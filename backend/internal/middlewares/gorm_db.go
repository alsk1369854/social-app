package middlewares

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const CONTEXT_KEY_GORM_DB string = "CONTEXT_KEY:GORM_DB"

func WarpGORMDBHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		SetContentGORMDB(ctx, db)
		ctx.Next()
	}
}

func GetContentGORMDB(ctx *gin.Context) (*gorm.DB, error) {
	value, exists := ctx.Get(CONTEXT_KEY_GORM_DB)
	if !exists {
		return nil, errors.New("GORM DB not found in context")
	}
	db, ok := value.(*gorm.DB)
	if !ok {
		return nil, errors.New("GORM DB type assertion failed")
	}
	return db, nil
}

func SetContentGORMDB(ctx *gin.Context, db *gorm.DB) {
	ctx.Set(CONTEXT_KEY_GORM_DB, db)
}

func TransactionGORMDB(ctx *gin.Context, fn func(*gorm.DB) error) error {
	db, err := GetContentGORMDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		SetContentGORMDB(ctx, tx)
		defer SetContentGORMDB(ctx, db)
		return fn(tx)
	})
}
