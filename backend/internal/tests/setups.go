package tests

import (
	"backend/internal/database"
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/servers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTestDB(dbFile string) (*gorm.DB, func()) {
	os.Remove(dbFile)
	db := database.SetupSQLite(&database.SQLiteConfig{
		DBFile:    dbFile, // ":memery:"
		EnableLog: false,
	})

	cleanup := func() {
		os.Remove(dbFile)
	}

	return db, cleanup
}

func SetupTestContext(dbFile string) (*gin.Context, *gorm.DB, func()) {
	db, cleanup := SetupTestDB(dbFile)

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set(middlewares.CONTEXT_KEY_GORM_DB, db)

	return ctx, db, cleanup
}

func SetupTestServer(dbFile string) (*gin.Engine, *gin.RouterGroup, *gin.Context, *gorm.DB, func()) {
	ctx, db, cleanup := SetupTestContext(dbFile)

	gin.SetMode(gin.TestMode)
	server, apiRouter := servers.SetupGin(&servers.GinConfig{
		DB: db,
	})

	return server, apiRouter, ctx, db, cleanup
}

// Required PostRouter.Bind
func SetupTestPost(server *gin.Engine, accessToken string) (*models.PostCreateResponse, error) {
	httpUtils := pkg.NewHTTPUtils()

	createPostReqBody := &models.PostCreateRequest{
		ImageURL: nil,
		Content:  "這是一個測試 Post",
		Tags:     []string{"測試", "Post"},
	}
	createPostReqBodyBuf, _ := httpUtils.ToJSONBuffer(createPostReqBody)
	req, _ := http.NewRequest("POST", "/api/post", createPostReqBodyBuf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", accessToken)
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)
	createPostRespBody := &models.PostCreateResponse{}
	if err := json.Unmarshal(recorder.Body.Bytes(), createPostRespBody); err != nil {
		return nil, err
	}
	return createPostRespBody, nil
}

// Required UserRouter.Bind
func SetupTestUser(server *gin.Engine) (*models.UserRegisterRequest, *models.UserLoginResponse, error) {
	httpUtils := pkg.NewHTTPUtils()

	// 1. 創建一個新用戶
	username := pkg.GetRandomString(5)
	registerReqBody := &models.UserRegisterRequest{
		Username: username,
		Email:    username + "@example.com",
		Password: "password123",
	}
	registerReqBodyBuf, _ := httpUtils.ToJSONBuffer(registerReqBody)
	req, _ := http.NewRequest("POST", "/api/user/register", registerReqBodyBuf)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)
	registerRespBody := &models.UserRegisterResponse{}
	if err := json.Unmarshal(recorder.Body.Bytes(), registerRespBody); err != nil {
		return nil, nil, err
	}

	// 2. 嘗試登入獲取 token
	loginReq := &models.UserLoginRequest{
		Email:    registerReqBody.Email,
		Password: registerReqBody.Password,
	}
	loginReqBuf, _ := httpUtils.ToJSONBuffer(loginReq)
	reqLogin, _ := http.NewRequest("POST", "/api/user/login", loginReqBuf)
	reqLogin.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	server.ServeHTTP(recorder, reqLogin)
	loginRespBody := &models.UserLoginResponse{}
	if err := json.Unmarshal(recorder.Body.Bytes(), loginRespBody); err != nil {
		return nil, nil, err
	}

	return registerReqBody, loginRespBody, nil
}
