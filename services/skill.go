package services

import (
	"context"
	"go-mongodb-api/models"
	"go-mongodb-api/repositories"
)

type SkillService struct {
	repo *repositories.SkillRepository
}

// NewSkillService creates a new skill service
func NewSkillService(repo *repositories.SkillRepository) *SkillService {
	return &SkillService{
		repo: repo,
	}
}

// GetAllSkills retrieves all skills with pagination, filtering, and sorting
func (s *SkillService) GetAllSkills(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Skill, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

// GetSkillByID retrieves a skill by ID
func (s *SkillService) GetSkillByID(ctx context.Context, id string) (*models.Skill, error) {
	return s.repo.GetByID(ctx, id)
}

// CreateSkill creates a new skill
func (s *SkillService) CreateSkill(ctx context.Context, skill *models.Skill) error {
	return s.repo.Create(ctx, skill)
}

// DeleteSkill deletes a skill by ID
func (s *SkillService) DeleteSkill(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
