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

type CandidateSkillHandler struct {
	service *services.CandidateSkillService
}

// NewCandidateSkillHandler creates a new candidate skill handler
func NewCandidateSkillHandler(service *services.CandidateSkillService) *CandidateSkillHandler {
	return &CandidateSkillHandler{
		service: service,
	}
}

// GetAllCandidateSkills handles GET /candidateskills request with pagination support
func (h *CandidateSkillHandler) GetAllCandidateSkills(w http.ResponseWriter, r *http.Request) {
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
		"candidate_id":      r.URL.Query().Get("candidate_id"),
		"skill_id":          r.URL.Query().Get("skill_id"),
		"proficiency_level": r.URL.Query().Get("proficiency_level"),
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Get candidate skills with pagination, filters, and sorting
	candidateSkills, total, err := h.service.GetAllCandidateSkills(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve candidate skills", http.StatusInternalServerError)
		return
	}

	// Build paginated response
	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       candidateSkills,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetCandidateSkillByID handles GET /candidateskills/{id} request
func (h *CandidateSkillHandler) GetCandidateSkillByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateSkillID := chi.URLParam(r, "id")

	candidateSkill, err := h.service.GetCandidateSkillByID(ctx, candidateSkillID)
	if err != nil {
		http.Error(w, "Candidate skill not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(candidateSkill); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetCandidateSkillsByCandidateID handles GET /candidates/{candidateId}/skills request
func (h *CandidateSkillHandler) GetCandidateSkillsByCandidateID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := chi.URLParam(r, "candidateId")

	candidateSkills, err := h.service.GetCandidateSkillsByCandidateID(ctx, candidateID)
	if err != nil {
		http.Error(w, "Failed to retrieve candidate skills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(candidateSkills); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateCandidateSkill handles POST /candidateskills request
func (h *CandidateSkillHandler) CreateCandidateSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var candidateSkill models.CandidateSkill
	err := json.NewDecoder(r.Body).Decode(&candidateSkill)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	validationErrors := helpers.ValidateStruct(candidateSkill)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	candidateSkill.CreatedTime = time.Now()
	candidateSkill.UpdatedTime = time.Now()
	candidateSkill.CreatedBy = "system"
	candidateSkill.UpdatedBy = "system"

	err = h.service.CreateCandidateSkill(ctx, &candidateSkill)
	if err != nil {
		http.Error(w, "Failed to create candidate skill", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(candidateSkill); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

// UpdateCandidateSkillProficiencyLevel handles PUT /candidateskills/{id} request
func (h *CandidateSkillHandler) UpdateCandidateSkillProficiencyLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateSkillID := chi.URLParam(r, "id")

	var request struct {
		ProficiencyLevel string `json:"proficiency_level"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateCandidateSkillProficiencyLevel(ctx, candidateSkillID, request.ProficiencyLevel)
	if err != nil {
		http.Error(w, "Failed to update candidate skill", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// DeleteCandidateSkill handles DELETE /candidateskills/{id} request
func (h *CandidateSkillHandler) DeleteCandidateSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateSkillID := chi.URLParam(r, "id")

	err := h.service.DeleteCandidateSkill(ctx, candidateSkillID)
	if err != nil {
		http.Error(w, "Candidate skill not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
