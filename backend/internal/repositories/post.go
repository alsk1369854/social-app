package repositories

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/pkg"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type PostRepository struct {
	ErrorUtils *pkg.ErrorUtils

	UserRepository *UserRepository
}

var postRepositoryOnce sync.Once
var postRepository *PostRepository

func NewPostRepository() *PostRepository {
	postRepositoryOnce.Do(func() {
		postRepository = &PostRepository{
			ErrorUtils: pkg.NewErrorUtils(),

			UserRepository: NewUserRepository(),
		}
	})
	return postRepository
}

func (r *PostRepository) GetByID(ctx *gin.Context, postID uuid.UUID) (*models.Post, error) {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, err
	}

	post := &models.Post{}
	if err := db.Model(post).
		Preload("Author").
		Preload("Tags").
		Preload("Likes").
		Where("id = ?", postID).
		First(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepository) GetListByKeywords(ctx *gin.Context, keywords []string, pagination *models.Pagination) ([]models.Post, uint, error) {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, 0, err
	}

	db = db.Model(&models.Post{}).
		Joins("LEFT JOIN users AS author ON author.id = posts.author_id").
		Joins("LEFT JOIN post_to_tag ON post_to_tag.post_id = posts.id").
		Joins("LEFT JOIN tags ON tags.id = post_to_tag.tag_id").
		Preload("Author").
		Preload("Tags").
		Preload("Likes").
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Table: "posts", Name: "created_at"}, Desc: true},
		}}).
		Group("posts.id")

	if len(keywords) > 0 {
		first := "%" + keywords[0] + "%"
		db = db.Where(
			"author.username ILIKE ? OR content ILIKE ? OR tags.name ILIKE ?",
			first, first, first,
		)
		for _, keyword := range keywords[1:] {
			k := "%" + keyword + "%"
			db = db.Or(
				"author.username ILIKE ? OR content ILIKE ? OR tags.name ILIKE ?",
				k, k, k,
			)
		}
	}

	totalCount := int64(0)
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if pagination != nil {
		if pagination.Limit <= 0 {
			return nil, 0, r.ErrorUtils.ServerInternalError("invalid pagination parameters")
		}
		db = db.Offset(int(pagination.Offset)).Limit(int(pagination.Limit))
	}

	posts := []models.Post{}
	if err := db.Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, uint(totalCount), nil
}

func (r *PostRepository) GetList(ctx *gin.Context, pagination *models.Pagination) ([]models.Post, uint, error) {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, 0, err
	}

	db = db.Model(&models.Post{}).
		Preload("Author").
		Preload("Tags").
		Preload("Likes").
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "created_at"}, Desc: true},
		}})

	totalCount := int64(0)
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if pagination != nil {
		if pagination.Limit <= 0 {
			return nil, 0, r.ErrorUtils.ServerInternalError("invalid pagination parameters")
		}
		db = db.Offset(int(pagination.Offset)).Limit(int(pagination.Limit))
	}

	posts := []models.Post{}
	if err := db.Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, uint(totalCount), nil
}

func (r *PostRepository) Create(ctx *gin.Context, postBases []models.PostBase, tags [][]models.Tag) ([]models.Post, error) {
	if len(postBases) != len(tags) {
		return nil, r.ErrorUtils.ServerInternalError("postBases and tags length mismatch")
	}

	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, err
	}

	// Create posts
	posts := make([]models.Post, len(postBases))
	for i, postBase := range postBases {
		posts[i] = models.Post{
			TableModel: models.TableModel{ID: uuid.New()},
			PostBase:   postBase,
		}
	}
	if err := db.Create(&posts).Error; err != nil {
		return nil, err
	}

	// Associate tags with posts
	for i := range posts {
		if err := db.Model(&posts[i]).Association("Tags").Append(tags[i]); err != nil {
			return nil, err
		}
	}

	if err := db.Model(&models.Post{}).Preload("Tags").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) LikedByUser(ctx *gin.Context, postID uuid.UUID, userID uuid.UUID) error {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return err
	}

	post, err := r.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	user, err := r.UserRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := db.Model(post).Association("Likes").Append(user); err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) GetPostsByAuthorID(ctx *gin.Context, AuthorID uuid.UUID, pagination *models.Pagination) ([]models.Post, uint, error) {
	db, err := middlewares.GetContentGORMDB(ctx)
	if err != nil {
		return nil, 0, err
	}

	db = db.Model(&models.Post{}).
		Where(&models.Post{PostBase: models.PostBase{AuthorID: AuthorID}}).
		Preload("Author").
		Preload("Tags").
		Preload("Likes").
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: "created_at"},
			Desc:   true,
		})

	totalCount := int64(0)
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if pagination != nil {
		if pagination.Limit <= 0 {
			return nil, 0, r.ErrorUtils.ServerInternalError("invalid pagination parameters")
		}
		db = db.Offset(int(pagination.Offset)).Limit(int(pagination.Limit))
	}

	posts := []models.Post{}
	if err := db.Find(&posts).Error; err != nil {
		return nil, 0, r.ErrorUtils.ServerInternalError(err.Error())
	}
	return posts, uint(totalCount), nil
}
