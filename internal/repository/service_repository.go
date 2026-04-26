package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/ent"
	"github.com/matheusgosk8/book-me-server/ent/service"
	"github.com/matheusgosk8/book-me-server/internal/dto"
)

type ServiceRepository struct {
	client *ent.Client
}

func NewServiceRepository(client *ent.Client) *ServiceRepository {
	return &ServiceRepository{client: client}
}

// Recebe o DTO e o ID do provedor (vindo do Middleware)
func (r *ServiceRepository) CreateService(ctx context.Context, input dto.ServiceRequestDTO, providerID uuid.UUID) (*ent.Service, error) {
	return r.client.Service.
		Create().
		SetTitle(input.Title).
		SetDescription(input.Description).
		SetPriceBase(input.PriceBase).
		SetPriceType(service.PriceType(input.PriceType)). // Link com o Enum do Ent
		SetDurationMinutes(input.DurationMinutes).
		SetProviderID(providerID). // O ID que pegamos no Token
		SetCategoryID(input.CategoryID).
		Save(ctx)
}

func (r *ServiceRepository) ListServices(ctx context.Context) ([]*ent.Service, error) {
	return r.client.Service.
		Query().
		WithCategory().
		WithProvider().
		All(ctx)
}
