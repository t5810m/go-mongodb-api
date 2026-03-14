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

func TestLocationAvailabilityHandler_GetAllLocationAvailabilities_Success(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	items := []models.LocationAvailability{{ID: bson.NewObjectID(), Title: "Remote"}}
	mockSvc.On("GetAllLocationAvailabilities", mock.Anything, 1, 10, mock.Anything, "", "").Return(items, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/locationavailabilities", nil)
	w := httptest.NewRecorder()

	h.GetAllLocationAvailabilities(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLocationAvailabilityHandler_GetAllLocationAvailabilities_Error(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	mockSvc.On("GetAllLocationAvailabilities", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.LocationAvailability{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/locationavailabilities", nil)
	w := httptest.NewRecorder()

	h.GetAllLocationAvailabilities(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLocationAvailabilityHandler_GetLocationAvailabilityByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	id := bson.NewObjectID()
	item := &models.LocationAvailability{ID: id, Title: "Remote"}
	mockSvc.On("GetLocationAvailabilityByID", mock.Anything, id.Hex()).Return(item, nil)

	r := httptest.NewRequest(http.MethodGet, "/locationavailabilities/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetLocationAvailabilityByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLocationAvailabilityHandler_GetLocationAvailabilityByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	mockSvc.On("GetLocationAvailabilityByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/locationavailabilities/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetLocationAvailabilityByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLocationAvailabilityHandler_CreateLocationAvailability_Success(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	mockSvc.On("CreateLocationAvailability", mock.Anything, mock.AnythingOfType("*models.LocationAvailability")).Return(nil)

	body := `{"title":"On-site"}`
	r := httptest.NewRequest(http.MethodPost, "/locationavailabilities", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateLocationAvailability(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLocationAvailabilityHandler_CreateLocationAvailability_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/locationavailabilities", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreateLocationAvailability(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLocationAvailabilityHandler_CreateLocationAvailability_ValidationError(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	body := `{"title":"A"}`
	r := httptest.NewRequest(http.MethodPost, "/locationavailabilities", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateLocationAvailability(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLocationAvailabilityHandler_UpdateLocationAvailability_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/locationavailabilities/some-id", bytes.NewBufferString("not-json"))
	r = addChiURLParam(r, "id", "some-id")
	w := httptest.NewRecorder()

	h.UpdateLocationAvailability(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLocationAvailabilityHandler_DeleteLocationAvailability_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	mockSvc.On("DeleteLocationAvailability", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/locationavailabilities/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteLocationAvailability(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLocationAvailabilityHandler_UpdateLocationAvailability_Success(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	id := bson.NewObjectID()
	updated := &models.LocationAvailability{ID: id, Title: "Hybrid"}
	mockSvc.On("UpdateLocationAvailability", mock.Anything, id.Hex(), mock.AnythingOfType("*models.LocationAvailability")).Return(updated, nil)

	body := `{"title":"Hybrid"}`
	r := httptest.NewRequest(http.MethodPut, "/locationavailabilities/"+id.Hex(), bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", id.Hex())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateLocationAvailability(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLocationAvailabilityHandler_UpdateLocationAvailability_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	mockSvc.On("UpdateLocationAvailability", mock.Anything, "bad-id", mock.AnythingOfType("*models.LocationAvailability")).Return(nil, errors.New("not found"))

	body := `{"title":"Hybrid"}`
	r := httptest.NewRequest(http.MethodPut, "/locationavailabilities/bad-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.UpdateLocationAvailability(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLocationAvailabilityHandler_DeleteLocationAvailability_Success(t *testing.T) {
	mockSvc := new(mocks.MockLocationAvailabilityService)
	h := handlers.NewLocationAvailabilityHandler(mockSvc)

	mockSvc.On("DeleteLocationAvailability", mock.Anything, "item-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/locationavailabilities/item-id", nil)
	r = addChiURLParam(r, "id", "item-id")
	w := httptest.NewRecorder()

	h.DeleteLocationAvailability(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}
