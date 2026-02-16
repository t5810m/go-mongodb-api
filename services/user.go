package services

import (
	"context"
	"go-mongodb-api/models"
	"go-mongodb-api/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
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
	return s.repo.Create(ctx, user)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
