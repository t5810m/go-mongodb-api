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

func TestSkillService_GetAllSkills(t *testing.T) {
	mockRepo := new(mocks.MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	expected := []models.Skill{{ID: bson.NewObjectID(), Name: "Go"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	skills, total, err := svc.GetAllSkills(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, skills, 1)
	mockRepo.AssertExpectations(t)
}

func TestSkillService_GetSkillByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	id := bson.NewObjectID()
	expected := &models.Skill{ID: id, Name: "Go"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	skill, err := svc.GetSkillByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "Go", skill.Name)
	mockRepo.AssertExpectations(t)
}

func TestSkillService_GetSkillByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	skill, err := svc.GetSkillByID(context.Background(), "bad-id")
	assert.Error(t, err)
	assert.Nil(t, skill)
	mockRepo.AssertExpectations(t)
}

func TestSkillService_CreateSkill(t *testing.T) {
	mockRepo := new(mocks.MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	skill := &models.Skill{Name: "Python"}
	mockRepo.On("Create", mock.Anything, skill).Return(nil)

	err := svc.CreateSkill(context.Background(), skill)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSkillService_DeleteSkill(t *testing.T) {
	mockRepo := new(mocks.MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "skill-id").Return(nil)

	err := svc.DeleteSkill(context.Background(), "skill-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSkillService_UpdateSkill_Success(t *testing.T) {
	mockRepo := new(mocks.MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	id := bson.NewObjectID()
	input := &models.Skill{Name: "Go", Description: "Programming language"}
	expected := &models.Skill{ID: id, Name: "Go", Description: "Programming language"}
	mockRepo.On("Update", mock.Anything, id.Hex(), input).Return(expected, nil)

	result, err := svc.UpdateSkill(context.Background(), id.Hex(), input)
	assert.NoError(t, err)
	assert.Equal(t, "Go", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestSkillService_UpdateSkill_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	input := &models.Skill{Name: "Go"}
	mockRepo.On("Update", mock.Anything, "bad-id", input).Return(nil, errors.New("not found"))

	result, err := svc.UpdateSkill(context.Background(), "bad-id", input)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
