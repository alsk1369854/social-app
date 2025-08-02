package repositories

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		if err := db.
			Where(models.Tag{TagBase: models.TagBase{Name: tagBase.Name}}).
			FirstOrCreate(&tags[i]).
			Error; err != nil {
			return nil, err
		}
	}

	return tags, nil
}
