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

type JobSkillHandler struct {
	service *services.JobSkillService
}

// NewJobSkillHandler creates a new job skill handler
func NewJobSkillHandler(service *services.JobSkillService) *JobSkillHandler {
	return &JobSkillHandler{
		service: service,
	}
}

// GetAllJobSkills handles GET /jobskills request with pagination support
func (h *JobSkillHandler) GetAllJobSkills(w http.ResponseWriter, r *http.Request) {
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
		"job_id":                     r.URL.Query().Get("job_id"),
		"skill_id":                   r.URL.Query().Get("skill_id"),
		"proficiency_level_required": r.URL.Query().Get("proficiency_level_required"),
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Get job skills with pagination, filters, and sorting
	jobSkills, total, err := h.service.GetAllJobSkills(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve job skills", http.StatusInternalServerError)
		return
	}

	// Build paginated response
	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       jobSkills,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetJobSkillByID handles GET /jobskills/{id} request
func (h *JobSkillHandler) GetJobSkillByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobSkillID := chi.URLParam(r, "id")

	jobSkill, err := h.service.GetJobSkillByID(ctx, jobSkillID)
	if err != nil {
		http.Error(w, "Job skill not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobSkill); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetJobSkillsByJobID handles GET /jobs/{jobId}/skills request
func (h *JobSkillHandler) GetJobSkillsByJobID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobID := chi.URLParam(r, "jobId")

	jobSkills, err := h.service.GetJobSkillsByJobID(ctx, jobID)
	if err != nil {
		http.Error(w, "Failed to retrieve job skills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobSkills); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateJobSkill handles POST /jobskills request
func (h *JobSkillHandler) CreateJobSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var jobSkill models.JobSkill
	err := json.NewDecoder(r.Body).Decode(&jobSkill)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	validationErrors := helpers.ValidateStruct(jobSkill)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	jobSkill.CreatedTime = time.Now()
	jobSkill.UpdatedTime = time.Now()
	jobSkill.CreatedBy = "system"
	jobSkill.UpdatedBy = "system"

	err = h.service.CreateJobSkill(ctx, &jobSkill)
	if err != nil {
		http.Error(w, "Failed to create job skill", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(jobSkill); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

// UpdateJobSkillProficiencyLevel handles PUT /jobskills/{id} request
func (h *JobSkillHandler) UpdateJobSkillProficiencyLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobSkillID := chi.URLParam(r, "id")

	var request struct {
		ProficiencyLevelRequired string `json:"proficiency_level_required"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateJobSkillProficiencyLevel(ctx, jobSkillID, request.ProficiencyLevelRequired)
	if err != nil {
		http.Error(w, "Failed to update job skill", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// DeleteJobSkill handles DELETE /jobskills/{id} request
func (h *JobSkillHandler) DeleteJobSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobSkillID := chi.URLParam(r, "id")

	err := h.service.DeleteJobSkill(ctx, jobSkillID)
	if err != nil {
		http.Error(w, "Job skill not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
