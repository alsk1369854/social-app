package tests

import (
	"backend/internal/database"
	"backend/internal/middlewares"
	"backend/internal/servers"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTestDB() (*gorm.DB, func()) {
	db := database.SetupSQLite(&database.SQLiteConfig{
		DBFile:    ":memory:",
		EnableLog: false,
	})

	cleanup := func() {
		// os.Remove("./test.db")
	}

	return db, cleanup
}

func SetupTestContext() (*gin.Context, func()) {
	db, cleanup := SetupTestDB()

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set(middlewares.CONTEXT_KEY_GORM_DB, db)

	return ctx, cleanup
}

func SetupTestServer() (*gin.Engine, *gin.RouterGroup, func()) {
	db, cleanup := SetupTestDB()

	gin.SetMode(gin.ReleaseMode)
	server, apiRouter := servers.SetupGin(&servers.GinConfig{
		DB: db,
	})

	return server, apiRouter, cleanup
}
