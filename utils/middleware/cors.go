package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSConfig() fiber.Handler {
	allowOrigins := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS"))
	if allowOrigins == "" {
		allowOrigins = "http://localhost:5173"
	}

	allowMethods := strings.TrimSpace(os.Getenv("CORS_ALLOW_METHODS"))
	if allowMethods == "" {
		allowMethods = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	}

	allowHeaders := strings.TrimSpace(os.Getenv("CORS_ALLOW_HEADERS"))
	if allowHeaders == "" {
		allowHeaders = "Origin,Content-Type,Accept,Authorization,X-API-Key"
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: strings.EqualFold(os.Getenv("CORS_ALLOW_CREDENTIALS"), "true"),
	})
}
