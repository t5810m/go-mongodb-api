package services

import (
	"context"
	"fmt"
	"go-mongodb-api/models"
	"go-mongodb-api/repositories"
)

type JobService struct {
	repo          *repositories.JobRepository
	recruiterRepo *repositories.RecruiterRepository
	companyRepo   *repositories.CompanyRepository
	categoryRepo  *repositories.JobCategoryRepository
}

// NewJobService creates a new job service
func NewJobService(repo *repositories.JobRepository) *JobService {
	return &JobService{
		repo: repo,
	}
}

// NewJobServiceWithDeps creates a new job service with dependencies
func NewJobServiceWithDeps(repo *repositories.JobRepository, recruiterRepo *repositories.RecruiterRepository, companyRepo *repositories.CompanyRepository, categoryRepo *repositories.JobCategoryRepository) *JobService {
	return &JobService{
		repo:          repo,
		recruiterRepo: recruiterRepo,
		companyRepo:   companyRepo,
		categoryRepo:  categoryRepo,
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

// GetJobsByCompany retrieves jobs by company ID
func (s *JobService) GetJobsByCompany(ctx context.Context, companyID string) ([]models.Job, error) {
	return s.repo.GetByCompanyID(ctx, companyID)
}

// CreateJob creates a new job
func (s *JobService) CreateJob(ctx context.Context, job *models.Job) error {
	// Validate recruiter exists
	if s.recruiterRepo != nil {
		_, err := s.recruiterRepo.GetByID(ctx, job.RecruiterID.Hex())
		if err != nil {
			return fmt.Errorf("recruiter not found")
		}
	}

	// Validate company exists
	if s.companyRepo != nil {
		_, err := s.companyRepo.GetByID(ctx, job.CompanyID.Hex())
		if err != nil {
			return fmt.Errorf("company not found")
		}
	}

	// Validate category exists
	if s.categoryRepo != nil {
		_, err := s.categoryRepo.GetByID(ctx, job.CategoryID.Hex())
		if err != nil {
			return fmt.Errorf("job category not found")
		}
	}

	return s.repo.Create(ctx, job)
}

// DeleteJob deletes a job by ID
func (s *JobService) DeleteJob(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
