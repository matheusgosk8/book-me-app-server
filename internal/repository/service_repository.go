package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/ent"
	"github.com/matheusgosk8/book-me-server/ent/category"
	"github.com/matheusgosk8/book-me-server/ent/service"
	"github.com/matheusgosk8/book-me-server/internal/dto"
)

type ServiceRepository struct {
	client *ent.Client
}

// NewServiceRepository cria uma nova instância do repositório de serviços
func NewServiceRepository(client *ent.Client) *ServiceRepository {
	return &ServiceRepository{client: client}
}

// CreateService persiste um novo serviço com suporte a localização e endereço nulo
func (r *ServiceRepository) CreateService(ctx context.Context, input dto.ServiceRequestDTO, providerID uuid.UUID) (*ent.Service, error) {
	return r.client.Service.
		Create().
		SetTitle(input.Title).
		SetDescription(input.Description).
		SetPriceBase(input.PriceBase).
		SetPriceType(service.PriceType(input.PriceType)).
		SetDurationMinutes(input.DurationMinutes).
		SetIsActive(input.IsActive).
		SetCategoryID(input.CategoryID).
		SetProviderID(providerID).
		SetIsInPlace(input.IsInPlace).
		SetNillableAddressID(input.AddressID). // Define o ID do endereço apenas se não for nil
		Save(ctx)
}

// ListServices busca serviços com paginação, filtros e pré-carregamento de relações (Edges)
func (r *ServiceRepository) ListServices(ctx context.Context, limit, offset int, categoryID *uuid.UUID) ([]*ent.Service, error) {
	query := r.client.Service.
		Query().
		WithCategory().
		WithProvider()

	// Filtro por Categoria caso o ID seja fornecido
	if categoryID != nil {
		query.Where(service.HasCategoryWith(category.ID(*categoryID)))
	}

	return query.
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(service.FieldCreatedAt)). // Ordenação por data de criação (mais novos primeiro)
		All(ctx)
}

// DeleteService remove um serviço do banco de dados pelo seu UUID
func (r *ServiceRepository) DeleteService(ctx context.Context, id uuid.UUID) error {
	return r.client.Service.
		DeleteOneID(id).
		Exec(ctx)
}