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

func TestApplicationHandler_GetAllApplications_Success(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	apps := []models.Application{{ID: bson.NewObjectID(), Status: "applied"}}
	mockSvc.On("GetAllApplications", mock.Anything, 1, 10, mock.Anything, "", "").Return(apps, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/applications", nil)
	w := httptest.NewRecorder()

	h.GetAllApplications(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestApplicationHandler_GetAllApplications_Error(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	mockSvc.On("GetAllApplications", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.Application{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/applications", nil)
	w := httptest.NewRecorder()

	h.GetAllApplications(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestApplicationHandler_GetApplicationByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	id := bson.NewObjectID()
	app := &models.Application{ID: id, Status: "applied"}
	mockSvc.On("GetApplicationByID", mock.Anything, id.Hex()).Return(app, nil)

	r := httptest.NewRequest(http.MethodGet, "/applications/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetApplicationByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestApplicationHandler_GetApplicationByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	mockSvc.On("GetApplicationByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/applications/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetApplicationByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestApplicationHandler_GetApplicationsByJobID(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	jobID := bson.NewObjectID()
	apps := []models.Application{{Status: "applied"}}
	mockSvc.On("GetApplicationsByJobID", mock.Anything, jobID.Hex()).Return(apps, nil)

	r := httptest.NewRequest(http.MethodGet, "/jobs/"+jobID.Hex()+"/applications", nil)
	r = addChiURLParam(r, "jobId", jobID.Hex())
	w := httptest.NewRecorder()

	h.GetApplicationsByJobID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestApplicationHandler_GetApplicationsByUserID(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	userID := bson.NewObjectID()
	apps := []models.Application{{Status: "rejected"}}
	mockSvc.On("GetApplicationsByUserID", mock.Anything, userID.Hex()).Return(apps, nil)

	r := httptest.NewRequest(http.MethodGet, "/users/"+userID.Hex()+"/applications", nil)
	r = addChiURLParam(r, "userId", userID.Hex())
	w := httptest.NewRecorder()

	h.GetApplicationsByUserID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestApplicationHandler_CreateApplication_Success(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	mockSvc.On("CreateApplication", mock.Anything, mock.AnythingOfType("*models.Application")).Return(nil)

	jobID := bson.NewObjectID()
	userID := bson.NewObjectID()
	body := `{"job_id":"` + jobID.Hex() + `","user_id":"` + userID.Hex() + `","status":"applied"}`
	r := httptest.NewRequest(http.MethodPost, "/applications", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateApplication(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestApplicationHandler_CreateApplication_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/applications", bytes.NewBufferString("bad-json"))
	w := httptest.NewRecorder()

	h.CreateApplication(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestApplicationHandler_CreateApplication_ServiceError(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	mockSvc.On("CreateApplication", mock.Anything, mock.AnythingOfType("*models.Application")).Return(errors.New("job not found"))

	jobID := bson.NewObjectID()
	userID := bson.NewObjectID()
	body := `{"job_id":"` + jobID.Hex() + `","user_id":"` + userID.Hex() + `","status":"applied"}`
	r := httptest.NewRequest(http.MethodPost, "/applications", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateApplication(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestApplicationHandler_UpdateApplicationStatus_Success(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	mockSvc.On("UpdateApplicationStatus", mock.Anything, "app-id", "accepted").Return(nil)

	body := `{"status":"accepted"}`
	r := httptest.NewRequest(http.MethodPut, "/applications/app-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "app-id")
	w := httptest.NewRecorder()

	h.UpdateApplicationStatus(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestApplicationHandler_UpdateApplicationStatus_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/applications/app-id", bytes.NewBufferString("bad"))
	r = addChiURLParam(r, "id", "app-id")
	w := httptest.NewRecorder()

	h.UpdateApplicationStatus(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestApplicationHandler_DeleteApplication_Success(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	mockSvc.On("DeleteApplication", mock.Anything, "app-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/applications/app-id", nil)
	r = addChiURLParam(r, "id", "app-id")
	w := httptest.NewRecorder()

	h.DeleteApplication(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestApplicationHandler_DeleteApplication_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockApplicationService)
	h := handlers.NewApplicationHandler(mockSvc)

	mockSvc.On("DeleteApplication", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/applications/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteApplication(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
