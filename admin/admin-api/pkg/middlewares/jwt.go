package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"admin-api/internal/admin/auth"
	jwtauth "admin-api/pkg/common/auth"
	types "admin-api/pkg/share"
)

// skipAuthRoutes — public endpoints that don't need a token
var skipAuthRoutes = []string{
	"/api/v1/admin/auth/login",
	"/api/v1/admin/auth/refresh",
}

func shouldSkip(path string) bool {
	for _, r := range skipAuthRoutes {
		if strings.HasPrefix(path, r) {
			return true
		}
	}
	return false
}

// NewJwtMiddleware validates every request using our ValidateToken
func NewJwtMiddleware(app *fiber.App, db *sqlx.DB, rdb *redis.Client) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		secret = "change-me-in-production"
	}

	// Single middleware for both HTTP and WebSocket
	app.Use(func(c fiber.Ctx) error {
		// Skip auth for public routes
		if shouldSkip(c.Path()) {
			return c.Next()
		}

		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success":     false,
				"message":     "Missing or invalid Authorization header",
				"status_code": 4001,
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate using our JWT lib (HMAC check + signature + expiry)
		claims, err := jwtauth.ValidateToken(tokenString, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success":     false,
				"message":     "Invalid or expired token",
				"status_code": 4002,
			})
		}

		// Build UserContext from claims
		uCtx := types.UserContext{
			UserID:       claims.UserID,
			UserName:     claims.UserName,
			LoginSession: claims.LoginSession,
			RoleID:       claims.RoleID,
			UserAgent:    c.Get("User-Agent", "unknown"),
			Ip:           c.IP(),
		}

		// Check session via Redis (sub-ms)
		sv := auth.NewAuthServiceImpl(db, rdb)
		if _, err := sv.CheckSession(uCtx.LoginSession, uCtx.UserID); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success":     false,
				"message":     "Session expired or invalid",
				"status_code": 4003,
			})
		}

		// Inject user context for downstream handlers
		c.Locals("UserContext", uCtx)
		return c.Next()
	})
}
