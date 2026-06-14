package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/ent"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/dto"
	"github.com/matheusgosk8/book-me-server/internal/models"
	cmdrepo "github.com/matheusgosk8/book-me-server/internal/repository"
)

type ServiceBusiness struct {
	repo *cmdrepo.ServiceRepository
}

func NewServiceBusiness() *ServiceBusiness {
	return &ServiceBusiness{
		repo: cmdrepo.NewServiceRepository(db.Client),
	}
}

type ListInput struct {
	Pagination *models.Pagination
	CategoryID *uuid.UUID
}

type ListOutput struct {
	Services []*ent.Service
	Meta     *models.Meta
}

func (s *ServiceBusiness) ListServices(ctx context.Context, input ListInput) (*ListOutput, error) {
	out, err := s.repo.ListServices(ctx, input.Pagination.PerPage, (input.Pagination.Page-1)*input.Pagination.PerPage, input.CategoryID)
	if err != nil {
		return nil, err
	}

	meta := &models.Meta{
		Total:   out.Total,
		Page:    input.Pagination.Page,
		PerPage: input.Pagination.PerPage,
	}

	return &ListOutput{Services: out.Data, Meta: meta}, nil
}

// ListMyServices: lista serviços do provider via cmd repository padrão (retorna data + total)
func (s *ServiceBusiness) ListMyServices(ctx context.Context, providerID uuid.UUID, pagination *models.Pagination) (*ListOutput, error) {
	repo := cmdrepo.NewListServices()
	out, err := repo.ListMyServices(ctx, providerID, pagination)
	if err != nil {
		return nil, err
	}

	meta := &models.Meta{
		Total:   out.Total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	}

	return &ListOutput{Services: out.Data, Meta: meta}, nil
}

func (s *ServiceBusiness) CreateService(ctx context.Context, input dto.ServiceRequestDTO, providerID uuid.UUID) (*ent.Service, error) {
	if input.IsInPlace && (input.AddressID == nil || *input.AddressID == uuid.Nil) {
		return nil, errors.New("address_id é obrigatório para atendimento no local do prestador")
	}
	return s.repo.CreateService(ctx, input, providerID)
}

func (s *ServiceBusiness) UpdateService(ctx context.Context, serviceID, providerID uuid.UUID, input dto.ServicePatchDTO) (*ent.Service, error) {
	// Se o client está definindo is_in_place = true via patch, certifique-se de que trouxe address_id
	if input.IsInPlace != nil && *input.IsInPlace {
		if input.AddressID == nil || *input.AddressID == uuid.Nil {
			return nil, errors.New("address_id é obrigatório para serviços realizados no local")
		}
	}
	return s.repo.UpdateService(ctx, serviceID, providerID, input)
}

func (s *ServiceBusiness) DeleteService(ctx context.Context, serviceID, providerID uuid.UUID) error {
	return s.repo.DeleteService(ctx, serviceID, providerID)
}
