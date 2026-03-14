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

func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	expected := []models.User{{ID: bson.NewObjectID(), FirstName: "Alice"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	users, total, err := svc.GetAllUsers(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, users, 1)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	id := bson.NewObjectID()
	expected := &models.User{ID: id, Email: "alice@example.com"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	user, err := svc.GetUserByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, expected.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, errors.New("not found"))

	user, err := svc.GetUserByID(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	expected := &models.User{Email: "alice@example.com"}
	mockRepo.On("GetByEmail", mock.Anything, "alice@example.com").Return(expected, nil)

	user, err := svc.GetUserByEmail(context.Background(), "alice@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "alice@example.com", user.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	user := &models.User{
		FirstName: "Alice",
		Password:  "plaintext123",
	}
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	err := svc.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	// Password should have been hashed
	assert.NotEqual(t, "plaintext123", user.Password)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	user := &models.User{Password: "password123"}
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("db error"))

	err := svc.CreateUser(context.Background(), user)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "some-id").Return(nil)

	err := svc.DeleteUser(context.Background(), "some-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_Error(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "bad-id").Return(errors.New("not found"))

	err := svc.DeleteUser(context.Background(), "bad-id")
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	id := bson.NewObjectID()
	input := &models.User{FirstName: "Alice", LastName: "Updated", Active: true}
	expected := &models.User{ID: id, FirstName: "Alice", LastName: "Updated"}
	mockRepo.On("Update", mock.Anything, id.Hex(), input).Return(expected, nil)

	result, err := svc.UpdateUser(context.Background(), id.Hex(), input)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", result.FirstName)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	input := &models.User{FirstName: "Alice"}
	mockRepo.On("Update", mock.Anything, "bad-id", input).Return(nil, errors.New("not found"))

	result, err := svc.UpdateUser(context.Background(), "bad-id", input)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_WithPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo)

	id := bson.NewObjectID()
	input := &models.User{FirstName: "Alice", Password: "newpassword123"}
	expected := &models.User{ID: id, FirstName: "Alice"}
	mockRepo.On("Update", mock.Anything, id.Hex(), mock.AnythingOfType("*models.User")).Return(expected, nil)

	result, err := svc.UpdateUser(context.Background(), id.Hex(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEqual(t, "newpassword123", input.Password)
	mockRepo.AssertExpectations(t)
}
