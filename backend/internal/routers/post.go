package routers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/pkg"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type PostRouter struct {
	JWTUtils  *pkg.JWTUtils
	AuthUtils *pkg.AuthUtils
}

var postRouterOnce sync.Once
var postRouter *PostRouter

func NewPostRouter() *PostRouter {
	postRouterOnce.Do(func() {
		postRouter = &PostRouter{
			JWTUtils:  pkg.NewJWTUtils(),
			AuthUtils: pkg.NewAuthUtils(),
		}
	})
	return postRouter
}

func (r *PostRouter) Bind(_router *gin.RouterGroup) {
	router := _router.Group("/post")
	// POST
	{
		router.POST("",
			middlewares.VerifyAccessToken(func(authHeader string) (jwt.MapClaims, bool) {
				claims, err := r.JWTUtils.ParseToken(authHeader, nil)
				if err != nil {
					return nil, false
				}
				return claims, true
			}),
			r.CreatePost,
		)
	}

	// GET
	{
		// router.GET("/user/:userID", r.GetPostsByUser)
	}
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
	reqBody := &models.PostCreateRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "invalid request body"})
		return
	}

	if reqBody.Content == "" || len(reqBody.Content) > 300 {
		ctx.JSON(400, models.ErrorResponse{Error: "content characters must be between 0 and 300"})
		return
	}

	tagIDs := make([]uuid.UUID, len(reqBody.Tags))
	for i, _ := range reqBody.Tags {
		tagIDs[i] = uuid.New() // Assuming tags are created or fetched elsewhere
	}
	respBody := models.PostCreateResponse{
		ID:        uuid.New(),
		AuthorID:  reqBody.AuthorID,
		ImageURL:  reqBody.ImageURL,
		Content:   reqBody.Content,
		TagIDs:    tagIDs,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	ctx.JSON(200, respBody)
}

// func (r *PostRouter) GetPostsByUser(ctx *gin.Context) {
// 	userID := ctx.Param("userID")

// 	ctx.JSON(200, gin.H{"message": "Get posts by user", "userID": userID})
// }
