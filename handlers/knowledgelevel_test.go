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

func TestKnowledgeLevelHandler_GetAllKnowledgeLevels_Success(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	levels := []models.KnowledgeLevel{{ID: bson.NewObjectID(), Title: "Beginner"}}
	mockSvc.On("GetAllKnowledgeLevels", mock.Anything, 1, 10, mock.Anything, "", "").Return(levels, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/knowledgelevels", nil)
	w := httptest.NewRecorder()

	h.GetAllKnowledgeLevels(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestKnowledgeLevelHandler_GetAllKnowledgeLevels_Error(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	mockSvc.On("GetAllKnowledgeLevels", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.KnowledgeLevel{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/knowledgelevels", nil)
	w := httptest.NewRecorder()

	h.GetAllKnowledgeLevels(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestKnowledgeLevelHandler_GetKnowledgeLevelByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	id := bson.NewObjectID()
	level := &models.KnowledgeLevel{ID: id, Title: "Beginner"}
	mockSvc.On("GetKnowledgeLevelByID", mock.Anything, id.Hex()).Return(level, nil)

	r := httptest.NewRequest(http.MethodGet, "/knowledgelevels/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetKnowledgeLevelByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestKnowledgeLevelHandler_GetKnowledgeLevelByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	mockSvc.On("GetKnowledgeLevelByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/knowledgelevels/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetKnowledgeLevelByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestKnowledgeLevelHandler_CreateKnowledgeLevel_Success(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	mockSvc.On("CreateKnowledgeLevel", mock.Anything, mock.AnythingOfType("*models.KnowledgeLevel")).Return(nil)

	body := `{"title":"Advanced"}`
	r := httptest.NewRequest(http.MethodPost, "/knowledgelevels", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateKnowledgeLevel(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestKnowledgeLevelHandler_CreateKnowledgeLevel_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/knowledgelevels", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreateKnowledgeLevel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestKnowledgeLevelHandler_CreateKnowledgeLevel_ValidationError(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	body := `{"title":"A"}`
	r := httptest.NewRequest(http.MethodPost, "/knowledgelevels", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateKnowledgeLevel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestKnowledgeLevelHandler_UpdateKnowledgeLevel_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/knowledgelevels/some-id", bytes.NewBufferString("not-json"))
	r = addChiURLParam(r, "id", "some-id")
	w := httptest.NewRecorder()

	h.UpdateKnowledgeLevel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestKnowledgeLevelHandler_DeleteKnowledgeLevel_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	mockSvc.On("DeleteKnowledgeLevel", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/knowledgelevels/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteKnowledgeLevel(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestKnowledgeLevelHandler_UpdateKnowledgeLevel_Success(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	id := bson.NewObjectID()
	updated := &models.KnowledgeLevel{ID: id, Title: "Expert"}
	mockSvc.On("UpdateKnowledgeLevel", mock.Anything, id.Hex(), mock.AnythingOfType("*models.KnowledgeLevel")).Return(updated, nil)

	body := `{"title":"Expert"}`
	r := httptest.NewRequest(http.MethodPut, "/knowledgelevels/"+id.Hex(), bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", id.Hex())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateKnowledgeLevel(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestKnowledgeLevelHandler_UpdateKnowledgeLevel_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	mockSvc.On("UpdateKnowledgeLevel", mock.Anything, "bad-id", mock.AnythingOfType("*models.KnowledgeLevel")).Return(nil, errors.New("not found"))

	body := `{"title":"Expert"}`
	r := httptest.NewRequest(http.MethodPut, "/knowledgelevels/bad-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.UpdateKnowledgeLevel(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestKnowledgeLevelHandler_DeleteKnowledgeLevel_Success(t *testing.T) {
	mockSvc := new(mocks.MockKnowledgeLevelService)
	h := handlers.NewKnowledgeLevelHandler(mockSvc)

	mockSvc.On("DeleteKnowledgeLevel", mock.Anything, "level-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/knowledgelevels/level-id", nil)
	r = addChiURLParam(r, "id", "level-id")
	w := httptest.NewRecorder()

	h.DeleteKnowledgeLevel(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}
