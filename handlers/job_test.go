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

func TestJobHandler_GetAllJobs_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobService)
	h := handlers.NewJobHandler(mockSvc)

	jobs := []models.Job{{ID: bson.NewObjectID(), Title: "Go Developer"}}
	mockSvc.On("GetAllJobs", mock.Anything, 1, 10, mock.Anything, "", "").Return(jobs, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	w := httptest.NewRecorder()

	h.GetAllJobs(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobHandler_GetAllJobs_Error(t *testing.T) {
	mockSvc := new(mocks.MockJobService)
	h := handlers.NewJobHandler(mockSvc)

	mockSvc.On("GetAllJobs", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.Job{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	w := httptest.NewRecorder()

	h.GetAllJobs(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJobHandler_GetJobByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobService)
	h := handlers.NewJobHandler(mockSvc)

	id := bson.NewObjectID()
	job := &models.Job{ID: id, Title: "Go Developer"}
	mockSvc.On("GetJobByID", mock.Anything, id.Hex()).Return(job, nil)

	r := httptest.NewRequest(http.MethodGet, "/jobs/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetJobByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobHandler_GetJobByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockJobService)
	h := handlers.NewJobHandler(mockSvc)

	mockSvc.On("GetJobByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/jobs/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetJobByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJobHandler_GetJobsByUser(t *testing.T) {
	mockSvc := new(mocks.MockJobService)
	h := handlers.NewJobHandler(mockSvc)

	userID := bson.NewObjectID()
	jobs := []models.Job{{Title: "SWE"}}
	mockSvc.On("GetJobsByUser", mock.Anything, userID.Hex()).Return(jobs, nil)

	r := httptest.NewRequest(http.MethodGet, "/users/"+userID.Hex()+"/jobs", nil)
	r = addChiURLParam(r, "userId", userID.Hex())
	w := httptest.NewRecorder()

	h.GetJobsByUser(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobHandler_CreateJob_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobService)
	h := handlers.NewJobHandler(mockSvc)

	mockSvc.On("CreateJob", mock.Anything, mock.AnythingOfType("*models.Job")).Return(nil)

	userID := bson.NewObjectID()
	categoryID := bson.NewObjectID()
	body := `{
		"title":"Go Developer",
		"description":"We need a Go developer with at least 3 years of experience",
		"user_id":"` + userID.Hex() + `",
		"category_id":"` + categoryID.Hex() + `",
		"location":"New York",
		"job_type":"full-time",
		"salary_min":80000,
		"salary_max":120000,
		"status":"active"
	}`
	r := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateJob(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobHandler_CreateJob_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockJobService)
	h := handlers.NewJobHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreateJob(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJobHandler_CreateJob_SalaryValidation(t *testing.T) {
	mockSvc := new(mocks.MockJobService)
	h := handlers.NewJobHandler(mockSvc)

	recruiterID := bson.NewObjectID()
	companyID := bson.NewObjectID()
	categoryID := bson.NewObjectID()
	body := `{
		"title":"Go Developer",
		"description":"We need a Go developer with at least 3 years of experience",
		"recruiter_id":"` + recruiterID.Hex() + `",
		"company_id":"` + companyID.Hex() + `",
		"category_id":"` + categoryID.Hex() + `",
		"location":"New York",
		"job_type":"full-time",
		"salary_min":120000,
		"salary_max":80000,
		"status":"active"
	}`
	r := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateJob(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJobHandler_DeleteJob_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobService)
	h := handlers.NewJobHandler(mockSvc)

	mockSvc.On("DeleteJob", mock.Anything, "job-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/jobs/job-id", nil)
	r = addChiURLParam(r, "id", "job-id")
	w := httptest.NewRecorder()

	h.DeleteJob(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobHandler_DeleteJob_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockJobService)
	h := handlers.NewJobHandler(mockSvc)

	mockSvc.On("DeleteJob", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/jobs/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteJob(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
