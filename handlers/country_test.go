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

func TestCountryHandler_GetAllCountries_Success(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	countries := []models.Country{{ID: bson.NewObjectID(), Name: "Serbia"}}
	mockSvc.On("GetAllCountries", mock.Anything, 1, 10, mock.Anything, "", "").Return(countries, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/countries", nil)
	w := httptest.NewRecorder()

	h.GetAllCountries(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCountryHandler_GetAllCountries_Error(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	mockSvc.On("GetAllCountries", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.Country{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/countries", nil)
	w := httptest.NewRecorder()

	h.GetAllCountries(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCountryHandler_GetCountryByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	id := bson.NewObjectID()
	country := &models.Country{ID: id, Name: "Serbia"}
	mockSvc.On("GetCountryByID", mock.Anything, id.Hex()).Return(country, nil)

	r := httptest.NewRequest(http.MethodGet, "/countries/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetCountryByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCountryHandler_GetCountryByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	mockSvc.On("GetCountryByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/countries/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetCountryByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCountryHandler_CreateCountry_Success(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	mockSvc.On("CreateCountry", mock.Anything, mock.AnythingOfType("*models.Country")).Return(nil)

	body := `{"name":"Germany"}`
	r := httptest.NewRequest(http.MethodPost, "/countries", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateCountry(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCountryHandler_CreateCountry_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/countries", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreateCountry(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCountryHandler_CreateCountry_ValidationError(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	body := `{"name":"A"}`
	r := httptest.NewRequest(http.MethodPost, "/countries", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateCountry(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCountryHandler_UpdateCountry_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/countries/some-id", bytes.NewBufferString("not-json"))
	r = addChiURLParam(r, "id", "some-id")
	w := httptest.NewRecorder()

	h.UpdateCountry(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCountryHandler_DeleteCountry_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	mockSvc.On("DeleteCountry", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/countries/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteCountry(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCountryHandler_UpdateCountry_Success(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	id := bson.NewObjectID()
	updated := &models.Country{ID: id, Name: "France"}
	mockSvc.On("UpdateCountry", mock.Anything, id.Hex(), mock.AnythingOfType("*models.Country")).Return(updated, nil)

	body := `{"name":"France"}`
	r := httptest.NewRequest(http.MethodPut, "/countries/"+id.Hex(), bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", id.Hex())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateCountry(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCountryHandler_UpdateCountry_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	mockSvc.On("UpdateCountry", mock.Anything, "bad-id", mock.AnythingOfType("*models.Country")).Return(nil, errors.New("not found"))

	body := `{"name":"France"}`
	r := httptest.NewRequest(http.MethodPut, "/countries/bad-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.UpdateCountry(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCountryHandler_DeleteCountry_Success(t *testing.T) {
	mockSvc := new(mocks.MockCountryService)
	h := handlers.NewCountryHandler(mockSvc)

	mockSvc.On("DeleteCountry", mock.Anything, "country-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/countries/country-id", nil)
	r = addChiURLParam(r, "id", "country-id")
	w := httptest.NewRecorder()

	h.DeleteCountry(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}
