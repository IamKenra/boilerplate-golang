package auth

import (
	"errors"
	"strings"

	"boilerplate-golang/utils/jwt"

	"github.com/gofiber/fiber/v2"
)

func ValidateAccessToken(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return errors.New("authorization header is required")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return errors.New("invalid authorization header format")
	}
	_, err := jwt.ValidateAccessToken(tokenParts[1])
	if err != nil {
		return errors.New("invalid or expired access token")
	}

	return nil
}

func GetUserFromXMember(ctx *fiber.Ctx) (string, error) {
	xMember := ctx.Get("X-Member")
	if xMember == "" {
		return "", errors.New("X-Member header is required")
	}
	return xMember, nil
}
