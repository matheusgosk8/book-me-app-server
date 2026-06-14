package address

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/ent"
	repositories "github.com/matheusgosk8/book-me-server/internal/repository"
)

type AddressService struct{}

func NewAddressService() *AddressService {
	return &AddressService{}
}

type ListInput struct {
	UserID uuid.UUID
}

type ListOutput struct {
	Addresses []*ent.Address
}

func (s *AddressService) ListMyAddresses(ctx context.Context, input ListInput) (*ListOutput, error) {
	repo := repositories.NewListMyAddresses()
	addrs, err := repo.ListMyAddresses(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar endereços: %w", err)
	}
	return &ListOutput{Addresses: addrs}, nil
}
