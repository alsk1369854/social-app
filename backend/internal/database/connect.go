package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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

func ConnectToPostgres(cfg PostgresConfig) (*gorm.DB, error) {
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

	return gorm.Open(postgres.Open(dsn), gormCfg)
}

func ConnectToSQLite(dbFile string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbFile)

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	return gorm.Open(sqlite.Open(dsn), gormCfg)
}
