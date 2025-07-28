package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Test API
// @Summary Test API
// @Tags Test
// @Accept application/json
// @Produce application/json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/ping [get]
func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// @title Social APP API
// @version 1.0
// @description Social APP API
// @securityDefinitions.apiKey bearerToken
// @in header
// @name Authorization
// @basePath /api
func main() {
	addr := "0.0.0.0:8080"
	engin := gin.Default()

	api := engin.Group("/api")
	{
		api.GET("/ping", ping)
	}

	// swagger docs
	// http://localhost:8080/docs
	// docs := engin.Group("/docs")
	// {
	// 	docs.GET("", func(ctx *gin.Context) {
	// 		ctx.Redirect(http.StatusPermanentRedirect, "/docs/index.html")
	// 	})
	// 	docs.GET("*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }

	// http://localhost:8080/swagger/index.html
	engin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := engin.Run(addr); err != nil {
		log.Fatal(err)
	}
}
