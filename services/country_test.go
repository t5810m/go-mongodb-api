package services_test

import (
	"context"
	"errors"
	"testing"

	"go-mongodb-api/mocks"
	"go-mongodb-api/models"
	"go-mongodb-api/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCountryService_GetAllCountries(t *testing.T) {
	mockRepo := new(mocks.MockCountryRepository)
	svc := services.NewCountryService(mockRepo)

	expected := []models.Country{{ID: bson.NewObjectID(), Name: "Serbia"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	countries, total, err := svc.GetAllCountries(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, countries, 1)
	mockRepo.AssertExpectations(t)
}

func TestCountryService_GetCountryByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockCountryRepository)
	svc := services.NewCountryService(mockRepo)

	id := bson.NewObjectID()
	expected := &models.Country{ID: id, Name: "Serbia"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	country, err := svc.GetCountryByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "Serbia", country.Name)
	mockRepo.AssertExpectations(t)
}

func TestCountryService_GetCountryByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockCountryRepository)
	svc := services.NewCountryService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	country, err := svc.GetCountryByID(context.Background(), "bad-id")
	assert.Error(t, err)
	assert.Nil(t, country)
	mockRepo.AssertExpectations(t)
}

func TestCountryService_CreateCountry(t *testing.T) {
	mockRepo := new(mocks.MockCountryRepository)
	svc := services.NewCountryService(mockRepo)

	country := &models.Country{Name: "Germany"}
	mockRepo.On("Create", mock.Anything, country).Return(nil)

	err := svc.CreateCountry(context.Background(), country)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCountryService_DeleteCountry(t *testing.T) {
	mockRepo := new(mocks.MockCountryRepository)
	svc := services.NewCountryService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "country-id").Return(nil)

	err := svc.DeleteCountry(context.Background(), "country-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCountryService_UpdateCountry_Success(t *testing.T) {
	mockRepo := new(mocks.MockCountryRepository)
	svc := services.NewCountryService(mockRepo)

	id := bson.NewObjectID()
	input := &models.Country{Name: "France"}
	expected := &models.Country{ID: id, Name: "France"}
	mockRepo.On("Update", mock.Anything, id.Hex(), input).Return(expected, nil)

	result, err := svc.UpdateCountry(context.Background(), id.Hex(), input)
	assert.NoError(t, err)
	assert.Equal(t, "France", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCountryService_UpdateCountry_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockCountryRepository)
	svc := services.NewCountryService(mockRepo)

	input := &models.Country{Name: "France"}
	mockRepo.On("Update", mock.Anything, "bad-id", input).Return(nil, errors.New("not found"))

	result, err := svc.UpdateCountry(context.Background(), "bad-id", input)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
