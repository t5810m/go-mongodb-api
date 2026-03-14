package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-mongodb-api/handlers"
	"go-mongodb-api/mocks"
	"go-mongodb-api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestJobCategoryHandler_GetAllJobCategories_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	cats := []models.JobCategory{{ID: bson.NewObjectID(), Name: "Engineering"}}
	mockSvc.On("GetAllJobCategories", mock.Anything, 1, 10, mock.Anything, "", "").Return(cats, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/jobcategories", nil)
	w := httptest.NewRecorder()

	h.GetAllJobCategories(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobCategoryHandler_GetAllJobCategories_Error(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	mockSvc.On("GetAllJobCategories", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.JobCategory{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/jobcategories", nil)
	w := httptest.NewRecorder()

	h.GetAllJobCategories(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJobCategoryHandler_GetJobCategoryByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	id := bson.NewObjectID()
	cat := &models.JobCategory{ID: id, Name: "Design"}
	mockSvc.On("GetJobCategoryByID", mock.Anything, id.Hex()).Return(cat, nil)

	r := httptest.NewRequest(http.MethodGet, "/jobcategories/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetJobCategoryByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobCategoryHandler_GetJobCategoryByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	mockSvc.On("GetJobCategoryByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/jobcategories/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetJobCategoryByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJobCategoryHandler_CreateJobCategory_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	mockSvc.On("CreateJobCategory", mock.Anything, mock.AnythingOfType("*models.JobCategory")).Return(nil)

	body := `{"name":"Engineering","description":"All engineering related jobs in tech sector"}`
	r := httptest.NewRequest(http.MethodPost, "/jobcategories", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateJobCategory(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobCategoryHandler_CreateJobCategory_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/jobcategories", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreateJobCategory(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJobCategoryHandler_CreateJobCategory_ValidationError(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	body := `{"name":"En"}`
	r := httptest.NewRequest(http.MethodPost, "/jobcategories", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateJobCategory(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJobCategoryHandler_DeleteJobCategory_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	mockSvc.On("DeleteJobCategory", mock.Anything, "cat-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/jobcategories/cat-id", nil)
	r = addChiURLParam(r, "id", "cat-id")
	w := httptest.NewRecorder()

	h.DeleteJobCategory(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobCategoryHandler_DeleteJobCategory_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	mockSvc.On("DeleteJobCategory", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/jobcategories/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteJobCategory(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJobCategoryHandler_UpdateJobCategory_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	id := bson.NewObjectID()
	updated := &models.JobCategory{ID: id, Name: "Engineering Updated"}
	mockSvc.On("UpdateJobCategory", mock.Anything, id.Hex(), mock.AnythingOfType("*models.JobCategory")).Return(updated, nil)

	body := `{"name":"Engineering Updated","description":"Updated description for engineering jobs"}`
	r := httptest.NewRequest(http.MethodPut, "/jobcategories/"+id.Hex(), bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", id.Hex())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateJobCategory(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobCategoryHandler_UpdateJobCategory_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	mockSvc.On("UpdateJobCategory", mock.Anything, "bad-id", mock.AnythingOfType("*models.JobCategory")).Return(nil, errors.New("not found"))

	body := `{"name":"Engineering"}`
	r := httptest.NewRequest(http.MethodPut, "/jobcategories/bad-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.UpdateJobCategory(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJobCategoryHandler_UpdateJobCategory_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockJobCategoryService)
	h := handlers.NewJobCategoryHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/jobcategories/some-id", bytes.NewBufferString("not-json"))
	r = addChiURLParam(r, "id", "some-id")
	w := httptest.NewRecorder()

	h.UpdateJobCategory(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
