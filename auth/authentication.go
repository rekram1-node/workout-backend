package auth

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/rekram1-node/workout-backend/models"
)

type JwtCustomClaims struct {
	Name string
	UUID string
	jwt.RegisteredClaims
}

func CreateAccessToken(user *models.User, secret string) (string, error) {
	exp := time.Now().Add(time.Hour * 72)
	claims := &JwtCustomClaims{
		user.Username,
		user.UUID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func CreateRefreshToken(user *models.User, secret string) (string, error) {
	exp := time.Now().Add(time.Hour * 72)
	claims := &JwtCustomClaims{
		user.Username,
		user.UUID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func IsAuthorized(token string, secret string) (bool, error) {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func ReadUUIDFromToken(token string, secret string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token %w", err)
	}

	if claims, ok := t.Claims.(*JwtCustomClaims); ok && t.Valid {
		return claims.UUID, nil
	}

	return "", fmt.Errorf("invalid token or missing claims")
}
