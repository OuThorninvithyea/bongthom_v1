package middlewares

import (

	// Commnuity packages
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// Internal packages
	"admin-api/internal/admin/auth"
	response "admin-api/pkg/http"
	types "admin-api/pkg/share"
	"admin-api/pkg/utls"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func NewJwtMiddleware(app *fiber.App, db_pool *sqlx.DB, redis *redis.Client) {
	errs := godotenv.Load()
	if errs != nil {
		log.Println("No .env file, using environment variables")
	}
	secret_key := os.Getenv("JWT_SECRET_KEY")

	app.Use(func(c fiber.Ctx) error {
		if websocketUpgrade := c.Get("Upgrade"); websocketUpgrade == "websocket" {
			webSocketProtocol := c.Get("Sec-WebSocket-Protocol")
			if webSocketProtocol == "" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Missing WebSocket protocol for authentication",
				})
			}

			parts := strings.Split(webSocketProtocol, ",")
			if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Bearer" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid WebSocket protocol authentication format",
				})
			}

			tokenString := strings.TrimSpace(parts[1])

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret_key), nil
			})
			if err != nil || !token.Valid {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid or expired JWT token",
				})
			}

			c.Locals("jwt_data", token)
			c.Set("Sec-WebSocket-Protocol", "Bearer")
			return c.Next()
		}

		// Skip auth for login
		if c.Path() == "/api/v1/admin/auth/login" {
			return c.Next()
		}

		return jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(secret_key)},
			SuccessHandler: func(c fiber.Ctx) error {
				if token := jwtware.FromContext(c); token != nil {
					c.Locals("jwt_data", token)
				}
				return c.Next()
			},
		})(c)
	})

	app.Use(func(c fiber.Ctx) error {
		// Skip session check for login
		if c.Path() == "/api/v1/admin/auth/login" {
			return c.Next()
		}

		user_token := c.Locals("jwt_data").(*jwt.Token)
		uclaim := user_token.Claims.(jwt.MapClaims)

		return handleUserContext(c, uclaim, db_pool, redis)
	})
}

func handleUserContext(c fiber.Ctx, uclaim jwt.MapClaims, db *sqlx.DB, redis *redis.Client) error {

	login_session, ok := uclaim["login_session"].(string)
	if !ok || login_session == "" {
		smg_error := response.NewResponseError(
			utls.Translate("loginSessionMissing", nil, c),
			-500,
			fmt.Errorf("%s", utls.Translate("loginSessionMissing", nil, c)),
		)
		return c.Status(http.StatusUnprocessableEntity).JSON(smg_error)
	}

	uCtx := types.UserContext{
		UserID:       int64(uclaim["user_id"].(float64)),
		UserName:     uclaim["user_name"].(string),
		RoleID:       int(uclaim["role_id"].(float64)),
		LoginSession: uclaim["login_session"].(string),
		Exp:          time.Unix(int64(uclaim["exp"].(float64)), 0),
		UserAgent:    c.Get("User-Agent", "unknown"),
		IP:           c.IP(),
	}
	c.Locals("UserContext", uCtx)

	sv := auth.NewAuthService(db, redis)
	success, err := sv.CheckSession(login_session, uCtx.UserID)
	if err != nil || !success {
		smg_error := response.NewResponseError(
			utls.Translate("loginSessionInvalid", nil, c),
			-500,
			fmt.Errorf("%s", utls.Translate("loginSessionInvalid", nil, c)),
		)
		return c.Status(http.StatusUnprocessableEntity).JSON(smg_error)
	}

	return c.Next()
}
