package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/dto"
	"github.com/matheusgosk8/book-me-server/internal/repository"
)

// ListServices busca serviços com suporte a paginação e filtro por categoria
func ListServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// 1. Tratamento de Query Parameters (Paginação e Filtros)
	query := r.URL.Query()

	// Define limite de itens (padrão 10)
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	// Define a página atual (padrão 1)
	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	// Filtro opcional por category_id
	var categoryID *uuid.UUID
	if catStr := query.Get("category_id"); catStr != "" {
		if id, err := uuid.Parse(catStr); err == nil {
			categoryID = &id
		}
	}

	// 2. Chamada ao Repository
	repo := repository.NewServiceRepository(db.Client)
	services, err := repo.ListServices(r.Context(), limit, offset, categoryID)
	if err != nil {
		http.Error(w, "Erro ao buscar serviços: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Mapeamento para DTO (Output limpo para o React Native)
	response := make([]dto.ServiceResponseDTO, 0)
	for _, s := range services {
		response = append(response, dto.ServiceResponseDTO{
			ID:              s.ID,
			Title:           s.Title,
			Description:     s.Description,
			PriceBase:       s.PriceBase,
			PriceType:       string(s.PriceType),
			DurationMinutes: s.DurationMinutes,
			IsActive:        s.IsActive,
			Provider: dto.ProviderShortResponse{
				ID:    s.Edges.Provider.ID,
				Email: s.Edges.Provider.Email,
			},
			Category: dto.CategoryShortResponse{
				ID:   s.Edges.Category.ID,
				Name: s.Edges.Category.Name,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateServiceHandler lida com a criação de novos serviços vinculados ao usuário logado
func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Recupera o user_id do contexto (setado pelo Middleware de Auth)
	providerIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Erro ao identificar usuário", http.StatusUnauthorized)
		return
	}
	providerID, _ := uuid.Parse(providerIDStr)

	// 2. Decodifica o payload de entrada
	var input dto.ServiceRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// 3. Persistência no banco via Repository
	repo := repository.NewServiceRepository(db.Client)
	newService, err := repo.CreateService(r.Context(), input, providerID)
	if err != nil {
		http.Error(w, "Erro ao criar serviço: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Mapeamento para DTO de retorno 
response := dto.ServiceResponseDTO{
    ID:              newService.ID,
    Title:           newService.Title,
    Description:     newService.Description,
    PriceBase:       newService.PriceBase,
    PriceType:       string(newService.PriceType),
    DurationMinutes: newService.DurationMinutes,
    IsActive:        newService.IsActive,
    Provider: dto.ProviderShortResponse{
        ID:    providerID, 
        Email: "",
    },
    Category: dto.CategoryShortResponse{
        ID:   input.CategoryID,
        Name: "",
    },
}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}