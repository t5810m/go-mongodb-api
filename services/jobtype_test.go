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

func TestJobTypeService_GetAllJobTypes(t *testing.T) {
	mockRepo := new(mocks.MockJobTypeRepository)
	svc := services.NewJobTypeService(mockRepo)

	expected := []models.JobType{{ID: bson.NewObjectID(), Title: "Full-time"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	jobTypes, total, err := svc.GetAllJobTypes(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, jobTypes, 1)
	mockRepo.AssertExpectations(t)
}

func TestJobTypeService_GetJobTypeByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockJobTypeRepository)
	svc := services.NewJobTypeService(mockRepo)

	id := bson.NewObjectID()
	expected := &models.JobType{ID: id, Title: "Full-time"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	jobType, err := svc.GetJobTypeByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "Full-time", jobType.Title)
	mockRepo.AssertExpectations(t)
}

func TestJobTypeService_GetJobTypeByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockJobTypeRepository)
	svc := services.NewJobTypeService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	jobType, err := svc.GetJobTypeByID(context.Background(), "bad-id")
	assert.Error(t, err)
	assert.Nil(t, jobType)
	mockRepo.AssertExpectations(t)
}

func TestJobTypeService_CreateJobType(t *testing.T) {
	mockRepo := new(mocks.MockJobTypeRepository)
	svc := services.NewJobTypeService(mockRepo)

	jobType := &models.JobType{Title: "Part-time"}
	mockRepo.On("Create", mock.Anything, jobType).Return(nil)

	err := svc.CreateJobType(context.Background(), jobType)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestJobTypeService_DeleteJobType(t *testing.T) {
	mockRepo := new(mocks.MockJobTypeRepository)
	svc := services.NewJobTypeService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "jobtype-id").Return(nil)

	err := svc.DeleteJobType(context.Background(), "jobtype-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestJobTypeService_UpdateJobType_Success(t *testing.T) {
	mockRepo := new(mocks.MockJobTypeRepository)
	svc := services.NewJobTypeService(mockRepo)

	id := bson.NewObjectID()
	input := &models.JobType{Title: "Contract"}
	expected := &models.JobType{ID: id, Title: "Contract"}
	mockRepo.On("Update", mock.Anything, id.Hex(), input).Return(expected, nil)

	result, err := svc.UpdateJobType(context.Background(), id.Hex(), input)
	assert.NoError(t, err)
	assert.Equal(t, "Contract", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestJobTypeService_UpdateJobType_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockJobTypeRepository)
	svc := services.NewJobTypeService(mockRepo)

	input := &models.JobType{Title: "Contract"}
	mockRepo.On("Update", mock.Anything, "bad-id", input).Return(nil, errors.New("not found"))

	result, err := svc.UpdateJobType(context.Background(), "bad-id", input)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
