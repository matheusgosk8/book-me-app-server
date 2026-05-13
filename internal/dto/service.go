package dto

import "github.com/google/uuid"

// ServiceResponseDTO: Objeto de saída para o Front-end
type ServiceResponseDTO struct {
	ID              uuid.UUID              `json:"id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	PriceBase       float64                `json:"price_base"`
	PriceType       string                 `json:"price_type"`
	DurationMinutes int                    `json:"duration_minutes"`
	IsActive        bool                   `json:"is_active"`
	IsInPlace       bool                   `json:"is_in_place"`
	Provider        ProviderShortResponse  `json:"provider"`
	Category        CategoryShortResponse  `json:"category"`
}

// ServiceRequestDTO: Objeto de entrada para criação/edição
type ServiceRequestDTO struct {
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	PriceBase       float64    `json:"price_base"`
	PriceType       string     `json:"price_type"`
	DurationMinutes int        `json:"duration_minutes"`
	CategoryID      uuid.UUID  `json:"category_id"`
	IsActive        bool       `json:"is_active"`
	IsInPlace       bool       `json:"is_in_place"` 
	AddressID       *uuid.UUID `json:"address_id"`   
}

// Estruturas auxiliares 
type ProviderShortResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type CategoryShortResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}