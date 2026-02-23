package middleware

import "github.com/gofiber/fiber/v2"

const (
	PUBLIC  string = "public"
	PRIVATE string = "private"
)

func PublicMiddleware(c *fiber.Ctx) error {
	c.Locals("requestSource", PUBLIC)
	return c.Next()
}

func PrivateMiddleware(c *fiber.Ctx) error {
	c.Locals("requestSource", PRIVATE)
	return c.Next()
}
