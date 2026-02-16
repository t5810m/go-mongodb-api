package services

import (
	"context"
	"fmt"
	"go-mongodb-api/models"
	"go-mongodb-api/repositories"
)

type CandidateSkillService struct {
	repo          *repositories.CandidateSkillRepository
	candidateRepo *repositories.CandidateRepository
	skillRepo     *repositories.SkillRepository
}

// NewCandidateSkillService creates a new candidate skill service
func NewCandidateSkillService(repo *repositories.CandidateSkillRepository) *CandidateSkillService {
	return &CandidateSkillService{
		repo: repo,
	}
}

// NewCandidateSkillServiceWithDeps creates a new candidate skill service with dependencies
func NewCandidateSkillServiceWithDeps(repo *repositories.CandidateSkillRepository, candidateRepo *repositories.CandidateRepository, skillRepo *repositories.SkillRepository) *CandidateSkillService {
	return &CandidateSkillService{
		repo:          repo,
		candidateRepo: candidateRepo,
		skillRepo:     skillRepo,
	}
}

// GetAllCandidateSkills retrieves all candidate skills with pagination and optional filtering
func (s *CandidateSkillService) GetAllCandidateSkills(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.CandidateSkill, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

// GetCandidateSkillByID retrieves a candidate skill by ID
func (s *CandidateSkillService) GetCandidateSkillByID(ctx context.Context, id string) (*models.CandidateSkill, error) {
	return s.repo.GetByID(ctx, id)
}

// GetCandidateSkillsByCandidateID retrieves all skills for a specific candidate
func (s *CandidateSkillService) GetCandidateSkillsByCandidateID(ctx context.Context, candidateID string) ([]models.CandidateSkill, error) {
	return s.repo.GetByCandidateID(ctx, candidateID)
}

// CreateCandidateSkill creates a new candidate skill
func (s *CandidateSkillService) CreateCandidateSkill(ctx context.Context, candidateSkill *models.CandidateSkill) error {
	// Validate candidate exists
	if s.candidateRepo != nil {
		_, err := s.candidateRepo.GetByID(ctx, candidateSkill.CandidateID.Hex())
		if err != nil {
			return fmt.Errorf("candidate not found")
		}
	}

	// Validate skill exists
	if s.skillRepo != nil {
		_, err := s.skillRepo.GetByID(ctx, candidateSkill.SkillID.Hex())
		if err != nil {
			return fmt.Errorf("skill not found")
		}
	}

	return s.repo.Create(ctx, candidateSkill)
}

// UpdateCandidateSkillProficiencyLevel updates the proficiency level of a candidate skill
func (s *CandidateSkillService) UpdateCandidateSkillProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error {
	return s.repo.UpdateProficiencyLevel(ctx, id, proficiencyLevel)
}

// DeleteCandidateSkill deletes a candidate skill by ID
func (s *CandidateSkillService) DeleteCandidateSkill(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
