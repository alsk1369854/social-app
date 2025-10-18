package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "backend/docs"
	"backend/internal/database"
	"backend/internal/models"
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
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// default values
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "28080"
	}
	debug := false
	for _, temp := range []string{"true", "1", "TRUE", "True"} {
		if temp == strings.ToLower(os.Getenv("DEBUG_MODE")) {
			debug = true
			break
		}
	}

	// Command line flags for host and port
	flag.StringVar(&host, "host", host, "Host for the server")
	flag.StringVar(&port, "port", port, "Port for the server")
	flag.BoolVar(&debug, "debug", debug, "Enable debug mode")
	flag.Parse()

	fmt.Printf("%v, %v, %v,\n", host, port, debug)

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db := database.SetupPostgres(&database.PostgresConfig{
		Host:      os.Getenv("DB_HOST"),
		Port:      os.Getenv("DB_PORT"),
		User:      os.Getenv("DB_USER"),
		Password:  os.Getenv("DB_PASSWORD"),
		DBName:    os.Getenv("DB_NAME"),
		Timezone:  "Asia/Taipei",
		EnableLog: debug,
	})

	// Setup Gin server
	server, apiRouter := servers.SetupGin(&servers.GinConfig{DB: db})
	routers.NewCityRouter().Bind(apiRouter)
	routers.NewUserRouter().Bind(apiRouter)
	routers.NewPostRouter().Bind(apiRouter)
	routers.NewCommentRouter().Bind(apiRouter)

	// Setup AI Router
	routers.NewAIRouter(&models.AIModelConfigs{
		ChatModel: models.AIModelConfig{
			APIKey:    os.Getenv("OPENAI_API_KEY"),
			BaseURL:   os.Getenv("OPENAI_BASE_URL"),
			ModelName: os.Getenv("OPENAI_CHAT_MODEL"),
		},
	}).Bind(apiRouter)

	// Start the server
	log.Printf("Swagger docs available at http://%s:%s/swagger/index.html\n", host, port)
	if err := server.Run(host + ":" + port); err != nil {
		log.Fatal(err)
	}
}
