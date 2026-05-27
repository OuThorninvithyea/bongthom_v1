package auth

import (
	// Commnunity Pacakges
	"github.com/gofiber/fiber/v3"

	// Internal Packages
	"admin-api/pkg/utls"
)

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
