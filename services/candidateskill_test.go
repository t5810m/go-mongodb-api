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

func TestCandidateSkillService_GetAllCandidateSkills(t *testing.T) {
	mockRepo := new(mocks.MockCandidateSkillRepository)
	svc := services.NewCandidateSkillService(mockRepo, nil, nil)

	expected := []models.CandidateSkill{{ID: bson.NewObjectID(), ProficiencyLevel: "expert"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	skills, total, err := svc.GetAllCandidateSkills(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, skills, 1)
	mockRepo.AssertExpectations(t)
}

func TestCandidateSkillService_GetCandidateSkillByID(t *testing.T) {
	mockRepo := new(mocks.MockCandidateSkillRepository)
	svc := services.NewCandidateSkillService(mockRepo, nil, nil)

	id := bson.NewObjectID()
	expected := &models.CandidateSkill{ID: id, ProficiencyLevel: "intermediate"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	cs, err := svc.GetCandidateSkillByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "intermediate", cs.ProficiencyLevel)
	mockRepo.AssertExpectations(t)
}

func TestCandidateSkillService_GetByUserID(t *testing.T) {
	mockRepo := new(mocks.MockCandidateSkillRepository)
	svc := services.NewCandidateSkillService(mockRepo, nil, nil)

	userID := bson.NewObjectID()
	expected := []models.CandidateSkill{{ProficiencyLevel: "beginner"}}
	mockRepo.On("GetByUserID", mock.Anything, userID.Hex()).Return(expected, nil)

	skills, err := svc.GetCandidateSkillsByUserID(context.Background(), userID.Hex())
	assert.NoError(t, err)
	assert.Len(t, skills, 1)
	mockRepo.AssertExpectations(t)
}

func TestCandidateSkillService_CreateCandidateSkill_Success(t *testing.T) {
	mockRepo := new(mocks.MockCandidateSkillRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSkillRepo := new(mocks.MockSkillRepository)
	svc := services.NewCandidateSkillService(mockRepo, mockUserRepo, mockSkillRepo)

	userID := bson.NewObjectID()
	skillID := bson.NewObjectID()
	cs := &models.CandidateSkill{UserID: userID, SkillID: skillID, ProficiencyLevel: "advanced"}
	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(&models.User{}, nil)
	mockSkillRepo.On("GetByID", mock.Anything, skillID.Hex()).Return(&models.Skill{}, nil)
	mockRepo.On("Create", mock.Anything, cs).Return(nil)

	err := svc.CreateCandidateSkill(context.Background(), cs)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockSkillRepo.AssertExpectations(t)
}

func TestCandidateSkillService_CreateCandidateSkill_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockCandidateSkillRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSkillRepo := new(mocks.MockSkillRepository)
	svc := services.NewCandidateSkillService(mockRepo, mockUserRepo, mockSkillRepo)

	userID := bson.NewObjectID()
	cs := &models.CandidateSkill{UserID: userID}
	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(nil, errors.New("not found"))

	err := svc.CreateCandidateSkill(context.Background(), cs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestCandidateSkillService_CreateCandidateSkill_SkillNotFound(t *testing.T) {
	mockRepo := new(mocks.MockCandidateSkillRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSkillRepo := new(mocks.MockSkillRepository)
	svc := services.NewCandidateSkillService(mockRepo, mockUserRepo, mockSkillRepo)

	userID := bson.NewObjectID()
	skillID := bson.NewObjectID()
	cs := &models.CandidateSkill{UserID: userID, SkillID: skillID}
	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(&models.User{}, nil)
	mockSkillRepo.On("GetByID", mock.Anything, skillID.Hex()).Return(nil, errors.New("not found"))

	err := svc.CreateCandidateSkill(context.Background(), cs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "skill not found")
}

func TestCandidateSkillService_UpdateProficiencyLevel(t *testing.T) {
	mockRepo := new(mocks.MockCandidateSkillRepository)
	svc := services.NewCandidateSkillService(mockRepo, nil, nil)

	mockRepo.On("UpdateProficiencyLevel", mock.Anything, "cs-id", "expert").Return(nil)

	err := svc.UpdateCandidateSkillProficiencyLevel(context.Background(), "cs-id", "expert")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCandidateSkillService_DeleteCandidateSkill(t *testing.T) {
	mockRepo := new(mocks.MockCandidateSkillRepository)
	svc := services.NewCandidateSkillService(mockRepo, nil, nil)

	mockRepo.On("Delete", mock.Anything, "cs-id").Return(nil)

	err := svc.DeleteCandidateSkill(context.Background(), "cs-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
