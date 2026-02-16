package services

import (
	"context"
	"go-mongodb-api/models"
	"go-mongodb-api/repositories"
)

type RecruiterService struct {
	repo *repositories.RecruiterRepository
}

// NewRecruiterService creates a new recruiter service
func NewRecruiterService(repo *repositories.RecruiterRepository) *RecruiterService {
	return &RecruiterService{
		repo: repo,
	}
}

// GetAllRecruiters retrieves all recruiters with pagination, filtering, and sorting
func (s *RecruiterService) GetAllRecruiters(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Recruiter, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

// GetRecruiterByID retrieves a recruiter by ID
func (s *RecruiterService) GetRecruiterByID(ctx context.Context, id string) (*models.Recruiter, error) {
	return s.repo.GetByID(ctx, id)
}

// CreateRecruiter creates a new recruiter
func (s *RecruiterService) CreateRecruiter(ctx context.Context, recruiter *models.Recruiter) error {
	return s.repo.Create(ctx, recruiter)
}

// DeleteRecruiter deletes a recruiter by ID
func (s *RecruiterService) DeleteRecruiter(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
