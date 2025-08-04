package repositories

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository struct{}

var tagRepositoryOnce sync.Once
var tagRepository *TagRepository

func NewTagRepository() *TagRepository {
	tagRepositoryOnce.Do(func() {
		tagRepository = &TagRepository{}
	})
	return tagRepository
}

func (r *TagRepository) Create(ctx *gin.Context, tagBases []models.TagBase) ([]models.Tag, error) {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]models.Tag, len(tagBases))
	for i, tagBase := range tagBases {
		tags[i] = models.Tag{
			TableModel: models.TableModel{
				ID: uuid.New(),
			},
			TagBase: tagBase,
		}
	}
	if err := db.Create(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *TagRepository) GetByName(ctx *gin.Context, name string) (*models.Tag, error) {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, err
	}

	tag := &models.Tag{}
	if err := db.Where("name = ?", name).First(tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return tag, nil
}

func (r *TagRepository) GetByIDs(ctx *gin.Context, ids []uuid.UUID) ([]models.Tag, error) {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, err
	}

	tags := []models.Tag{}
	if err := db.Where("id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}
