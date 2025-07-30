package middlewares

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const CONTEXT_KEY_GORM_DB string = "CONTEXT_KEY_GORM_DB"

func SetGORMDB(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(CONTEXT_KEY_GORM_DB, db)
		ctx.Next()
	}
}
