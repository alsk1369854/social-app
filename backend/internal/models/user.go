package models

import "github.com/google/uuid"

type User struct {
	TableModel
	UserBase
}

type UserBase struct {
	Username       string `gorm:"not null"`
	Email          string `gorm:"uniqueIndex;not null"`
	HashedPassword string `gorm:"not null"`
	Age            *int64
	AddressID      *uuid.UUID
	Address        *Address `gorm:"foreignKey:AddressID"`
}

// User Register structs
type UserRegisterRequest struct {
	Username string                      `json:"username" binding:"required"`
	Email    string                      `json:"email" binding:"required"`
	Password string                      `json:"password" binding:"required"`
	Age      *int64                      `json:"age"`
	Address  *UserRegisterRequestAddress `json:"address"`
}

type UserRegisterRequestAddress struct {
	CityID uuid.UUID `json:"cityID" binding:"required"`
	Street string    `json:"street" binding:"required"`
}

type UserRegisterResponse struct {
	ID       uuid.UUID                    `json:"id"`
	Username string                       `json:"username"`
	Email    string                       `json:"email"`
	Age      *int64                       `json:"age"`
	Address  *UserRegisterResponseAddress `json:"address"`
}

type UserRegisterResponseAddress struct {
	CityID uuid.UUID `json:"cityID"`
	Street string    `json:"street"`
}

// User Login structs
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	AccessToken string    `json:"accessToken"`
}
