package services

import (
	"context"
	"fmt"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo interfaces.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetAllUsers retrieves all users with pagination, filtering, and sorting
func (s *UserService) GetAllUsers(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.User, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashed)
	return s.repo.Create(ctx, user)
}

// UpdateUser updates a user's allowed fields
func (s *UserService) UpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error) {
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = string(hashed)
	}
	return s.repo.Update(ctx, id, user)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
