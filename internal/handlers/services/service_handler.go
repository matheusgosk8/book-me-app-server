package services

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/internal/dto"
	"github.com/matheusgosk8/book-me-server/internal/models"
	service "github.com/matheusgosk8/book-me-server/internal/service/services"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	vld "github.com/matheusgosk8/book-me-server/internal/validator"
)

func ListServices(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	perPage := utils.QueryInt(r, "perPage", 10)
	page := utils.QueryInt(r, "page", 1)

	var categoryID *uuid.UUID
	if catStr := query.Get("category_id"); catStr != "" {
		if id, err := uuid.Parse(catStr); err == nil {
			categoryID = &id
		}
	}

	svc := service.NewServiceBusiness()
	out, err := svc.ListServices(r.Context(), service.ListInput{Pagination: &models.Pagination{Page: page, PerPage: perPage}, CategoryID: categoryID})
	if err != nil {
		utils.ServerError(w, http.StatusInternalServerError, err)
		return
	}

	if out != nil && out.Meta != nil {
		utils.ServerResponse(w, out.Services, *out.Meta)
	} else if out != nil {
		utils.ServerResponse(w, out.Services)
	} else {
		utils.ServerResponse(w, []any{})
	}
}

func ListMyServices(w http.ResponseWriter, r *http.Request) {
	userType, _ := r.Context().Value("user_type").(string)

	if userType != "provider" && userType != "collaborator" {
		if userType == "customer" {
			utils.ServerError(w, http.StatusForbidden, errors.New("Acesso negado: clientes não podem acessar este recurso"))
		} else {
			utils.ServerError(w, http.StatusForbidden, errors.New("Acesso negado: apenas providers e colaboradores podem acessar seus serviços"))
		}
		return
	}

	providerIDStr, _ := r.Context().Value("user_id").(string)
	providerID, _ := uuid.Parse(providerIDStr)

	perPage := utils.QueryInt(r, "perPage", 10)
	page := utils.QueryInt(r, "page", 1)

	svc := service.NewServiceBusiness()
	out, err := svc.ListMyServices(r.Context(), providerID, &models.Pagination{Page: page, PerPage: perPage})
	if err != nil {
		utils.ServerError(w, http.StatusInternalServerError, err)
		return
	}

	if out != nil && out.Meta != nil {
		utils.ServerResponse(w, out.Services, *out.Meta)
	} else if out != nil {
		utils.ServerResponse(w, out.Services)
	} else {
		utils.ServerResponse(w, []any{})
	}
}

func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	userType, _ := r.Context().Value("user_type").(string)
	if userType != "provider" && userType != "collaborator" {
		if userType == "customer" {
			utils.ServerError(w, http.StatusForbidden, errors.New("Acesso negado: clientes não podem criar serviços"))
		} else {
			utils.ServerError(w, http.StatusForbidden, errors.New("Acesso negado: apenas providers e colaboradores podem criar serviços"))
		}
		return
	}

	providerIDStr, _ := r.Context().Value("user_id").(string)
	providerID, _ := uuid.Parse(providerIDStr)

	// parse into validation DTO to provide friendly validation messages
	input, err := utils.BodyParser[vld.ServiceDTO](r)
	if err != nil {
		utils.RequestErrorHandler(w, err)
		return
	}

	if msgs := vld.ValidateService(*input); msgs != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.ServerResponse(w, msgs)
		return
	}

	// map validated input to internal DTO
	var svcReq dto.ServiceRequestDTO
	svcReq.Title = input.Title
	svcReq.Description = input.Description
	svcReq.PriceBase = input.PriceBase
	svcReq.PriceType = input.PriceType
	svcReq.DurationMinutes = input.DurationMinutes
	svcReq.IsActive = input.IsActive
	svcReq.IsInPlace = input.IsInPlace
	if id, err := uuid.Parse(input.CategoryID); err == nil {
		svcReq.CategoryID = id
	}
	if input.AddressID != nil {
		if aid, err := uuid.Parse(*input.AddressID); err == nil {
			svcReq.AddressID = &aid
		}
	}

	log.Printf("Payload do serviço: %+v", svcReq)

	svc := service.NewServiceBusiness()
	newService, err := svc.CreateService(r.Context(), svcReq, providerID)
	if err != nil {
		// map Postgres-specific errors to friendly messages when possible
		msg, status := utils.HandlePGError(err)
		utils.ServerError(w, status, errors.New(msg))
		return
	}

	utils.ServerSuccess(w, http.StatusCreated, "created", newService)
}

func UpdateServiceHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Payload de edição: %+v", r.Context())

	userType, _ := r.Context().Value("user_type").(string)
	if userType != "provider" && userType != "collaborator" {
		if userType == "customer" {
			utils.ServerError(w, http.StatusForbidden, errors.New("Acesso negado: clientes não podem editar serviços"))
		} else {
			utils.ServerError(w, http.StatusForbidden, errors.New("Acesso negado: apenas providers e colaboradores podem editar serviços"))
		}
		return
	}

	providerIDStr, _ := r.Context().Value("user_id").(string)
	providerID, _ := uuid.Parse(providerIDStr)

	serviceIDStr := chi.URLParam(r, "id")
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		utils.ServerError(w, http.StatusBadRequest, err)
		return
	}

	input, err := utils.BodyParser[vld.ServicePatchDTO](r)
	if err != nil {
		utils.RequestErrorHandler(w, err)
		return
	}

	// Validate patch input (checks only fields provided)
	if msgs := vld.ValidateServicePatch(*input); msgs != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.ServerResponse(w, msgs)
		return
	}

	// map to DTO de patch usado pelo service/repo
	var svcPatch dto.ServicePatchDTO
	if input.Title != nil {
		svcPatch.Title = input.Title
	}
	if input.Description != nil {
		svcPatch.Description = input.Description
	}
	if input.PriceBase != nil {
		svcPatch.PriceBase = input.PriceBase
	}
	if input.PriceType != nil {
		svcPatch.PriceType = input.PriceType
	}
	if input.DurationMinutes != nil {
		svcPatch.DurationMinutes = input.DurationMinutes
	}
	if input.CategoryID != nil {
		if id, err := uuid.Parse(*input.CategoryID); err == nil {
			svcPatch.CategoryID = &id
		}
	}
	if input.IsActive != nil {
		svcPatch.IsActive = input.IsActive
	}
	if input.IsInPlace != nil {
		svcPatch.IsInPlace = input.IsInPlace
	}
	if input.AddressID != nil {
		if *input.AddressID == "" {
			// explicit empty -> set nil pointer
			nilptr := (*uuid.UUID)(nil)
			svcPatch.AddressID = nilptr
		} else if aid, err := uuid.Parse(*input.AddressID); err == nil {
			svcPatch.AddressID = &aid
		}
	}

	svc := service.NewServiceBusiness()
	updated, err := svc.UpdateService(r.Context(), serviceID, providerID, svcPatch)
	if err != nil {
		msg, status := utils.HandlePGError(err)
		utils.ServerError(w, status, errors.New(msg))
		return
	}

	utils.ServerResponse(w, updated)
}

func DeleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	userType, _ := r.Context().Value("user_type").(string)
	if userType != "provider" && userType != "collaborator" {
		if userType == "customer" {
			utils.ServerError(w, http.StatusForbidden, errors.New("Acesso negado: clientes não podem deletar serviços"))
		} else {
			utils.ServerError(w, http.StatusForbidden, errors.New("Acesso negado: apenas providers e colaboradores podem deletar serviços"))
		}
		return
	}

	providerIDStr, _ := r.Context().Value("user_id").(string)
	providerID, _ := uuid.Parse(providerIDStr)

	serviceIDStr := chi.URLParam(r, "id")
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		utils.ServerError(w, http.StatusBadRequest, err)
		return
	}

	svc := service.NewServiceBusiness()
	if err := svc.DeleteService(r.Context(), serviceID, providerID); err != nil {
		msg, status := utils.HandlePGError(err)
		utils.ServerError(w, status, errors.New(msg))
		return
	}

	utils.ServerSuccess[any](w, http.StatusOK, "Serviço deletado com sucesso", nil)
}
