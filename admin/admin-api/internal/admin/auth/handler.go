package auth

import (
	// Commnuity Packages
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"

	// Internal Packages
	constants "admin-api/pkg/constants"
	response "admin-api/pkg/http"
	"admin-api/pkg/translate"
	"admin-api/pkg/utls"
)

type AuthHandler struct {
	Services AuthService
}

func NewAuthHandler(a *fiber.App, db *sqlx.DB) *AuthHandler {
	s := NewAuthServiceImpl(db)
	return &AuthHandler{
		Services: s,
	}
}

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

// func (a *AuthHandler) CheckStats(c fiber.Ctx) error  {
// 	if e != nil {
// 		c.Status(fiber.StatusInternalServerError)
// 		return  c.JSON(fiber.Map {
// 			"success": false,
// 			"message": "msg_failed",
// 			"status_code": 3000,
// 			"data": "The stats info is no showing. Please tryagain later",
// 		})
// 	} else {
// 		return c.JSON(fiber.Map {
// 			"success": true,
// 			"message": "Status are showing",
// 			"data":
// 		})
// 	}
// }
