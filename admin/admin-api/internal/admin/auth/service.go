package auth

import (
	//Community pacakges
	"github.com/jmoiron/sqlx"
	error_responses "admin-api/pkg/responses"
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

func (s *AuthServiceImpl) Login(username, password string) (*AuthLoginReponse, *error_responses.ErrorResponse) {
	return s.Repo.Login(username, password)
}
