package auth

import (
	// Community Pacakges
	"github.com/jmoiron/sqlx"
)

type AuthRepo interface {
	Login(username string, password string) (error, bool)
}

type AuthRepoImpl struct {
	db *sqlx.DB
}

func NewAuthRepoImpl(db *sqlx.DB) *AuthRepoImpl {
	return &AuthRepoImpl{db: db}
}

func (r *AuthRepoImpl) Login(username string, password string) (error, bool) {
	var count int
	// db is pool connecitons
	err := r.db.Get(
		&count,
		"SELECT COUNT(*) FROM users WHERE username = $1 AND password = $2",
		username, password,
	)
	if err != nil {
		return err, false
	}
	return nil, count > 0
}
