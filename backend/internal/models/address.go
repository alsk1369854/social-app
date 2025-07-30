package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	CityID uint   `json:"-"`
	City   City   `json:"City" gorm:"foreignKey:CityID"`
	Street string `json:"street"`
}
