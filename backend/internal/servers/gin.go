package servers

import (
	"backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type GinConfig struct {
	DB    *gorm.DB
	Debug bool
}

func SetupGin(cfg *GinConfig) (*gin.Engine, *gin.RouterGroup) {
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()
	apiRouter := server.Group("/api")
	apiRouter.Use(middlewares.WarpGORMDBHandler(cfg.DB))
	// allow CORS for development
	if cfg.Debug {
		apiRouter.Use(func(ctx *gin.Context) {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		})
	}

	// Add API documentation, e.g http://<user-host>/swagger/index.html
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return server, apiRouter
}
