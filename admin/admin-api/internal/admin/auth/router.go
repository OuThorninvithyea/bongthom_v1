package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type AuthRoute struct {
	handler *AuthHandler
}

func NewAuthRoute(a *fiber.App, db *sqlx.DB, rdb *redis.Client) *AuthRoute {
	h := NewAuthHandler(a, db, rdb)
	a.Post("api/v1/admin/auth/logins", h.Login)
	return &AuthRoute{
		handler: h,
	}
}
