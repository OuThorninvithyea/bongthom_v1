package auth

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
	types "admin-api/pkg/share"
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

	rs, e := a.Services.Login(req)
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

// checking sessions user is its still valid for frontend validations and flutter

func (a *AuthHandler) Me(c fiber.Ctx) error {
	userContext, ok := c.Locals("UserContext").(types.UserContext)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.NewResponseError(
				"Unauthorized",
				constants.Generic_invalid,
				nil,
			),
		)
	}
	data := a.Services.Me(userContext)

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse("Session valid", constants.Generic_success, data),
	)

}
