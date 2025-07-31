package models

type Comment struct {
	TableModel
	CommentBase
}

type CommentBase struct {
	PostID   string  `gorm:"type:uuid;not null"`
	UserID   string  `gorm:"type:uuid;not null"`
	Content  string  `gorm:"type:text;not null"`
	ParentID *string `gorm:"type:uuid"`
}
