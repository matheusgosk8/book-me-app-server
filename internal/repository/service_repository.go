package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/ent"
	"github.com/matheusgosk8/book-me-server/ent/category"
	"github.com/matheusgosk8/book-me-server/ent/service"
	"github.com/matheusgosk8/book-me-server/ent/user"
	"github.com/matheusgosk8/book-me-server/internal/dto"
)

type ServiceRepository struct {
	client *ent.Client
}

func NewServiceRepository(client *ent.Client) *ServiceRepository {
	return &ServiceRepository{client: client}
}

// CreateService: Cria um novo serviço vinculado a um prestador
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
		SetNillableAddressID(input.AddressID).
		Save(ctx)
}

// ListServices: Lista serviços com filtros de categoria e paginação
func (r *ServiceRepository) ListServices(ctx context.Context, limit, offset int, categoryID *uuid.UUID) ([]*ent.Service, error) {
	query := r.client.Service.Query().WithCategory().WithProvider()

	if categoryID != nil {
		query.Where(service.HasCategoryWith(category.ID(*categoryID)))
	}

	return query.
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(service.FieldCreatedAt)).
		All(ctx)
}

// UpdateService: Atualiza um serviço garantindo que o executor seja o dono (Provider)
func (r *ServiceRepository) UpdateService(ctx context.Context, id uuid.UUID, providerID uuid.UUID, input dto.ServiceRequestDTO) (*ent.Service, error) {
	// Verificamos a existência e propriedade antes de atualizar
	exists, err := r.client.Service.Query().
		Where(
			service.ID(id),
			service.HasProviderWith(user.ID(providerID)),
		).Exist(ctx)

	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("serviço não encontrado ou você não tem permissão para editá-lo")
	}

	return r.client.Service.
		UpdateOneID(id).
		SetTitle(input.Title).
		SetDescription(input.Description).
		SetPriceBase(input.PriceBase).
		SetPriceType(service.PriceType(input.PriceType)).
		SetDurationMinutes(input.DurationMinutes).
		SetIsActive(input.IsActive).
		SetCategoryID(input.CategoryID).
		SetIsInPlace(input.IsInPlace).
		SetNillableAddressID(input.AddressID).
		Save(ctx)
}

// DeleteService: Remove o serviço validando o proprietário via ID do Provider
func (r *ServiceRepository) DeleteService(ctx context.Context, id uuid.UUID, providerID uuid.UUID) error {
	// Delete().Where() retorna (int, error). Descartamos o count para retornar apenas o error.
	_, err := r.client.Service.
		Delete().
		Where(
			service.ID(id),
			service.HasProviderWith(user.ID(providerID)),
		).
		Exec(ctx)

	return err
}