package servers

import (
	"backend/internal/middlewares"

	"github.com/gin-contrib/cors"
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
	if cfg.Debug {
		// allow CORS for development
		apiRouter.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, // 允许的前端地址
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
		}))
		apiRouter.OPTIONS("/*path", func(c *gin.Context) {
			c.AbortWithStatus(204) // Status code 204 No Content
		})
	}

	// Add API documentation, e.g http://<user-host>/swagger/index.html
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return server, apiRouter
}
