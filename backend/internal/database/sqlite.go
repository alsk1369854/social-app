package database

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLiteConfig struct {
	DBFile    string
	EnableLog bool
}

func SetupSQLite(cfg *SQLiteConfig) *gorm.DB {
	if cfg.DBFile == "" {
		cfg.DBFile = "default.db"
	}

	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", cfg.DBFile)

	gormCfg := &gorm.Config{}
	if cfg.EnableLog {
		gormCfg.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(sqlite.Open(dsn), gormCfg)
	if err != nil {
		log.Fatal("Failed to connect to SQLite database:", err)
	}
	if err := Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	return db
}
