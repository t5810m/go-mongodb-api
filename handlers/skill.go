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

type SkillHandler struct {
	service *services.SkillService
}

// NewSkillHandler creates a new skill handler
func NewSkillHandler(service *services.SkillService) *SkillHandler {
	return &SkillHandler{
		service: service,
	}
}

// GetAllSkills handles GET /skills request with pagination support
func (h *SkillHandler) GetAllSkills(w http.ResponseWriter, r *http.Request) {
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
		"name": r.URL.Query().Get("name"),
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Get skills with pagination, filters, and sorting
	skills, total, err := h.service.GetAllSkills(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve skills", http.StatusInternalServerError)
		return
	}

	// Build paginated response
	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       skills,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetSkillByID handles GET /skills/{id} request
func (h *SkillHandler) GetSkillByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	skillID := chi.URLParam(r, "id")

	skill, err := h.service.GetSkillByID(ctx, skillID)
	if err != nil {
		http.Error(w, "Skill not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(skill); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateSkill handles POST /skills request
func (h *SkillHandler) CreateSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var skill models.Skill
	err := json.NewDecoder(r.Body).Decode(&skill)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	validationErrors := helpers.ValidateStruct(skill)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	skill.CreatedTime = time.Now()
	skill.UpdatedTime = time.Now()
	skill.CreatedBy = "system"
	skill.UpdatedBy = "system"

	err = h.service.CreateSkill(ctx, &skill)
	if err != nil {
		http.Error(w, "Failed to create skill", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(skill); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

// DeleteSkill handles DELETE /skills/{id} request
func (h *SkillHandler) DeleteSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	skillID := chi.URLParam(r, "id")

	err := h.service.DeleteSkill(ctx, skillID)
	if err != nil {
		http.Error(w, "Skill not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
