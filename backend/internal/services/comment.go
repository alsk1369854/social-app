package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentService struct {
	CommentRepository *repositories.CommentRepository
}

var commentServiceOnce sync.Once
var commentService *CommentService

func NewCommentService() *CommentService {
	commentServiceOnce.Do(func() {
		commentService = &CommentService{
			CommentRepository: repositories.NewCommentRepository(),
		}
	})
	return commentService
}

func (s *CommentService) Create(ctx *gin.Context, commentBases []models.CommentBase) ([]models.Comment, error) {
	return s.CommentRepository.Create(ctx, commentBases)
}

func (s *CommentService) GetByID(ctx *gin.Context, commentID uuid.UUID) (*models.Comment, error) {
	return s.CommentRepository.GetByID(ctx, commentID)
}

func (s *CommentService) GetListByPostID(ctx *gin.Context, postID uuid.UUID) ([]models.Comment, error) {
	return s.CommentRepository.GetListByPostID(ctx, postID)
}
