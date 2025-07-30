package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CityService struct {
	CityRepository *repositories.CityRepository
}

var cityOnce sync.Once
var cityService *CityService

func NewCityService() *CityService {
	cityOnce.Do(func() {
		cityService = &CityService{
			CityRepository: repositories.NewCityRepository(),
		}
	})
	return cityService
}

func (s *CityService) GetByID(ctx *gin.Context, cityID uuid.UUID) (*models.City, error) {
	return s.CityRepository.GetByID(ctx, cityID)
}

func (s *CityService) GetAll(ctx *gin.Context) ([]models.City, error) {
	return s.CityRepository.GetAll(ctx)
}

func (s *CityService) Create(ctx *gin.Context, cityBaseSlice []models.CityBase) ([]models.City, error) {
	return s.CityRepository.Create(ctx, cityBaseSlice)
}
