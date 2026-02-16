package services

import (
	"context"
	"go-mongodb-api/models"
	"go-mongodb-api/repositories"
)

type CandidateService struct {
	repo *repositories.CandidateRepository
}

// NewCandidateService creates a new candidate service
func NewCandidateService(repo *repositories.CandidateRepository) *CandidateService {
	return &CandidateService{
		repo: repo,
	}
}

// GetAllCandidates retrieves all candidates with pagination, filtering, and sorting
func (s *CandidateService) GetAllCandidates(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Candidate, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

func (s *CandidateService) GetCandidateByID(ctx context.Context, id string) (*models.Candidate, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CandidateService) CreateCandidate(ctx context.Context, candidate *models.Candidate) error {
	return s.repo.Create(ctx, candidate)
}

// DeleteCandidate deletes a candidate by ID
func (s *CandidateService) DeleteCandidate(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
