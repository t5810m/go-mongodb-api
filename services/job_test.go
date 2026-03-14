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

func TestJobService_GetAllJobs(t *testing.T) {
	mockRepo := new(mocks.MockJobRepository)
	svc := services.NewJobService(mockRepo, nil, nil)

	expected := []models.Job{{ID: bson.NewObjectID(), Title: "Dev"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	jobs, total, err := svc.GetAllJobs(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, jobs, 1)
	mockRepo.AssertExpectations(t)
}

func TestJobService_GetJobByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockJobRepository)
	svc := services.NewJobService(mockRepo, nil, nil)

	id := bson.NewObjectID()
	expected := &models.Job{ID: id, Title: "Dev"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	job, err := svc.GetJobByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "Dev", job.Title)
	mockRepo.AssertExpectations(t)
}

func TestJobService_GetJobsByUser(t *testing.T) {
	mockRepo := new(mocks.MockJobRepository)
	svc := services.NewJobService(mockRepo, nil, nil)

	userID := bson.NewObjectID()
	expected := []models.Job{{Title: "SWE"}}
	mockRepo.On("GetByUserID", mock.Anything, userID.Hex()).Return(expected, nil)

	jobs, err := svc.GetJobsByUser(context.Background(), userID.Hex())
	assert.NoError(t, err)
	assert.Len(t, jobs, 1)
	mockRepo.AssertExpectations(t)
}

func TestJobService_CreateJob_Success(t *testing.T) {
	mockRepo := new(mocks.MockJobRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCategoryRepo := new(mocks.MockJobCategoryRepository)
	svc := services.NewJobService(mockRepo, mockUserRepo, mockCategoryRepo)

	userID := bson.NewObjectID()
	categoryID := bson.NewObjectID()
	job := &models.Job{UserID: userID, CategoryID: categoryID}
	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(&models.User{}, nil)
	mockCategoryRepo.On("GetByID", mock.Anything, categoryID.Hex()).Return(&models.JobCategory{}, nil)
	mockRepo.On("Create", mock.Anything, job).Return(nil)

	err := svc.CreateJob(context.Background(), job)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestJobService_CreateJob_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockJobRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCategoryRepo := new(mocks.MockJobCategoryRepository)
	svc := services.NewJobService(mockRepo, mockUserRepo, mockCategoryRepo)

	userID := bson.NewObjectID()
	job := &models.Job{UserID: userID}
	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(nil, errors.New("not found"))

	err := svc.CreateJob(context.Background(), job)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestJobService_CreateJob_CategoryNotFound(t *testing.T) {
	mockRepo := new(mocks.MockJobRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCategoryRepo := new(mocks.MockJobCategoryRepository)
	svc := services.NewJobService(mockRepo, mockUserRepo, mockCategoryRepo)

	userID := bson.NewObjectID()
	categoryID := bson.NewObjectID()
	job := &models.Job{UserID: userID, CategoryID: categoryID}
	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(&models.User{}, nil)
	mockCategoryRepo.On("GetByID", mock.Anything, categoryID.Hex()).Return(nil, errors.New("not found"))

	err := svc.CreateJob(context.Background(), job)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "job category not found")
}

func TestJobService_DeleteJob(t *testing.T) {
	mockRepo := new(mocks.MockJobRepository)
	svc := services.NewJobService(mockRepo, nil, nil)

	mockRepo.On("Delete", mock.Anything, "job-id").Return(nil)

	err := svc.DeleteJob(context.Background(), "job-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
