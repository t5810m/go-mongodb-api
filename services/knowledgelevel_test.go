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

func TestKnowledgeLevelService_GetAllKnowledgeLevels(t *testing.T) {
	mockRepo := new(mocks.MockKnowledgeLevelRepository)
	svc := services.NewKnowledgeLevelService(mockRepo)

	expected := []models.KnowledgeLevel{{ID: bson.NewObjectID(), Title: "Beginner"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	levels, total, err := svc.GetAllKnowledgeLevels(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, levels, 1)
	mockRepo.AssertExpectations(t)
}

func TestKnowledgeLevelService_GetKnowledgeLevelByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockKnowledgeLevelRepository)
	svc := services.NewKnowledgeLevelService(mockRepo)

	id := bson.NewObjectID()
	expected := &models.KnowledgeLevel{ID: id, Title: "Beginner"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	level, err := svc.GetKnowledgeLevelByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "Beginner", level.Title)
	mockRepo.AssertExpectations(t)
}

func TestKnowledgeLevelService_GetKnowledgeLevelByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockKnowledgeLevelRepository)
	svc := services.NewKnowledgeLevelService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	level, err := svc.GetKnowledgeLevelByID(context.Background(), "bad-id")
	assert.Error(t, err)
	assert.Nil(t, level)
	mockRepo.AssertExpectations(t)
}

func TestKnowledgeLevelService_CreateKnowledgeLevel(t *testing.T) {
	mockRepo := new(mocks.MockKnowledgeLevelRepository)
	svc := services.NewKnowledgeLevelService(mockRepo)

	level := &models.KnowledgeLevel{Title: "Advanced"}
	mockRepo.On("Create", mock.Anything, level).Return(nil)

	err := svc.CreateKnowledgeLevel(context.Background(), level)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestKnowledgeLevelService_DeleteKnowledgeLevel(t *testing.T) {
	mockRepo := new(mocks.MockKnowledgeLevelRepository)
	svc := services.NewKnowledgeLevelService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "level-id").Return(nil)

	err := svc.DeleteKnowledgeLevel(context.Background(), "level-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestKnowledgeLevelService_UpdateKnowledgeLevel_Success(t *testing.T) {
	mockRepo := new(mocks.MockKnowledgeLevelRepository)
	svc := services.NewKnowledgeLevelService(mockRepo)

	id := bson.NewObjectID()
	input := &models.KnowledgeLevel{Title: "Expert"}
	expected := &models.KnowledgeLevel{ID: id, Title: "Expert"}
	mockRepo.On("Update", mock.Anything, id.Hex(), input).Return(expected, nil)

	result, err := svc.UpdateKnowledgeLevel(context.Background(), id.Hex(), input)
	assert.NoError(t, err)
	assert.Equal(t, "Expert", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestKnowledgeLevelService_UpdateKnowledgeLevel_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockKnowledgeLevelRepository)
	svc := services.NewKnowledgeLevelService(mockRepo)

	input := &models.KnowledgeLevel{Title: "Expert"}
	mockRepo.On("Update", mock.Anything, "bad-id", input).Return(nil, errors.New("not found"))

	result, err := svc.UpdateKnowledgeLevel(context.Background(), "bad-id", input)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
