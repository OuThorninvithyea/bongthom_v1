package auth

import (
	// Commnuity Packages
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"

	// Internal Packages
	"admin-api/pkg/utls"
)

type AuthHandler struct {
	Services AuthService
}

func NewAuthHandler(a *fiber.App, db *sqlx.DB,) *AuthHandler {
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
		return err
	}

	e, rs := a.Services.Login(req.Username, req.Password)
	// conditions for database error
	if e != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"success":     false,
			"message":     "msg_failed",
			"status_code": 3000,
			"data":        "The JWT could not be generated. Please try again later",
		})
	} else {
		return c.JSON(fiber.Map{
			"success":     true,
			"message":     "Login successful.",
			"status_code": 9000,
			"data":        AuthResponse{IsSucces: rs},
		})
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
