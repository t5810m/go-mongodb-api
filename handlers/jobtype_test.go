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

func TestJobTypeHandler_GetAllJobTypes_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	jobTypes := []models.JobType{{ID: bson.NewObjectID(), Title: "Full-time"}}
	mockSvc.On("GetAllJobTypes", mock.Anything, 1, 10, mock.Anything, "", "").Return(jobTypes, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/jobtypes", nil)
	w := httptest.NewRecorder()

	h.GetAllJobTypes(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobTypeHandler_GetAllJobTypes_Error(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	mockSvc.On("GetAllJobTypes", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.JobType{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/jobtypes", nil)
	w := httptest.NewRecorder()

	h.GetAllJobTypes(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJobTypeHandler_GetJobTypeByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	id := bson.NewObjectID()
	jobType := &models.JobType{ID: id, Title: "Full-time"}
	mockSvc.On("GetJobTypeByID", mock.Anything, id.Hex()).Return(jobType, nil)

	r := httptest.NewRequest(http.MethodGet, "/jobtypes/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetJobTypeByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobTypeHandler_GetJobTypeByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	mockSvc.On("GetJobTypeByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/jobtypes/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetJobTypeByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJobTypeHandler_CreateJobType_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	mockSvc.On("CreateJobType", mock.Anything, mock.AnythingOfType("*models.JobType")).Return(nil)

	body := `{"title":"Part-time"}`
	r := httptest.NewRequest(http.MethodPost, "/jobtypes", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateJobType(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobTypeHandler_CreateJobType_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/jobtypes", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreateJobType(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJobTypeHandler_CreateJobType_ValidationError(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	body := `{"title":"A"}`
	r := httptest.NewRequest(http.MethodPost, "/jobtypes", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateJobType(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJobTypeHandler_UpdateJobType_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/jobtypes/some-id", bytes.NewBufferString("not-json"))
	r = addChiURLParam(r, "id", "some-id")
	w := httptest.NewRecorder()

	h.UpdateJobType(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJobTypeHandler_DeleteJobType_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	mockSvc.On("DeleteJobType", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/jobtypes/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteJobType(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJobTypeHandler_UpdateJobType_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	id := bson.NewObjectID()
	updated := &models.JobType{ID: id, Title: "Contract"}
	mockSvc.On("UpdateJobType", mock.Anything, id.Hex(), mock.AnythingOfType("*models.JobType")).Return(updated, nil)

	body := `{"title":"Contract"}`
	r := httptest.NewRequest(http.MethodPut, "/jobtypes/"+id.Hex(), bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", id.Hex())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateJobType(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobTypeHandler_UpdateJobType_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	mockSvc.On("UpdateJobType", mock.Anything, "bad-id", mock.AnythingOfType("*models.JobType")).Return(nil, errors.New("not found"))

	body := `{"title":"Contract"}`
	r := httptest.NewRequest(http.MethodPut, "/jobtypes/bad-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.UpdateJobType(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJobTypeHandler_DeleteJobType_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobTypeService)
	h := handlers.NewJobTypeHandler(mockSvc)

	mockSvc.On("DeleteJobType", mock.Anything, "jobtype-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/jobtypes/jobtype-id", nil)
	r = addChiURLParam(r, "id", "jobtype-id")
	w := httptest.NewRecorder()

	h.DeleteJobType(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}
