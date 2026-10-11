package middleware

import (
	"strings"
	jwtutil "svc-boilerplate-fiber-golang/utils/jwt"
	"svc-boilerplate-fiber-golang/utils/logger"
	"svc-boilerplate-fiber-golang/utils/response"

	"github.com/gofiber/fiber/v2"
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
		claims, err := jwtutil.ValidateAccessToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired JWT",
			})
		}

		c.Locals("user_uuid", claims.UserUUID)
		c.Locals("username", claims.Username)
		c.Locals("authority", claims.Authority)

		return c.Next()
	}
}
