package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-mongodb-api/handlers"
	"go-mongodb-api/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Login_Success(t *testing.T) {
	mockSvc := new(mocks.MockAuthService)
	h := handlers.NewAuthHandler(mockSvc)

	mockSvc.On("Login", mock.Anything, "alice@example.com", "password123").Return("jwt-token", nil)

	body := `{"email":"alice@example.com","password":"password123"}`
	r := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "jwt-token", resp["token"])
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockSvc := new(mocks.MockAuthService)
	h := handlers.NewAuthHandler(mockSvc)

	mockSvc.On("Login", mock.Anything, "alice@example.com", "wrong").Return("", errors.New("invalid credentials"))

	body := `{"email":"alice@example.com","password":"wrong"}`
	r := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Login(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Login_MissingBody(t *testing.T) {
	mockSvc := new(mocks.MockAuthService)
	h := handlers.NewAuthHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.Login(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_MissingEmail(t *testing.T) {
	mockSvc := new(mocks.MockAuthService)
	h := handlers.NewAuthHandler(mockSvc)

	body := `{"password":"password123"}`
	r := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Login(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_MissingPassword(t *testing.T) {
	mockSvc := new(mocks.MockAuthService)
	h := handlers.NewAuthHandler(mockSvc)

	body := `{"email":"alice@example.com"}`
	r := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Login(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockSvc := new(mocks.MockAuthService)
	h := handlers.NewAuthHandler(mockSvc)

	mockSvc.On("Register", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	body := `{"first_name":"Alice","last_name":"Smith","email":"alice@example.com","password":"password123","role":"candidate"}`
	r := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Register_MissingBody(t *testing.T) {
	mockSvc := new(mocks.MockAuthService)
	h := handlers.NewAuthHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.Register(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_MissingRole(t *testing.T) {
	mockSvc := new(mocks.MockAuthService)
	h := handlers.NewAuthHandler(mockSvc)

	body := `{"email":"alice@example.com","password":"password123"}`
	r := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Register(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_ServiceError(t *testing.T) {
	mockSvc := new(mocks.MockAuthService)
	h := handlers.NewAuthHandler(mockSvc)

	mockSvc.On("Register", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("db error"))

	body := `{"first_name":"Alice","last_name":"Smith","email":"alice@example.com","password":"password123","role":"candidate"}`
	r := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}
