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

type JobTypeHandler struct {
	service interfaces.JobTypeService
}

func NewJobTypeHandler(service interfaces.JobTypeService) *JobTypeHandler {
	return &JobTypeHandler{service: service}
}

func (h *JobTypeHandler) GetAllJobTypes(w http.ResponseWriter, r *http.Request) {
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
		"title": r.URL.Query().Get("title"),
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	jobTypes, total, err := h.service.GetAllJobTypes(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve job types", http.StatusInternalServerError)
		return
	}

	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       jobTypes,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *JobTypeHandler) GetJobTypeByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobTypeID := chi.URLParam(r, "id")

	jobType, err := h.service.GetJobTypeByID(ctx, jobTypeID)
	if err != nil {
		http.Error(w, "Job type not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobType); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *JobTypeHandler) CreateJobType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var jobType models.JobType
	err := json.NewDecoder(r.Body).Decode(&jobType)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validationErrors := helpers.ValidateStruct(jobType)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	jobType.CreatedTime = time.Now()
	jobType.UpdatedTime = time.Now()
	jobType.CreatedBy = "system"
	jobType.UpdatedBy = "system"

	err = h.service.CreateJobType(ctx, &jobType)
	if err != nil {
		http.Error(w, "Failed to create job type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(jobType); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

func (h *JobTypeHandler) UpdateJobType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobTypeID := chi.URLParam(r, "id")

	var jobType models.JobType
	if err := json.NewDecoder(r.Body).Decode(&jobType); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jobType.UpdatedTime = time.Now()
	jobType.UpdatedBy = "system"

	updated, err := h.service.UpdateJobType(ctx, jobTypeID, &jobType)
	if err != nil {
		http.Error(w, "Job type not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *JobTypeHandler) DeleteJobType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobTypeID := chi.URLParam(r, "id")

	err := h.service.DeleteJobType(ctx, jobTypeID)
	if err != nil {
		http.Error(w, "Job type not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
