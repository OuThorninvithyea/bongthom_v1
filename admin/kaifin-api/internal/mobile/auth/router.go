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
	a.Post("/api/v1/mobile/auth/login", h.Login)
	return &AuthRoute{
		handler: h,
	}
}

func NewProtectedAuthRoute(a *fiber.App, db *sqlx.DB, rdb *redis.Client) *AuthRoute {
	h := NewAuthHandler(a, db, rdb)
	a.Get("/api/v1/mobile/auth/me", h.Me)
	return &AuthRoute{
		handler: h,
	}
}
