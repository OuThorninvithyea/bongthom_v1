package auth

import "github.com/jmoiron/sqlx"

type AuthServiceImpl struct {
	Repo AuthRepo
}

type AuthService interface {
	Login(username string, password string) (error, bool)
}

func NewAuthServiceImpl(db *sqlx.DB) *AuthServiceImpl {
	r := NewAuthRepoImpl(db)
	return &AuthServiceImpl{
		Repo: r,
	}
}

func (s *AuthServiceImpl) Login(username, password string) (error, bool) {
	return s.Repo.Login(username, password)
}
