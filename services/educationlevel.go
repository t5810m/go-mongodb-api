package services

import (
	"context"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type EducationLevelService struct {
	repo interfaces.EducationLevelRepository
}

func NewEducationLevelService(repo interfaces.EducationLevelRepository) *EducationLevelService {
	return &EducationLevelService{repo: repo}
}

func (s *EducationLevelService) GetAllEducationLevels(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.EducationLevel, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

func (s *EducationLevelService) GetEducationLevelByID(ctx context.Context, id string) (*models.EducationLevel, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EducationLevelService) CreateEducationLevel(ctx context.Context, educationLevel *models.EducationLevel) error {
	return s.repo.Create(ctx, educationLevel)
}

func (s *EducationLevelService) UpdateEducationLevel(ctx context.Context, id string, educationLevel *models.EducationLevel) (*models.EducationLevel, error) {
	return s.repo.Update(ctx, id, educationLevel)
}

func (s *EducationLevelService) DeleteEducationLevel(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
