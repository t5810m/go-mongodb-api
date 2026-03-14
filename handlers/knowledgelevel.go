package handlers

import (
	"encoding/json"
	"go-mongodb-api/helpers"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type KnowledgeLevelHandler struct {
	service interfaces.KnowledgeLevelService
}

func NewKnowledgeLevelHandler(service interfaces.KnowledgeLevelService) *KnowledgeLevelHandler {
	return &KnowledgeLevelHandler{service: service}
}

func (h *KnowledgeLevelHandler) GetAllKnowledgeLevels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	filters := map[string]string{
		"title": r.URL.Query().Get("title"),
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	knowledgeLevels, total, err := h.service.GetAllKnowledgeLevels(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve knowledge levels", http.StatusInternalServerError)
		return
	}

	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       knowledgeLevels,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *KnowledgeLevelHandler) GetKnowledgeLevelByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	knowledgeLevelID := chi.URLParam(r, "id")

	knowledgeLevel, err := h.service.GetKnowledgeLevelByID(ctx, knowledgeLevelID)
	if err != nil {
		http.Error(w, "Knowledge level not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(knowledgeLevel); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *KnowledgeLevelHandler) CreateKnowledgeLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var knowledgeLevel models.KnowledgeLevel
	err := json.NewDecoder(r.Body).Decode(&knowledgeLevel)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validationErrors := helpers.ValidateStruct(knowledgeLevel)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	knowledgeLevel.CreatedTime = time.Now()
	knowledgeLevel.UpdatedTime = time.Now()
	knowledgeLevel.CreatedBy = "system"
	knowledgeLevel.UpdatedBy = "system"

	err = h.service.CreateKnowledgeLevel(ctx, &knowledgeLevel)
	if err != nil {
		http.Error(w, "Failed to create knowledge level", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(knowledgeLevel); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

func (h *KnowledgeLevelHandler) UpdateKnowledgeLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	knowledgeLevelID := chi.URLParam(r, "id")

	var knowledgeLevel models.KnowledgeLevel
	if err := json.NewDecoder(r.Body).Decode(&knowledgeLevel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	knowledgeLevel.UpdatedTime = time.Now()
	knowledgeLevel.UpdatedBy = "system"

	updated, err := h.service.UpdateKnowledgeLevel(ctx, knowledgeLevelID, &knowledgeLevel)
	if err != nil {
		http.Error(w, "Knowledge level not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *KnowledgeLevelHandler) DeleteKnowledgeLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	knowledgeLevelID := chi.URLParam(r, "id")

	err := h.service.DeleteKnowledgeLevel(ctx, knowledgeLevelID)
	if err != nil {
		http.Error(w, "Knowledge level not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
