package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostService struct {
	PostRepository *repositories.PostRepository
}

var postServiceOnce sync.Once
var postService *PostService

func NewPostService() *PostService {
	postServiceOnce.Do(func() {
		postService = &PostService{
			PostRepository: repositories.NewPostRepository(),
		}
	})
	return postService
}

func (s *PostService) GetPostByID(ctx *gin.Context, postID uuid.UUID) (*models.Post, error) {
	return s.PostRepository.GetPostByID(ctx, postID)
}

func (s *PostService) Create(ctx *gin.Context, postBases []models.PostBase) ([]models.Post, error) {
	return s.PostRepository.Create(ctx, postBases)
}
