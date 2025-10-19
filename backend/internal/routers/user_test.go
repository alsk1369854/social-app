package routers

import (
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/tests"
	"encoding/json"
	"log"
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

	t.Run("Login", func(t *testing.T) {
		t.Run("登入失敗 - 電子郵件用戶不存在", func(t *testing.T) {
			body := &models.UserLoginRequest{
				Email:    pkg.GetRandomString(5),
				Password: "passwerd123",
			}
			buf, _ := httpUtils.ToJSONBuffer(body)
			req, _ := http.NewRequest("POST", "/api/user/login", buf)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			response := &models.ErrorResponse{}
			err := json.Unmarshal(recorder.Body.Bytes(), response)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, 400, recorder.Code, "應該回傳 400 表示登入失敗")
			assert.Equal(t, "email not found", response.Error, "Error message should match")
		})

		t.Run("登入失敗 - 密碼錯誤", func(t *testing.T) {
			// 1. 創建一個新用戶
			username := pkg.GetRandomString(5)
			payload := models.UserRegisterRequest{

				Email:    username + "@example.com",
				Password: "password123",
			}
			buf, _ := httpUtils.ToJSONBuffer(payload)
			req, _ := http.NewRequest("POST", "/api/user/register", buf)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, 200, recorder.Code, "應該回傳 200 表示註冊成功")

			// 2. 嘗試使用錯誤的密碼登入
			loginPayload := models.UserLoginRequest{
				Email:    payload.Email,
				Password: "wrongpassword",
			}
			loginBuf, _ := httpUtils.ToJSONBuffer(loginPayload)
			loginReq, _ := http.NewRequest("POST", "/api/user/login", loginBuf)
			loginReq.Header.Set("Content-Type", "application/json")
			loginRecorder := httptest.NewRecorder()
			server.ServeHTTP(loginRecorder, loginReq)
			assert.Equal(t, 400, loginRecorder.Code, "應該回傳 400 表示登入失敗")
			response := &models.ErrorResponse{}
			err := json.Unmarshal(loginRecorder.Body.Bytes(), response)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, "incorrect email or password", response.Error)
		})

		t.Run("成功登入", func(t *testing.T) {
			// 1. 創建一個新用戶
			username := pkg.GetRandomString(5)
			payload := models.UserRegisterRequest{

				Email:    username + "@example.com",
				Password: "password123",
			}
			buf, _ := httpUtils.ToJSONBuffer(payload)
			req, _ := http.NewRequest("POST", "/api/user/register", buf)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, 200, recorder.Code, "應該回傳 200 表示註冊成功")

			// 2. 嘗試登入
			loginPayload := models.UserLoginRequest{
				Email:    payload.Email,
				Password: payload.Password,
			}
			loginBuf, _ := httpUtils.ToJSONBuffer(loginPayload)
			loginReq, _ := http.NewRequest("POST", "/api/user/login", loginBuf)
			loginReq.Header.Set("Content-Type", "application/json")
			loginRecorder := httptest.NewRecorder()
			server.ServeHTTP(loginRecorder, loginReq)
			assert.Equal(t, 200, loginRecorder.Code, "應該回傳 200 表示登入成功")
			// response := &models.UserLoginResponse{}
			response := &models.UserLoginResponse{}
			log.Println(loginRecorder.Body.String())
			err := json.Unmarshal(loginRecorder.Body.Bytes(), response)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.NotEmpty(t, response.AccessToken, "Token should not be empty")
			assert.Equal(t, response.Username, payload.Username)
			assert.Equal(t, response.Email, payload.Email)
			assert.NotEmpty(t, response.ID, "User ID should not be empty")
		})

		t.Run("登入失敗 - 錯誤的電子郵件或密碼", func(t *testing.T) {
			// 嘗試登入不存在的用戶
			loginPayload := models.UserLoginRequest{
				Email:    "nonexistent@example.com",
				Password: "wrongpassword",
			}
			loginBuf, _ := httpUtils.ToJSONBuffer(loginPayload)
			loginReq, _ := http.NewRequest("POST", "/api/user/login", loginBuf)
			loginReq.Header.Set("Content-Type", "application/json")
			loginRecorder := httptest.NewRecorder()
			server.ServeHTTP(loginRecorder, loginReq)

			assert.Equal(t, 400, loginRecorder.Code, "應該回傳 400 表示登入失敗")
		})

	})

	t.Run("Register", func(t *testing.T) {
		t.Run("成功註冊", func(t *testing.T) {
			username := pkg.GetRandomString(5)
			payload := models.UserRegisterRequest{

				Email:    username + "@example.com",
				Password: "password123",
			}
			buf, _ := httpUtils.ToJSONBuffer(payload)
			req, _ := http.NewRequest("POST", "/api/user/register", buf)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			expectedResponse := &models.UserRegisterResponse{}
			err := json.Unmarshal(recorder.Body.Bytes(), expectedResponse)

			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, 200, recorder.Code)
			assert.NotEmpty(t, expectedResponse.ID, "User ID should not be empty")
			assert.Equal(t, expectedResponse.Username, payload.Username)
			assert.Equal(t, expectedResponse.Email, payload.Email)
		})

		t.Run("註冊包含地址會員", func(t *testing.T) {
			citySlice, _ := userRouter.CityService.GetAll(ctx)
			username := pkg.GetRandomString(5)
			payload := models.UserRegisterRequest{

				Email:    username + "@example.com",
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
			username := pkg.GetRandomString(5)
			payload := models.UserRegisterRequest{

				Email:    username + "@example.com",
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
			username := pkg.GetRandomString(5)
			payload := models.UserRegisterRequest{

				Email:    username + "@example.com",
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
			username := pkg.GetRandomString(5)
			payload := models.UserRegisterRequest{

				Email:    username + "@example.com",
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

		t.Run("密碼長度不符合 6~12 個字元", func(t *testing.T) {
			t.Run("少於 6 個字元", func(t *testing.T) {
				username := pkg.GetRandomString(5)
				payload := models.UserRegisterRequest{

					Email:    username + "@example.com",
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
				username := pkg.GetRandomString(5)
				payload := models.UserRegisterRequest{

					Email:    username + "@example.com",
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
