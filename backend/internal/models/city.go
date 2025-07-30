package models

import "github.com/google/uuid"

type City struct {
	TableModel
	CityBase
}

type CityBase struct {
	Name string
}

type CityGetAllResponse []CityGetAllResponseItem

type CityGetAllResponseItem struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
