package database

import (
	"backend/internal/models"
	"backend/internal/pkg"
	"errors"
	"os"
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
	if adminEmail, adminPassword := os.Getenv("ADMIN_EMAIL"), os.Getenv("ADMIN_PASSWORD"); adminEmail != "" && adminPassword != "" {
		_, err := createAdminUser(db, adminEmail, adminPassword)
		if err != nil {
			return err
		}
	}

	// 創建訪客帳號
	_, err := createGuestUser(db, "temp@temp.com", "temp@temp")
	if err != nil {
		return err
	}

	// 預設創建，台北到台南的 city 資料
	createTaiwanCitys(db)

	return nil
}

func createGuestUser(db *gorm.DB, email string, password string) (*models.User, error) {
	cryptoUtils := pkg.NewCryptoUtils()

	// 驗證 email 格式，提取 username
	re := regexp.MustCompile(`^(.+?)@.+\..+$`)
	matches := re.FindStringSubmatch(email)
	if len(matches) < 2 {
		return nil, errors.New("invalid email format")
	}

	// 檢查是否已存在相同 email 的使用者
	user := &models.User{}
	if err := db.Where("email = ?", email).First(user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 已存在訪客帳號，直接回傳
	if user.ID != uuid.Nil {
		return user, nil
	}

	// 建立訪客帳號
	newUser := &models.User{
		TableModel: models.TableModel{ID: uuid.New()},
		UserBase: models.UserBase{
			Username: "訪客",
			Email:    email,
			HashedPassword: cryptoUtils.GeneratePasswordHash(&pkg.CryptoUtilsPasswordHashInput{
				Email:    email,
				Password: password,
			}),
			Role: models.RoleNormalCustomer,
		},
	}
	if err := db.Create(newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
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

	// 刪除現有的管理員帳號
	admins := []models.User{}
	if err := db.Where("role = ?", models.RoleAdmin).Find(&admins).Error; err != nil {
		return nil, err
	}
	for _, admin := range admins {
		posts := []models.Post{}
		if err := db.Where("author_id = ?", admin.ID).Find(&posts).Error; err != nil {
			return nil, err
		}
		for _, post := range posts {
			if err := db.Where("post_id = ?", post.ID).Delete(&models.Comment{}).Error; err != nil {
				return nil, err
			}
			if err := db.Model(&post).Association("Tags").Clear(); err != nil {
				return nil, err
			}
			if err := db.Model(&post).Association("Likes").Clear(); err != nil {
				return nil, err
			}
			if err := db.Where("id = ?", post.ID).Delete(&models.Post{}).Error; err != nil {
				return nil, err
			}
		}
		if err := db.Where("id = ?", admin.ID).Delete(&models.User{}).Error; err != nil {
			return nil, err
		}
	}

	// 檢查是否已存在相同 email 的使用者
	user := &models.User{}
	if err := db.Where("email = ?", email).First(user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 建立管理帳號
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
