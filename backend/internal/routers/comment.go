package routers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/services"
	"sync"

	"github.com/gin-gonic/gin"
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
		router.POST("",
			middlewares.VerifyAccessToken(middlewares.ParseJWTAccessToken),
			r.CreateComment,
		)
	}
	// // GET
	// {
	// 	router.GET("/post/:postID/offset/:offset/limit/:limit", r.GetCommentsByPostID)
	// }
}

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
