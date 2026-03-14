package services

import (
	"context"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type JobCategoryService struct {
	repo interfaces.JobCategoryRepository
}

// NewJobCategoryService creates a new job category service
func NewJobCategoryService(repo interfaces.JobCategoryRepository) *JobCategoryService {
	return &JobCategoryService{repo: repo}
}

// GetAllJobCategories retrieves all job categories with pagination support
func (s *JobCategoryService) GetAllJobCategories(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobCategory, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

// GetJobCategoryByID retrieves a job category by ID
func (s *JobCategoryService) GetJobCategoryByID(ctx context.Context, id string) (*models.JobCategory, error) {
	return s.repo.GetByID(ctx, id)
}

// CreateJobCategory creates a new job category
func (s *JobCategoryService) CreateJobCategory(ctx context.Context, jobCategory *models.JobCategory) error {
	return s.repo.Create(ctx, jobCategory)
}

// UpdateJobCategory updates a job category's allowed fields
func (s *JobCategoryService) UpdateJobCategory(ctx context.Context, id string, jobCategory *models.JobCategory) (*models.JobCategory, error) {
	return s.repo.Update(ctx, id, jobCategory)
}

// DeleteJobCategory deletes a job category by ID (only if no jobs use it)
func (s *JobCategoryService) DeleteJobCategory(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
