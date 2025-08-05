package repositories

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentRepository struct {
}

var commentRepositoryOnce sync.Once
var commentRepository *CommentRepository

func NewCommentRepository() *CommentRepository {
	commentRepositoryOnce.Do(func() {
		commentRepository = &CommentRepository{}
	})
	return commentRepository
}

func (r *CommentRepository) Create(ctx *gin.Context, commentBase []models.CommentBase) ([]models.Comment, error) {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, err
	}

	comment := make([]models.Comment, len(commentBase))
	for i, base := range commentBase {
		comment[i] = models.Comment{
			TableModel: models.TableModel{
				ID: uuid.New(),
			},
			CommentBase: base,
		}
	}
	if err := db.Create(comment).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepository) GetByID(ctx *gin.Context, commentID uuid.UUID) (*models.Comment, error) {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, err
	}

	var comment models.Comment
	if err := db.First(&comment, "id = ?", commentID).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) GetListByPostID(ctx *gin.Context, postID uuid.UUID) ([]models.Comment, error) {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, err
	}

	comments := []models.Comment{}
	if err := db.Where("post_id = ?", postID).
		Order("created_at ASC").
		Preload("User").
		Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}
