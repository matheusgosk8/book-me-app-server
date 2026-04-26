package dto

import (
	"github.com/google/uuid"
)

// ServiceRequestDTO: O que o App envia (Input)
type ServiceRequestDTO struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description"`
	PriceBase       float64   `json:"price_base" binding:"required"`
	PriceType       string    `json:"price_type" binding:"required,oneof=fixed hourly"`
	DurationMinutes int       `json:"duration_minutes" binding:"required"`
	CategoryID      uuid.UUID `json:"category_id" binding:"required"`
}

// Sub-DTOs para não expor dados sensíveis do User
type ProviderShortResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type CategoryShortResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// ServiceResponseDTO: O que o servidor devolve (Output limpo)
type ServiceResponseDTO struct {
	ID              uuid.UUID             `json:"id"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	PriceBase       float64               `json:"price_base"`
	PriceType       string                `json:"price_type"`
	DurationMinutes int                   `json:"duration_minutes"`
	IsActive        bool                  `json:"is_active"`
	Provider        ProviderShortResponse `json:"provider"`
	Category        CategoryShortResponse `json:"category"`
}