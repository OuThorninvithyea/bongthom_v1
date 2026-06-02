package user

import (

	// Community pacakges
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	error_responses "admin-api/pkg/responses"
)

type UserRepo interface {
	FindAll(page, perPage int) ([]User, int, *error_responses.ErrorResponse)
	FindByID(id int64) (*User, *error_responses.ErrorResponse)
	FindByUserName(userName string) (*User, *error_responses.ErrorResponse)
	Create(user *User) *error_responses.ErrorResponse
	Update(id int64, updates map[string]any) (*User, *error_responses.ErrorResponse)
	SoftDelete(id int64, deletedBy int64) *error_responses.ErrorResponse
}

type UserRepoImpl struct {
	db *sqlx.DB
}

func NewUserRepoImpl(db *sqlx.DB) UserRepo {
	return &UserRepoImpl{db: db}
}

func (r *UserRepoImpl) FindAll(page, perPage int) ([]User, int, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	var total int
	r.db.Get(&total, `SELECT COUNT(*) FROM tbl_users WHERE deleted_at IS NULL`)

	offset := (page - 1) * perPage
	var users []User
	err := r.db.Select(&users,
		`SELECT id, user_name, first_name, last_name, email, role_name, role_id, is_admin,
		 login_session, last_login, currency_id, language_id, status_id, created_at, updated_at
		 FROM tbl_users
		 WHERE deleted_at IS NULL
		 ORDER BY id ASC
		 LIMIT $1 OFFSET $2`,
		perPage, offset,
	)
	if err != nil {
		return nil, 0, msg.NewErrorResponse("database_error", err)
	}

	return users, total, nil
}

func (r *UserRepoImpl) FindByID(id int64) (*User, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	var user User
	err := r.db.Get(&user,
		`SELECT * FROM tbl_users WHERE id = $1 AND deleted_at IS NULL LIMIT 1`, id,
	)
	if err != nil {
		return nil, msg.NewErrorResponse("user_not_found", err)
	}
	return &user, nil
}

func (r *UserRepoImpl) FindByUserName(userName string) (*User, *error_responses.ErrorResponse) {
	var user User
	err := r.db.Get(&user,
		`SELECT * FROM tbl_users WHERE user_name = $1 LIMIT 1`, userName,
	)
	if err != nil {
		return nil, nil // not found is OK — caller checks
	}
	return &user, nil
}

func (r *UserRepoImpl) Create(user *User) *error_responses.ErrorResponse {
	msg := error_responses.ErrorResponse{}

	err := r.db.QueryRow(
		`INSERT INTO tbl_users (first_name, last_name, user_name, email, password, role_name, role_id, is_admin, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, created_at`,
		user.FirstName, user.LastName, user.UserName, user.Email,
		user.Password, user.RoleName, user.RoleID, user.IsAdmin, user.CreatedBy,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return msg.NewErrorResponse("database_error", err)
	}
	return nil
}

func (r *UserRepoImpl) Update(id int64, updates map[string]any) (*User, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	setClauses := []string{}
	args := []any{}
	i := 1
	for col, val := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}

	if len(setClauses) == 0 {
		return nil, msg.NewErrorResponse("no_updates_provided", fmt.Errorf("empty update"))
	}

	query := fmt.Sprintf(
		`UPDATE tbl_users SET %s WHERE id = $%d AND deleted_at IS NULL RETURNING *`,
		strings.Join(setClauses, ", "), i,
	)
	args = append(args, id)

	var user User
	err := r.db.Get(&user, query, args...)
	if err != nil {
		return nil, msg.NewErrorResponse("database_error", err)
	}
	return &user, nil
}

func (r *UserRepoImpl) SoftDelete(id int64, deletedBy int64) *error_responses.ErrorResponse {
	msg := error_responses.ErrorResponse{}

	result, err := r.db.Exec(
		`UPDATE tbl_users SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`,
		deletedBy, id,
	)
	if err != nil {
		return msg.NewErrorResponse("database_error", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return msg.NewErrorResponse("user_not_found", fmt.Errorf("user %d not found", id))
	}
	return nil
}
