package database

import (
	"gorm.io/gorm"
)

// Migrate
// param db *gorm.DB Database connect
// return error
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(); err != nil {
		return err
	}
	return nil
}
