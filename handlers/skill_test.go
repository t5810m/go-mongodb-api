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

func TestSkillHandler_GetAllSkills_Success(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	skills := []models.Skill{{ID: bson.NewObjectID(), Name: "Go"}}
	mockSvc.On("GetAllSkills", mock.Anything, 1, 10, mock.Anything, "", "").Return(skills, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/skills", nil)
	w := httptest.NewRecorder()

	h.GetAllSkills(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestSkillHandler_GetAllSkills_Error(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	mockSvc.On("GetAllSkills", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.Skill{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/skills", nil)
	w := httptest.NewRecorder()

	h.GetAllSkills(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSkillHandler_GetSkillByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	id := bson.NewObjectID()
	skill := &models.Skill{ID: id, Name: "Python"}
	mockSvc.On("GetSkillByID", mock.Anything, id.Hex()).Return(skill, nil)

	r := httptest.NewRequest(http.MethodGet, "/skills/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetSkillByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestSkillHandler_GetSkillByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	mockSvc.On("GetSkillByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/skills/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetSkillByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSkillHandler_CreateSkill_Success(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	mockSvc.On("CreateSkill", mock.Anything, mock.AnythingOfType("*models.Skill")).Return(nil)

	body := `{"name":"Kubernetes","description":"Container orchestration tool"}`
	r := httptest.NewRequest(http.MethodPost, "/skills", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateSkill(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestSkillHandler_CreateSkill_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/skills", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreateSkill(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSkillHandler_CreateSkill_ValidationError(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	body := `{"name":"K"}`
	r := httptest.NewRequest(http.MethodPost, "/skills", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateSkill(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSkillHandler_DeleteSkill_Success(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	mockSvc.On("DeleteSkill", mock.Anything, "skill-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/skills/skill-id", nil)
	r = addChiURLParam(r, "id", "skill-id")
	w := httptest.NewRecorder()

	h.DeleteSkill(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestSkillHandler_DeleteSkill_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	mockSvc.On("DeleteSkill", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/skills/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteSkill(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSkillHandler_UpdateSkill_Success(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	id := bson.NewObjectID()
	updated := &models.Skill{ID: id, Name: "Go", Description: "Programming language"}
	mockSvc.On("UpdateSkill", mock.Anything, id.Hex(), mock.AnythingOfType("*models.Skill")).Return(updated, nil)

	body := `{"name":"Go","description":"Programming language"}`
	r := httptest.NewRequest(http.MethodPut, "/skills/"+id.Hex(), bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", id.Hex())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateSkill(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestSkillHandler_UpdateSkill_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	mockSvc.On("UpdateSkill", mock.Anything, "bad-id", mock.AnythingOfType("*models.Skill")).Return(nil, errors.New("not found"))

	body := `{"name":"Go"}`
	r := httptest.NewRequest(http.MethodPut, "/skills/bad-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.UpdateSkill(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSkillHandler_UpdateSkill_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockSkillService)
	h := handlers.NewSkillHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/skills/some-id", bytes.NewBufferString("not-json"))
	r = addChiURLParam(r, "id", "some-id")
	w := httptest.NewRecorder()

	h.UpdateSkill(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
