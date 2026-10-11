package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"svc-boilerplate-fiber-golang/valueobject"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserUUID  string `json:"user_uuid"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Authority string `json:"authority"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user valueobject.User) (string, int64, error) {
	secret := getJWTSecret()
	if secret == "" {
		return "", 0, errors.New("JWT_SECRET is not configured")
	}
	duration := getAccessTokenDuration()
	expirationTime := time.Now().Add(duration)

	claims := &Claims{
		UserUUID:  user.UUID,
		Username:  user.Username,
		Name:      user.Identity.Name,
		Authority: user.Authority,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(duration.Seconds()), nil
}

func GenerateRefreshTokenExpiry() time.Time {
	duration := getRefreshTokenDuration()
	return time.Now().Add(duration)
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	secret := getJWTSecret()
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func getJWTSecret() string {
	return strings.TrimSpace(os.Getenv("JWT_SECRET"))
}

func getAccessTokenDuration() time.Duration {
	durationStr := os.Getenv("ACCESS_TOKEN_DURATION")
	if durationStr == "" {
		return 15 * time.Minute
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 15 * time.Minute
	}

	return duration
}

func getRefreshTokenDuration() time.Duration {
	durationStr := os.Getenv("REFRESH_TOKEN_DURATION")
	if durationStr == "" {
		return 7 * 24 * time.Hour // 7 days
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 7 * 24 * time.Hour
	}

	return duration
}

// GenerateRandomToken generates a random token for refresh token
func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken creates a SHA256 hash of the given token
func HashToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
