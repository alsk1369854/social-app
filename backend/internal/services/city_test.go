package services

import (
	"backend/internal/models"
	"backend/internal/tests"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewCityService(t *testing.T) {
	t.Run("單例模式測試", func(t *testing.T) {
		repo1 := NewCityService()
		repo2 := NewCityService()

		assert.Same(t, repo1, repo2, "應該返回相同的實例")
		assert.NotNil(t, repo1, "實例不應該為 nil")
	})
}

func TestCityServiceGetByID(t *testing.T) {
	ctx, cleanup := tests.SetupTestContext()
	defer cleanup()

	t.Run("成功找到城市", func(t *testing.T) {
		cityID := uint(3)
		expectedCity := &models.City{
			Model: gorm.Model{ID: cityID},
			Name:  "台北市",
		}

		repo := NewCityService()
		city, err := repo.GetByID(ctx, cityID)

		assert.NoError(t, err)
		assert.NotNil(t, city)
		assert.Equal(t, expectedCity.ID, city.ID)
		assert.Equal(t, expectedCity.Name, city.Name)
	})

	t.Run("城市不存在", func(t *testing.T) {
		cityID := uint(999)

		repo := NewCityService()
		city, err := repo.GetByID(ctx, cityID)

		assert.Error(t, err)
		assert.Nil(t, city)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestCityServiceGetAll(t *testing.T) {
	ctx, cleanup := tests.SetupTestContext()
	defer cleanup()

	t.Run("成功獲取所有城市", func(t *testing.T) {

		repo := NewCityService()
		cities, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, cities)
		assert.Greater(t, len(cities), 0, "應該有城市資料")
		for _, city := range cities {
			assert.NotEmpty(t, city.Name, "城市名稱不應該為空")
		}
	})
}
