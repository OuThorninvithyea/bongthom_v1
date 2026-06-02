package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"admin-api/internal/admin/auth"
	jwtauth "admin-api/pkg/common/auth"
	response "admin-api/pkg/http"
	types "admin-api/pkg/share"
	"admin-api/pkg/utls"
)

// NewJwtMiddleware validates every request using our ValidateToken
func NewJwtMiddleware(app *fiber.App, db *sqlx.DB, redis *redis.Client) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		panic("JWT_SECRET_KEY not set in .env")
	}

	app.Use(func(c fiber.Ctx) error {
		// WebSocket auth — extract from Sec-WebSocket-Protocol
		if c.Get("Upgrade") == "websocket" {
			protocol := c.Get("Sec-WebSocket-Protocol")
			if protocol == "" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Missing WebSocket protocol for authentication",
				})
			}

			parts := strings.Split(protocol, ",")
			if len(parts) < 2 {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid WebSocket auth format",
				})
			}

			tokenString := strings.TrimPrefix(strings.TrimSpace(parts[1]), "Bearer ")
			claims, err := jwtauth.ValidateToken(tokenString, secret)
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid or expired token",
				})
			}

			return handleContext(c, claims, db, redis)
		}

		// HTTP auth — extract from Authorization header
		if c.Path() == "/api/v1/admin/auth/login" {
			return c.Next()
		}

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

		return handleContext(c, claims, db, redis)
	})
}

func handleContext(c fiber.Ctx, claims *jwtauth.Claims, db *sqlx.DB, redis *redis.Client) error {
	// Validate login session exists
	if claims.LoginSession == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(
			response.NewResponseError(
				utls.Translate("loginSessionMissing", nil, c),
				4003,
				fmt.Errorf("missing login session"),
			),
		)
	}

	// Build user context from JWT claims
	uCtx := types.UserContext{
		UserID:       claims.UserID,
		UserName:     claims.UserName,
		LoginSession: claims.LoginSession,
		RoleID:       claims.RoleID,
		UserAgent:    c.Get("User-Agent", "unknown"),
		Ip:           c.IP(),
	}

	// Check session via Redis (sub-ms)
	sv := auth.NewAuthServiceImpl(db, redis)
	if _, err := sv.CheckSession(uCtx.LoginSession, uCtx.UserID); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success":     false,
			"message":     "Session expired or invalid",
			"status_code": 4003,
		})
	}

	// Inject user context for downstream handlers
	c.Locals("UserContext", uCtx)
	return c.Next()
}
