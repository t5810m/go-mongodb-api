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

func TestJobSkillHandler_GetAllJobSkills_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	skills := []models.JobSkill{{ID: bson.NewObjectID(), ProficiencyLevelRequired: "expert"}}
	mockSvc.On("GetAllJobSkills", mock.Anything, 1, 10, mock.Anything, "", "").Return(skills, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/jobskills", nil)
	w := httptest.NewRecorder()

	h.GetAllJobSkills(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobSkillHandler_GetAllJobSkills_Error(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	mockSvc.On("GetAllJobSkills", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.JobSkill{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/jobskills", nil)
	w := httptest.NewRecorder()

	h.GetAllJobSkills(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJobSkillHandler_GetJobSkillByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	id := bson.NewObjectID()
	js := &models.JobSkill{ID: id, ProficiencyLevelRequired: "intermediate"}
	mockSvc.On("GetJobSkillByID", mock.Anything, id.Hex()).Return(js, nil)

	r := httptest.NewRequest(http.MethodGet, "/jobskills/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetJobSkillByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobSkillHandler_GetJobSkillByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	mockSvc.On("GetJobSkillByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/jobskills/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetJobSkillByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJobSkillHandler_GetJobSkillsByJobID(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	jobID := bson.NewObjectID()
	skills := []models.JobSkill{{ProficiencyLevelRequired: "advanced"}}
	mockSvc.On("GetJobSkillsByJobID", mock.Anything, jobID.Hex()).Return(skills, nil)

	r := httptest.NewRequest(http.MethodGet, "/jobs/"+jobID.Hex()+"/skills", nil)
	r = addChiURLParam(r, "jobId", jobID.Hex())
	w := httptest.NewRecorder()

	h.GetJobSkillsByJobID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobSkillHandler_CreateJobSkill_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	mockSvc.On("CreateJobSkill", mock.Anything, mock.AnythingOfType("*models.JobSkill")).Return(nil)

	jobID := bson.NewObjectID()
	skillID := bson.NewObjectID()
	body := `{"job_id":"` + jobID.Hex() + `","skill_id":"` + skillID.Hex() + `","proficiency_level_required":"advanced"}`
	r := httptest.NewRequest(http.MethodPost, "/jobskills", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateJobSkill(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobSkillHandler_CreateJobSkill_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/jobskills", bytes.NewBufferString("bad-json"))
	w := httptest.NewRecorder()

	h.CreateJobSkill(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJobSkillHandler_UpdateJobSkillProficiencyLevel_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	mockSvc.On("UpdateJobSkillProficiencyLevel", mock.Anything, "js-id", "expert").Return(nil)

	body := `{"proficiency_level_required":"expert"}`
	r := httptest.NewRequest(http.MethodPut, "/jobskills/js-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "js-id")
	w := httptest.NewRecorder()

	h.UpdateJobSkillProficiencyLevel(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobSkillHandler_UpdateJobSkillProficiencyLevel_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/jobskills/js-id", bytes.NewBufferString("bad"))
	r = addChiURLParam(r, "id", "js-id")
	w := httptest.NewRecorder()

	h.UpdateJobSkillProficiencyLevel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJobSkillHandler_DeleteJobSkill_Success(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	mockSvc.On("DeleteJobSkill", mock.Anything, "js-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/jobskills/js-id", nil)
	r = addChiURLParam(r, "id", "js-id")
	w := httptest.NewRecorder()

	h.DeleteJobSkill(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestJobSkillHandler_DeleteJobSkill_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockJobSkillService)
	h := handlers.NewJobSkillHandler(mockSvc)

	mockSvc.On("DeleteJobSkill", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/jobskills/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteJobSkill(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
