package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/dto"
	"github.com/matheusgosk8/book-me-server/internal/repository"
)

// ListServices busca todos os serviços cadastrados e mapeia para o DTO de resposta
func ListServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	repo := repository.NewServiceRepository(db.Client)

	services, err := repo.ListServices(r.Context())
	if err != nil {
		http.Error(w, "Erro ao buscar serviços: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// MAPEAMENTO: Transforma a lista do Ent na lista do DTO (Sem dados sensíveis)
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

func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	providerIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Erro ao identificar usuário", http.StatusUnauthorized)
		return
	}
	providerID, _ := uuid.Parse(providerIDStr)

	var input dto.ServiceRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	repo := repository.NewServiceRepository(db.Client)
	newService, err := repo.CreateService(r.Context(), input, providerID)
	if err != nil {
		http.Error(w, "Erro ao criar serviço: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// MAPEAMENTO: Também no Create para manter o padrão de saída
	response := dto.ServiceResponseDTO{
		ID:              newService.ID,
		Title:           newService.Title,
		Description:     newService.Description,
		PriceBase:       newService.PriceBase,
		PriceType:       string(newService.PriceType),
		DurationMinutes: newService.DurationMinutes,
		IsActive:        newService.IsActive,
		Provider: dto.ProviderShortResponse{
			ID:    newService.Edges.Provider.ID,
			Email: newService.Edges.Provider.Email,
		},
		Category: dto.CategoryShortResponse{
			ID:   newService.Edges.Category.ID,
			Name: newService.Edges.Category.Name,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}