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

type ArticleHandler struct {
	service interfaces.ArticleService
}

func NewArticleHandler(service interfaces.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) GetAllArticles(w http.ResponseWriter, r *http.Request) {
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
		"name": r.URL.Query().Get("name"),
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	articles, total, err := h.service.GetAllArticles(ctx, page, limit, filters, sort, order)
	if err != nil {
		http.Error(w, "Failed to retrieve articles", http.StatusInternalServerError)
		return
	}

	pagination := helpers.NewPagination(page, limit)
	pagination.SetTotal(total)

	response := helpers.PaginatedResponse{
		Data:       articles,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *ArticleHandler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articleID := chi.URLParam(r, "id")

	article, err := h.service.GetArticleByID(ctx, articleID)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(article); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var article models.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validationErrors := helpers.ValidateStruct(article)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(helpers.ErrorResponse{Errors: validationErrors}); err != nil {
			log.Printf("error encoding validation error response: %v", err)
		}
		return
	}

	article.CreatedTime = time.Now()
	article.UpdatedTime = time.Now()
	article.CreatedBy = "system"
	article.UpdatedBy = "system"

	err = h.service.CreateArticle(ctx, &article)
	if err != nil {
		http.Error(w, "Failed to create article", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(article); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articleID := chi.URLParam(r, "id")

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	article.UpdatedTime = time.Now()
	article.UpdatedBy = "system"

	updated, err := h.service.UpdateArticle(ctx, articleID, &article)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articleID := chi.URLParam(r, "id")

	err := h.service.DeleteArticle(ctx, articleID)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
