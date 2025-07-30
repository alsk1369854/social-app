package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "backend/docs"
	"backend/internal/database"

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
	// Command line flags for host and port
	var host string
	var port string
	flag.StringVar(&host, "host", "0.0.0.0", "Host for the server")
	flag.StringVar(&port, "port", "8080", "Port for the server")
	flag.Parse()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Connect to SQLite database
	db, err := database.ConnectToSQLLine(os.Getenv("SQLITE_DATABASE"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database.Migrate(db)

	// Build API
	engin := gin.Default()
	api := engin.Group("/api")
	{
		api.GET("/ping", ping)
	}
	// Add API documentation, e.g http://localhost:8080/swagger/index.html
	engin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Printf("Server is running on %s:%s", host, port)
	if err := engin.Run(host + ":" + port); err != nil {
		log.Fatal(err)
	}
}
