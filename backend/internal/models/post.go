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
}
