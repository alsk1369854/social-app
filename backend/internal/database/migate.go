package database

import (
	"backend/internal/models"
	"backend/internal/pkg"
	"errors"
	"regexp"

	"github.com/google/uuid"
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
		&models.Post{},
		&models.Comment{},
		&models.Tag{},
	); err != nil {
		return err
	}

	// 創建管理員帳號
	_, err := createAdminUser(db, "admin@example.com", "admin@admin")
	if err != nil {
		return err
	}

	// 預設創建，台北到台南的 city 資料
	createTaiwanCitys(db)

	return nil
}

func createAdminUser(db *gorm.DB, email string, password string) (*models.User, error) {
	cryptoUtils := pkg.NewCryptoUtils()

	// 使用正則表達式提取使用者名稱
	extractUsernameRegex := regexp.MustCompile(`^(?P<username>.+)@.+\..+$`)
	if !extractUsernameRegex.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	// 刪除已存在的管理員帳號
	ordAdminUser := &models.User{}
	if err := db.Where(&models.User{UserBase: models.UserBase{Email: email}}).
		Select("id").First(ordAdminUser).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if ordAdminUser.ID != uuid.Nil {
		if err := db.Where("author_id = ?", ordAdminUser.ID).Delete(&models.Post{}).Error; err != nil {
			return nil, err
		}
		if err := db.Where("id = ?", ordAdminUser.ID).Delete(&models.User{}).Error; err != nil {
			return nil, err
		}
	}

	// 創建新的管理員帳號
	matches := extractUsernameRegex.FindStringSubmatch(email)
	username := matches[extractUsernameRegex.SubexpIndex("username")]
	hashedPassword := cryptoUtils.GeneratePasswordHash(&pkg.CryptoUtilsPasswordHashInput{
		Email:    email,
		Password: password,
	})
	adminUser := &models.User{
		TableModel: models.TableModel{
			ID: uuid.New(),
		},
		UserBase: models.UserBase{
			Username:       username,
			Email:          email,
			HashedPassword: hashedPassword,
		},
	}
	if err := db.Create(adminUser).Error; err != nil {
		return nil, err
	}
	return adminUser, nil
}

func createTaiwanCitys(db *gorm.DB) error {
	cityNames := []string{"基隆市", "嘉義市", "台北市", "嘉義縣", "新北市", "台南市", "桃園縣", "高雄市", "新竹市", "屏東縣", "新竹縣", "台東縣", "苗栗縣", "花蓮縣", "台中市", "宜蘭縣", "彰化縣", "澎湖縣", "南投縣", "金門縣", "雲林縣", "連江縣"}
	for _, name := range cityNames {
		if db.Where(&models.City{CityBase: models.CityBase{Name: name}}).First(&models.City{}).RowsAffected != 0 {
			continue
		}
		city := &models.City{
			TableModel: models.TableModel{
				ID: uuid.New(),
			},
			CityBase: models.CityBase{
				Name: name,
			},
		}
		if err := db.Create(city).Error; err != nil {
			return err
		}
	}
	return nil
}
