package services

import (
	"context"
	"fmt"
	"go-mongodb-api/models"
	"go-mongodb-api/repositories"
)

type ApplicationService struct {
	repo          *repositories.ApplicationRepository
	jobRepo       *repositories.JobRepository
	candidateRepo *repositories.CandidateRepository
}

// NewApplicationService creates a new application service
func NewApplicationService(repo *repositories.ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		repo: repo,
	}
}

// NewApplicationServiceWithDeps creates a new application service with dependencies
func NewApplicationServiceWithDeps(repo *repositories.ApplicationRepository, jobRepo *repositories.JobRepository, candidateRepo *repositories.CandidateRepository) *ApplicationService {
	return &ApplicationService{
		repo:          repo,
		jobRepo:       jobRepo,
		candidateRepo: candidateRepo,
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

// GetApplicationsByCandidateID retrieves all applications from a specific candidate
func (s *ApplicationService) GetApplicationsByCandidateID(ctx context.Context, candidateID string) ([]models.Application, error) {
	return s.repo.GetByCandidateID(ctx, candidateID)
}

// CreateApplication creates a new application
func (s *ApplicationService) CreateApplication(ctx context.Context, application *models.Application) error {
	// Validate job exists
	if s.jobRepo != nil {
		_, err := s.jobRepo.GetByID(ctx, application.JobID.Hex())
		if err != nil {
			return fmt.Errorf("job not found")
		}
	}

	// Validate candidate exists
	if s.candidateRepo != nil {
		_, err := s.candidateRepo.GetByID(ctx, application.CandidateID.Hex())
		if err != nil {
			return fmt.Errorf("candidate not found")
		}
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
