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

func TestJobSkillService_GetAllJobSkills(t *testing.T) {
	mockRepo := new(mocks.MockJobSkillRepository)
	svc := services.NewJobSkillService(mockRepo, nil, nil)

	expected := []models.JobSkill{{ID: bson.NewObjectID(), ProficiencyLevelRequired: "advanced"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	skills, total, err := svc.GetAllJobSkills(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, skills, 1)
	mockRepo.AssertExpectations(t)
}

func TestJobSkillService_GetJobSkillByID(t *testing.T) {
	mockRepo := new(mocks.MockJobSkillRepository)
	svc := services.NewJobSkillService(mockRepo, nil, nil)

	id := bson.NewObjectID()
	expected := &models.JobSkill{ID: id, ProficiencyLevelRequired: "beginner"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	js, err := svc.GetJobSkillByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "beginner", js.ProficiencyLevelRequired)
	mockRepo.AssertExpectations(t)
}

func TestJobSkillService_GetJobSkillsByJobID(t *testing.T) {
	mockRepo := new(mocks.MockJobSkillRepository)
	svc := services.NewJobSkillService(mockRepo, nil, nil)

	jobID := bson.NewObjectID()
	expected := []models.JobSkill{{ProficiencyLevelRequired: "expert"}}
	mockRepo.On("GetByJobID", mock.Anything, jobID.Hex()).Return(expected, nil)

	skills, err := svc.GetJobSkillsByJobID(context.Background(), jobID.Hex())
	assert.NoError(t, err)
	assert.Len(t, skills, 1)
	mockRepo.AssertExpectations(t)
}

func TestJobSkillService_CreateJobSkill_Success(t *testing.T) {
	mockRepo := new(mocks.MockJobSkillRepository)
	mockJobRepo := new(mocks.MockJobRepository)
	mockSkillRepo := new(mocks.MockSkillRepository)
	svc := services.NewJobSkillService(mockRepo, mockJobRepo, mockSkillRepo)

	jobID := bson.NewObjectID()
	skillID := bson.NewObjectID()
	js := &models.JobSkill{JobID: jobID, SkillID: skillID, ProficiencyLevelRequired: "intermediate"}
	mockJobRepo.On("GetByID", mock.Anything, jobID.Hex()).Return(&models.Job{}, nil)
	mockSkillRepo.On("GetByID", mock.Anything, skillID.Hex()).Return(&models.Skill{}, nil)
	mockRepo.On("Create", mock.Anything, js).Return(nil)

	err := svc.CreateJobSkill(context.Background(), js)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockJobRepo.AssertExpectations(t)
	mockSkillRepo.AssertExpectations(t)
}

func TestJobSkillService_CreateJobSkill_JobNotFound(t *testing.T) {
	mockRepo := new(mocks.MockJobSkillRepository)
	mockJobRepo := new(mocks.MockJobRepository)
	mockSkillRepo := new(mocks.MockSkillRepository)
	svc := services.NewJobSkillService(mockRepo, mockJobRepo, mockSkillRepo)

	jobID := bson.NewObjectID()
	js := &models.JobSkill{JobID: jobID}
	mockJobRepo.On("GetByID", mock.Anything, jobID.Hex()).Return(nil, errors.New("not found"))

	err := svc.CreateJobSkill(context.Background(), js)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "job not found")
}

func TestJobSkillService_CreateJobSkill_SkillNotFound(t *testing.T) {
	mockRepo := new(mocks.MockJobSkillRepository)
	mockJobRepo := new(mocks.MockJobRepository)
	mockSkillRepo := new(mocks.MockSkillRepository)
	svc := services.NewJobSkillService(mockRepo, mockJobRepo, mockSkillRepo)

	jobID := bson.NewObjectID()
	skillID := bson.NewObjectID()
	js := &models.JobSkill{JobID: jobID, SkillID: skillID}
	mockJobRepo.On("GetByID", mock.Anything, jobID.Hex()).Return(&models.Job{}, nil)
	mockSkillRepo.On("GetByID", mock.Anything, skillID.Hex()).Return(nil, errors.New("not found"))

	err := svc.CreateJobSkill(context.Background(), js)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "skill not found")
}

func TestJobSkillService_UpdateProficiencyLevel(t *testing.T) {
	mockRepo := new(mocks.MockJobSkillRepository)
	svc := services.NewJobSkillService(mockRepo, nil, nil)

	mockRepo.On("UpdateProficiencyLevel", mock.Anything, "js-id", "expert").Return(nil)

	err := svc.UpdateJobSkillProficiencyLevel(context.Background(), "js-id", "expert")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestJobSkillService_DeleteJobSkill(t *testing.T) {
	mockRepo := new(mocks.MockJobSkillRepository)
	svc := services.NewJobSkillService(mockRepo, nil, nil)

	mockRepo.On("Delete", mock.Anything, "js-id").Return(nil)

	err := svc.DeleteJobSkill(context.Background(), "js-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
