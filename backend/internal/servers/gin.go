package servers

import (
	"backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type GinConfig struct {
	DB   *gorm.DB
	Host string
	Port string
}

func SetupGin(cfg *GinConfig) (*gin.Engine, *gin.RouterGroup) {
	server := gin.Default()
	apiRouter := server.Group("/api")
	apiRouter.Use(middlewares.SetGORMDB(cfg.DB))

	// Add API documentation, e.g http://localhost:8080/swagger/index.html
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return server, apiRouter
}
