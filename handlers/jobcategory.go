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

type JobCategoryHandler struct {
	service *services.JobCategoryService
}

// NewJobCategoryHandler creates a new job category handler
func NewJobCategoryHandler(service *services.JobCategoryService) *JobCategoryHandler {
	return &JobCategoryHandler{
		service: service,
	}
}

// GetAllJobCategories handles GET /jobcategories request with pagination support
func (h *JobCategoryHandler) GetAllJobCategories(w http.ResponseWriter, r *http.Request) {
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
		"description": r.URL.Query().Get("description"),
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Get job categories with pagination, filters, and sorting
	jobCategories, total, err := h.service.GetAllJobCategories(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve job categories", http.StatusInternalServerError)
		return
	}

	// Build paginated response
	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       jobCategories,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetJobCategoryByID handles GET /jobcategories/{id} request
func (h *JobCategoryHandler) GetJobCategoryByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobCategoryID := chi.URLParam(r, "id")

	jobCategory, err := h.service.GetJobCategoryByID(ctx, jobCategoryID)
	if err != nil {
		http.Error(w, "Job category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobCategory); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateJobCategory handles POST /jobcategories request
func (h *JobCategoryHandler) CreateJobCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var jobCategory models.JobCategory
	err := json.NewDecoder(r.Body).Decode(&jobCategory)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	validationErrors := helpers.ValidateStruct(jobCategory)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	jobCategory.CreatedTime = time.Now()
	jobCategory.UpdatedTime = time.Now()
	jobCategory.CreatedBy = "system"
	jobCategory.UpdatedBy = "system"

	err = h.service.CreateJobCategory(ctx, &jobCategory)
	if err != nil {
		http.Error(w, "Failed to create job category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(jobCategory); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

// DeleteJobCategory handles DELETE /jobcategories/{id} request
func (h *JobCategoryHandler) DeleteJobCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobCategoryID := chi.URLParam(r, "id")

	err := h.service.DeleteJobCategory(ctx, jobCategoryID)
	if err != nil {
		http.Error(w, "Job category not found or has related jobs", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
