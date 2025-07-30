package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddressService struct {
	AddressRepository *repositories.AddressRepository
}

var addressOnce sync.Once
var addressService *AddressService

func NewAddressService() *AddressService {
	addressOnce.Do(func() {
		addressService = &AddressService{
			AddressRepository: repositories.NewAddressRepository(),
		}
	})
	return addressService
}

func (s *AddressService) GetByID(ctx *gin.Context, addressID uuid.UUID) (*models.Address, error) {
	return s.AddressRepository.GetByID(ctx, addressID)
}

func (s *AddressService) Create(ctx *gin.Context, addressBaseSlice []models.AddressBase) ([]models.Address, error) {
	return s.AddressRepository.Create(ctx, addressBaseSlice)
}

func (s *AddressService) DeleteByID(ctx *gin.Context, addressID uuid.UUID) error {
	return s.AddressRepository.DeleteByID(ctx, addressID)
}
