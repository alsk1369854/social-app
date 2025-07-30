package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService struct {
	UserRepository *repositories.UserRepository

	CityService    *repositories.CityRepository
	AddressService *repositories.AddressRepository
}

var userOnce sync.Once
var userService *UserService

func NewUserService() *UserService {
	userOnce.Do(func() {
		userService = &UserService{
			UserRepository: repositories.NewUserRepository(),
			CityService:    repositories.NewCityRepository(),
			AddressService: repositories.NewAddressRepository(),
		}
	})
	return userService
}

func (s *UserService) GetByID(ctx *gin.Context, userID uuid.UUID) (*models.User, error) {
	return s.UserRepository.GetByID(ctx, userID)
}

func (s *UserService) Register(ctx *gin.Context, userData models.UserRegisterRequest) (*models.User, error) {
	var addressID *uuid.UUID
	if userData.Address != nil {
		cityID, err := uuid.Parse(userData.Address.CityID)
		if err != nil {
			return nil, errors.New("invalid city ID")
		}
		if _, err := s.CityService.GetByID(ctx, cityID); err != nil {
			return nil, errors.New("city not found")
		}
		addressSlice, err := s.AddressService.Create(ctx, []models.AddressBase{{
			CityID: cityID,
			Street: userData.Address.Street,
		}})
		if err != nil {
			return nil, errors.New("invalid address data")
		}
		addressID = &(addressSlice[0].ID)
	}

	passwordHash := userData.Password // Assume password hashing is done here
	userBase := models.UserBase{
		Username:     userData.Username,
		Email:        userData.Email,
		PasswordHash: passwordHash,
		Age:          userData.Age,
		AddressID:    addressID,
	}

	userSlice, err := s.UserRepository.Create(ctx, []models.UserBase{userBase})
	if err != nil {
		return nil, err
	}
	return &userSlice[0], nil
}

func (s *UserService) Create(ctx *gin.Context, userBaseSlice []models.UserBase) ([]models.User, error) {
	return s.UserRepository.Create(ctx, userBaseSlice)
}

func (s *UserService) DeleteByID(ctx *gin.Context, userID uuid.UUID) error {
	return s.UserRepository.DeleteByID(ctx, userID)
}
