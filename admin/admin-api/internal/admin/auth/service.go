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
	"golang.org/x/crypto/bcrypt"

	// internal pacakges
	jwtauth "admin-api/pkg/common/auth"
	"admin-api/pkg/logs"
	error_responses "admin-api/pkg/responses"
)

type AuthService interface {
	Login(username string, password string) (*AuthLoginReponse, *error_responses.ErrorResponse)
	CheckSession(loginSession string, userID int64) (bool, *error_responses.ErrorResponse)
}

type AuthServiceImpl struct {
	Repo  AuthRepo
	Redis *redis.Client
}

func NewAuthServiceImpl(db *sqlx.DB, rdb *redis.Client) *AuthServiceImpl {
	r := NewAuthRepoImpl(db)
	return &AuthServiceImpl{
		Repo:  r,
		Redis: rdb,
	}
}

func NewAuthService(db *sqlx.DB, rdb *redis.Client) AuthService {
	return NewAuthServiceImpl(db, rdb)
}

func (s *AuthServiceImpl) Login(username string, password string) (*AuthLoginReponse, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	// Step 1: find user (repo only does DB work)
	user, err := s.Repo.Login(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, msg.NewErrorResponse("invalid_credentials", err)
	}

	// Step 2: generate login session UUID
	loginSession := uuid.New().String()

	// Step 3: read JWT secret and expiry
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		secret = "change-me-in-production"
	}

	jwtDuration := 15 * time.Minute // default
	if v := os.Getenv("JWT_EXPIRE"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			jwtDuration = d
		}
	}

	// Step 5: generate access token (business logic — lives in service)
	accessToken, _, jerr := jwtauth.GenerateToken(
		user.ID, user.UserName, user.RoleID, loginSession,
		secret, jwtDuration,
	)
	if jerr != nil {
		return nil, msg.NewErrorResponse("token_generation_failed", jerr)
	}

	// Step 5: store login session in Redis (TTL matches JWT expiry)
	if err := s.Redis.Set(context.Background(),
		fmt.Sprintf("session:%d", user.ID), loginSession,
		jwtDuration,
	).Err(); err != nil {
		logs.NewCustomLog("redis_session_set_failed", err.Error(), "warn")
	}

	// Step 6: store login session in database (audit trail)
	if err := s.Repo.UpdateLoginSession(user.ID, loginSession); err != nil {
		return nil, err
	}

	var au AuthLoginReponse
	au.Auth.Token = accessToken
	au.Auth.TokenType = "jwt"
	return &au, nil
}

func (s *AuthServiceImpl) CheckSession(loginSession string, userID int64) (bool, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	// Fast path: Redis
	key := fmt.Sprintf("session:%d", userID)
	stored, redisErr := s.Redis.Get(context.Background(), key).Result()
	if redisErr == nil && stored == loginSession {
		return true, nil
	}

	// Slow path: PostgreSQL fallback (Redis down or key missing)
	dbErr := s.Repo.CheckDatabaseLoginSession(userID, loginSession)
	if dbErr == nil {
		return true, nil
	}

	return false, msg.NewErrorResponse("invalid_session", fmt.Errorf("session mismatch"))
}
