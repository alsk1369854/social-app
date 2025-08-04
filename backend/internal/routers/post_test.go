package routers

import (
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/tests"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostRouter(t *testing.T) {
	httpUtils := pkg.NewHTTPUtils()

	server, apiRouter, _, _, cleanup := tests.SetupTestServer("test_post_router.db")
	defer cleanup()

	postRouter := NewPostRouter()
	postRouter.Bind(apiRouter)
	userRouter := NewUserRouter()
	userRouter.Bind(apiRouter)

	// 1. 創建一個新用戶
	username := pkg.GetRandomString(5)
	reqReqRegisterBody := &models.UserRegisterRequest{
		Username: username,
		Email:    username + "@example.com",
		Password: "password123",
	}
	bugReqRegisterBody, _ := httpUtils.ToJSONBuffer(reqReqRegisterBody)
	reqRegister, _ := http.NewRequest("POST", "/api/user/register", bugReqRegisterBody)
	reqRegister.Header.Set("Content-Type", "application/json")
	recorderRegister := httptest.NewRecorder()
	server.ServeHTTP(recorderRegister, reqRegister)
	assert.Equal(t, 200, recorderRegister.Code, "應該回傳 200 表示註冊成功")
	respRegisterBody := &models.UserRegisterResponse{}
	err := json.Unmarshal(recorderRegister.Body.Bytes(), respRegisterBody)
	assert.NoError(t, err, "Response should be valid JSON")

	// 2. 嘗試登入獲取 token
	reqLoginBody := &models.UserLoginRequest{
		Email:    reqReqRegisterBody.Email,
		Password: reqReqRegisterBody.Password,
	}
	bufReqLoginBody, _ := httpUtils.ToJSONBuffer(reqLoginBody)
	reqLogin, _ := http.NewRequest("POST", "/api/user/login", bufReqLoginBody)
	reqLogin.Header.Set("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	server.ServeHTTP(recorderLogin, reqLogin)
	respLoginBody := &models.UserLoginResponse{}
	err = json.Unmarshal(recorderLogin.Body.Bytes(), respLoginBody)
	assert.NoError(t, err, "Response should be valid JSON")
	assert.Equal(t, 200, recorderLogin.Code, "應該回傳 200 表示登入成功")
	assert.NotEmpty(t, respLoginBody.AccessToken, "Access token should not be empty")

	t.Run("獲取作者 Posts", func(t *testing.T) {

		t.Run("失敗 - 使用者不存在", func(t *testing.T) {
			authorID := uuid.New().String()
			reqGetPostsByAuthorID, _ := http.NewRequest("GET", "/api/post/author/"+authorID+"/offset/0/limit/10", nil)
			recorderGetPostsByAuthorID := httptest.NewRecorder()
			server.ServeHTTP(recorderGetPostsByAuthorID, reqGetPostsByAuthorID)
			assert.Equal(t, 404, recorderGetPostsByAuthorID.Code, "應該回傳 404 表示使用者不存在")
			respGetPostsByAuthorIDBody := &models.ErrorResponse{}
			err := json.Unmarshal(recorderGetPostsByAuthorID.Body.Bytes(), respGetPostsByAuthorIDBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, "author not found", respGetPostsByAuthorIDBody.Error, "Error message should match")
		})

		t.Run("成功獲取 Posts", func(t *testing.T) {
			// 新增一個 Post 以便後續測試
			reqCreatePostBody := &models.PostCreateRequest{
				AuthorID: respRegisterBody.ID,
				ImageURL: nil,
				Content:  "這是一個測試 Post",
				Tags:     []string{"測試", "Post"},
			}
			bufReqCreatePostBody, _ := httpUtils.ToJSONBuffer(reqCreatePostBody)
			reqCreatePost, _ := http.NewRequest("POST", "/api/post", bufReqCreatePostBody)
			reqCreatePost.Header.Set("Content-Type", "application/json")
			reqCreatePost.Header.Set("Authorization", respLoginBody.AccessToken)
			recorderCreatePost := httptest.NewRecorder()
			server.ServeHTTP(recorderCreatePost, reqCreatePost)
			assert.Equal(t, 200, recorderCreatePost.Code, "應該回傳 200 表示創建 Post 成功")
			respCreatePostBody := &models.PostCreateResponse{}
			err = json.Unmarshal(recorderCreatePost.Body.Bytes(), respCreatePostBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, reqCreatePostBody.AuthorID, respCreatePostBody.AuthorID, "Post AuthorID should match the created user ID")
			assert.Equal(t, reqCreatePostBody.Content, respCreatePostBody.Content, "Post content should match the request content")
			assert.Equal(t, len(reqCreatePostBody.Tags), len(respCreatePostBody.TagIDs), "Post tags should match the request tags")
			assert.NotEmpty(t, respCreatePostBody.ID, "Post ID should not be empty")
			assert.NotEmpty(t, respCreatePostBody.CreatedAt, "Post CreatedAt should not be empty")
			assert.NotEmpty(t, respCreatePostBody.UpdatedAt, "Post UpdatedAt should not be empty")

			// 1. 使用 token 獲取作者的 Posts
			reqGetPostsByAuthorID, _ := http.NewRequest("GET", "/api/post/author/"+respLoginBody.ID.String()+"/offset/0/limit/10", nil)
			reqGetPostsByAuthorID.Header.Set("Authorization", respLoginBody.AccessToken)
			recorderGetPostsByAuthorID := httptest.NewRecorder()
			server.ServeHTTP(recorderGetPostsByAuthorID, reqGetPostsByAuthorID)
			assert.Equal(t, 200, recorderGetPostsByAuthorID.Code, "應該回傳 200 表示獲取 Posts 成功")
			respGetPostsByAuthorIDBody := &models.PaginationResponse[models.PostGetPostsByAuthorIDResponseItem]{}
			err = json.Unmarshal(recorderGetPostsByAuthorID.Body.Bytes(), respGetPostsByAuthorIDBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.NotEmpty(t, respGetPostsByAuthorIDBody.Data, "Posts data should not be empty")
			assert.Equal(t, respLoginBody.ID, respGetPostsByAuthorIDBody.Data[0].AuthorID, "Post AuthorID should match the logged-in user ID")
			assert.NotEmpty(t, respGetPostsByAuthorIDBody.Data[0].ID, "Post ID should not be empty")
			assert.NotEmpty(t, respGetPostsByAuthorIDBody.Data[0].CreatedAt, "Post CreatedAt should not be empty")
			assert.NotEmpty(t, respGetPostsByAuthorIDBody.Data[0].UpdatedAt, "Post UpdatedAt should not be empty")
			assert.NotEmpty(t, respGetPostsByAuthorIDBody.Data[0].Tags, "Post Tags should not be empty")
			assert.NotNil(t, respGetPostsByAuthorIDBody.Data[0].LikedCount, "Post LikedCount should not be empty")
		})
	})

	t.Run("創建 Post", func(t *testing.T) {

		t.Run("創建失敗 - 內容為空", func(t *testing.T) {
			reqCreatePostBody := &models.PostCreateRequest{
				AuthorID: respLoginBody.ID,
				ImageURL: nil,
				Content:  "",
				Tags:     []string{"測試", "Post"},
			}
			bufReqCreatePostBody, _ := httpUtils.ToJSONBuffer(reqCreatePostBody)
			reqCreatePost, _ := http.NewRequest("POST", "/api/post", bufReqCreatePostBody)
			reqCreatePost.Header.Set("Content-Type", "application/json")
			reqCreatePost.Header.Set("Authorization", respLoginBody.AccessToken)
			recorderCreatePost := httptest.NewRecorder()
			server.ServeHTTP(recorderCreatePost, reqCreatePost)
			assert.Equal(t, 400, recorderCreatePost.Code, "應該回傳 400 表示請求錯誤")
			respCreatePostBody := &models.ErrorResponse{}
			err := json.Unmarshal(recorderCreatePost.Body.Bytes(), respCreatePostBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, "invalid request body", respCreatePostBody.Error, "Error message should match")
		})

		t.Run("創建失敗 - 內容超過 300 字元", func(t *testing.T) {
			reqCreatePostBody := &models.PostCreateRequest{
				AuthorID: respLoginBody.ID,
				ImageURL: nil,
				Content:  pkg.GetRandomString(301), // 超過 300 字元
				Tags:     []string{"測試", "Post"},
			}
			bufReqCreatePostBody, _ := httpUtils.ToJSONBuffer(reqCreatePostBody)
			reqCreatePost, _ := http.NewRequest("POST", "/api/post", bufReqCreatePostBody)
			reqCreatePost.Header.Set("Content-Type", "application/json")
			reqCreatePost.Header.Set("Authorization", respLoginBody.AccessToken)
			recorderCreatePost := httptest.NewRecorder()
			server.ServeHTTP(recorderCreatePost, reqCreatePost)
			assert.Equal(t, 400, recorderCreatePost.Code, "應該回傳 400 表示請求錯誤")
			respCreatePostBody := &models.ErrorResponse{}
			err := json.Unmarshal(recorderCreatePost.Body.Bytes(), respCreatePostBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, "content characters must be between 0 and 300", respCreatePostBody.Error, "Error message should match")
		})

		t.Run("創建失敗 - 缺少 Authorization", func(t *testing.T) {
			reqCreatePostBody := &models.PostCreateRequest{
				AuthorID: respLoginBody.ID,
				ImageURL: nil,
				Content:  "這是一個測試 Post",
				Tags:     []string{"測試", "Post"},
			}
			bufReqCreatePostBody, _ := httpUtils.ToJSONBuffer(reqCreatePostBody)
			reqCreatePost, _ := http.NewRequest("POST", "/api/post", bufReqCreatePostBody)
			reqCreatePost.Header.Set("Content-Type", "application/json")
			recorderCreatePost := httptest.NewRecorder()
			server.ServeHTTP(recorderCreatePost, reqCreatePost)
			assert.Equal(t, 401, recorderCreatePost.Code, "應該回傳 401 表示未授權")
			respCreatePostBody := &models.ErrorResponse{}
			err := json.Unmarshal(recorderCreatePost.Body.Bytes(), respCreatePostBody)
			assert.NoError(t, err, "Response should be valid JSON")
		})

		t.Run("成功創建 Post", func(t *testing.T) {

			// 3. 使用 token 創建 Post
			reqCreatePostBody := &models.PostCreateRequest{
				AuthorID: respRegisterBody.ID,
				ImageURL: nil,
				Content:  "這是一個測試 Post",
				Tags:     []string{"測試", "Post"},
			}
			bufReqCreatePostBody, _ := httpUtils.ToJSONBuffer(reqCreatePostBody)
			reqCreatePost, _ := http.NewRequest("POST", "/api/post", bufReqCreatePostBody)
			reqCreatePost.Header.Set("Content-Type", "application/json")
			reqCreatePost.Header.Set("Authorization", respLoginBody.AccessToken)
			recorderCreatePost := httptest.NewRecorder()
			server.ServeHTTP(recorderCreatePost, reqCreatePost)
			assert.Equal(t, 200, recorderCreatePost.Code, "應該回傳 200 表示創建 Post 成功")
			respCreatePostBody := &models.PostCreateResponse{}
			err = json.Unmarshal(recorderCreatePost.Body.Bytes(), respCreatePostBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, reqCreatePostBody.AuthorID, respCreatePostBody.AuthorID, "Post AuthorID should match the created user ID")
			assert.Equal(t, reqCreatePostBody.Content, respCreatePostBody.Content, "Post content should match the request content")
			assert.Equal(t, len(reqCreatePostBody.Tags), len(respCreatePostBody.TagIDs), "Post tags should match the request tags")
			assert.NotEmpty(t, respCreatePostBody.ID, "Post ID should not be empty")
			assert.NotEmpty(t, respCreatePostBody.CreatedAt, "Post CreatedAt should not be empty")
			assert.NotEmpty(t, respCreatePostBody.UpdatedAt, "Post UpdatedAt should not be empty")
		})

	})

}
