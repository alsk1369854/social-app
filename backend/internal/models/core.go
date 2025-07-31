package models

import "github.com/google/uuid"

type TableModel struct {
	ID uuid.UUID `gorm:"primaryKey"`
}

type JWTClaimsData struct {
	UserID uuid.UUID
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
