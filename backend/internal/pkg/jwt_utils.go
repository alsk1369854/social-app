package pkg

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type JWTUtils struct {
	DefaultEnvKey string
}

var jwtUtilsOnce sync.Once
var jwtUtils *JWTUtils

func NewJWTUtils() *JWTUtils {
	jwtUtilsOnce.Do(func() {
		jwtUtils = &JWTUtils{
			DefaultEnvKey: "JWT_SECRET",
		}
	})
	return jwtUtils
}

func (u *JWTUtils) GenerateToken(data any, secret string) (string, error) {
	claims := jwt.MapClaims{
		"data": data,
		"iat":  jwt.NewNumericDate(time.Now()),
		"exp":  jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (u *JWTUtils) ParseToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
