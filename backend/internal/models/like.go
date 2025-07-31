package models

type Like struct {
	TableModel
	LikeBase
}

type LikeBase struct {
	PostID string `gorm:"type:uuid;not null"`
	UserID string `gorm:"type:uuid;not null"`
}
