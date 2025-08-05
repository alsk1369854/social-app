package middlewares

import (
	"backend/internal/models"
	"backend/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func ParseJWTAccessToken(authHeader string) (jwt.MapClaims, bool) {
	claims, err := pkg.NewJWTUtils().ParseToken(authHeader, nil)
	if err != nil {
		return nil, false
	}
	return claims, true
}

func GetContentAccessTokenData(ctx *gin.Context) (*models.JWTClaimsData, error) {
	errorUtils := pkg.NewErrorUtils()

	value, exists := ctx.Get(CONTEXT_KEY_ACCESS_TOKEN_DATA)
	if !exists {
		return nil, errorUtils.ServerInternalError("AccessToken Data not found in context")
	}
	tokenData, ok := value.(jwt.MapClaims)
	if !ok {
		return nil, errorUtils.ServerInternalError("AccessToken Data not found in context, type assertion failed")
	}

	// 解析 Token 中的數據
	userID, err := uuid.Parse(tokenData["data"].(map[string]any)["UserID"].(string))
	if err != nil {
		return nil, errorUtils.ServerInternalError("failed to parse user ID from token")
	}

	result := &models.JWTClaimsData{
		UserID: userID,
	}
	return result, nil
}
