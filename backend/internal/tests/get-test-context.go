package tests

import (
	"backend/internal/database"
	"backend/internal/middlewares"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func truncatePostgresAllTables(db *gorm.DB) {
	var tableNames []string
	err := db.Raw(`
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
	`).Scan(&tableNames).Error
	if err != nil {
		fmt.Println("❌ Failed to list tables:", err)
		return
	}

	// // 關閉 FK 檢查避免錯誤
	// if err := db.Exec("SET session_replication_role = 'replica'").Error; err != nil {
	// 	fmt.Println("❌ Failed to disable FK checks:", err)
	// 	return
	// }

	// for _, table := range tableNames {
	// 	if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE \"%s\" RESTART IDENTITY CASCADE", table)).Error; err != nil {
	// 		fmt.Printf("❌ Failed to truncate table %s: %v\n", table, err)
	// 	}
	// }
	for _, table := range tableNames {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE \"%s\" CASCADE", table)).Error; err != nil {
			fmt.Printf("❌ Failed to truncate table %s: %v\n", table, err)
		}
	}

	// // 恢復 FK 檢查
	// if err := db.Exec("SET session_replication_role = 'origin'").Error; err != nil {
	// 	fmt.Println("❌ Failed to re-enable FK checks:", err)
	// }
}

func GetTestContext() (*gin.Context, func()) {
	// Load environment variables from .env file
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("Failed to load .env file: %v", err)
	// }

	// db, err := database.ConnectToPostgres(database.PostgresConfig{
	// 	Host:      os.Getenv("DB_HOST"),
	// 	Port:      os.Getenv("DB_PORT"),
	// 	User:      os.Getenv("DB_USER"),
	// 	Password:  os.Getenv("DB_PASSWORD"),
	// 	DBName:    os.Getenv("DB_NAME"),
	// 	Timezone:  "Asia/Taipei",
	// 	EnableLog: true,
	// })

	db, err := database.ConnectToSQLite("./test.db")
	if err != nil {
		log.Fatalf("Failed to create GORM database: %v", err)
	}
	database.Migrate(db)
	cleanup := func() {
		// truncatePostgresAllTables(db)
		// remove test database file
		os.Remove("./test.db")
	}

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("GET", "/test", nil)
	ctx.Request = req
	ctx.Set(middlewares.CONTEXT_KEY_GORM_DB, db)

	return ctx, cleanup
}
