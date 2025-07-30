package repositories

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressRepository struct{}

var addressOnce sync.Once
var addressInstance *AddressRepository

func NewAddressRepository() *AddressRepository {
	addressOnce.Do(func() {
		addressInstance = &AddressRepository{}
	})
	return addressInstance
}

func (r *AddressRepository) GetByID(ctx *gin.Context, addressID uuid.UUID) (*models.Address, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)

	address := &models.Address{}
	if err := db.Model(&models.Address{}).
		Where(&models.Address{TableModel: models.TableModel{ID: addressID}}).
		First(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}

func (r *AddressRepository) Create(ctx *gin.Context, addressBaseSlice []models.AddressBase) ([]models.Address, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)

	addressSlice := make([]models.Address, len(addressBaseSlice))
	for i, addressBase := range addressBaseSlice {
		addressSlice[i] = models.Address{
			TableModel:  models.TableModel{ID: uuid.New()},
			AddressBase: addressBase,
		}
	}
	if err := db.Create(addressSlice).Error; err != nil {
		return nil, err
	}
	return addressSlice, nil
}

func (r *AddressRepository) DeleteByID(ctx *gin.Context, addressID uuid.UUID) error {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)

	if err := db.Model(&models.Address{}).
		Where(&models.Address{TableModel: models.TableModel{ID: addressID}}).
		Delete(&models.Address{}).Error; err != nil {
		return err
	}
	return nil
}
