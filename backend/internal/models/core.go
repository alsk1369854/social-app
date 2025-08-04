package models

import "github.com/google/uuid"

type TableModel struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updatedAt"`
}

type JWTClaimsData struct {
	UserID uuid.UUID
}

type Pagination struct {
	Offset uint
	Limit  uint
}

type PaginationResponse[T any] struct {
	Data       []T         `json:"data"`
	TotalCount uint        `json:"totalCount"`
	Pagination *Pagination `json:"pagination"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
