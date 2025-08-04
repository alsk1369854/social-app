package models

import "github.com/google/uuid"

type Post struct {
	TableModel
	PostBase
}

type PostBase struct {
	AuthorID uuid.UUID `gorm:"not null"`
	Author   *User     `gorm:"foreignKey:AuthorID"`
	ImageURL *string
	Content  string  `gorm:"not null"`
	Tags     []*Tag  `gorm:"many2many:post_to_tag;"`
	Likes    []*User `gorm:"many2many:post_to_user;"`
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

// Post GetPostsByAuthorID structs
type PostGetPostsByAuthorIDResponseItem struct {
	ID         uuid.UUID                               `json:"id"`
	AuthorID   uuid.UUID                               `json:"authorID"`
	ImageURL   *string                                 `json:"imageURL"`
	Content    string                                  `json:"content"`
	CreatedAt  string                                  `json:"createdAt"`
	UpdatedAt  string                                  `json:"updatedAt"`
	Tags       []PostGetPostsByAuthorIDResponseItemTag `json:"tags"`
	LikedCount uint                                    `json:"likedCount"`
}

type PostGetPostsByAuthorIDResponseItemTag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
