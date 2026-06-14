package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/ent"
	"github.com/matheusgosk8/book-me-server/ent/category"
	"github.com/matheusgosk8/book-me-server/ent/service"
	"github.com/matheusgosk8/book-me-server/ent/user"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/dto"
	models "github.com/matheusgosk8/book-me-server/internal/models"
)

type ListServices struct{}

func NewListServices() *ListServices {
	return &ListServices{}
}

type ListServicesOutput struct {
	Data  []*ent.Service
	Total int
}

func (s *ListServices) ListMyServices(ctx context.Context, providerId uuid.UUID, pagination *models.Pagination) (*ListServicesOutput, error) {
	q := db.Client.Service.Query().Where(service.HasProviderWith(user.ID(providerId)))

	total, err := q.Count(ctx)
	if err != nil {
		return nil, err
	}

	data, err := q.
		Limit(pagination.PerPage).
		Offset((pagination.Page - 1) * pagination.PerPage).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &ListServicesOutput{Data: data, Total: total}, nil
}

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
func (r *ServiceRepository) ListServices(ctx context.Context, limit, offset int, categoryID *uuid.UUID) (*ListServicesOutput, error) {
	query := r.client.Service.Query().WithCategory().WithProvider()

	if categoryID != nil {
		query.Where(service.HasCategoryWith(category.ID(*categoryID)))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	data, err := query.
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(service.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &ListServicesOutput{Data: data, Total: total}, nil
}

// UpdateService: Atualiza um serviço garantindo que o executor seja o dono (Provider)
func (r *ServiceRepository) UpdateService(ctx context.Context, id uuid.UUID, providerID uuid.UUID, input dto.ServicePatchDTO) (*ent.Service, error) {
	// Verifica a existência e propriedade antes de atualizar
	exists, err := r.client.Service.Query().
		Where(
			service.ID(id),
			service.HasProviderWith(user.ID(providerID)),
		).Exist(ctx)

	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("serviço não encontrado ou acesso negado")
	}

	// Executa o Update aplicando apenas os campos presentes no patch
	upd := r.client.Service.UpdateOneID(id)

	if input.Title != nil {
		upd = upd.SetTitle(*input.Title)
	}
	if input.Description != nil {
		upd = upd.SetDescription(*input.Description)
	}
	if input.PriceBase != nil {
		upd = upd.SetPriceBase(*input.PriceBase)
	}
	if input.PriceType != nil {
		upd = upd.SetPriceType(service.PriceType(*input.PriceType))
	}
	if input.DurationMinutes != nil {
		upd = upd.SetDurationMinutes(*input.DurationMinutes)
	}
	if input.IsActive != nil {
		upd = upd.SetIsActive(*input.IsActive)
	}
	if input.CategoryID != nil {
		upd = upd.SetCategoryID(*input.CategoryID)
	}
	if input.IsInPlace != nil {
		upd = upd.SetIsInPlace(*input.IsInPlace)
	}
	if input.AddressID != nil {
		// SetNillableAddressID accepts *uuid.UUID; passing the pointer will set or nil accordingly
		upd = upd.SetNillableAddressID(input.AddressID)
	}

	return upd.Save(ctx)
}

// DeleteService: Remove o serviço validando o proprietário via ID do Provider
func (r *ServiceRepository) DeleteService(ctx context.Context, id uuid.UUID, providerID uuid.UUID) error {
	// Delete().Where() retorna (int, error).
	_, err := r.client.Service.
		Delete().
		Where(
			service.ID(id),
			service.HasProviderWith(user.ID(providerID)),
		).
		Exec(ctx)

	return err
}
