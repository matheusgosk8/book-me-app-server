package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ServiceDTO struct {
	Title string `json:"title" validate:"required"
    "`
	Description     string  `json:"description" validate:"required"`
	PriceBase       float64 `json:"price_base" validate:"required,gt=0"`
	PriceType       string  `json:"price_type" validate:"required"`
	DurationMinutes int     `json:"duration_minutes" validate:"required,gt=0"`
	CategoryID      string  `json:"category_id" validate:"required"`
	IsActive        bool    `json:"is_active"`
	IsInPlace       bool    `json:"is_in_place"`
	AddressID       *string `json:"address_id"`
}

func ValidateService(s ServiceDTO) map[string]string {
	v := validator.New()
	err := v.Struct(s)
	if err == nil {
		// cross-field validation: if IsInPlace then AddressID required
		if s.IsInPlace && (s.AddressID == nil || *s.AddressID == "") {
			return map[string]string{"address_id": "address_id é obrigatório quando is_in_place=true"}
		}
		return nil
	}

	verrs := err.(validator.ValidationErrors)
	msgs := map[string]string{}
	for _, e := range verrs {
		switch e.Tag() {
		case "required":
			msgs[e.Field()] = "campo obrigatório"
		case "gt":
			msgs[e.Field()] = fmt.Sprintf("deve ser maior que %s", e.Param())
		default:
			msgs[e.Field()] = e.Error()
		}
	}

	// Re-check cross-field to append message if needed
	if s.IsInPlace && (s.AddressID == nil || *s.AddressID == "") {
		msgs["address_id"] = "address_id é obrigatório quando is_in_place=true"
	}

	return msgs
}
