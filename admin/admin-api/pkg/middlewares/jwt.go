package middlewares

import (

	// Community pacakges
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	// Internal pacakges
	"admin-api/internal/admin/auth"
	jwtauth "admin-api/pkg/common/auth"
	response "admin-api/pkg/http"
	types "admin-api/pkg/share"
	"admin-api/pkg/utls"
)

func NewJwtMiddleware(app *fiber.App, db *sqlx.DB, redis *redis.Client) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		panic("JWT_SECRET_KEY not set")
	}

	app.Use(func(c fiber.Ctx) error {
		// WebSocket auth
		if c.Get("Upgrade") == "websocket" {
			protocol := c.Get("Sec-WebSocket-Protocol")
			if protocol == "" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Missing WebSocket protocol for authentication",
				})
			}

			parts := strings.Split(protocol, ",")
			if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Bearer" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid WebSocket protocol authentication format",
				})
			}

			tokenString := strings.TrimSpace(parts[1])

			claims, err := jwtauth.ValidateToken(tokenString, secret)
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid or expired JWT token",
				})
			}

			c.Locals("jwt_claims", claims)
			c.Set("Sec-WebSocket-Protocol", "Bearer")

			return handleUserContext(c, claims, db, redis)
		}

		// Skip auth for login
		if c.Path() == "/api/v1/admin/auth/login" {
			return c.Next()
		}

		// HTTP auth
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"success":     false,
				"message":     "Missing or invalid Authorization header",
				"status_code": 4001,
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtauth.ValidateToken(tokenString, secret)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"success":     false,
				"message":     "Invalid or expired token",
				"status_code": 4002,
			})
		}

		c.Locals("jwt_claims", claims)
		return handleUserContext(c, claims, db, redis)
	})
}

func handleUserContext(c fiber.Ctx, claims *jwtauth.Claims, db *sqlx.DB, redis *redis.Client) error {
	loginSession := claims.LoginSession
	if loginSession == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(
			response.NewResponseError(
				utls.Translate("loginSessionMissing", nil, c),
				4003,
				fmt.Errorf("missing login session"),
			),
		)
	}

	uCtx := types.UserContext{
		UserID:       claims.UserID,
		UserName:     claims.UserName,
		RoleID:       claims.RoleID,
		LoginSession: claims.LoginSession,
		Exp:          time.Now(),
		UserAgent:    c.Get("User-Agent", "unknown"),
		Ip:           c.IP(),
	}
	c.Locals("UserContext", uCtx)

	sv := auth.NewAuthService(db, redis)
	if success, err := sv.CheckSession(loginSession, claims.UserID); err != nil || !success {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success":     false,
			"message":     "Session expired or invalid",
			"status_code": 4003,
		})
	}

	return c.Next()
}
