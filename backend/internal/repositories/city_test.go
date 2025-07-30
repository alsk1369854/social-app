package repositories

import (
	"backend/internal/models"
	"backend/internal/tests"
	"testing"

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
	ctx, cleanup := tests.GetTestContext()
	defer cleanup()

	t.Run("成功找到城市", func(t *testing.T) {
		cityID := uint(3)
		expectedCity := &models.City{
			Model: gorm.Model{ID: cityID},
			Name:  "台北市",
		}

		repo := NewCityRepository()
		city, err := repo.GetByID(ctx, cityID)

		assert.NoError(t, err)
		assert.NotNil(t, city)
		assert.Equal(t, expectedCity.ID, city.ID)
		assert.Equal(t, expectedCity.Name, city.Name)
	})

	t.Run("城市不存在", func(t *testing.T) {
		cityID := uint(999)

		repo := NewCityRepository()
		city, err := repo.GetByID(ctx, cityID)

		assert.Error(t, err)
		assert.Nil(t, city)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

// func TestCityRepository_GetAll(t *testing.T) {
// 	gormDB, mock, cleanup := setupCityTestDB(t)
// 	defer cleanup()

// 	repo := NewCityRepository()
// 	ctx := setupCityGinContext(gormDB)

// 	t.Run("成功獲取所有城市", func(t *testing.T) {
// 		expectedCities := []models.City{
// 			{Model: gorm.Model{ID: 1}, Name: "台北市"},
// 			{Model: gorm.Model{ID: 2}, Name: "新北市"},
// 			{Model: gorm.Model{ID: 3}, Name: "桃園市"},
// 		}

// 		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
// 			AddRow(expectedCities[0].ID, expectedCities[0].Name, "2024-01-01 00:00:00", "2024-01-01 00:00:00", nil).
// 			AddRow(expectedCities[1].ID, expectedCities[1].Name, "2024-01-01 00:00:00", "2024-01-01 00:00:00", nil).
// 			AddRow(expectedCities[2].ID, expectedCities[2].Name, "2024-01-01 00:00:00", "2024-01-01 00:00:00", nil)

// 		mock.ExpectQuery(`SELECT \* FROM "cities" WHERE "cities"\."deleted_at" IS NULL`).
// 			WillReturnRows(rows)

// 		cities, err := repo.GetAll(ctx)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, cities)
// 		assert.Len(t, cities, 3)
// 		assert.Equal(t, expectedCities[0].Name, cities[0].Name)
// 		assert.Equal(t, expectedCities[1].Name, cities[1].Name)
// 		assert.Equal(t, expectedCities[2].Name, cities[2].Name)
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})

// 	t.Run("沒有城市資料", func(t *testing.T) {
// 		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"})

// 		mock.ExpectQuery(`SELECT \* FROM "cities" WHERE "cities"\."deleted_at" IS NULL`).
// 			WillReturnRows(rows)

// 		cities, err := repo.GetAll(ctx)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, cities)
// 		assert.Len(t, cities, 0)
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})

// 	t.Run("資料庫查詢錯誤", func(t *testing.T) {
// 		expectedError := gorm.ErrInvalidDB

// 		mock.ExpectQuery(`SELECT \* FROM "cities" WHERE "cities"\."deleted_at" IS NULL`).
// 			WillReturnError(expectedError)

// 		cities, err := repo.GetAll(ctx)

// 		assert.Error(t, err)
// 		assert.Nil(t, cities)
// 		assert.Equal(t, expectedError, err)
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})
// }

// func TestCityRepository_ContextWithoutDB(t *testing.T) {
// 	repo := NewCityRepository()

// 	gin.SetMode(gin.TestMode)
// 	w := httptest.NewRecorder()
// 	ctx, _ := gin.CreateTestContext(w)
// 	req, _ := http.NewRequest("GET", "/test", nil)
// 	ctx.Request = req
// 	// 注意：這裡不設置 GORM DB 到 context

// 	t.Run("GetByID 沒有 DB context 應該 panic", func(t *testing.T) {
// 		assert.Panics(t, func() {
// 			repo.GetByID(ctx, 1)
// 		}, "應該因為沒有 DB context 而 panic")
// 	})

// 	t.Run("GetAll 沒有 DB context 應該 panic", func(t *testing.T) {
// 		assert.Panics(t, func() {
// 			repo.GetAll(ctx)
// 		}, "應該因為沒有 DB context 而 panic")
// 	})
// }
