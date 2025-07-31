package services

import (
	"backend/internal/models"
	"backend/internal/tests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCityService(t *testing.T) {
	repo := NewCityService()
	ctx, _, cleanup := tests.SetupTestContext("test_city_service.db")
	defer cleanup()

	t.Run("單例模式測試", func(t *testing.T) {
		repo2 := NewCityService()
		assert.Same(t, repo, repo2, "應該返回相同的實例")
		assert.NotNil(t, repo, "實例不應該為 nil")
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("成功創建城市", func(t *testing.T) {
			cityBaseSlice := []models.CityBase{
				{Name: "city1"},
				{Name: "city2"},
			}
			cities, err := repo.Create(ctx, cityBaseSlice)
			assert.NoError(t, err)
			assert.Len(t, cities, 2)
			assert.Equal(t, "city1", cities[0].Name)
			assert.Equal(t, "city2", cities[1].Name)
		})

		t.Run("創建城市失敗 - 名稱重複", func(t *testing.T) {
			cityBase := models.CityBase{Name: "city3"}
			repo.Create(ctx, []models.CityBase{cityBase})
			_, err := repo.Create(ctx, []models.CityBase{cityBase})
			assert.Error(t, err)
		})
	})

	t.Run("GetByID", func(*testing.T) {
		t.Run("成功找到城市", func(t *testing.T) {
			citySlice, err := repo.Create(ctx, []models.CityBase{{Name: "city-getbyid-1"}})
			assert.NotNil(t, citySlice)
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
			assert.Nil(t, city)
			assert.Error(t, err)
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
