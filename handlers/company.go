package handlers

import (
	"encoding/json"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"
	"go-mongodb-api/services"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type CompanyHandler struct {
	service *services.CompanyService
}

// NewCompanyHandler creates a new company handler
func NewCompanyHandler(service *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		service: service,
	}
}

// GetAllCompanies handles GET /companies request with pagination support
func (h *CompanyHandler) GetAllCompanies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Convert to integers with defaults
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

	// Parse search filters
	filters := map[string]string{
		"name":        r.URL.Query().Get("name"),
		"location":    r.URL.Query().Get("location"),
		"country":     r.URL.Query().Get("country"),
		"city":        r.URL.Query().Get("city"),
		"postal_code": r.URL.Query().Get("postal_code"),
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Get companies with pagination, filters, and sorting
	companies, total, err := h.service.GetAllCompanies(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve companies", http.StatusInternalServerError)
		return
	}

	// Build paginated response
	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       companies,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetCompanyByID handles GET /companies/{id} request
func (h *CompanyHandler) GetCompanyByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	companyID := chi.URLParam(r, "id")

	company, err := h.service.GetCompanyByID(ctx, companyID)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(company); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateCompany handles POST /companies request
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var company models.Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	validationErrors := helpers.ValidateStruct(company)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	company.CreatedTime = time.Now()
	company.UpdatedTime = time.Now()
	company.CreatedBy = "system"
	company.UpdatedBy = "system"

	err = h.service.CreateCompany(ctx, &company)
	if err != nil {
		http.Error(w, "Failed to create company", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(company); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

// DeleteCompany handles DELETE /companies/{id} request
func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	companyID := chi.URLParam(r, "id")

	err := h.service.DeleteCompany(ctx, companyID)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
