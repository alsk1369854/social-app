package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string
	Email        string
	PasswordHash string
	Age          *int64
	AddressID    *uint
	Address      *Address `gorm:"foreignKey:AddressID"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=12"`
	Age      *int64 `json:"age"`
	Address  *struct {
		CityID uint   `json:"cityID" binding:"required"`
		Street string `json:"street" binding:"required"`
	} `json:"address"`
}
