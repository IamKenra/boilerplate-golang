package middleware

import (
	"os"
	"strings"
	"boilerplate-golang/utils/logger"
	"boilerplate-golang/utils/response"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.Error("Authorization header is missing")
			return response.Unauthorized(c, "Authorization header is missing")
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Error("Authorization header format is invalid")
			return response.Unauthorized(c, "Authorization header format is invalid")
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired JWT",
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("uuid", claims["uuid"])
		c.Locals("username", claims["username"])

		return c.Next()
	}
}
