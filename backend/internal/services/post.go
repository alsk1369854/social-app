package services

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/repositories"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostService struct {
	ErrorUtils *pkg.ErrorUtils

	PostRepository *repositories.PostRepository

	TagService *TagService
}

var postServiceOnce sync.Once
var postService *PostService

func NewPostService() *PostService {
	postServiceOnce.Do(func() {
		postService = &PostService{
			ErrorUtils: pkg.NewErrorUtils(),

			PostRepository: repositories.NewPostRepository(),

			TagService: NewTagService(),
		}
	})
	return postService
}

func (s *PostService) GetByID(ctx *gin.Context, postID uuid.UUID) (*models.Post, error) {
	return s.PostRepository.GetByID(ctx, postID)
}

func (s *PostService) Create(ctx *gin.Context, postBases []models.PostBase, tags [][]models.Tag) ([]models.Post, error) {
	return s.PostRepository.Create(ctx, postBases, tags)
}

func (s *PostService) CreatePostWithTags(ctx *gin.Context, postBase models.PostBase, tagBases []models.TagBase) (*models.Post, error) {

	var post *models.Post
	if err := middlewares.TransactionGORMDB(ctx, func() error {
		tags, err := s.TagService.CreateIfNotExist(ctx, tagBases)
		if err != nil {
			return err
		}

		posts, err := s.PostRepository.Create(ctx, []models.PostBase{postBase}, [][]models.Tag{tags})
		if err != nil {
			return err
		}
		post = &posts[0]
		return nil
	}); err != nil {
		return nil, s.ErrorUtils.ServerInternalError(err.Error())
	}

	return post, nil
}

func (s *PostService) GetPostsByAuthorID(ctx *gin.Context, AuthorID uuid.UUID, pagination *models.Pagination) ([]models.Post, uint, error) {
	return s.PostRepository.GetPostsByAuthorID(ctx, AuthorID, pagination)
}

func (s *PostService) LikedByUser(ctx *gin.Context, postID uuid.UUID, userID uuid.UUID) error {
	return s.PostRepository.LikedByUser(ctx, postID, userID)
}

func (s *PostService) GetList(ctx *gin.Context, pagination *models.Pagination) ([]models.Post, uint, error) {
	return s.PostRepository.GetList(ctx, pagination)
}

func (s *PostService) GetListByKeywords(ctx *gin.Context, keywords []string, pagination *models.Pagination) ([]models.Post, uint, error) {
	return s.PostRepository.GetListByKeywords(ctx, keywords, pagination)
}
