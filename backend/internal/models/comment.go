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
	PostID   uuid.UUID  `json:"postID" binding:"required"`
	Content  string     `json:"content" binding:"required"`
	ParentID *uuid.UUID `json:"parentID"`
}

type CommentCreateResponse struct {
	ID       uuid.UUID  `json:"id"`
	PostID   uuid.UUID  `json:"postID"`
	Content  string     `json:"content"`
	ParentID *uuid.UUID `json:"parentID"`
	UserID   uuid.UUID  `json:"userID"`
}

// Get Comments By PostID
type CommentGetListByPostIDResponseItem struct {
	ID          uuid.UUID                            `json:"id"`
	PostID      uuid.UUID                            `json:"postID"`
	Content     string                               `json:"content"`
	ParentID    *uuid.UUID                           `json:"parentID"`
	UserID      uuid.UUID                            `json:"userID"`
	UserName    string                               `json:"userName"`
	CreatedAt   string                               `json:"createdAt"`
	UpdatedAt   string                               `json:"updatedAt"`
	SubComments []CommentGetListByPostIDResponseItem `json:"subComments"`
}
