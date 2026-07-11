package user

import (

	// Community pacakges
	"kaifin-api/pkg/share"
	"kaifin-api/pkg/utls"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
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

type UserShowRequest struct {
	PageOption share.Paging   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts      []share.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters    []share.Filter `json:"filters,omitempty" query:"filters"`
	Search     string         `json:"q,omitempty" query:"q"`
	CurrencyID int            `json:"currency_id,omitempty" query:"currency_id"`
}

func (u *UserShowRequest) bind(c fiber.Ctx, v *utls.Validator) error {

	if err := c.Bind().Query(u); err != nil {
		return err
	}

	for i := range u.Filters {
		// request takes &filters[index][value]= int,
		value := c.Query(fmt.Sprintf("filters[%d][value]", i))
		if intValue, err := strconv.Atoi(value); err == nil {
			u.Filters[i].Value = intValue
		} else if boolValue, err := strconv.ParseBool(value); err == nil {
			u.Filters[i].Value = boolValue
		} else {
			u.Filters[i].Value = value
		}
	}

	if u.Search == "" {
		u.Search = c.Query("q")
	}
	if u.CurrencyID == 0 {
		if v := strings.TrimSpace(c.Query("currency_id")); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				u.CurrencyID = n
			}
		}
	}

	if err := v.Validate(u); err != nil {
		return err
	}
	return nil
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
	Users []User `json:"users"`
	Total int
}

func (r *CreateUserRequest) bind(c fiber.Ctx, v *utls.Validator) error {
	if err := c.Bind().Body(&r); err != nil {
		return err
	}
	if err := v.Validate(r); err != nil {
		return err
	}
	return nil
}

type PagingRequest struct {
	Page    int `query:"page" validate:"min=1"`
	PerPage int `query:"per_page" validate:"min=1,max=100"`
}
