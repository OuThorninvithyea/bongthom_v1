package auth

import (

	// Community pacakges
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	// internal pacakges
	jwtauth "admin-api/pkg/common/auth"
	error_responses "admin-api/pkg/responses"
)

type AuthServiceImpl struct {
	Repo  AuthRepo
	Redis *redis.Client
}

type AuthService interface {
	Login(username string, password string) (*AuthLoginReponse, *error_responses.ErrorResponse)
	CheckRedisSession(loginSession string, userID int64) (*Auth, *error_responses.ErrorResponse)
}

func NewAuthServiceImpl(db *sqlx.DB, rdb *redis.Client) *AuthServiceImpl {
	r := NewAuthRepoImpl(db)
	return &AuthServiceImpl{
		Repo:  r,
		Redis: rdb,
	}
}

func (s *AuthServiceImpl) Login(username string, password string) (*AuthLoginReponse, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	// Step 1: find user (repo only does DB work)
	user, err := s.Repo.Login(username, password)
	if err != nil {
		return nil, err
	}

	// Step 2: generate login session UUID
	loginSession := uuid.New().String()

	// Step 3: read JWT secret
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		secret = "change-me-in-production"
	}

	// Step 5: generate access token (business logic — lives in service)
	accessToken, _, jerr := jwtauth.GenerateToken(
		user.ID, user.UserName, user.RoleID, loginSession,
		secret, 15*time.Minute,
	)
	if jerr != nil {
		return nil, msg.NewErrorResponse("token_generation_failed", jerr)
	}

	// Step 5: store login session in Redis
	s.Redis.Set(context.Background(),
		fmt.Sprintf("session:%d", user.ID), loginSession,
		0, // no expiry — lives until manual delete
	)

	// Step 6: store login session in database (audit trail)
	if err := s.Repo.UpdateLoginSession(user.ID, loginSession); err != nil {
		return nil, err
	}

	var au AuthLoginReponse
	au.Auth.Token = accessToken
	au.Auth.TokenType = "jwt"
	return &au, nil
}

func (s *AuthServiceImpl) CheckRedisSession(loginSession string, userID int64) (*Auth, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	// Fast path: Redis
	key := fmt.Sprintf("session:%d", userID)
	stored, redisErr := s.Redis.Get(context.Background(), key).Result()
	if redisErr == nil && stored == loginSession {
		return &Auth{ID: userID}, nil
	}

	// Slow path: PostgreSQL fallback (Redis down or key missing)
	user, dbErr := s.Repo.CheckDatabaseLoginSession(userID, loginSession)
	if dbErr == nil {
		return user, nil
	}

	return nil, msg.NewErrorResponse("invalid_session", fmt.Errorf("session mismatch"))
}
