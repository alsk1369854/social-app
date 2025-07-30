package database

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

// Migrate
// param db *gorm.DB Database connect
// return error
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.City{},
		&models.Address{},
		&models.User{},
	); err != nil {
		return err
	}

	// 預設創建 台北到台南的 city 資料
	cityNames := []string{"基隆市", "嘉義市", "台北市", "嘉義縣", "新北市", "台南市", "桃園縣", "高雄市", "新竹市", "屏東縣", "新竹縣", "台東縣", "苗栗縣", "花蓮縣", "台中市", "宜蘭縣", "彰化縣", "澎湖縣", "南投縣", "金門縣", "雲林縣", "連江縣"}
	for _, name := range cityNames {
		city := &models.City{Name: name}
		if err := db.FirstOrCreate(city, city).Error; err != nil {
			return err
		}
	}
	return nil
}
