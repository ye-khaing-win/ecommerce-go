package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	UserID string
	Email  string
	Role   string
	jwt.RegisteredClaims
}

func SignJWT(claims CustomClaims, secret string, expires time.Duration) (string, error) {
	now := time.Now()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(expires)),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))

	return signed, err
}
