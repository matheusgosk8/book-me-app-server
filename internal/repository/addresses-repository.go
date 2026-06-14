package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/ent"
	"github.com/matheusgosk8/book-me-server/ent/address"
	"github.com/matheusgosk8/book-me-server/ent/user"
	"github.com/matheusgosk8/book-me-server/internal/db"
)

type ListMyAddresses struct{}

func NewListMyAddresses() *ListMyAddresses {
	return &ListMyAddresses{}
}

func (s *ListMyAddresses) ListMyAddresses(ctx context.Context, userID uuid.UUID) ([]*ent.Address, error) {
	return db.Client.Address.
		Query().
		Where(address.HasUserWith(user.IDEQ(userID))).
		All(ctx)
}
