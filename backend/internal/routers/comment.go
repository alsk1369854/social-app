package routers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/services"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentRouter struct {
	ErrorUtils *pkg.ErrorUtils

	CommentService *services.CommentService
	PostService    *services.PostService
}

var commentRouterOnce sync.Once
var commentRouter *CommentRouter

func NewCommentRouter() *CommentRouter {
	commentRouterOnce.Do(func() {
		commentRouter = &CommentRouter{
			ErrorUtils: pkg.NewErrorUtils(),

			CommentService: services.NewCommentService(),
			PostService:    services.NewPostService(),
		}
	})
	return commentRouter
}
func (r *CommentRouter) Bind(_router *gin.RouterGroup) {
	router := _router.Group("/comment")
	// POST
	{
		// @Summary Create a new comment
		router.POST("",
			middlewares.VerifyAccessToken(middlewares.ParseJWTAccessToken),
			r.CreateComment,
		)
	}
	// GET
	{
		// @Summary Get comments by post ID
		router.GET("/list/post/:postID", r.GetCommentsByPostID)
	}
}

// @Tags Comment
// @Summary Get comments by post ID
// @Accept application/json
// @Produce application/json
// @Param postID path string true "Post ID"
// @Success 200 {array} models.CommentGetListByPostIDResponseItem
// @Failure 400 {object} models.ErrorResponse "Invalid post ID format"
// @Failure 404 {object} models.ErrorResponse "Post not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/comment/list/post/{postID} [get]
func (r *CommentRouter) GetCommentsByPostID(ctx *gin.Context) {
	// 解析路由參數
	postID, err := uuid.Parse(ctx.Param("postID"))
	if err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "Invalid post ID format"})
		return
	}

	// 檢查 Post 是否存在
	if _, err := r.PostService.GetByID(ctx, postID); err != nil {
		ctx.JSON(404, models.ErrorResponse{Error: "Post not found"})
		return
	}

	// 獲取評論
	comments, err := r.CommentService.GetListByPostID(ctx, postID)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{Error: "Failed to retrieve comments"})
		return
	}

	// 組織評論資料
	rootComments := make([]models.CommentGetListByPostIDResponseItem, 0)
	subComments := make([]models.CommentGetListByPostIDResponseItem, 0)
	parentCommentsMap := make(map[uuid.UUID][]models.CommentGetListByPostIDResponseItem)
	for _, comment := range comments {
		item := models.CommentGetListByPostIDResponseItem{
			ID:          comment.ID,
			PostID:      comment.PostID,
			Content:     comment.Content,
			ParentID:    comment.ParentID,
			UserID:      comment.UserID,
			UserName:    (*comment.User).Username,
			CreatedAt:   time.Unix(comment.CreatedAt, 0).Format(time.RFC3339),
			UpdatedAt:   time.Unix(comment.UpdatedAt, 0).Format(time.RFC3339),
			SubComments: []models.CommentGetListByPostIDResponseItem{},
		}
		if item.ParentID == nil {
			rootComments = append(rootComments, item)
			parentCommentsMap[item.ID] = item.SubComments
			continue
		} else {
			subComments = append(subComments, item)
		}
	}
	for _, subComment := range subComments {
		parentCommentsMap[*subComment.ParentID] = append(parentCommentsMap[*subComment.ParentID], subComment)
	}
	ctx.JSON(200, rootComments)
}

// @Tags Comment
// @Summary Create a new comment
// @Security AccessToken
// @Accept application/json
// @Produce application/json
// @Param comment body models.CommentCreateRequest true "Comment data"
// @Success 200 {object} models.CommentCreateResponse
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 404 {object} models.ErrorResponse "Post or parent comment not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/comment [post]
func (r *CommentRouter) CreateComment(ctx *gin.Context) {
	// 解析請求體
	commentCreateRequest := &models.CommentCreateRequest{}
	if err := ctx.ShouldBindJSON(commentCreateRequest); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// 檢查 Post 是否存在
	if _, err := r.PostService.GetByID(ctx, commentCreateRequest.PostID); err != nil {
		ctx.JSON(404, models.ErrorResponse{Error: "Post not found"})
		return
	}

	// 如果 ParentID 不為空，檢查父評論是否存在
	if commentCreateRequest.ParentID != nil {
		if _, err := r.CommentService.GetByID(ctx, *commentCreateRequest.ParentID); err != nil {
			ctx.JSON(404, models.ErrorResponse{Error: "Parent comment not found"})
			return
		}
	}

	// 獲取登入使用者資料
	tokenData, err := middlewares.GetContentAccessTokenData(ctx)
	if err != nil {
		err = r.ErrorUtils.ServerInternalError(err.Error())
		ctx.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	// 創建評論
	commentBase := models.CommentBase{
		PostID:   commentCreateRequest.PostID,
		UserID:   tokenData.UserID,
		Content:  commentCreateRequest.Content,
		ParentID: commentCreateRequest.ParentID,
	}
	comments, err := r.CommentService.Create(ctx, []models.CommentBase{commentBase})
	if err != nil {
		err = r.ErrorUtils.ServerInternalError(err.Error())
		ctx.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}
	if len(comments) == 0 {
		ctx.JSON(500, models.ErrorResponse{Error: "Failed to create comment"})
		return
	}
	comment := comments[0]

	// 構建回應資料
	respData := models.CommentCreateResponse{
		ID:       comment.ID,
		PostID:   comment.PostID,
		Content:  comment.Content,
		ParentID: comment.ParentID,
		UserID:   tokenData.UserID,
	}
	ctx.JSON(200, respData)
}
