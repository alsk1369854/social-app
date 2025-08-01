package models

type Tag struct {
	TableModel
	TagBase
}

type TagBase struct {
	Name  string  `gorm:"not null;unique"`
	Posts []*Post `gorm:"many2many:post_to_tag;"`
}
