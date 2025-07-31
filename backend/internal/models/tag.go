package models

type Tag struct {
	TableModel
	TagBase
}

type TagBase struct {
	Name string `gorm:"not null;unique"`
}
