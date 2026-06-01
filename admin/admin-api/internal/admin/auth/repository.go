package auth

// REPOSITORY LAYER — ONLY talks to the database.
// No business logic. No HTTP. No token generation.
// Job: accept parameters → run SQL → return results.
//
// Example: "Find user WHERE user_name = 'admin' AND password = 'admin123'"
//          → SELECT → returns *Auth or error

import (

	// Commnunity pacakges
	"github.com/jmoiron/sqlx"

	// Internal pacakges
	error_responses "admin-api/pkg/responses"
)

type AuthRepo interface {
	Login(username string, password string) (*Auth, *error_responses.ErrorResponse)
	UpdateLoginSession(userID int64, loginSession string) *error_responses.ErrorResponse
}

type AuthRepoImpl struct {
	db *sqlx.DB
}

func NewAuthRepoImpl(db *sqlx.DB) AuthRepo {
	return &AuthRepoImpl{db: db}
}

func (r *AuthRepoImpl) Login(username string, password string) (*Auth, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

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

	return &user, nil
}

func (r *AuthRepoImpl) UpdateLoginSession(userID int64, loginSession string) *error_responses.ErrorResponse {
	msg := error_responses.ErrorResponse{}
	_, err := r.db.Exec(
		`UPDATE tbl_users SET login_session = $1, last_login = NOW() WHERE id = $2`,
		loginSession, userID,
	)
	if err != nil {
		return msg.NewErrorResponse("database_error", err)
	}
	return nil
}
