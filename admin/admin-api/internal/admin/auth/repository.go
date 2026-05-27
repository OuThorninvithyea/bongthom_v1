package auth

import (
	// Community Pacakges
	"fmt"
	
	// Internal packages
	error_responses "admin-api/pkg/responses"
	"github.com/jmoiron/sqlx"
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
	var count int
	msg := error_responses.ErrorResponse{}
	// db is pool connecitons
	err := r.db.Get(
		&count,
		"SELECT COUNT(*) FROM users WHERE username = $1 AND password = $2",
		username, password,
	)
	
	if err != nil || count == 0 {
		return nil, msg.NewErrorResponse("user_not_found", fmt.Errorf("Failed to generate UUID. Please try again later."))
	}
	
	var au AuthLoginReponse
	au.Auth.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
	au.Auth.TokenType = "jwt"
	return &au, nil
}
