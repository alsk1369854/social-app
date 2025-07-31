package repositories

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct{}

var postRepositoryOnce sync.Once
var postRepository *PostRepository

func NewPostRepository() *PostRepository {
	postRepositoryOnce.Do(func() {
		postRepository = &PostRepository{}
	})
	return postRepository
}

func (r *PostRepository) GetPostByID(ctx *gin.Context, postID uuid.UUID) (*models.Post, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)
	post := &models.Post{}
	if err := db.Model(post).
		Where(&models.Post{TableModel: models.TableModel{ID: postID}}).
		First(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}
