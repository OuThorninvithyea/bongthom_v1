package auth

import (
	// Community Packages
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

type AuthRoute struct {
	Handler *AuthHandler
}

func NewAuthRoute(a *fiber.App, db *sqlx.DB) *AuthRoute {
	h := NewAuthHandler(a, db)
	a.Post("api/v1/admin/auth/login", h.Login)
	return &AuthRoute{
		Handler: h,
	}
}
