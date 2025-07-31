package repositories

import (
	"backend/internal/models"
	"backend/internal/tests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCityRepository(t *testing.T) {
	repo := NewCityRepository()
	ctx, _, cleanup := tests.SetupTestContext("test_city_repository.db")
	defer cleanup()

	t.Run("單例模式測試", func(t *testing.T) {
		repo1 := NewCityRepository()
		assert.Same(t, repo, repo1, "應該返回相同的實例")
		assert.NotNil(t, repo1, "實例不應該為 nil")
	})

	t.Run("GetByID", func(t *testing.T) {
		t.Run("成功找到城市", func(t *testing.T) {
			citySlice, err := repo.Create(ctx, []models.CityBase{{Name: "getbyid-city-1"}})
			assert.NoError(t, err)
			expectedCity := citySlice[0]

			city, err := repo.GetByID(ctx, expectedCity.ID)
			assert.NoError(t, err)
			assert.NotNil(t, city)
			assert.Equal(t, expectedCity.ID, city.ID)
			assert.Equal(t, expectedCity.Name, city.Name)
		})

		t.Run("城市不存在", func(t *testing.T) {
			city, err := repo.GetByID(ctx, uuid.New())

			assert.Error(t, err)
			assert.Nil(t, city)
			assert.Equal(t, gorm.ErrRecordNotFound, err)
		})
	})

	t.Run("GetAll", func(t *testing.T) {
		t.Run("成功獲取所有城市", func(t *testing.T) {
			cities, err := repo.GetAll(ctx)

			assert.NoError(t, err)
			assert.NotNil(t, cities)
			assert.Greater(t, len(cities), 0, "應該有城市資料")
			for _, city := range cities {
				assert.NotEmpty(t, city.Name, "城市名稱不應該為空")
			}
		})
	})
}
