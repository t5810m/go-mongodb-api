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

func TestJobCategoryService_GetAllJobCategories(t *testing.T) {
	mockRepo := new(mocks.MockJobCategoryRepository)
	svc := services.NewJobCategoryService(mockRepo)

	expected := []models.JobCategory{{ID: bson.NewObjectID(), Name: "Engineering"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	cats, total, err := svc.GetAllJobCategories(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, cats, 1)
	mockRepo.AssertExpectations(t)
}

func TestJobCategoryService_GetJobCategoryByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockJobCategoryRepository)
	svc := services.NewJobCategoryService(mockRepo)

	id := bson.NewObjectID()
	expected := &models.JobCategory{ID: id, Name: "Engineering"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	cat, err := svc.GetJobCategoryByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "Engineering", cat.Name)
	mockRepo.AssertExpectations(t)
}

func TestJobCategoryService_GetJobCategoryByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockJobCategoryRepository)
	svc := services.NewJobCategoryService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	cat, err := svc.GetJobCategoryByID(context.Background(), "bad-id")
	assert.Error(t, err)
	assert.Nil(t, cat)
	mockRepo.AssertExpectations(t)
}

func TestJobCategoryService_CreateJobCategory(t *testing.T) {
	mockRepo := new(mocks.MockJobCategoryRepository)
	svc := services.NewJobCategoryService(mockRepo)

	cat := &models.JobCategory{Name: "Design"}
	mockRepo.On("Create", mock.Anything, cat).Return(nil)

	err := svc.CreateJobCategory(context.Background(), cat)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestJobCategoryService_DeleteJobCategory(t *testing.T) {
	mockRepo := new(mocks.MockJobCategoryRepository)
	svc := services.NewJobCategoryService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "cat-id").Return(nil)

	err := svc.DeleteJobCategory(context.Background(), "cat-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestJobCategoryService_UpdateJobCategory_Success(t *testing.T) {
	mockRepo := new(mocks.MockJobCategoryRepository)
	svc := services.NewJobCategoryService(mockRepo)

	id := bson.NewObjectID()
	input := &models.JobCategory{Name: "Engineering Updated", Description: "Updated description"}
	expected := &models.JobCategory{ID: id, Name: "Engineering Updated"}
	mockRepo.On("Update", mock.Anything, id.Hex(), input).Return(expected, nil)

	result, err := svc.UpdateJobCategory(context.Background(), id.Hex(), input)
	assert.NoError(t, err)
	assert.Equal(t, "Engineering Updated", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestJobCategoryService_UpdateJobCategory_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockJobCategoryRepository)
	svc := services.NewJobCategoryService(mockRepo)

	input := &models.JobCategory{Name: "Engineering"}
	mockRepo.On("Update", mock.Anything, "bad-id", input).Return(nil, errors.New("not found"))

	result, err := svc.UpdateJobCategory(context.Background(), "bad-id", input)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
