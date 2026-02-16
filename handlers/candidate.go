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

type CandidateHandler struct {
	service *services.CandidateService
}

// NewCandidateHandler creates a new candidate handler
func NewCandidateHandler(service *services.CandidateService) *CandidateHandler {
	return &CandidateHandler{
		service: service,
	}
}

// GetAllCandidates handles GET /candidates request
func (h *CandidateHandler) GetAllCandidates(w http.ResponseWriter, r *http.Request) {
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
		"location":   r.URL.Query().Get("location"),
	}

	// Parse sort parameters
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Get Candidates with pagination, filters, and sorting
	candidates, total, err := h.service.GetAllCandidates(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve candidates", http.StatusInternalServerError)
		return
	}

	// Build paginated response
	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       candidates,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

// GetCandidateByID handles GET /candidates/{id} request
func (h *CandidateHandler) GetCandidateByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := chi.URLParam(r, "id")

	candidate, err := h.service.GetCandidateByID(ctx, candidateID)
	if err != nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(candidate); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateCandidate handles POST /candidates request
func (h *CandidateHandler) CreateCandidate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var candidate models.Candidate
	err := json.NewDecoder(r.Body).Decode(&candidate)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	validationErrors := helpers.ValidateStruct(candidate)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	candidate.CreatedTime = time.Now()
	candidate.UpdatedTime = time.Now()
	candidate.CreatedBy = "system"
	candidate.UpdatedBy = "system"

	err = h.service.CreateCandidate(ctx, &candidate)
	if err != nil {
		http.Error(w, "Failed to create candidate", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(candidate); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

// DeleteCandidate handles DELETE /candidates/{id} request
func (h *CandidateHandler) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := chi.URLParam(r, "id")

	err := h.service.DeleteCandidate(ctx, candidateID)
	if err != nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
