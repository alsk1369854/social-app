package models

import "github.com/google/uuid"

type Comment struct {
	TableModel
	CommentBase
}

type CommentBase struct {
	PostID   uuid.UUID  `gorm:"type:uuid;not null"`
	Post     *Post      `gorm:"foreignKey:PostID"`
	UserID   uuid.UUID  `gorm:"type:uuid;not null"`
	User     *User      `gorm:"foreignKey:UserID"`
	Content  string     `gorm:"type:text;not null"`
	ParentID *uuid.UUID `gorm:"type:uuid"`
}

// Create Comment structs
type CommentCreateRequest struct {
	PostID   uuid.UUID  `json:"post_id" binding:"required"`
	Content  string     `json:"content" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type CommentCreateResponse struct {
	ID       uuid.UUID  `json:"id"`
	PostID   uuid.UUID  `json:"post_id"`
	Content  string     `json:"content"`
	ParentID *uuid.UUID `json:"parent_id"`
	UserID   uuid.UUID  `json:"user_id"`
}
