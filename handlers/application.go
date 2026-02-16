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

type ApplicationHandler struct {
	service *services.ApplicationService
}

// NewApplicationHandler creates a new application handler
func NewApplicationHandler(service *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		service: service,
	}
}

// GetAllApplications handles GET /applications request with pagination support
func (h *ApplicationHandler) GetAllApplications(w http.ResponseWriter, r *http.Request) {
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
		"status":       r.URL.Query().Get("status"),
		"job_id":       r.URL.Query().Get("job_id"),
		"candidate_id": r.URL.Query().Get("candidate_id"),
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Get applications with pagination, filters, and sorting
	applications, total, err := h.service.GetAllApplications(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	// Build paginated response
	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       applications,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetApplicationByID handles GET /applications/{id} request
func (h *ApplicationHandler) GetApplicationByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	applicationID := chi.URLParam(r, "id")

	application, err := h.service.GetApplicationByID(ctx, applicationID)
	if err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(application); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetApplicationsByJobID handles GET /jobs/{jobId}/applications request
func (h *ApplicationHandler) GetApplicationsByJobID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobID := chi.URLParam(r, "jobId")

	applications, err := h.service.GetApplicationsByJobID(ctx, jobID)
	if err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(applications); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetApplicationsByCandidateID handles GET /candidates/{candidateId}/applications request
func (h *ApplicationHandler) GetApplicationsByCandidateID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := chi.URLParam(r, "candidateId")

	applications, err := h.service.GetApplicationsByCandidateID(ctx, candidateID)
	if err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(applications); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateApplication handles POST /applications request
func (h *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var application models.Application
	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	validationErrors := helpers.ValidateStruct(application)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	application.AppliedTime = time.Now()
	application.UpdatedTime = time.Now()
	application.CreatedBy = "system"
	application.UpdatedBy = "system"

	err = h.service.CreateApplication(ctx, &application)
	if err != nil {
		http.Error(w, "Failed to create application", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(application); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

// UpdateApplicationStatus handles PUT /applications/{id} request
func (h *ApplicationHandler) UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	applicationID := chi.URLParam(r, "id")

	var request struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateApplicationStatus(ctx, applicationID, request.Status)
	if err != nil {
		http.Error(w, "Failed to update application", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// DeleteApplication handles DELETE /applications/{id} request
func (h *ApplicationHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	applicationID := chi.URLParam(r, "id")

	err := h.service.DeleteApplication(ctx, applicationID)
	if err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
