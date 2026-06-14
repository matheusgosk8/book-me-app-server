package validator

import (
	"fmt"

	"github.com/google/uuid"
)

// DTO para patch -- campos opcionais
type ServicePatchDTO struct {
	Title           *string  `json:"title,omitempty"`
	Description     *string  `json:"description,omitempty"`
	PriceBase       *float64 `json:"price_base,omitempty"`
	PriceType       *string  `json:"price_type,omitempty"`
	DurationMinutes *int     `json:"duration_minutes,omitempty"`
	CategoryID      *string  `json:"category_id,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
	IsInPlace       *bool    `json:"is_in_place,omitempty"`
	AddressID       *string  `json:"address_id,omitempty"`
}

// ValidateServicePatch valida apenas campos presentes
func ValidateServicePatch(s ServicePatchDTO) map[string]string {
	msgs := map[string]string{}

	if s.PriceBase != nil {
		if *s.PriceBase <= 0 {
			msgs["price_base"] = "deve ser maior que 0"
		}
	}
	if s.DurationMinutes != nil {
		if *s.DurationMinutes <= 0 {
			msgs["duration_minutes"] = "deve ser maior que 0"
		}
	}
	if s.CategoryID != nil {
		if _, err := uuid.Parse(*s.CategoryID); err != nil {
			msgs["category_id"] = "category_id inválido"
		}
	}
	if s.AddressID != nil {
		if *s.AddressID != "" {
			if _, err := uuid.Parse(*s.AddressID); err != nil {
				msgs["address_id"] = "address_id inválido"
			}
		}
	}
	if s.PriceType != nil {
		pt := *s.PriceType
		if pt != "fixed" && pt != "hourly" {
			msgs["price_type"] = fmt.Sprintf("price_type inválido: %s", pt)
		}
	}

	// cross-field: se is_in_place estiver sendo setado para true, exigir address_id estar presente (no patch)
	if s.IsInPlace != nil && *s.IsInPlace {
		if s.AddressID == nil || *s.AddressID == "" {
			msgs["address_id"] = "address_id é obrigatório quando is_in_place=true"
		}
	}

	if len(msgs) == 0 {
		return nil
	}
	return msgs
}
