package repositories

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CityRepository struct{}

var cityOnce sync.Once
var cityRepository *CityRepository

func NewCityRepository() *CityRepository {
	cityOnce.Do(func() {
		cityRepository = &CityRepository{}
	})
	return cityRepository
}

func (r *CityRepository) GetByID(ctx *gin.Context, cityID uint) (*models.City, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)

	city := &models.City{}
	if err := db.Model(city).Where(&models.City{Model: gorm.Model{ID: cityID}}).First(city).Error; err != nil {
		return nil, err
	}
	return city, nil
}

func (r *CityRepository) GetAll(ctx *gin.Context) ([]models.City, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)

	cities := []models.City{}
	if err := db.Model(&models.City{}).Find(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}
