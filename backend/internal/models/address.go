package models

import "github.com/google/uuid"

type Address struct {
	TableModel
	AddressBase
}

type AddressBase struct {
	CityID uuid.UUID `gorm:"type:uuid"`
	City   *City     `gorm:"foreignKey:CityID"`
	Street string
}
