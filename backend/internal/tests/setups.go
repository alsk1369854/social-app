package tests

import (
	"backend/internal/database"
	"backend/internal/middlewares"
	"backend/internal/servers"
	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTestDB(dbFile string) (*gorm.DB, func()) {
	os.Remove(dbFile)
	db := database.SetupSQLite(&database.SQLiteConfig{
		DBFile:    dbFile, // ":memery:"
		EnableLog: false,
	})

	cleanup := func() {
		os.Remove(dbFile)
	}

	return db, cleanup
}

func SetupTestContext(dbFile string) (*gin.Context, *gorm.DB, func()) {
	db, cleanup := SetupTestDB(dbFile)

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set(middlewares.CONTEXT_KEY_GORM_DB, db)

	return ctx, db, cleanup
}

func SetupTestServer(dbFile string) (*gin.Engine, *gin.RouterGroup, *gin.Context, *gorm.DB, func()) {
	ctx, db, cleanup := SetupTestContext(dbFile)

	gin.SetMode(gin.TestMode)
	server, apiRouter := servers.SetupGin(&servers.GinConfig{
		DB: db,
	})

	return server, apiRouter, ctx, db, cleanup
}
