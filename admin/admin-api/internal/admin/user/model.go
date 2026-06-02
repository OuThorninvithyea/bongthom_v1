package user

import (

	// Community pacakges
	"time"
)

// User maps to tbl_users
type User struct {
	ID           int64      `json:"id" db:"id"`
	FirstName    *string    `json:"first_name" db:"first_name"`
	LastName     *string    `json:"last_name" db:"last_name"`
	UserName     string     `json:"user_name" db:"user_name"`
	Email        *string    `json:"email" db:"email"`
	Password     string     `json:"-" db:"password"`
	RoleName     string     `json:"role_name" db:"role_name"`
	RoleID       int        `json:"role_id" db:"role_id"`
	IsAdmin      bool       `json:"is_admin" db:"is_admin"`
	LoginSession *string    `json:"-" db:"login_session"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
	CurrencyID   *int       `json:"currency_id" db:"currency_id"`
	LanguageID   *int       `json:"language_id" db:"language_id"`
	StatusID     *int       `json:"status_id" db:"status_id"`
	Order        *int       `json:"order" db:"order"`
	CreatedBy    *int64     `json:"-" db:"created_by"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedBy    *int64     `json:"-" db:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
	DeletedBy    *int64     `json:"-" db:"deleted_by"`
	DeletedAt    *time.Time `json:"-" db:"deleted_at"`
}

// CreateUserRequest — fields needed to create a user
type CreateUserRequest struct {
	UserName  string  `json:"user_name" validate:"required,min=4"`
	Password  string  `json:"password" validate:"required,min=6"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	RoleID    int     `json:"role_id" validate:"required"`
	RoleName  string  `json:"role_name" validate:"required"`
}

// UpdateUserRequest — partial update
type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	RoleID    *int    `json:"role_id"`
	RoleName  *string `json:"role_name"`
	StatusID  *int    `json:"status_id"`
}

// UserResponse — safe fields returned to client
type UserResponse struct {
	ID        int64      `json:"id"`
	UserName  string     `json:"user_name"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	Email     *string    `json:"email,omitempty"`
	RoleName  string     `json:"role_name"`
	RoleID    int        `json:"role_id"`
	IsAdmin   bool       `json:"is_admin"`
	StatusID  int        `json:"status_id"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

	type PagingRequest struct {
	Page    int `query:"page" validate:"min=1"`
	PerPage int `query:"per_page" validate:"min=1,max=100"`
}
