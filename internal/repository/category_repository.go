package repository

import (
	"context"
	"github.com/matheusgosk8/book-me-server/ent"
	"github.com/matheusgosk8/book-me-server/ent/category"
)

type CategoryRepository struct {
	client *ent.Client
}


func NewCategoryRepository(client *ent.Client) *CategoryRepository {
	return &CategoryRepository{client: client}
}

func (r *CategoryRepository) ListAll(ctx context.Context) ([]*ent.Category, error) {
	return r.client.Category.
		Query().
		Order(ent.Asc(category.FieldName)).
		All(ctx)
}