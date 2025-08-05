package routers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/services"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostRouter struct {
	JWTUtils  *pkg.JWTUtils
	AuthUtils *pkg.AuthUtils

	PostService *services.PostService
	TagService  *services.TagService
	UserService *services.UserService
}

var postRouterOnce sync.Once
var postRouter *PostRouter

func NewPostRouter() *PostRouter {
	postRouterOnce.Do(func() {
		postRouter = &PostRouter{
			JWTUtils:  pkg.NewJWTUtils(),
			AuthUtils: pkg.NewAuthUtils(),

			PostService: services.NewPostService(),
			TagService:  services.NewTagService(),
			UserService: services.NewUserService(),
		}
	})
	return postRouter
}

func (r *PostRouter) Bind(_router *gin.RouterGroup) {
	router := _router.Group("/post")
	// POST
	{
		router.POST("",
			middlewares.VerifyAccessToken(middlewares.ParseJWTAccessToken),
			r.CreatePost,
		)
	}
	// GET
	{
		router.GET("/author/:authorID/offset/:offset/limit/:limit", r.GetPostsByAuthorID)
	}
	//PUT
	{
		// router.PUT("/like/:postID",
		// 	middlewares.VerifyAccessToken(middlewares.ParseJWTAccessToken),
		// 	r.LikedPostByUser,
		// )
	}
}

// // @title Post API
// // @Summary Like a post by user
// // @Tags Post
// // @Security AccessToken
// // @Accept text/plain
// // @Produce application/json
// // @Param postID path string true "Post ID"
// // @Success 200 {object} models.SuccessResponse
// // @Failure 400 {object} models.ErrorResponse
// // @Failure 404 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Router /api/post/like/{postID} [put]
// func (r *PostRouter) LikedPostByUser(ctx *gin.Context) {
// 	postID, err := uuid.Parse(ctx.Param("postID"))
// 	if err != nil {
// 		ctx.JSON(400, models.ErrorResponse{Error: "invalid post ID"})
// 		return
// 	}
// 	// 檢查貼文是否存在
// 	if _, err := r.PostService.GetByID(ctx, postID); err != nil {
// 		ctx.JSON(404, models.ErrorResponse{Error: "post not found"})
// 		return
// 	}

// 	// 從 Token 中取得用戶 ID
// 	claims, err := middlewares.GetContentAccessTokenData(ctx)
// 	if err != nil {
// 		ctx.JSON(500, models.ErrorResponse{Error: "failed to get user ID from token"})
// 		return
// 	}

// 	// 添加 用戶對貼文的喜歡
// 	if err := r.PostService.LikedByUser(ctx, postID, claims.UserID); err != nil {
// 		ctx.JSON(500, models.ErrorResponse{Error: err.Error()})
// 		return
// 	}

// 	ctx.JSON(200, models.SuccessResponse{Success: true})
// }

// @title Post API
// @Summary Get posts by author ID
// @Tags Post
// @Accept text/plain
// @Produce application/json
// @Param authorID path string true "Author ID"
// @Param offset path string true "Offset"
// @Param limit path string true "Limit"
// @Success 200 {object} models.PaginationResponse[models.PostGetPostsByAuthorIDResponseItem]
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/post/author/{authorID}/offset/{offset}/limit/{limit} [get]
func (r *PostRouter) GetPostsByAuthorID(ctx *gin.Context) {
	// 解析路由參數
	authorID, err := uuid.Parse(ctx.Param("authorID"))
	if err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "invalid author ID"})
		return
	}
	offset, err := strconv.ParseUint(ctx.Param("offset"), 10, 64)
	if err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "invalid offset"})
		return
	}
	limit, err := strconv.ParseUint(ctx.Param("limit"), 10, 64)
	if err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "invalid limit"})
		return
	}

	// 檢查用戶是否存在
	user, err := r.UserService.GetByID(ctx, authorID)
	if err != nil {
		ctx.JSON(404, models.ErrorResponse{Error: "author not found"})
		return
	}
	// 獲取用戶的貼文
	pagination := &models.Pagination{
		Offset: uint(offset),
		Limit:  uint(limit),
	}
	posts, totalCount, err := r.PostService.GetPostsByAuthorID(ctx, user.ID, pagination)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	// 構建回應
	responseData := make([]models.PostGetPostsByAuthorIDResponseItem, len(posts))
	for i, post := range posts {
		tags := make([]models.PostGetPostsByAuthorIDResponseItemTag, len(post.Tags))
		for j, tag := range post.Tags {
			tags[j] = models.PostGetPostsByAuthorIDResponseItemTag{Name: tag.Name}
		}
		responseData[i] = models.PostGetPostsByAuthorIDResponseItem{
			ID:         post.ID,
			AuthorID:   post.AuthorID,
			ImageURL:   post.ImageURL,
			Content:    post.Content,
			CreatedAt:  time.Unix(post.CreatedAt, 0).Format(time.RFC3339),
			UpdatedAt:  time.Unix(post.UpdatedAt, 0).Format(time.RFC3339),
			Tags:       tags,
			LikedCount: uint(len(post.Likes)),
		}
	}
	ctx.JSON(200, models.PaginationResponse[models.PostGetPostsByAuthorIDResponseItem]{
		Data:       responseData,
		TotalCount: totalCount,
		Pagination: pagination,
	})
}

// @title Post API
// @Summary Create a post
// @Tags Post
// @Security AccessToken
// @Accept application/json
// @Produce application/json
// @Param post body models.PostCreateRequest true "Post create request"
// @Success 200 {object} models.PostCreateResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/post [post]
func (r *PostRouter) CreatePost(ctx *gin.Context) {
	// 解析請求體
	reqBody := &models.PostCreateRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "invalid request body"})
		return
	}
	if reqBody.Content == "" || len(reqBody.Content) > 300 {
		ctx.JSON(400, models.ErrorResponse{Error: "content characters must be between 0 and 300"})
		return
	}

	// 獲取 Token Data
	tokenData, err := middlewares.GetContentAccessTokenData(ctx)
	if err != nil {
		err = r.PostService.ErrorUtils.ServerInternalError(err.Error())
		ctx.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	// 創建 Post
	tagBases := make([]models.TagBase, len(reqBody.Tags))
	for i, tagName := range reqBody.Tags {
		tagBases[i] = models.TagBase{Name: strings.Trim(tagName, " ")}
	}
	postBase := models.PostBase{
		AuthorID: tokenData.UserID,
		ImageURL: reqBody.ImageURL,
		Content:  reqBody.Content,
	}
	post, err := r.PostService.CreatePostWithTags(ctx, postBase, tagBases)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	// 構建回應
	tagIDs := make([]uuid.UUID, len(post.Tags))
	for i, tag := range post.Tags {
		tagIDs[i] = tag.ID
	}
	respBody := models.PostCreateResponse{
		ID:        post.ID,
		AuthorID:  post.AuthorID,
		ImageURL:  post.ImageURL,
		Content:   post.Content,
		TagIDs:    tagIDs,
		CreatedAt: time.Unix(post.CreatedAt, 0).Format(time.RFC3339),
		UpdatedAt: time.Unix(post.UpdatedAt, 0).Format(time.RFC3339),
	}
	ctx.JSON(200, respBody)
}
