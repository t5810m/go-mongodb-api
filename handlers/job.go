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

type JobHandler struct {
	service *services.JobService
}

// NewJobHandler creates a new job handler
func NewJobHandler(service *services.JobService) *JobHandler {
	return &JobHandler{
		service: service,
	}
}

// GetAllJobs handles GET /jobs request with pagination and search support
func (h *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse pagination parameters
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
		"title":       r.URL.Query().Get("title"),
		"description": r.URL.Query().Get("description"),
		"location":    r.URL.Query().Get("location"),
		"job_type":    r.URL.Query().Get("job_type"),
		"status":      r.URL.Query().Get("status"),
	}

	// Parse sort parameters
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Get jobs with pagination, filters, and sorting
	jobs, total, err := h.service.GetAllJobs(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve jobs", http.StatusInternalServerError)
		return
	}

	// Build paginated response
	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       jobs,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetJobByID handles GET /jobs/:id request
func (h *JobHandler) GetJobByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobID := chi.URLParam(r, "id")

	job, err := h.service.GetJobByID(ctx, jobID)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(job); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetJobsByCompany handles GET /companies/:companyId/jobs request
func (h *JobHandler) GetJobsByCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	companyID := chi.URLParam(r, "companyId")

	jobs, err := h.service.GetJobsByCompany(ctx, companyID)
	if err != nil {
		http.Error(w, "Failed to retrieve jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateJob handles POST /jobs request
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var job models.Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	validationErrors := helpers.ValidateStruct(job)

	// Validate cross-field constraint: salary_max must be >= salary_min
	if job.SalaryMax < job.SalaryMin {
		validationErrors = append(validationErrors, helpers.ValidationError{
			Field:   "salary_max",
			Message: "must be greater than or equal to salary_min",
		})
	}

	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	job.CreatedTime = time.Now()
	job.UpdatedTime = time.Now()
	job.CreatedBy = "system"
	job.UpdatedBy = "system"

	err = h.service.CreateJob(ctx, &job)
	if err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(job); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

// DeleteJob handles DELETE /jobs/{id} request
func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobID := chi.URLParam(r, "id")

	err := h.service.DeleteJob(ctx, jobID)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
