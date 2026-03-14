package services

import (
	"context"
	"fmt"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type CandidateSkillService struct {
	repo      interfaces.CandidateSkillRepository
	userRepo  interfaces.UserRepository
	skillRepo interfaces.SkillRepository
}

// NewCandidateSkillService creates a new candidate skill service
func NewCandidateSkillService(repo interfaces.CandidateSkillRepository, userRepo interfaces.UserRepository, skillRepo interfaces.SkillRepository) *CandidateSkillService {
	return &CandidateSkillService{
		repo:      repo,
		userRepo:  userRepo,
		skillRepo: skillRepo,
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

// GetCandidateSkillsByUserID retrieves all skills for a specific user
func (s *CandidateSkillService) GetCandidateSkillsByUserID(ctx context.Context, userID string) ([]models.CandidateSkill, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// CreateCandidateSkill creates a new candidate skill
func (s *CandidateSkillService) CreateCandidateSkill(ctx context.Context, candidateSkill *models.CandidateSkill) error {
	if _, err := s.userRepo.GetByID(ctx, candidateSkill.UserID.Hex()); err != nil {
		return fmt.Errorf("user not found")
	}

	if _, err := s.skillRepo.GetByID(ctx, candidateSkill.SkillID.Hex()); err != nil {
		return fmt.Errorf("skill not found")
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
