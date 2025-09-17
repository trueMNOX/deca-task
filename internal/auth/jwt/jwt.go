package jwt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatalf("failed to get JWT_SECRET: environment variable is not set")
	}
	return []byte(secret)
}

func TokenTTL() time.Duration {
	if v := os.Getenv("JWT_EXPIRE_IN"); v != "" {
		if m, err := strconv.Atoi(v); err == nil && m > 0 {
			return time.Duration(m) * time.Minute
		}
	}
	return 24 * time.Hour
}

func GenerateToken(userId uint) (string, error) {
	now := time.Now()
	ttl := TokenTTL()

	claims := jwt.MapClaims{
		"user_id": userId,
		"iat":     jwt.NewNumericDate(now),
		"exp":     jwt.NewNumericDate(now.Add(ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JwtSecret())
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signedToken, nil
}
func ParseToken(tokenString string) (*jwt.MapClaims, error) {
	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	if tokenString == "" {
		return nil, errors.New("empty token")
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return JwtSecret(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return &claims, nil
}

func GetUserId(claims jwt.MapClaims) (uint, error) {
	v, ok := claims["user_id"]
	if !ok {
		return 0, errors.New("user_id claim not found")
	}
	switch t := v.(type) {
	case float64:
		return uint(t), nil
	case float32:
		return uint(t), nil
	case int:
		return uint(t), nil
	case int64:
		return uint(t), nil
	case string:
		if u, err := strconv.ParseUint(t, 10, 64); err == nil {
			return uint(u), nil
		} else {
			return 0, fmt.Errorf("invalid user_id string: %w", err)
		}
	default:
		return 0, fmt.Errorf("unexpected user_id type %T", v)
	}
}
