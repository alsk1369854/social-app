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

	NewUserRouter().Bind(apiRouter)
	NewPostRouter().Bind(apiRouter)

	// 1. 創建一個新用戶
	userData, loginData, err := tests.SetupTestUser(server)
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.NotNil(t, loginData)

	t.Run("獲取作者 Posts", func(t *testing.T) {

		t.Run("失敗 - 使用者不存在", func(t *testing.T) {
			authorID := uuid.New().String()
			reqGetPostsByAuthorID, _ := http.NewRequest("GET", "/api/post/list/author/"+authorID+"/offset/0/limit/10", nil)
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
				ImageURL: nil,
				Content:  "這是一個測試 Post #測試 #Post",
			}
			bufReqCreatePostBody, _ := httpUtils.ToJSONBuffer(reqCreatePostBody)
			reqCreatePost, _ := http.NewRequest("POST", "/api/post", bufReqCreatePostBody)
			reqCreatePost.Header.Set("Content-Type", "application/json")
			reqCreatePost.Header.Set("Authorization", loginData.AccessToken)
			recorderCreatePost := httptest.NewRecorder()
			server.ServeHTTP(recorderCreatePost, reqCreatePost)
			assert.Equal(t, 200, recorderCreatePost.Code, "應該回傳 200 表示創建 Post 成功")
			respCreatePostBody := &models.PostCreateResponse{}
			err = json.Unmarshal(recorderCreatePost.Body.Bytes(), respCreatePostBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, reqCreatePostBody.Content, respCreatePostBody.Content, "Post content should match the request content")
			assert.Equal(t, 2, len(respCreatePostBody.TagIDs), "Post tags should match the request tags")
			assert.NotEmpty(t, respCreatePostBody.ID, "Post ID should not be empty")
			assert.NotEmpty(t, respCreatePostBody.CreatedAt, "Post CreatedAt should not be empty")
			assert.NotEmpty(t, respCreatePostBody.UpdatedAt, "Post UpdatedAt should not be empty")

			// 1. 使用 token 獲取作者的 Posts
			reqGetPostsByAuthorID, _ := http.NewRequest("GET", "/api/post/list/author/"+loginData.ID.String()+"/offset/0/limit/10", nil)
			reqGetPostsByAuthorID.Header.Set("Authorization", loginData.AccessToken)
			recorderGetPostsByAuthorID := httptest.NewRecorder()
			server.ServeHTTP(recorderGetPostsByAuthorID, reqGetPostsByAuthorID)
			assert.Equal(t, 200, recorderGetPostsByAuthorID.Code, "應該回傳 200 表示獲取 Posts 成功")
			respGetPostsByAuthorIDBody := &models.PaginationResponse[models.PostGetPostsByAuthorIDResponseItem]{}
			err = json.Unmarshal(recorderGetPostsByAuthorID.Body.Bytes(), respGetPostsByAuthorIDBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.NotEmpty(t, respGetPostsByAuthorIDBody.Data, "Posts data should not be empty")
			assert.Equal(t, loginData.ID, respGetPostsByAuthorIDBody.Data[0].AuthorID, "Post AuthorID should match the logged-in user ID")
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
				ImageURL: nil,
				Content:  "",
			}
			bufReqCreatePostBody, _ := httpUtils.ToJSONBuffer(reqCreatePostBody)
			reqCreatePost, _ := http.NewRequest("POST", "/api/post", bufReqCreatePostBody)
			reqCreatePost.Header.Set("Content-Type", "application/json")
			reqCreatePost.Header.Set("Authorization", loginData.AccessToken)
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
				ImageURL: nil,
				Content:  pkg.GetRandomString(301) + "  #測試 #Post", // 超過 300 字元
			}
			bufReqCreatePostBody, _ := httpUtils.ToJSONBuffer(reqCreatePostBody)
			reqCreatePost, _ := http.NewRequest("POST", "/api/post", bufReqCreatePostBody)
			reqCreatePost.Header.Set("Content-Type", "application/json")
			reqCreatePost.Header.Set("Authorization", loginData.AccessToken)
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
				ImageURL: nil,
				Content:  "這是一個測試 Post #測試 #Post",
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
				ImageURL: nil,
				Content:  "這是一個測試 Post #測試 #Post",
			}
			bufReqCreatePostBody, _ := httpUtils.ToJSONBuffer(reqCreatePostBody)
			reqCreatePost, _ := http.NewRequest("POST", "/api/post", bufReqCreatePostBody)
			reqCreatePost.Header.Set("Content-Type", "application/json")
			reqCreatePost.Header.Set("Authorization", loginData.AccessToken)
			recorderCreatePost := httptest.NewRecorder()
			server.ServeHTTP(recorderCreatePost, reqCreatePost)
			assert.Equal(t, 200, recorderCreatePost.Code, "應該回傳 200 表示創建 Post 成功")
			respCreatePostBody := &models.PostCreateResponse{}
			err = json.Unmarshal(recorderCreatePost.Body.Bytes(), respCreatePostBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, reqCreatePostBody.Content, respCreatePostBody.Content, "Post content should match the request content")
			assert.Equal(t, 2, len(respCreatePostBody.TagIDs), "Post tags should match the request tags")
			assert.NotEmpty(t, respCreatePostBody.ID, "Post ID should not be empty")
			assert.NotEmpty(t, respCreatePostBody.CreatedAt, "Post CreatedAt should not be empty")
			assert.NotEmpty(t, respCreatePostBody.UpdatedAt, "Post UpdatedAt should not be empty")
		})

	})

}
