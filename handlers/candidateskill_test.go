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

func TestCandidateSkillHandler_GetAllCandidateSkills_Success(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	skills := []models.CandidateSkill{{ID: bson.NewObjectID(), ProficiencyLevel: "advanced"}}
	mockSvc.On("GetAllCandidateSkills", mock.Anything, 1, 10, mock.Anything, "", "").Return(skills, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/candidateskills", nil)
	w := httptest.NewRecorder()

	h.GetAllCandidateSkills(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCandidateSkillHandler_GetAllCandidateSkills_Error(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	mockSvc.On("GetAllCandidateSkills", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.CandidateSkill{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/candidateskills", nil)
	w := httptest.NewRecorder()

	h.GetAllCandidateSkills(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCandidateSkillHandler_GetCandidateSkillByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	id := bson.NewObjectID()
	cs := &models.CandidateSkill{ID: id, ProficiencyLevel: "expert"}
	mockSvc.On("GetCandidateSkillByID", mock.Anything, id.Hex()).Return(cs, nil)

	r := httptest.NewRequest(http.MethodGet, "/candidateskills/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetCandidateSkillByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCandidateSkillHandler_GetCandidateSkillByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	mockSvc.On("GetCandidateSkillByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/candidateskills/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetCandidateSkillByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCandidateSkillHandler_GetCandidateSkillsByUserID(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	userID := bson.NewObjectID()
	skills := []models.CandidateSkill{{ProficiencyLevel: "beginner"}}
	mockSvc.On("GetCandidateSkillsByUserID", mock.Anything, userID.Hex()).Return(skills, nil)

	r := httptest.NewRequest(http.MethodGet, "/users/"+userID.Hex()+"/skills", nil)
	r = addChiURLParam(r, "userId", userID.Hex())
	w := httptest.NewRecorder()

	h.GetCandidateSkillsByUserID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCandidateSkillHandler_CreateCandidateSkill_Success(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	mockSvc.On("CreateCandidateSkill", mock.Anything, mock.AnythingOfType("*models.CandidateSkill")).Return(nil)

	userID := bson.NewObjectID()
	skillID := bson.NewObjectID()
	body := `{"user_id":"` + userID.Hex() + `","skill_id":"` + skillID.Hex() + `","proficiency_level":"advanced"}`
	r := httptest.NewRequest(http.MethodPost, "/candidateskills", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateCandidateSkill(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCandidateSkillHandler_CreateCandidateSkill_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/candidateskills", bytes.NewBufferString("bad-json"))
	w := httptest.NewRecorder()

	h.CreateCandidateSkill(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCandidateSkillHandler_UpdateCandidateSkillProficiencyLevel_Success(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	mockSvc.On("UpdateCandidateSkillProficiencyLevel", mock.Anything, "cs-id", "expert").Return(nil)

	body := `{"proficiency_level":"expert"}`
	r := httptest.NewRequest(http.MethodPut, "/candidateskills/cs-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "cs-id")
	w := httptest.NewRecorder()

	h.UpdateCandidateSkillProficiencyLevel(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCandidateSkillHandler_UpdateCandidateSkillProficiencyLevel_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/candidateskills/cs-id", bytes.NewBufferString("bad"))
	r = addChiURLParam(r, "id", "cs-id")
	w := httptest.NewRecorder()

	h.UpdateCandidateSkillProficiencyLevel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCandidateSkillHandler_DeleteCandidateSkill_Success(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	mockSvc.On("DeleteCandidateSkill", mock.Anything, "cs-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/candidateskills/cs-id", nil)
	r = addChiURLParam(r, "id", "cs-id")
	w := httptest.NewRecorder()

	h.DeleteCandidateSkill(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCandidateSkillHandler_DeleteCandidateSkill_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockCandidateSkillService)
	h := handlers.NewCandidateSkillHandler(mockSvc)

	mockSvc.On("DeleteCandidateSkill", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/candidateskills/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteCandidateSkill(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
