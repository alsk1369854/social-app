package routers

import (
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/tests"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRouter(t *testing.T) {
	httpUtils := pkg.NewHTTPUtils()
	server, apiRouter, ctx, _, cleanup := tests.SetupTestServer("test_user_router.db")
	defer cleanup()
	userRouter := NewUserRouter()
	userRouter.Bind(apiRouter)

	t.Run("Register", func(t *testing.T) {
		t.Run("Successful Registration", func(t *testing.T) {
			payload := models.UserRegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			}
			buf, _ := httpUtils.ToJSONBuffer(payload)
			req, _ := http.NewRequest("POST", "/api/user/register", buf)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			expectedStatus := http.StatusOK
			expectedResponse := &models.UserRegisterResponse{}
			err := json.Unmarshal(recorder.Body.Bytes(), expectedResponse)

			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, expectedStatus, recorder.Code)
			assert.NotEmpty(t, expectedResponse.ID, "User ID should not be empty")
			assert.Equal(t, expectedResponse.Username, payload.Username)
			assert.Equal(t, expectedResponse.Email, payload.Email)
		})

		t.Run("註冊包含地址會員", func(t *testing.T) {
			citySlice, _ := userRouter.CityService.GetAll(ctx)
			payload := models.UserRegisterRequest{
				Username: "testuser2",
				Email:    "test@test.com",
				Password: "password123",
				Address: &models.UserRegisterRequestAddress{
					CityID: citySlice[0].ID,
					Street: "123 Test St",
				},
			}
			buf, _ := httpUtils.ToJSONBuffer(payload)
			req, _ := http.NewRequest("POST", "/api/user/register", buf)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			expectedStatus := http.StatusOK
			expectedResponse := &models.UserRegisterResponse{}
			err := json.Unmarshal(recorder.Body.Bytes(), expectedResponse)

			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, expectedStatus, recorder.Code)
			assert.NotEmpty(t, expectedResponse.ID, "User ID should not be empty")
			assert.Equal(t, expectedResponse.Username, payload.Username)
			assert.Equal(t, expectedResponse.Email, payload.Email)
			assert.Equal(t, expectedResponse.Address.CityID, payload.Address.CityID)
			assert.Equal(t, expectedResponse.Address.Street, payload.Address.Street)
		})

		t.Run("註冊會員 - 包含年齡", func(t *testing.T) {
			payload := models.UserRegisterRequest{
				Username: "testuser3",
				Email:    "test3@example.com",
				Password: "password123",
				Age:      pkg.GetPointer(int64(30)),
			}
			buf, _ := httpUtils.ToJSONBuffer(payload)
			req, _ := http.NewRequest("POST", "/api/user/register", buf)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			expectedStatus := http.StatusOK
			expectedResponse := &models.UserRegisterResponse{}
			err := json.Unmarshal(recorder.Body.Bytes(), expectedResponse)

			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, expectedStatus, recorder.Code)
			assert.NotEmpty(t, expectedResponse.ID, "User ID should not be empty")
			assert.Equal(t, expectedResponse.Username, payload.Username)
			assert.Equal(t, expectedResponse.Email, payload.Email)
			assert.Equal(t, expectedResponse.Age, payload.Age)
		})

		t.Run("註冊會員 - 包含地址和年齡", func(t *testing.T) {
			citySlice, _ := userRouter.CityService.GetAll(ctx)
			payload := models.UserRegisterRequest{
				Username: "test4123",
				Email:    "test4123@example.com",
				Password: "password123",
				Age:      pkg.GetPointer(int64(30)),
				Address: &models.UserRegisterRequestAddress{
					CityID: citySlice[0].ID,
					Street: "123 Test St",
				},
			}
			buf, _ := httpUtils.ToJSONBuffer(payload)
			req, _ := http.NewRequest("POST", "/api/user/register", buf)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			expectedStatus := http.StatusOK
			expectedResponse := &models.UserRegisterResponse{}
			err := json.Unmarshal(recorder.Body.Bytes(), expectedResponse)

			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, expectedStatus, recorder.Code)
			assert.NotEmpty(t, expectedResponse.ID, "User ID should not be empty")
			assert.Equal(t, expectedResponse.Username, payload.Username)
			assert.Equal(t, expectedResponse.Email, payload.Email)
			assert.Equal(t, expectedResponse.Age, payload.Age)
			assert.Equal(t, expectedResponse.Address.CityID, payload.Address.CityID)
			assert.Equal(t, expectedResponse.Address.Street, payload.Address.Street)
		})

		t.Run("註冊失敗 - 重複的電子郵件", func(t *testing.T) {
			// 準備重複的 email payload
			payload := models.UserRegisterRequest{
				Username: "existinguser",
				Email:    "testuserasdf3@example.com",
				Password: "password123",
			}

			// === 1. 先註冊一次使用者（預期成功）===
			payloadBuf, _ := httpUtils.ToJSONBuffer(payload)
			req, _ := http.NewRequest("POST", "/api/user/register", payloadBuf)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, 200, recorder.Code, "第一次註冊應該成功")

			// === 2. 再次註冊相同 email（預期失敗）===
			payloadBuf, _ = httpUtils.ToJSONBuffer(payload)
			req, _ = http.NewRequest("POST", "/api/user/register", payloadBuf)
			recorder = httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			assert.Equal(t, 400, recorder.Code, "應該回傳 400 表示 email 重複")

			var resp models.ErrorResponse
			err := json.Unmarshal(recorder.Body.Bytes(), &resp)
			require.NoError(t, err, "回傳應為合法 JSON")
			assert.Equal(t, "email already exists", resp.Error)
		})

		t.Run("重複名稱", func(t *testing.T) {
			// 準備重複的 username payload
			payload := models.UserRegisterRequest{
				Username: "duplicate-name-testuser",
				Email:    "duplicate-name-testuser@example.com",
				Password: "password123",
			}
			// === 1. 先註冊一次使用者（預期成功）===
			payloadBuf, _ := httpUtils.ToJSONBuffer(payload)
			req, _ := http.NewRequest("POST", "/api/user/register", payloadBuf)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, 200, recorder.Code, "第一次註冊應該成功")

			// === 2. 再次註冊相同 username（預期失敗）===
			payload.Email = "duplicate-name-testuser-new@example.com"
			payloadBuf, _ = httpUtils.ToJSONBuffer(payload)
			req, _ = http.NewRequest("POST", "/api/user/register", payloadBuf)
			recorder = httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			assert.Equal(t, 400, recorder.Code, "應該回傳 400 表示 username 重複")

			var resp models.ErrorResponse
			err := json.Unmarshal(recorder.Body.Bytes(), &resp)
			require.NoError(t, err, "回傳應為合法 JSON")
			assert.Equal(t, "username already exists", resp.Error)
		})

		t.Run("密碼長度不符合 6~12 個字元", func(t *testing.T) {
			t.Run("少於 6 個字元", func(t *testing.T) {
				payload := models.UserRegisterRequest{
					Username: "shortpassword",
					Email:    "shortpassword@example.com",
					Password: "123",
				}

				buf, _ := httpUtils.ToJSONBuffer(payload)
				req, _ := http.NewRequest("POST", "/api/user/register", buf)
				req.Header.Set("Content-Type", "application/json")
				recorder := httptest.NewRecorder()
				server.ServeHTTP(recorder, req)

				assert.Equal(t, 400, recorder.Code, "應該回傳 400 表示密碼長度不符合")

				var resp models.ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				require.NoError(t, err, "回傳應為合法 JSON")
				assert.Equal(t, "password length must be between 6 and 12 characters", resp.Error)
			})

			t.Run("超過 12 個字元", func(t *testing.T) {
				payload := models.UserRegisterRequest{
					Username: "longpassword",
					Email:    "longpassword@example.com",
					Password: "1234567890123",
				}

				buf, _ := httpUtils.ToJSONBuffer(payload)
				req, _ := http.NewRequest("POST", "/api/user/register", buf)
				req.Header.Set("Content-Type", "application/json")
				recorder := httptest.NewRecorder()
				server.ServeHTTP(recorder, req)

				assert.Equal(t, 400, recorder.Code, "應該回傳 400 表示密碼長度不符合")

				var resp models.ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				require.NoError(t, err, "回傳應為合法 JSON")
				assert.Equal(t, "password length must be between 6 and 12 characters", resp.Error)
			})
		})
	})
}
