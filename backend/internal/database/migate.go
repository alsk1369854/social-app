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

	// 驗證 email 格式，提取 username
	re := regexp.MustCompile(`^(.+?)@.+\..+$`)
	matches := re.FindStringSubmatch(email)
	if len(matches) < 2 {
		return nil, errors.New("invalid email format")
	}
	username := matches[1]

	var user models.User
	err := db.Where("email = ?", email).First(&user).Error

	if err == nil {
		// 使用者已存在，更新密碼
		user.HashedPassword = cryptoUtils.GeneratePasswordHash(&pkg.CryptoUtilsPasswordHashInput{
			Email:    email,
			Password: password,
		})
		if err := db.Save(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 使用者不存在，建立新帳號
	newUser := &models.User{
		TableModel: models.TableModel{ID: uuid.New()},
		UserBase: models.UserBase{
			Username: username,
			Email:    email,
			HashedPassword: cryptoUtils.GeneratePasswordHash(&pkg.CryptoUtilsPasswordHashInput{
				Email:    email,
				Password: password,
			}),
		},
	}

	if err := db.Create(newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
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
