package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/dto"
	"github.com/matheusgosk8/book-me-server/internal/repository"
)

// ListServices: Atende GET /customer/services e GET /provider/services
func ListServices(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	var categoryID *uuid.UUID
	if catStr := query.Get("category_id"); catStr != "" {
		if id, err := uuid.Parse(catStr); err == nil {
			categoryID = &id
		}
	}

	repo := repository.NewServiceRepository(db.Client)
	services, err := repo.ListServices(r.Context(), limit, offset, categoryID)
	if err != nil {
		http.Error(w, "Erro ao buscar serviços", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

// CreateServiceHandler: POST /provider/services
func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Recupera o tipo do usuário do contexto injetado pelo middleware
	userType, _ := r.Context().Value("user_type").(string)

	// Validação de segurança: apenas prestadores ou colaboradores
	if userType != "provider" && userType != "collaborator" {
		http.Error(w, "Acesso negado: apenas prestadores podem gerenciar serviços", http.StatusForbidden)
		return
	}

	// Recupera o ID do prestador
	providerIDStr, _ := r.Context().Value("user_id").(string)
	providerID, _ := uuid.Parse(providerIDStr)

	var input dto.ServiceRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Regra de negócio: is_in_place exige address_id
	if input.IsInPlace && (input.AddressID == nil || *input.AddressID == uuid.Nil) {
		http.Error(w, "address_id é obrigatório para atendimento no local do prestador", http.StatusBadRequest)
		return
	}

	repo := repository.NewServiceRepository(db.Client)
	newService, err := repo.CreateService(r.Context(), input, providerID)
	if err != nil {
		http.Error(w, "Erro ao salvar serviço", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newService)
}

// UpdateServiceHandler: PUT /provider/services/{id}
// UpdateServiceHandler: PUT /provider/services/{id}
func UpdateServiceHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Validação de userType (Inclusão de Collaborator)
	userType, _ := r.Context().Value("user_type").(string)
	if userType != "provider" && userType != "collaborator" {
		http.Error(w, "Acesso negado", http.StatusForbidden)
		return
	}

    // 2. Extração do ID do executor para garantir posse do registro
	providerIDStr, _ := r.Context().Value("user_id").(string)
	providerID, _ := uuid.Parse(providerIDStr)

	serviceIDStr := chi.URLParam(r, "id")
	serviceID, err := uuid.Parse(serviceIDStr)
    if err != nil {
		http.Error(w, "ID do serviço inválido", http.StatusBadRequest)
		return
	}

	var input dto.ServiceRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 3. Validação da Regra de Negócio Crítica
	if input.IsInPlace && (input.AddressID == nil || *input.AddressID == uuid.Nil) {
		http.Error(w, "address_id é obrigatório para serviços realizados no local", http.StatusBadRequest)
		return
	}

	repo := repository.NewServiceRepository(db.Client)
    // Passamos o providerID para o repositório validar a posse
	updated, err := repo.UpdateService(r.Context(), serviceID, providerID, input)
	if err != nil {
		http.Error(w, "Erro ao atualizar: verifique se o serviço pertence a você", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updated)
}

// DeleteServiceHandler: DELETE /provider/services/{id}
func DeleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	userType, _ := r.Context().Value("user_type").(string)
	if userType != "provider" && userType != "collaborator" {
		http.Error(w, "Acesso negado", http.StatusForbidden)
		return
	}

	providerIDStr, _ := r.Context().Value("user_id").(string)
	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
		return
	}

	serviceIDStr := chi.URLParam(r, "id")
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	repo := repository.NewServiceRepository(db.Client)
	if err := repo.DeleteService(r.Context(), serviceID, providerID); err != nil {
		http.Error(w, "Erro ao deletar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
