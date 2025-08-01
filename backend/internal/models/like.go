package models

import "github.com/google/uuid"

type Like struct {
	TableModel
	LikeBase
}

type LikeBase struct {
	PostID uuid.UUID `gorm:"type:uuid;not null"`
	Post   *Post     `gorm:"foreignKey:PostID"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	User   *User     `gorm:"foreignKey:UserID"`
}
