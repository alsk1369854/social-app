package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	DBName    string
	Timezone  string
	EnableLog bool
}

func SetupPostgres(cfg *PostgresConfig) *gorm.DB {
	if cfg.Timezone == "" {
		cfg.Timezone = "Asia/Taipei"
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.Timezone,
	)

	gormCfg := &gorm.Config{}
	if cfg.EnableLog {
		gormCfg.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL database:", err)
	}
	if err := Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	return db
}
