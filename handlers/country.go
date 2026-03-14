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

type CountryHandler struct {
	service interfaces.CountryService
}

func NewCountryHandler(service interfaces.CountryService) *CountryHandler {
	return &CountryHandler{service: service}
}

func (h *CountryHandler) GetAllCountries(w http.ResponseWriter, r *http.Request) {
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
		"name": r.URL.Query().Get("name"),
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	countries, total, err := h.service.GetAllCountries(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve countries", http.StatusInternalServerError)
		return
	}

	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       countries,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *CountryHandler) GetCountryByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	countryID := chi.URLParam(r, "id")

	country, err := h.service.GetCountryByID(ctx, countryID)
	if err != nil {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(country); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *CountryHandler) CreateCountry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var country models.Country
	err := json.NewDecoder(r.Body).Decode(&country)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validationErrors := helpers.ValidateStruct(country)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	country.CreatedTime = time.Now()
	country.UpdatedTime = time.Now()
	country.CreatedBy = "system"
	country.UpdatedBy = "system"

	err = h.service.CreateCountry(ctx, &country)
	if err != nil {
		http.Error(w, "Failed to create country", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(country); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

func (h *CountryHandler) UpdateCountry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	countryID := chi.URLParam(r, "id")

	var country models.Country
	if err := json.NewDecoder(r.Body).Decode(&country); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	country.UpdatedTime = time.Now()
	country.UpdatedBy = "system"

	updated, err := h.service.UpdateCountry(ctx, countryID, &country)
	if err != nil {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *CountryHandler) DeleteCountry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	countryID := chi.URLParam(r, "id")

	err := h.service.DeleteCountry(ctx, countryID)
	if err != nil {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
