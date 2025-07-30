package models

import "gorm.io/gorm"

type City struct {
	gorm.Model
	Name string `json:"name"`
}

type CityGetAllResponse []struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
