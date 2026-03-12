package middleware

import (
	"os"
	"svc-boilerplate-fiber-golang/utils/response"

	"github.com/gofiber/fiber/v2"
)

func InternalApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" || apiKey != os.Getenv("INTERNAL_API_KEY") {
			return response.Unauthorized(c, "Unauthorized")
		}
		return c.Next()
	}
}
