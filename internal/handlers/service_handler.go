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

// ListServices: Busca serviços com paginação e filtro por categoria
func ListServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	// Configuração de paginação (Limit e Offset)
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	// Filtro opcional por ID da categoria
	var categoryID *uuid.UUID
	if catStr := query.Get("category_id"); catStr != "" {
		if id, err := uuid.Parse(catStr); err == nil {
			categoryID = &id
		}
	}

	repo := repository.NewServiceRepository(db.Client)
	services, err := repo.ListServices(r.Context(), limit, offset, categoryID)
	if err != nil {
		http.Error(w, "Erro ao buscar serviços: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Mapeamento para DTO de resposta
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
			IsInPlace:       s.IsInPlace, // Adicionado conforme regra do Notion
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

// CreateServiceHandler: Criação de serviços com validação de local e endereço
func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Verificação do tipo de utilizador (Segurança)
	userType, _ := r.Context().Value("user_type").(string)
	if userType != "provider" && userType != "collaborator" {
		http.Error(w, "Acesso negado: apenas prestadores podem criar serviços", http.StatusForbidden)
		return
	}

	// Recuperação do ID do prestador do contexto
	providerIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Erro ao identificar usuário", http.StatusUnauthorized)
		return
	}
	providerID, _ := uuid.Parse(providerIDStr)

	var input dto.ServiceRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Dados inválidos no JSON", http.StatusBadRequest)
		return
	}

	// Regra de Negócio: Se for no local do prestador, addressId é obrigatório
	if input.IsInPlace && (input.AddressID == nil || *input.AddressID == uuid.Nil) {
		http.Error(w, "address_id é obrigatório para serviços em local fixo", http.StatusBadRequest)
		return
	}

	repo := repository.NewServiceRepository(db.Client)
	newService, err := repo.CreateService(r.Context(), input, providerID)
	if err != nil {
		http.Error(w, "Erro ao salvar serviço: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Resposta de sucesso limpa (evita Panic de Edges vazias)
	response := dto.ServiceResponseDTO{
		ID:              newService.ID,
		Title:           newService.Title,
		Description:     newService.Description,
		PriceBase:       newService.PriceBase,
		PriceType:       string(newService.PriceType),
		DurationMinutes: newService.DurationMinutes,
		IsActive:        newService.IsActive,
		IsInPlace:       newService.IsInPlace,
		Provider: dto.ProviderShortResponse{
			ID: providerID,
		},
		Category: dto.CategoryShortResponse{
			ID: input.CategoryID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DeleteServiceHandler: Elimina o serviço (compatível com roteador padrão Go 1.22+)
func DeleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Uso do PathValue em vez do Chi para capturar o ID da URL
	serviceIDStr := r.PathValue("id")
	if serviceIDStr == "" {
		http.Error(w, "ID do serviço não fornecido", http.StatusBadRequest)
		return
	}

	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		http.Error(w, "ID de serviço inválido", http.StatusBadRequest)
		return
	}

	repo := repository.NewServiceRepository(db.Client)
	err = repo.DeleteService(r.Context(), serviceID)
	if err != nil {
		http.Error(w, "Erro ao eliminar serviço: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}