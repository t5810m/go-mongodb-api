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

func TestEducationLevelService_GetAllEducationLevels(t *testing.T) {
	mockRepo := new(mocks.MockEducationLevelRepository)
	svc := services.NewEducationLevelService(mockRepo)

	expected := []models.EducationLevel{{ID: bson.NewObjectID(), Title: "Bachelor"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	levels, total, err := svc.GetAllEducationLevels(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, levels, 1)
	mockRepo.AssertExpectations(t)
}

func TestEducationLevelService_GetEducationLevelByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockEducationLevelRepository)
	svc := services.NewEducationLevelService(mockRepo)

	id := bson.NewObjectID()
	expected := &models.EducationLevel{ID: id, Title: "Bachelor"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	level, err := svc.GetEducationLevelByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "Bachelor", level.Title)
	mockRepo.AssertExpectations(t)
}

func TestEducationLevelService_GetEducationLevelByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockEducationLevelRepository)
	svc := services.NewEducationLevelService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	level, err := svc.GetEducationLevelByID(context.Background(), "bad-id")
	assert.Error(t, err)
	assert.Nil(t, level)
	mockRepo.AssertExpectations(t)
}

func TestEducationLevelService_CreateEducationLevel(t *testing.T) {
	mockRepo := new(mocks.MockEducationLevelRepository)
	svc := services.NewEducationLevelService(mockRepo)

	level := &models.EducationLevel{Title: "Master"}
	mockRepo.On("Create", mock.Anything, level).Return(nil)

	err := svc.CreateEducationLevel(context.Background(), level)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEducationLevelService_DeleteEducationLevel(t *testing.T) {
	mockRepo := new(mocks.MockEducationLevelRepository)
	svc := services.NewEducationLevelService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "level-id").Return(nil)

	err := svc.DeleteEducationLevel(context.Background(), "level-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEducationLevelService_UpdateEducationLevel_Success(t *testing.T) {
	mockRepo := new(mocks.MockEducationLevelRepository)
	svc := services.NewEducationLevelService(mockRepo)

	id := bson.NewObjectID()
	input := &models.EducationLevel{Title: "PhD"}
	expected := &models.EducationLevel{ID: id, Title: "PhD"}
	mockRepo.On("Update", mock.Anything, id.Hex(), input).Return(expected, nil)

	result, err := svc.UpdateEducationLevel(context.Background(), id.Hex(), input)
	assert.NoError(t, err)
	assert.Equal(t, "PhD", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestEducationLevelService_UpdateEducationLevel_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockEducationLevelRepository)
	svc := services.NewEducationLevelService(mockRepo)

	input := &models.EducationLevel{Title: "PhD"}
	mockRepo.On("Update", mock.Anything, "bad-id", input).Return(nil, errors.New("not found"))

	result, err := svc.UpdateEducationLevel(context.Background(), "bad-id", input)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
