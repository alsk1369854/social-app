package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TagService struct {
	TagRepository *repositories.TagRepository
}

var tagService *TagService
var tagServiceOnce sync.Once

func NewTagService() *TagService {
	tagServiceOnce.Do(func() {
		tagService = &TagService{
			TagRepository: repositories.NewTagRepository(),
		}
	})
	return tagService
}

func (s *TagService) Create(ctx *gin.Context, tagBases []models.TagBase) ([]models.Tag, error) {
	return s.TagRepository.Create(ctx, tagBases)
}

func (s *TagService) CreateIfNotExist(ctx *gin.Context, tagBases []models.TagBase) ([]models.Tag, error) {
	// 檢查是否有已存在的 Tag
	result := make([]models.Tag, 0, len(tagBases))
	mustCreate := make([]models.TagBase, 0, len(tagBases))
	for _, tagBase := range tagBases {
		existingTag, err := s.TagRepository.GetByName(ctx, tagBase.Name)
		if err != nil {
			return nil, err
		}
		if existingTag != nil {
			result = append(result, *existingTag)
		} else {
			mustCreate = append(mustCreate, tagBase)
		}
	}
	if len(mustCreate) == 0 {
		return result, nil
	}

	// 創建不存在的 Tag
	tags, err := s.TagRepository.Create(ctx, mustCreate)
	if err != nil {
		return nil, err
	}
	result = append(result, tags...)

	return result, nil
}

func (s *TagService) GetByIDs(ctx *gin.Context, ids []uuid.UUID) ([]models.Tag, error) {
	return s.TagRepository.GetByIDs(ctx, ids)
}
