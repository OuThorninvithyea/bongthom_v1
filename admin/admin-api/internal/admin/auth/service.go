package auth

import (
	error_responses "admin-api/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type AuthServiceImpl struct {
	Repo AuthRepo
}

type AuthService interface {
	Login(username string, password string) (*AuthLoginReponse, *error_responses.ErrorResponse)
}

func NewAuthServiceImpl(db *sqlx.DB) *AuthServiceImpl {
	r := NewAuthRepoImpl(db)
	return &AuthServiceImpl{
		Repo: r,
	}
}

func (s *AuthServiceImpl) Login(username string, password string) (*AuthLoginReponse, *error_responses.ErrorResponse) {
	return s.Repo.Login(username, password)
}
