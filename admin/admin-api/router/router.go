package router

import (

	// Community pacakges
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

const MaxUploadBytes = 500 * 1024 * 1024 // 500 MB

func New() *fiber.App {
	f := fiber.New(fiber.Config{
		BodyLimit: MaxUploadBytes,
	})
	f.Use(logger.New()) // "Use" is an extentions from other libraries
	f.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "Accept-Language"},
		AllowMethods: []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},
	}))
	return f
}
