package categories

import (
	"context"
	"fmt"

	models "github.com/matheusgosk8/book-me-server/internal/models"

	"github.com/matheusgosk8/book-me-server/ent"
	repositories "github.com/matheusgosk8/book-me-server/internal/repository"
)

type CategoriesService struct{}

func NewCategoriesService() *CategoriesService {
	return &CategoriesService{}
}

type ListInput struct {
	Pagination *models.Pagination
}

type ListOutput struct {
	Categories []*ent.Category
	Meta       *models.Meta
}

func (s *CategoriesService) ListMyCategories(ctx context.Context, input ListInput) (*ListOutput, error) {
	repo := repositories.NewListMyCategories()
	cats, total, err := repo.ListMyCategories(ctx, input.Pagination)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar categorias: %w", err)
	}

	meta := &models.Meta{
		Total:   total,
		Page:    input.Pagination.Page,
		PerPage: input.Pagination.PerPage,
	}

	return &ListOutput{Categories: cats, Meta: meta}, nil
}
