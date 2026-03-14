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

type EducationLevelHandler struct {
	service interfaces.EducationLevelService
}

func NewEducationLevelHandler(service interfaces.EducationLevelService) *EducationLevelHandler {
	return &EducationLevelHandler{service: service}
}

func (h *EducationLevelHandler) GetAllEducationLevels(w http.ResponseWriter, r *http.Request) {
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

	educationLevels, total, err := h.service.GetAllEducationLevels(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve education levels", http.StatusInternalServerError)
		return
	}

	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       educationLevels,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *EducationLevelHandler) GetEducationLevelByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	educationLevelID := chi.URLParam(r, "id")

	educationLevel, err := h.service.GetEducationLevelByID(ctx, educationLevelID)
	if err != nil {
		http.Error(w, "Education level not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(educationLevel); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *EducationLevelHandler) CreateEducationLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var educationLevel models.EducationLevel
	err := json.NewDecoder(r.Body).Decode(&educationLevel)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validationErrors := helpers.ValidateStruct(educationLevel)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	educationLevel.CreatedTime = time.Now()
	educationLevel.UpdatedTime = time.Now()
	educationLevel.CreatedBy = "system"
	educationLevel.UpdatedBy = "system"

	err = h.service.CreateEducationLevel(ctx, &educationLevel)
	if err != nil {
		http.Error(w, "Failed to create education level", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(educationLevel); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

func (h *EducationLevelHandler) UpdateEducationLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	educationLevelID := chi.URLParam(r, "id")

	var educationLevel models.EducationLevel
	if err := json.NewDecoder(r.Body).Decode(&educationLevel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	educationLevel.UpdatedTime = time.Now()
	educationLevel.UpdatedBy = "system"

	updated, err := h.service.UpdateEducationLevel(ctx, educationLevelID, &educationLevel)
	if err != nil {
		http.Error(w, "Education level not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *EducationLevelHandler) DeleteEducationLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	educationLevelID := chi.URLParam(r, "id")

	err := h.service.DeleteEducationLevel(ctx, educationLevelID)
	if err != nil {
		http.Error(w, "Education level not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
