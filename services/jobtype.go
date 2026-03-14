package services

import (
	"context"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type JobTypeService struct {
	repo interfaces.JobTypeRepository
}

func NewJobTypeService(repo interfaces.JobTypeRepository) *JobTypeService {
	return &JobTypeService{repo: repo}
}

func (s *JobTypeService) GetAllJobTypes(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobType, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

func (s *JobTypeService) GetJobTypeByID(ctx context.Context, id string) (*models.JobType, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *JobTypeService) CreateJobType(ctx context.Context, jobType *models.JobType) error {
	return s.repo.Create(ctx, jobType)
}

func (s *JobTypeService) UpdateJobType(ctx context.Context, id string, jobType *models.JobType) (*models.JobType, error) {
	return s.repo.Update(ctx, id, jobType)
}

func (s *JobTypeService) DeleteJobType(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
