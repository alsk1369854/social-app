package routers

import (
	"backend/internal/models"
	"backend/internal/tests"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCityRouter(t *testing.T) {
	server, apiRouter, _, _, cleanup := tests.SetupTestServer("test_city_router.db")
	defer cleanup()
	cityRouter := NewCityRouter()
	cityRouter.Bind(apiRouter)

	t.Run("單例模式測試", func(t *testing.T) {
		cityRouter2 := NewCityRouter()
		assert.Same(t, cityRouter, cityRouter2, "應該返回相同的實例")
	})

	t.Run("GetAll", func(t *testing.T) {
		t.Run("成功獲取所有城市", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/city/all", nil)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			cities := models.CityGetAllResponse{}
			err := json.Unmarshal(recorder.Body.Bytes(), &cities)
			assert.NoError(t, err)
			assert.Greater(t, len(cities), 0, "城市列表應該包含至少一個城市")
		})
	})
}
