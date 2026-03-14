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

type LocationAvailabilityHandler struct {
	service interfaces.LocationAvailabilityService
}

func NewLocationAvailabilityHandler(service interfaces.LocationAvailabilityService) *LocationAvailabilityHandler {
	return &LocationAvailabilityHandler{service: service}
}

func (h *LocationAvailabilityHandler) GetAllLocationAvailabilities(w http.ResponseWriter, r *http.Request) {
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

	locationAvailabilities, total, err := h.service.GetAllLocationAvailabilities(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve location availabilities", http.StatusInternalServerError)
		return
	}

	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       locationAvailabilities,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *LocationAvailabilityHandler) GetLocationAvailabilityByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	locationAvailabilityID := chi.URLParam(r, "id")

	locationAvailability, err := h.service.GetLocationAvailabilityByID(ctx, locationAvailabilityID)
	if err != nil {
		http.Error(w, "Location availability not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locationAvailability); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *LocationAvailabilityHandler) CreateLocationAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var locationAvailability models.LocationAvailability
	err := json.NewDecoder(r.Body).Decode(&locationAvailability)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validationErrors := helpers.ValidateStruct(locationAvailability)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	locationAvailability.CreatedTime = time.Now()
	locationAvailability.UpdatedTime = time.Now()
	locationAvailability.CreatedBy = "system"
	locationAvailability.UpdatedBy = "system"

	err = h.service.CreateLocationAvailability(ctx, &locationAvailability)
	if err != nil {
		http.Error(w, "Failed to create location availability", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(locationAvailability); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

func (h *LocationAvailabilityHandler) UpdateLocationAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	locationAvailabilityID := chi.URLParam(r, "id")

	var locationAvailability models.LocationAvailability
	if err := json.NewDecoder(r.Body).Decode(&locationAvailability); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	locationAvailability.UpdatedTime = time.Now()
	locationAvailability.UpdatedBy = "system"

	updated, err := h.service.UpdateLocationAvailability(ctx, locationAvailabilityID, &locationAvailability)
	if err != nil {
		http.Error(w, "Location availability not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *LocationAvailabilityHandler) DeleteLocationAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	locationAvailabilityID := chi.URLParam(r, "id")

	err := h.service.DeleteLocationAvailability(ctx, locationAvailabilityID)
	if err != nil {
		http.Error(w, "Location availability not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
