package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-mongodb-api/handlers"
	"go-mongodb-api/mocks"
	"go-mongodb-api/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func addChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestUserHandler_GetAllUsers_Success(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	users := []models.User{{ID: bson.NewObjectID(), FirstName: "Alice", Email: "alice@example.com"}}
	mockSvc.On("GetAllUsers", mock.Anything, 1, 10, mock.Anything, "", "").Return(users, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	mockSvc.AssertExpectations(t)
}

func TestUserHandler_GetAllUsers_ServiceError(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	mockSvc.On("GetAllUsers", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.User{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUserHandler_GetUserByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	id := bson.NewObjectID()
	user := &models.User{ID: id, Email: "alice@example.com"}
	mockSvc.On("GetUserByID", mock.Anything, id.Hex()).Return(user, nil)

	r := httptest.NewRequest(http.MethodGet, "/users/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetUserByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUserHandler_GetUserByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	mockSvc.On("GetUserByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/users/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetUserByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	mockSvc.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	body := `{"first_name":"Alice","last_name":"Smith","email":"alice@example.com","password":"password123","role":"candidate"}`
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateUser(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUserHandler_CreateUser_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreateUser(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_CreateUser_ValidationError(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	// Missing required fields
	body := `{"first_name":"A"}`
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateUser(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_CreateUser_ServiceError(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	mockSvc.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("db error"))

	body := `{"first_name":"Alice","last_name":"Smith","email":"alice@example.com","password":"password123","role":"candidate"}`
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateUser(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUserHandler_DeleteUser_Success(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	mockSvc.On("DeleteUser", mock.Anything, "some-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/users/some-id", nil)
	r = addChiURLParam(r, "id", "some-id")
	w := httptest.NewRecorder()

	h.DeleteUser(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUserHandler_DeleteUser_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	mockSvc.On("DeleteUser", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/users/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteUser(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	id := bson.NewObjectID()
	updated := &models.User{ID: id, FirstName: "Alice", LastName: "Updated"}
	mockSvc.On("UpdateUser", mock.Anything, id.Hex(), mock.AnythingOfType("*models.User")).Return(updated, nil)

	body := `{"first_name":"Alice","last_name":"Updated","active":true}`
	r := httptest.NewRequest(http.MethodPut, "/users/"+id.Hex(), bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", id.Hex())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateUser(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	mockSvc.On("UpdateUser", mock.Anything, "bad-id", mock.AnythingOfType("*models.User")).Return(nil, errors.New("not found"))

	body := `{"first_name":"Alice"}`
	r := httptest.NewRequest(http.MethodPut, "/users/bad-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.UpdateUser(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_UpdateUser_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/users/some-id", bytes.NewBufferString("not-json"))
	r = addChiURLParam(r, "id", "some-id")
	w := httptest.NewRecorder()

	h.UpdateUser(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetAllUsers_WithPagination(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	mockSvc.On("GetAllUsers", mock.Anything, 2, 5, mock.Anything, "email", "asc").Return([]models.User{}, int64(0), nil)

	r := httptest.NewRequest(http.MethodGet, "/users?page=2&limit=5&sort=email&order=asc", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUserHandler_GetUserByID_ResponseBody(t *testing.T) {
	mockSvc := new(mocks.MockUserService)
	h := handlers.NewUserHandler(mockSvc)

	id := bson.NewObjectID()
	user := &models.User{ID: id, FirstName: "Alice", Email: "alice@example.com"}
	mockSvc.On("GetUserByID", mock.Anything, id.Hex()).Return(user, nil)

	r := httptest.NewRequest(http.MethodGet, "/users/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetUserByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.UserResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "alice@example.com", resp.Email)
}
