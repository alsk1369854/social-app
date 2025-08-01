package middlewares

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const CONTEXT_KEY_GORM_DB string = "CONTEXT_KEY:GORM_DB"

func SetGORMDB(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(CONTEXT_KEY_GORM_DB, db)
		ctx.Next()
	}
}

func GetContentGORMDB(ctx *gin.Context) (*gorm.DB, bool) {
	value, exists := ctx.Get(CONTEXT_KEY_GORM_DB)
	if !exists {
		return nil, false
	}
	db, ok := value.(*gorm.DB)
	if !ok {
		return nil, false
	}
	return db, true
}
