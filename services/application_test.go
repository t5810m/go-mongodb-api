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

func TestApplicationService_GetAllApplications(t *testing.T) {
	mockRepo := new(mocks.MockApplicationRepository)
	svc := services.NewApplicationService(mockRepo, nil, nil)

	expected := []models.Application{{ID: bson.NewObjectID(), Status: "applied"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	apps, total, err := svc.GetAllApplications(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, apps, 1)
	mockRepo.AssertExpectations(t)
}

func TestApplicationService_GetApplicationByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockApplicationRepository)
	svc := services.NewApplicationService(mockRepo, nil, nil)

	id := bson.NewObjectID()
	expected := &models.Application{ID: id, Status: "applied"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	app, err := svc.GetApplicationByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "applied", app.Status)
	mockRepo.AssertExpectations(t)
}

func TestApplicationService_GetApplicationsByJobID(t *testing.T) {
	mockRepo := new(mocks.MockApplicationRepository)
	svc := services.NewApplicationService(mockRepo, nil, nil)

	jobID := bson.NewObjectID()
	expected := []models.Application{{Status: "applied"}}
	mockRepo.On("GetByJobID", mock.Anything, jobID.Hex()).Return(expected, nil)

	apps, err := svc.GetApplicationsByJobID(context.Background(), jobID.Hex())
	assert.NoError(t, err)
	assert.Len(t, apps, 1)
	mockRepo.AssertExpectations(t)
}

func TestApplicationService_GetApplicationsByUserID(t *testing.T) {
	mockRepo := new(mocks.MockApplicationRepository)
	svc := services.NewApplicationService(mockRepo, nil, nil)

	userID := bson.NewObjectID()
	expected := []models.Application{{Status: "accepted"}}
	mockRepo.On("GetByUserID", mock.Anything, userID.Hex()).Return(expected, nil)

	apps, err := svc.GetApplicationsByUserID(context.Background(), userID.Hex())
	assert.NoError(t, err)
	assert.Len(t, apps, 1)
	mockRepo.AssertExpectations(t)
}

func TestApplicationService_CreateApplication_Success(t *testing.T) {
	mockRepo := new(mocks.MockApplicationRepository)
	mockJobRepo := new(mocks.MockJobRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	svc := services.NewApplicationService(mockRepo, mockJobRepo, mockUserRepo)

	jobID := bson.NewObjectID()
	userID := bson.NewObjectID()
	app := &models.Application{
		JobID:  jobID,
		UserID: userID,
		Status: "applied",
	}
	mockJobRepo.On("GetByID", mock.Anything, jobID.Hex()).Return(&models.Job{}, nil)
	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(&models.User{}, nil)
	mockRepo.On("Create", mock.Anything, app).Return(nil)

	err := svc.CreateApplication(context.Background(), app)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockJobRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestApplicationService_CreateApplication_JobNotFound(t *testing.T) {
	mockRepo := new(mocks.MockApplicationRepository)
	mockJobRepo := new(mocks.MockJobRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	svc := services.NewApplicationService(mockRepo, mockJobRepo, mockUserRepo)

	jobID := bson.NewObjectID()
	app := &models.Application{
		JobID:  jobID,
		UserID: bson.NewObjectID(),
		Status: "applied",
	}
	mockJobRepo.On("GetByID", mock.Anything, jobID.Hex()).Return(nil, errors.New("not found"))

	err := svc.CreateApplication(context.Background(), app)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "job not found")
	mockJobRepo.AssertExpectations(t)
}

func TestApplicationService_CreateApplication_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockApplicationRepository)
	mockJobRepo := new(mocks.MockJobRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	svc := services.NewApplicationService(mockRepo, mockJobRepo, mockUserRepo)

	jobID := bson.NewObjectID()
	userID := bson.NewObjectID()
	app := &models.Application{
		JobID:  jobID,
		UserID: userID,
		Status: "applied",
	}
	mockJobRepo.On("GetByID", mock.Anything, jobID.Hex()).Return(&models.Job{}, nil)
	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(nil, errors.New("not found"))

	err := svc.CreateApplication(context.Background(), app)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestApplicationService_UpdateApplicationStatus(t *testing.T) {
	mockRepo := new(mocks.MockApplicationRepository)
	svc := services.NewApplicationService(mockRepo, nil, nil)

	mockRepo.On("UpdateStatus", mock.Anything, "app-id", "accepted").Return(nil)

	err := svc.UpdateApplicationStatus(context.Background(), "app-id", "accepted")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestApplicationService_DeleteApplication(t *testing.T) {
	mockRepo := new(mocks.MockApplicationRepository)
	svc := services.NewApplicationService(mockRepo, nil, nil)

	mockRepo.On("Delete", mock.Anything, "app-id").Return(nil)

	err := svc.DeleteApplication(context.Background(), "app-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
