package models

import "github.com/google/uuid"

type User struct {
	TableModel
	UserBase
}

type UserBase struct {
	Username     string
	Email        string
	PasswordHash string
	Age          *int64
	AddressID    *uuid.UUID
	Address      *Address `gorm:"foreignKey:AddressID"`
}

type UserRegisterRequest struct {
	Username string                      `json:"username" binding:"required"`
	Email    string                      `json:"email" binding:"required"`
	Password string                      `json:"password" binding:"required,min=6,max=12"`
	Age      *int64                      `json:"age"`
	Address  *UserRegisterRequestAddress `json:"address"`
}

type UserRegisterRequestAddress struct {
	CityID string `json:"cityID" binding:"required"`
	Street string `json:"street" binding:"required"`
}

type UserRegisterResponse struct {
	Username string                       `json:"username"`
	Email    string                       `json:"email"`
	Age      *int64                       `json:"age"`
	Address  *UserRegisterResponseAddress `json:"address"`
}

type UserRegisterResponseAddress struct {
	CityID string `json:"cityID"`
	Street string `json:"street"`
}
