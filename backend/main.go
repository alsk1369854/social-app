package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "backend/docs"
	"backend/internal/database"
	"backend/internal/routers"
	"backend/internal/servers"
)

// @title Social APP API
// @version 1.0
// @description Social APP API
// @securityDefinitions.apiKey AccessToken
// @in header
// @name Authorization
// @basePath /api
func main() {
	// Command line flags for host and port
	var host string
	var port string
	var debug bool
	flag.StringVar(&host, "host", "0.0.0.0", "Host for the server")
	flag.StringVar(&port, "port", "28080", "Port for the server")
	flag.BoolVar(&debug, "debug", true, "Enable debug mode")
	flag.Parse()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Connect to database
	db := database.SetupPostgres(&database.PostgresConfig{
		Host:      os.Getenv("DB_HOST"),
		Port:      os.Getenv("DB_PORT"),
		User:      os.Getenv("DB_USER"),
		Password:  os.Getenv("DB_PASSWORD"),
		DBName:    os.Getenv("DB_NAME"),
		Timezone:  "Asia/Taipei",
		EnableLog: true,
	})

	// Setup Gin server
	server, apiRouter := servers.SetupGin(&servers.GinConfig{DB: db})
	routers.NewCityRouter().Bind(apiRouter)
	routers.NewUserRouter().Bind(apiRouter)
	routers.NewPostRouter().Bind(apiRouter)
	routers.NewCommentRouter().Bind(apiRouter)

	// Start the server
	log.Printf("Swagger docs available at http://%s:%s/swagger/index.html\n", host, port)
	if err := server.Run(host + ":" + port); err != nil {
		log.Fatal(err)
	}
}
