package repositories

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CityRepository struct{}

var cityOnce sync.Once
var cityInstance *CityRepository

func NewCityRepository() *CityRepository {
	cityOnce.Do(func() {
		cityInstance = &CityRepository{}
	})
	return cityInstance
}

func (r *CityRepository) GetByID(ctx *gin.Context, cityID uuid.UUID) (*models.City, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)

	city := &models.City{}
	if err := db.Model(city).
		Where(&models.City{TableModel: models.TableModel{ID: cityID}}).
		First(city).Error; err != nil {
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

func (r *CityRepository) Create(ctx *gin.Context, cityBaseSlice []models.CityBase) ([]models.City, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)

	citySlice := make([]models.City, len(cityBaseSlice))
	for i, cityBase := range cityBaseSlice {
		citySlice[i] = models.City{
			TableModel: models.TableModel{ID: uuid.New()},
			CityBase:   cityBase,
		}
	}
	if err := db.Create(citySlice).Error; err != nil {
		return nil, err
	}
	return citySlice, nil
}
