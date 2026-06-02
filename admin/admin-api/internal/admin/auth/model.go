package auth

import (
	// Commnunity Pacakges
	"time"
	
	"github.com/gofiber/fiber/v3"

	// Internal Packages
	"admin-api/pkg/utls"
)

// Auth maps to tbl_users
type Auth struct {
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
	Order        int        `json:"order" db:"order"`
	CreatedBy    int        `json:"created_by" db:"created_by"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedBy    int        `json:"updated_by" db:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
	DeletedBy    int        `json:"deleted_by" db:"deleted_by"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required"`
}

// functions bind purpose is to be able to add 2+ specail validate , or over standard,
// icase there is more to validated with database......
func (r *AuthRequest) bind(c fiber.Ctx, v *utls.Validator) error {
	if err := c.Bind().Body(&r); err != nil {
		return err
	}
	if err := v.Validate(r); err != nil {
		return err
	}
	return nil
}

type AuthResponse struct {
	// Tag field (is basicly functions)
	IsSucces bool   `json:"is_success"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthLoginReponse struct {
	Auth struct {
		Token     string `json:"token"`
		TokenType string `json:"token_type"`
	} `json:"auth"`
}
