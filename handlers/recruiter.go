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

type RecruiterHandler struct {
	service *services.RecruiterService
}

// NewRecruiterHandler creates a new recruiter handler
func NewRecruiterHandler(service *services.RecruiterService) *RecruiterHandler {
	return &RecruiterHandler{
		service: service,
	}
}

// GetAllRecruiters handles GET /recruiters request with pagination support
func (h *RecruiterHandler) GetAllRecruiters(w http.ResponseWriter, r *http.Request) {
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
		"first_name": r.URL.Query().Get("first_name"),
		"last_name":  r.URL.Query().Get("last_name"),
		"email":      r.URL.Query().Get("email"),
		"company_id": r.URL.Query().Get("company_id"),
	}

	// Parse sort parameters
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Get recruiters with pagination, filters, and sorting
	recruiters, total, err := h.service.GetAllRecruiters(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve recruiters", http.StatusInternalServerError)
		return
	}

	// Build paginated response
	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       recruiters,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetRecruiterByID handles GET /recruiters/:id request
func (h *RecruiterHandler) GetRecruiterByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	recruiterID := chi.URLParam(r, "id")

	recruiter, err := h.service.GetRecruiterByID(ctx, recruiterID)
	if err != nil {
		http.Error(w, "Recruiter not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recruiter); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateRecruiter handles POST /recruiters request
func (h *RecruiterHandler) CreateRecruiter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var recruiter models.Recruiter
	err := json.NewDecoder(r.Body).Decode(&recruiter)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	validationErrors := helpers.ValidateStruct(recruiter)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	recruiter.CreatedTime = time.Now()
	recruiter.UpdatedTime = time.Now()
	recruiter.CreatedBy = "system"
	recruiter.UpdatedBy = "system"

	err = h.service.CreateRecruiter(ctx, &recruiter)
	if err != nil {
		http.Error(w, "Failed to create recruiter", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(recruiter); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

// DeleteRecruiter handles DELETE /recruiters/{id} request
func (h *RecruiterHandler) DeleteRecruiter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	recruiterID := chi.URLParam(r, "id")

	err := h.service.DeleteRecruiter(ctx, recruiterID)
	if err != nil {
		http.Error(w, "Recruiter not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
