package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/repository"
)

// ListCategoriesHandler: Atende GET /public/categories
func ListCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewCategoryRepository(db.Client)
	
	// Busca todas as categorias no banco de dados
	categories, err := repo.ListAll(r.Context())
	if err != nil {
		http.Error(w, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}