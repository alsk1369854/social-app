package services

import (
	"backend/internal/models"
	"backend/internal/tests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddressService(t *testing.T) {
	service := NewAddressService()
	ctx, _, cleanup := tests.SetupTestContext("test_address_service.db")
	defer cleanup()

	t.Run("NewAddressService", func(t *testing.T) {
		service2 := NewAddressService()
		assert.Same(t, service, service2, "Expected singleton instance of AddressService")
	})

	t.Run("Create", func(t *testing.T) {
		addressBaseSlice := []models.AddressBase{
			{CityID: uuid.New(), Street: "123 Test St"},
			{CityID: uuid.New(), Street: "456 Example Ave"},
		}
		addresses, err := service.Create(ctx, addressBaseSlice)
		assert.NoError(t, err, "Expected no error when creating addresses")
		assert.Len(t, addresses, 2, "Expected two addresses to be created")
	})

	t.Run("GetByID not exist id", func(t *testing.T) {
		address, err := service.GetByID(ctx, uuid.New()) // Replace with a valid UUID for a real test
		assert.Nil(t, address, "Expected nil address for non-existent ID")
		assert.Error(t, err, "Expected error for non-existent address ID")
	})

	t.Run("GetByID exist id", func(t *testing.T) {
		// create a known address first
		addressBase := models.AddressBase{CityID: uuid.New(), Street: "789 Sample Rd"}
		addresses, _ := service.Create(ctx, []models.AddressBase{addressBase})
		expected := addresses[0]

		address, err := service.GetByID(ctx, expected.ID)
		assert.NotNil(t, address, "Expected non-nil address for existing ID")
		assert.NoError(t, err, "Expected no error for existing address ID")
		assert.Equal(t, expected.ID, address.ID, "Expected address ID to match")
	})

	t.Run("DeleteByID not exist id", func(t *testing.T) {
		address, err := service.DeleteByID(ctx, uuid.New()) // Replace with a valid UUID for a real test
		assert.Error(t, err, "Expected error when deleting non-existent address ID")
		assert.Nil(t, address, "Expected nil address for non-existent ID")
	})

	t.Run("DeleteByID exist id", func(t *testing.T) {
		// create a known address first
		addressBase := models.AddressBase{CityID: uuid.New(), Street: "123 Delete St"}
		addresses, _ := service.Create(ctx, []models.AddressBase{addressBase})
		expected := addresses[0]

		address, err := service.DeleteByID(ctx, expected.ID)
		assert.NoError(t, err, "Expected no error when deleting existing address ID")
		assert.NotNil(t, address, "Expected non-nil address for deleted ID")
		assert.Equal(t, expected.ID, address.ID, "Expected deleted address ID to match")
	})

}
