package repositories

import (
	"backend/internal/models"
	"backend/internal/tests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewCityRepository(t *testing.T) {
	t.Run("單例模式測試", func(t *testing.T) {
		repo1 := NewCityRepository()
		repo2 := NewCityRepository()

		assert.Same(t, repo1, repo2, "應該返回相同的實例")
		assert.NotNil(t, repo1, "實例不應該為 nil")
	})
}

func TestCityRepositoryGetByID(t *testing.T) {
	ctx, cleanup := tests.SetupTestContext()
	defer cleanup()

	t.Run("成功找到城市", func(t *testing.T) {
		repo := NewCityRepository()
		citySlice, err := repo.Create(ctx, []models.CityBase{{Name: "台北市"}})
		assert.NoError(t, err)
		expectedCity := citySlice[0]

		city, err := repo.GetByID(ctx, expectedCity.ID)
		assert.NoError(t, err)
		assert.NotNil(t, city)
		assert.Equal(t, expectedCity.ID, city.ID)
		assert.Equal(t, expectedCity.Name, city.Name)
	})

	t.Run("城市不存在", func(t *testing.T) {
		repo := NewCityRepository()
		city, err := repo.GetByID(ctx, uuid.New())

		assert.Error(t, err)
		assert.Nil(t, city)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestCityRepositoryGetAll(t *testing.T) {
	ctx, cleanup := tests.SetupTestContext()
	defer cleanup()

	t.Run("成功獲取所有城市", func(t *testing.T) {

		repo := NewCityRepository()
		cities, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, cities)
		assert.Greater(t, len(cities), 0, "應該有城市資料")
		for _, city := range cities {
			assert.NotEmpty(t, city.Name, "城市名稱不應該為空")
		}
	})
}
