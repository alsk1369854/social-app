package main

import (
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "backend/docs"
	"backend/internal/database"
	"backend/internal/middlewares"
	"backend/internal/routers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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

	// Connect to database
	db, err := database.ConnectToPostgres(database.PostgresConfig{
		Host:      os.Getenv("DB_HOST"),
		Port:      os.Getenv("DB_PORT"),
		User:      os.Getenv("DB_USER"),
		Password:  os.Getenv("DB_PASSWORD"),
		DBName:    os.Getenv("DB_NAME"),
		Timezone:  "Asia/Taipei",
		EnableLog: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Build API
	engin := gin.Default()
	api := engin.Group("/api")
	api.Use(middlewares.SetGORMDB(db))
	{
		routers.NewCityRouter().Bind(api)
		routers.NewUserRouter().Bind(api)
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
