package services

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/repositories"
	"sync"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService struct {
	UserRepository    *repositories.UserRepository
	AddressRepository *repositories.AddressRepository

	ErrorUtils *pkg.ErrorUtils
}

var userOnce sync.Once
var userService *UserService

func NewUserService() *UserService {
	userOnce.Do(func() {
		userService = &UserService{
			UserRepository:    repositories.NewUserRepository(),
			AddressRepository: repositories.NewAddressRepository(),

			ErrorUtils: pkg.NewErrorUtils(),
		}
	})
	return userService
}

func (s *UserService) GetByID(ctx *gin.Context, userID uuid.UUID) (*models.User, error) {
	return s.UserRepository.GetByID(ctx, userID)
}

func (s *UserService) GetByUsername(ctx *gin.Context, username string) (*models.User, error) {
	return s.UserRepository.GetByUsername(ctx, username)
}

func (s *UserService) GetByEmail(ctx *gin.Context, email string) (*models.User, error) {
	return s.UserRepository.GetByEmail(ctx, email)
}

func (s *UserService) Create(ctx *gin.Context, userBaseSlice []models.UserBase) ([]models.User, error) {
	return s.UserRepository.Create(ctx, userBaseSlice)
}

func (s *UserService) DeleteByID(ctx *gin.Context, userID uuid.UUID) error {
	return s.UserRepository.DeleteByID(ctx, userID)
}

func (s *UserService) CreateUserWithAddress(ctx *gin.Context, userBase *models.UserBase, addressBase *models.AddressBase) (*models.User, error) {

	var user *models.User
	if err := middlewares.TransactionGORMDB(ctx, func(db *gorm.DB) error {
		// 如果有提供地址資訊，則建立地址
		var address *models.Address
		if addressBase != nil {
			addresses, err := s.AddressRepository.Create(ctx, []models.AddressBase{*addressBase})
			if err != nil {
				return err
			}
			address = &addresses[0]
			userBase.AddressID = &address.ID
		}

		// 創建使用者
		users, err := s.UserRepository.Create(ctx, []models.UserBase{*userBase})
		if err != nil {
			return err
		}

		user = &users[0]
		user.Address = address
		return nil
	}); err != nil {
		return nil, s.ErrorUtils.ServerInternalError(err.Error())
	}

	return user, nil
}

// func (s *UserService) Register(ctx *gin.Context, userData *models.UserRegisterRequest) (*models.User, error) {

// 	// 檢查 email 是否存在
// 	if _, err := s.GetByEmail(ctx, userData.Email); err == nil {
// 		return nil, errors.New("email already exists")
// 	}
// 	// 檢查 username 是否存在
// 	if _, err := s.GetByUsername(ctx, userData.Username); err == nil {
// 		return nil, errors.New("username already exists")
// 	}
// 	// 檢查是否有填寫 address
// 	var addressID *uuid.UUID
// 	if userData.Address != nil {
// 		if _, err := s.CityService.GetByID(ctx, userData.Address.CityID); err != nil {
// 			return nil, errors.New("city not found")
// 		}
// 		addressSlice, err := s.AddressService.Create(ctx, []models.AddressBase{{
// 			CityID: userData.Address.CityID,
// 			Street: userData.Address.Street,
// 		}})
// 		if err != nil {
// 			return nil, errors.New("invalid address data")
// 		}
// 		addressID = &(addressSlice[0].ID)
// 	}
// 	// 建立 User 資料
// 	cryptoUtils := pkg.NewCryptoUtils()
// 	HashedPassword := cryptoUtils.GeneratePasswordHash(&pkg.CryptoUtilsPasswordHashInput{
// 		Email:    userData.Email,
// 		Username: userData.Username,
// 		Password: userData.Password,
// 	})
// 	userBase := models.UserBase{
// 		Username:       userData.Username,
// 		Email:          userData.Email,
// 		HashedPassword: HashedPassword,
// 		Age:            userData.Age,
// 		AddressID:      addressID,
// 	}
// 	users, err := s.UserRepository.Create(ctx, []models.UserBase{userBase})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 載入地址資訊
// 	result := &users[0]
// 	if addressID != nil {
// 		address, err := s.AddressService.GetByID(ctx, *addressID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		result.Address = address
// 	}
// 	return result, nil
// }
