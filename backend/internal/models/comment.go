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
