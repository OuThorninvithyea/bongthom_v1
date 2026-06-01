package auth

// HANDLER LAYER — HTTP only. No business logic. No SQL.
// Job: bind JSON → validate input → call service → translate errors → respond.
//
// Example: "POST /login {username, password} → is username ≥4 chars?
//          → call service.Login() → translate result → return JSON"

import (
	// Commnuity Packages
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	// Internal Packages
	constants "admin-api/pkg/constants"
	response "admin-api/pkg/http"
	"admin-api/pkg/translate"
	"admin-api/pkg/utls"
)

type AuthHandler struct {
	Services AuthService
}

func NewAuthHandler(a *fiber.App, db *sqlx.DB, rdb *redis.Client) *AuthHandler {
	s := NewAuthServiceImpl(db, rdb)
	return &AuthHandler{
		Services: s,
	}
}

// handler -> service -> repository

// payload = http body response // request
func (a *AuthHandler) Login(c fiber.Ctx) error {

	req := &AuthRequest{}
	// validator is dev package,
	v := utls.NewValidator()
	if err := req.bind(c, v); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fe := ve[0]
			msg, _ := translate.TranslateWithError(c, "validation_"+fe.Tag(),
				map[string]any{
					"Field": fe.Field(),
					"Param": fe.Param(),
				})

			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(msg, constants.Login_invalid, err),
			)
		}
		return err
	}

	rs, e := a.Services.Login(req.Username, req.Password)
	// conditions for database error
	if e != nil {
		msg, e_msg := translate.TranslateWithError(c, e.MessageID)
		if e_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(
					e_msg.Err.Error(),
					constants.Translate_Failed,
					e.Err,
				),
			)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(
				msg,
				constants.Login_failed,
				e.Err,
			),
		)

	} else {
		msg, e_msg := translate.TranslateWithError(c, "login_success")
		if e_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(
					e_msg.Err.Error(),
					constants.Translate_Failed,
					e_msg.Err,
				),
			)
		}
		return c.Status(fiber.StatusOK).JSON(
			response.NewResponse(msg, constants.Login_success, rs),
		)
	}

}
