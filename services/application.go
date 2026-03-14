package services

import (
	"context"
	"fmt"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type ApplicationService struct {
	repo     interfaces.ApplicationRepository
	jobRepo  interfaces.JobRepository
	userRepo interfaces.UserRepository
}

// NewApplicationService creates a new application service
func NewApplicationService(repo interfaces.ApplicationRepository, jobRepo interfaces.JobRepository, userRepo interfaces.UserRepository) *ApplicationService {
	return &ApplicationService{
		repo:     repo,
		jobRepo:  jobRepo,
		userRepo: userRepo,
	}
}

// GetAllApplications retrieves all applications with pagination and optional filtering
func (s *ApplicationService) GetAllApplications(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Application, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

// GetApplicationByID retrieves an application by ID
func (s *ApplicationService) GetApplicationByID(ctx context.Context, id string) (*models.Application, error) {
	return s.repo.GetByID(ctx, id)
}

// GetApplicationsByJobID retrieves all applications for a specific job
func (s *ApplicationService) GetApplicationsByJobID(ctx context.Context, jobID string) ([]models.Application, error) {
	return s.repo.GetByJobID(ctx, jobID)
}

// GetApplicationsByUserID retrieves all applications from a specific user
func (s *ApplicationService) GetApplicationsByUserID(ctx context.Context, userID string) ([]models.Application, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// CreateApplication creates a new application
func (s *ApplicationService) CreateApplication(ctx context.Context, application *models.Application) error {
	if _, err := s.jobRepo.GetByID(ctx, application.JobID.Hex()); err != nil {
		return fmt.Errorf("job not found")
	}

	if _, err := s.userRepo.GetByID(ctx, application.UserID.Hex()); err != nil {
		return fmt.Errorf("user not found")
	}

	return s.repo.Create(ctx, application)
}

// UpdateApplicationStatus updates the status of an application
func (s *ApplicationService) UpdateApplicationStatus(ctx context.Context, id string, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

// DeleteApplication deletes an application by ID
func (s *ApplicationService) DeleteApplication(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
