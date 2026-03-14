package services

import (
	"context"
	"fmt"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type JobService struct {
	repo         interfaces.JobRepository
	userRepo     interfaces.UserRepository
	categoryRepo interfaces.JobCategoryRepository
}

// NewJobService creates a new job service
func NewJobService(repo interfaces.JobRepository, userRepo interfaces.UserRepository, categoryRepo interfaces.JobCategoryRepository) *JobService {
	return &JobService{
		repo:         repo,
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
	}
}

// GetAllJobs retrieves all jobs with pagination, filtering, and sorting
func (s *JobService) GetAllJobs(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Job, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

// GetJobByID retrieves a job by ID
func (s *JobService) GetJobByID(ctx context.Context, id string) (*models.Job, error) {
	return s.repo.GetByID(ctx, id)
}

// GetJobsByUser retrieves jobs by user (recruiter) ID
func (s *JobService) GetJobsByUser(ctx context.Context, userID string) ([]models.Job, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// CreateJob creates a new job
func (s *JobService) CreateJob(ctx context.Context, job *models.Job) error {
	if _, err := s.userRepo.GetByID(ctx, job.UserID.Hex()); err != nil {
		return fmt.Errorf("user not found")
	}

	if _, err := s.categoryRepo.GetByID(ctx, job.CategoryID.Hex()); err != nil {
		return fmt.Errorf("job category not found")
	}

	return s.repo.Create(ctx, job)
}

// DeleteJob deletes a job by ID
func (s *JobService) DeleteJob(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
