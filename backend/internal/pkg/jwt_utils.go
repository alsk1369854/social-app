package pkg

import (
	"os"
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

func (u *JWTUtils) GenerateToken(data any, secret []byte) (string, error) {
	if secret == nil {
		secret = []byte(os.Getenv(u.DefaultEnvKey))
	}
	claims := jwt.MapClaims{
		"data": data,
		"iat":  jwt.NewNumericDate(time.Now()),
		"exp":  jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (u *JWTUtils) ParseToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	if secret == nil {
		secret = []byte(os.Getenv(u.DefaultEnvKey))
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
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
