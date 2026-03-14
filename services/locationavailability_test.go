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

func TestLocationAvailabilityService_GetAllLocationAvailabilities(t *testing.T) {
	mockRepo := new(mocks.MockLocationAvailabilityRepository)
	svc := services.NewLocationAvailabilityService(mockRepo)

	expected := []models.LocationAvailability{{ID: bson.NewObjectID(), Title: "Remote"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	items, total, err := svc.GetAllLocationAvailabilities(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, items, 1)
	mockRepo.AssertExpectations(t)
}

func TestLocationAvailabilityService_GetLocationAvailabilityByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockLocationAvailabilityRepository)
	svc := services.NewLocationAvailabilityService(mockRepo)

	id := bson.NewObjectID()
	expected := &models.LocationAvailability{ID: id, Title: "Remote"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	item, err := svc.GetLocationAvailabilityByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "Remote", item.Title)
	mockRepo.AssertExpectations(t)
}

func TestLocationAvailabilityService_GetLocationAvailabilityByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockLocationAvailabilityRepository)
	svc := services.NewLocationAvailabilityService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	item, err := svc.GetLocationAvailabilityByID(context.Background(), "bad-id")
	assert.Error(t, err)
	assert.Nil(t, item)
	mockRepo.AssertExpectations(t)
}

func TestLocationAvailabilityService_CreateLocationAvailability(t *testing.T) {
	mockRepo := new(mocks.MockLocationAvailabilityRepository)
	svc := services.NewLocationAvailabilityService(mockRepo)

	item := &models.LocationAvailability{Title: "On-site"}
	mockRepo.On("Create", mock.Anything, item).Return(nil)

	err := svc.CreateLocationAvailability(context.Background(), item)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLocationAvailabilityService_DeleteLocationAvailability(t *testing.T) {
	mockRepo := new(mocks.MockLocationAvailabilityRepository)
	svc := services.NewLocationAvailabilityService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "item-id").Return(nil)

	err := svc.DeleteLocationAvailability(context.Background(), "item-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLocationAvailabilityService_UpdateLocationAvailability_Success(t *testing.T) {
	mockRepo := new(mocks.MockLocationAvailabilityRepository)
	svc := services.NewLocationAvailabilityService(mockRepo)

	id := bson.NewObjectID()
	input := &models.LocationAvailability{Title: "Hybrid"}
	expected := &models.LocationAvailability{ID: id, Title: "Hybrid"}
	mockRepo.On("Update", mock.Anything, id.Hex(), input).Return(expected, nil)

	result, err := svc.UpdateLocationAvailability(context.Background(), id.Hex(), input)
	assert.NoError(t, err)
	assert.Equal(t, "Hybrid", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestLocationAvailabilityService_UpdateLocationAvailability_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockLocationAvailabilityRepository)
	svc := services.NewLocationAvailabilityService(mockRepo)

	input := &models.LocationAvailability{Title: "Hybrid"}
	mockRepo.On("Update", mock.Anything, "bad-id", input).Return(nil, errors.New("not found"))

	result, err := svc.UpdateLocationAvailability(context.Background(), "bad-id", input)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
