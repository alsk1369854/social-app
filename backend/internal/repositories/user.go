package repositories

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct{}

var userOnce sync.Once
var userRepository *UserRepository

func NewUserRepository() *UserRepository {
	userOnce.Do(func() {
		userRepository = &UserRepository{}
	})
	return userRepository
}

func (r *UserRepository) GetByID(ctx *gin.Context, userID uuid.UUID) (*models.User, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)

	user := &models.User{}
	if err := db.Model(user).
		Where(&models.User{TableModel: models.TableModel{ID: userID}}).
		First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByUsername(ctx *gin.Context, username string) (*models.User, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)
	user := &models.User{}
	if err := db.Model(user).
		Where(&models.User{UserBase: models.UserBase{Username: username}}).
		First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx *gin.Context, email string) (*models.User, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)
	user := &models.User{}
	if err := db.Model(user).
		Where(&models.User{UserBase: models.UserBase{Email: email}}).
		First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx *gin.Context, userBaseSlice []models.UserBase) ([]models.User, error) {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)

	userSlice := make([]models.User, len(userBaseSlice))
	for i, userBase := range userBaseSlice {
		userSlice[i] = models.User{
			TableModel: models.TableModel{ID: uuid.New()},
			UserBase:   userBase,
		}
	}
	if err := db.Create(userSlice).Error; err != nil {
		return nil, err
	}
	return userSlice, nil
}

func (r *UserRepository) DeleteByID(ctx *gin.Context, userID uuid.UUID) error {
	db := ctx.MustGet(middlewares.CONTEXT_KEY_GORM_DB).(*gorm.DB)

	if err := db.Model(&models.User{}).
		Where(&models.User{TableModel: models.TableModel{ID: userID}}).
		Delete(&models.User{}).Error; err != nil {
		return err
	}

	return nil
}
