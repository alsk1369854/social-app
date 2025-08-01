package middlewares

import (
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const CONTEXT_KEY_ACCESS_TOKEN_DATA string = "CONTEXT_KEY:ACCESS_TOKEN_DATA"

func VerifyAccessToken(validateToken func(authHeader string) (jwt.MapClaims, bool)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(401, models.ErrorResponse{Error: "authorization header is required"})
			ctx.Abort()
			return
		}
		tokenData, valid := validateToken(authHeader)
		if !valid {
			ctx.JSON(401, models.ErrorResponse{Error: "invalid token"})
			ctx.Abort()
			return
		}
		ctx.Set(CONTEXT_KEY_ACCESS_TOKEN_DATA, tokenData)

		ctx.Next()
	}
}

func GetContentAccessTokenData(ctx *gin.Context) (jwt.MapClaims, bool) {
	value, exists := ctx.Get(CONTEXT_KEY_ACCESS_TOKEN_DATA)
	if !exists {
		return nil, false
	}
	tokenData, ok := value.(jwt.MapClaims)
	if !ok {
		return nil, false
	}
	return tokenData, true
}
