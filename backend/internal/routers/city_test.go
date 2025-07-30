package routers

import (
	"backend/internal/models"
	"backend/internal/tests"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewCityRouter(t *testing.T) {
	t.Run("單例模式測試", func(t *testing.T) {
		router1 := NewCityRouter()
		router2 := NewCityRouter()

		assert.Same(t, router1, router2, "應該返回相同的實例")
		assert.NotNil(t, router1, "實例不應該為 nil")
		assert.NotNil(t, router1.CityService, "CityService 不應該為 nil")
	})
}

func TestCityRouterBind(t *testing.T) {
	t.Run("路由綁定測試", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()
		apiGroup := router.Group("/api")

		cityRouter := NewCityRouter()
		cityRouter.Bind(apiGroup)

		routes := router.Routes()
		assert.NotEmpty(t, routes, "應該有註冊的路由")
	})
}

func TestCityRouterGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功獲取所有城市", func(t *testing.T) {
		server, apiRouter, cleanup := tests.SetupTestServer()
		defer cleanup()
		cityRouter := NewCityRouter()
		cityRouter.Bind(apiRouter)

		req, _ := http.NewRequest("GET", "/api/city/all", nil)
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		cities := models.CityGetAllResponse{}
		err := json.Unmarshal(recorder.Body.Bytes(), &cities)
		assert.NoError(t, err)
		assert.Greater(t, len(cities), 0, "城市列表應該包含至少一個城市")
	})
}
