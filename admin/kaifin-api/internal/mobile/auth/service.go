package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	jwtauth "kaifin-api/pkg/common/auth"
	"kaifin-api/pkg/logs"
	error_responses "kaifin-api/pkg/responses"
	types "kaifin-api/pkg/share"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(usreq *AuthRequest) (*AuthLoginReponse, *error_responses.ErrorResponse)
	CheckSession(loginSession string, userID int64) (bool, *error_responses.ErrorResponse)
	ForceLogout(userID int64) *error_responses.ErrorResponse
	Me(userContext types.UserContext) *AuthMeResponse
}

type AuthServiceImpl struct {
	Repo  AuthRepo
	Redis *redis.Client
}

func NewAuthServiceImpl(db *sqlx.DB, rdb *redis.Client) AuthService {
	r := NewAuthRepoImpl(db)
	return &AuthServiceImpl{
		Repo:  r,
		Redis: rdb,
	}
}

func NewAuthService(db *sqlx.DB, rdb *redis.Client) AuthService {
	return NewAuthServiceImpl(db, rdb)
}

func (s *AuthServiceImpl) Login(ureq *AuthRequest) (*AuthLoginReponse, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	user, err := s.Repo.Login(ureq)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ureq.Password)); err != nil {
		return nil, msg.NewErrorResponse("invalid_credentials", err)
	}

	loginSession := uuid.New().String()

	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		secret = "change-me-in-production"
	}

	jwtDuration := 15 * time.Minute
	if v := os.Getenv("JWT_EXPIRE"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			jwtDuration = d
		}
	}

	accessToken, _, jerr := jwtauth.GenerateToken(
		user.ID, user.UserName, loginSession, user.RoleID,
		secret, jwtDuration,
	)
	if jerr != nil {
		return nil, msg.NewErrorResponse("token_generation_failed", jerr)
	}

	if err := s.Redis.Set(context.Background(),
		fmt.Sprintf("session:%d", user.ID), loginSession,
		jwtDuration,
	).Err(); err != nil {
		logs.NewCustomLog("redis_session_set_failed", err.Error(), "warn")
	}

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

	key := fmt.Sprintf("session:%d", userID)
	stored, redisErr := s.Redis.Get(context.Background(), key).Result()
	if redisErr == nil && stored == loginSession {
		return true, nil
	}

	dbErr := s.Repo.CheckDatabaseLoginSession(userID, loginSession)
	if dbErr == nil {
		return true, nil
	}

	return false, msg.NewErrorResponse("invalid_session", fmt.Errorf("session mismatch"))
}

func (s *AuthServiceImpl) ForceLogout(userID int64) *error_responses.ErrorResponse {
	key := fmt.Sprintf("session:%d", userID)
	_ = s.Redis.Del(context.Background(), key).Err()
	return s.Repo.ClearLoginSession(userID)
}

func (s *AuthServiceImpl) Me(userContext types.UserContext) *AuthMeResponse {
	return &AuthMeResponse{
		User: AuthMeUser{
			UserID:   userContext.UserID,
			Username: userContext.UserName,
			RoleID:   userContext.RoleID,
		},
	}
}
