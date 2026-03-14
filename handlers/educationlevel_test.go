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

func TestEducationLevelHandler_GetAllEducationLevels_Success(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	levels := []models.EducationLevel{{ID: bson.NewObjectID(), Title: "Bachelor"}}
	mockSvc.On("GetAllEducationLevels", mock.Anything, 1, 10, mock.Anything, "", "").Return(levels, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/educationlevels", nil)
	w := httptest.NewRecorder()

	h.GetAllEducationLevels(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestEducationLevelHandler_GetAllEducationLevels_Error(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	mockSvc.On("GetAllEducationLevels", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.EducationLevel{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/educationlevels", nil)
	w := httptest.NewRecorder()

	h.GetAllEducationLevels(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEducationLevelHandler_GetEducationLevelByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	id := bson.NewObjectID()
	level := &models.EducationLevel{ID: id, Title: "Bachelor"}
	mockSvc.On("GetEducationLevelByID", mock.Anything, id.Hex()).Return(level, nil)

	r := httptest.NewRequest(http.MethodGet, "/educationlevels/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetEducationLevelByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestEducationLevelHandler_GetEducationLevelByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	mockSvc.On("GetEducationLevelByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/educationlevels/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetEducationLevelByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEducationLevelHandler_CreateEducationLevel_Success(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	mockSvc.On("CreateEducationLevel", mock.Anything, mock.AnythingOfType("*models.EducationLevel")).Return(nil)

	body := `{"title":"Master"}`
	r := httptest.NewRequest(http.MethodPost, "/educationlevels", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateEducationLevel(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestEducationLevelHandler_CreateEducationLevel_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/educationlevels", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreateEducationLevel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEducationLevelHandler_CreateEducationLevel_ValidationError(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	body := `{"title":"A"}`
	r := httptest.NewRequest(http.MethodPost, "/educationlevels", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateEducationLevel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEducationLevelHandler_UpdateEducationLevel_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/educationlevels/some-id", bytes.NewBufferString("not-json"))
	r = addChiURLParam(r, "id", "some-id")
	w := httptest.NewRecorder()

	h.UpdateEducationLevel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEducationLevelHandler_DeleteEducationLevel_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	mockSvc.On("DeleteEducationLevel", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/educationlevels/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteEducationLevel(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEducationLevelHandler_UpdateEducationLevel_Success(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	id := bson.NewObjectID()
	updated := &models.EducationLevel{ID: id, Title: "PhD"}
	mockSvc.On("UpdateEducationLevel", mock.Anything, id.Hex(), mock.AnythingOfType("*models.EducationLevel")).Return(updated, nil)

	body := `{"title":"PhD"}`
	r := httptest.NewRequest(http.MethodPut, "/educationlevels/"+id.Hex(), bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", id.Hex())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateEducationLevel(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestEducationLevelHandler_UpdateEducationLevel_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	mockSvc.On("UpdateEducationLevel", mock.Anything, "bad-id", mock.AnythingOfType("*models.EducationLevel")).Return(nil, errors.New("not found"))

	body := `{"title":"PhD"}`
	r := httptest.NewRequest(http.MethodPut, "/educationlevels/bad-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.UpdateEducationLevel(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEducationLevelHandler_DeleteEducationLevel_Success(t *testing.T) {
	mockSvc := new(mocks.MockEducationLevelService)
	h := handlers.NewEducationLevelHandler(mockSvc)

	mockSvc.On("DeleteEducationLevel", mock.Anything, "level-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/educationlevels/level-id", nil)
	r = addChiURLParam(r, "id", "level-id")
	w := httptest.NewRecorder()

	h.DeleteEducationLevel(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}
