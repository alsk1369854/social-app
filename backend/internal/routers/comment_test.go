package routers

import (
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/tests"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommentRouter(t *testing.T) {
	httpUtils := pkg.NewHTTPUtils()

	server, apiRouter, _, _, cleanup := tests.SetupTestServer("test_comment_router")
	defer cleanup()

	NewCommentRouter().Bind(apiRouter)
	NewPostRouter().Bind(apiRouter)
	NewUserRouter().Bind(apiRouter)

	userData, loginData, err := tests.SetupTestUser(server)
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.NotNil(t, loginData)

	postData, err := tests.SetupTestPost(server, loginData.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, postData)

	t.Run("GetCommentsByPostID", func(t *testing.T) {
		t.Run("成功獲取評論列表", func(t *testing.T) {
			// 創建兩個評論
			for i := 0; i < 2; i++ {
				commentCreateRequest := &models.CommentCreateRequest{
					PostID:  postData.ID,
					Content: "這是一個測試評論 " + fmt.Sprint(i+1),
				}
				bufCommentCreateRequest, _ := httpUtils.ToJSONBuffer(commentCreateRequest)
				reqCreateComment, _ := http.NewRequest("POST", "/api/comment", bufCommentCreateRequest)
				reqCreateComment.Header.Set("Content-Type", "application/json")
				reqCreateComment.Header.Set("Authorization", loginData.AccessToken)
				recorder := httptest.NewRecorder()
				server.ServeHTTP(recorder, reqCreateComment)
				assert.Equal(t, 200, recorder.Code, "應該回傳 200 表示評論創建成功")
			}

			// 獲取評論列表
			reqGetComments, _ := http.NewRequest("GET", "/api/comment/list/post/"+postData.ID.String(), nil)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, reqGetComments)
			assert.Equal(t, 200, recorder.Code, "應該回傳 200 表示成功獲取評論列表")
			responseBody := make([]models.CommentGetListByPostIDResponseItem, 0)
			err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Greater(t, len(responseBody), 0, "應該至少有一條評論")
			for _, comment := range responseBody {
				assert.Equal(t, postData.ID, comment.PostID, "Post ID should match the created post")
				assert.NotEmpty(t, comment.Content, "Comment content should not be empty")
				assert.NotEmpty(t, comment.UserID, "User ID should not be empty")
			}
		})
	})

	t.Run("CreateComment", func(t *testing.T) {

		t.Run("成功創建評論", func(t *testing.T) {
			commentCreateRequest := &models.CommentCreateRequest{
				PostID:  postData.ID,
				Content: "這是一個測試評論",
			}
			bufCommentCreateRequest, _ := httpUtils.ToJSONBuffer(commentCreateRequest)
			reqCreateComment, _ := http.NewRequest("POST", "/api/comment", bufCommentCreateRequest)
			reqCreateComment.Header.Set("Content-Type", "application/json")
			reqCreateComment.Header.Set("Authorization", loginData.AccessToken)

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, reqCreateComment)
			responseBody := &models.CommentCreateResponse{}
			err := json.Unmarshal(recorder.Body.Bytes(), responseBody)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.Equal(t, 200, recorder.Code, "應該回傳 200 表示評論創建成功")
			assert.NotEmpty(t, responseBody.ID, "Comment ID should not be empty")
			assert.Equal(t, commentCreateRequest.Content, responseBody.Content, "Comment content should match the request")
			assert.Equal(t, postData.ID, responseBody.PostID, "Post ID should match the request")
			assert.Equal(t, loginData.ID, responseBody.UserID, "User ID should match the creator's ID")
			assert.Nil(t, responseBody.ParentID, "Parent ID should be nil for top-level comments")
		})
	})
}
