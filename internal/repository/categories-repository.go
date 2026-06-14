package repository

import (
	"context"

	"github.com/matheusgosk8/book-me-server/ent"
	"github.com/matheusgosk8/book-me-server/internal/db"
	models "github.com/matheusgosk8/book-me-server/internal/models"
)

type ListMyCategories struct{}

func NewListMyCategories() *ListMyCategories {
	return &ListMyCategories{}
}

func (s *ListMyCategories) ListMyCategories(ctx context.Context, pagination *models.Pagination) ([]*ent.Category, int, error) {
	q := db.Client.Category.Query()

	// total com os mesmos filtros aplicados (se houverem)
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	cats, err := q.
		Limit(pagination.PerPage).
		Offset((pagination.Page - 1) * pagination.PerPage).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return cats, total, nil
}
