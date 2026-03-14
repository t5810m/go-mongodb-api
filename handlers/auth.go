package handlers

import (
	"encoding/json"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
	service interfaces.AuthService
}

func NewAuthHandler(service interfaces.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(loginResponse{Token: token}); err != nil {
		log.Printf("error encoding login response: %v", err)
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" || user.Role == "" {
		http.Error(w, "email, password and role are required", http.StatusBadRequest)
		return
	}

	user.CreatedTime = time.Now()
	user.UpdatedTime = time.Now()

	if err := h.service.Register(r.Context(), &user); err != nil {
		http.Error(w, "registration failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user.ToResponse()); err != nil {
		log.Printf("error encoding register response: %v", err)
	}
}
