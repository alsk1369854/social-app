package models

import "github.com/google/uuid"

type PostToTag struct {
	TableModel
	PostToTagBase
}

type PostToTagBase struct {
	PostID uuid.UUID `gorm:"not null"`
	TagID  uuid.UUID `gorm:"not null"`
}
