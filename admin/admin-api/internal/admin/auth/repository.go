package auth

import (
	// Community packages
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	// Internal packages
	jwtauth "admin-api/pkg/common/auth"
	error_responses "admin-api/pkg/responses"
)

type AuthRepo interface {
	Login(username string, password string) (*AuthLoginReponse, *error_responses.ErrorResponse)
}

type AuthRepoImpl struct {
	db *sqlx.DB
}

func NewAuthRepoImpl(db *sqlx.DB) *AuthRepoImpl {
	return &AuthRepoImpl{db: db}
}

func (r *AuthRepoImpl) Login(username string, password string) (*AuthLoginReponse, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	// Step 1: find the user
	var user Auth
	err := r.db.Get(&user,
		`SELECT id, user_name, role_id
		 FROM tbl_users
		 WHERE user_name = $1 AND password = $2
		 LIMIT 1`,
		username, password,
	)
	if err != nil {
		return nil, msg.NewErrorResponse("invalid_credentials", err)
	}

	// Step 2: generate login session UUID
	loginSession := uuid.New().String()

	// Step 3: read JWT secret
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		secret = "change-me-in-production"
	}

	// Step 4: generate access token (24h, self-validating, no DB storage)
	accessToken, _, err := jwtauth.GenerateToken(
		user.ID, user.UserName, loginSession, user.RoleId,
		secret, 24*time.Hour,
	)
	if err != nil {
		return nil, msg.NewErrorResponse("token_generation_failed", err)
	}

	// Step 5: update last_login on user record
	_, _ = r.db.Exec(
		`UPDATE tbl_users SET login_session = $1, last_login = NOW() WHERE id = $2`,
		loginSession, user.ID,
	)

	var au AuthLoginReponse
	au.Auth.Token = accessToken
	au.Auth.TokenType = "jwt"
	return &au, nil
}
