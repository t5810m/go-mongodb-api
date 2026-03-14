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
	"golang.org/x/crypto/bcrypt"
)

func makeHashedPassword(plain string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockUserSvc := new(mocks.MockUserService)
	svc := services.NewAuthService(mockUserSvc, "test-secret")

	user := &models.User{
		ID:       bson.NewObjectID(),
		Email:    "alice@example.com",
		Password: makeHashedPassword("password123"),
		Role:     "candidate",
	}
	mockUserSvc.On("GetUserByEmail", mock.Anything, "alice@example.com").Return(user, nil)

	token, err := svc.Login(context.Background(), "alice@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockUserSvc.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockUserSvc := new(mocks.MockUserService)
	svc := services.NewAuthService(mockUserSvc, "test-secret")

	mockUserSvc.On("GetUserByEmail", mock.Anything, "nope@example.com").Return(nil, errors.New("not found"))

	token, err := svc.Login(context.Background(), "nope@example.com", "password123")
	assert.Error(t, err)
	assert.Empty(t, token)
	mockUserSvc.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockUserSvc := new(mocks.MockUserService)
	svc := services.NewAuthService(mockUserSvc, "test-secret")

	user := &models.User{
		ID:       bson.NewObjectID(),
		Email:    "alice@example.com",
		Password: makeHashedPassword("correct-password"),
		Role:     "recruiter",
	}
	mockUserSvc.On("GetUserByEmail", mock.Anything, "alice@example.com").Return(user, nil)

	token, err := svc.Login(context.Background(), "alice@example.com", "wrong-password")
	assert.Error(t, err)
	assert.Empty(t, token)
	mockUserSvc.AssertExpectations(t)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockUserSvc := new(mocks.MockUserService)
	svc := services.NewAuthService(mockUserSvc, "test-secret")

	mockUserSvc.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	user := &models.User{
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice@example.com",
		Password:  "password123",
		Role:      "candidate",
	}

	err := svc.Register(context.Background(), user)
	assert.NoError(t, err)
	assert.True(t, user.Active)
	mockUserSvc.AssertExpectations(t)
}

func TestAuthService_Register_ServiceError(t *testing.T) {
	mockUserSvc := new(mocks.MockUserService)
	svc := services.NewAuthService(mockUserSvc, "test-secret")

	mockUserSvc.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("db error"))

	user := &models.User{Email: "alice@example.com", Password: "password123", Role: "candidate"}
	err := svc.Register(context.Background(), user)
	assert.Error(t, err)
	mockUserSvc.AssertExpectations(t)
}
