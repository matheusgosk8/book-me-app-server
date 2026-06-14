package categories

import (
	"net/http"

	models "github.com/matheusgosk8/book-me-server/internal/models"
	service "github.com/matheusgosk8/book-me-server/internal/service/categories"
	"github.com/matheusgosk8/book-me-server/internal/utils"
)

func ListCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	pageSizeInt := utils.QueryInt(r, "perPage", 10)
	pageInt := utils.QueryInt(r, "page", 1)

	categoriesService := service.NewCategoriesService()
	categories, err := categoriesService.ListMyCategories(r.Context(), service.ListInput{Pagination: &models.Pagination{
		Page:    pageInt,
		PerPage: pageSizeInt,
	}})

	if err != nil {
		http.Error(w, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if categories != nil && categories.Meta != nil {
		utils.ServerResponse(w, categories.Categories, *categories.Meta)
	} else if categories != nil {
		utils.ServerResponse(w, categories.Categories)
	} else {
		utils.ServerResponse(w, []any{})
	}

}
