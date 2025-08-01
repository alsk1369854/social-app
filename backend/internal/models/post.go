package models

import "github.com/google/uuid"

type Post struct {
	TableModel
	PostBase
}

type PostBase struct {
	AuthorID uuid.UUID `gorm:"not null"`
	Author   *User     `gorm:"foreignKey:AuthorID"`
	ImageURL *string   `gorm:"default:null"`
	Content  string    `gorm:"not null"`
	Tags     []*Tag    `gorm:"many2many:post_to_tag;"`
}

// Post Create structs
type PostCreateRequest struct {
	AuthorID uuid.UUID `json:"authorID" binding:"required"`
	ImageURL *string   `json:"imageURL"`
	Content  string    `json:"content" binding:"required"`
	Tags     []string  `json:"tags" binding:"required"`
}

type PostCreateResponse struct {
	ID        uuid.UUID   `json:"id"`
	AuthorID  uuid.UUID   `json:"authorID"`
	ImageURL  *string     `json:"imageURL"`
	Content   string      `json:"content"`
	TagIDs    []uuid.UUID `json:"tagIDs"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
}

type PostCreateResponseTag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// type PostGetPostsByUserIDRequest struct{}
