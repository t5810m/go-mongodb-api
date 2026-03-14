package services

import (
	"context"
	"fmt"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type JobSkillService struct {
	repo      interfaces.JobSkillRepository
	jobRepo   interfaces.JobRepository
	skillRepo interfaces.SkillRepository
}

// NewJobSkillService creates a new job skill service
func NewJobSkillService(repo interfaces.JobSkillRepository, jobRepo interfaces.JobRepository, skillRepo interfaces.SkillRepository) *JobSkillService {
	return &JobSkillService{
		repo:      repo,
		jobRepo:   jobRepo,
		skillRepo: skillRepo,
	}
}

// GetAllJobSkills retrieves all job skills with pagination and optional filtering
func (s *JobSkillService) GetAllJobSkills(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobSkill, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

// GetJobSkillByID retrieves a job skill by ID
func (s *JobSkillService) GetJobSkillByID(ctx context.Context, id string) (*models.JobSkill, error) {
	return s.repo.GetByID(ctx, id)
}

// GetJobSkillsByJobID retrieves all skills required for a specific job
func (s *JobSkillService) GetJobSkillsByJobID(ctx context.Context, jobID string) ([]models.JobSkill, error) {
	return s.repo.GetByJobID(ctx, jobID)
}

// CreateJobSkill creates a new job skill
func (s *JobSkillService) CreateJobSkill(ctx context.Context, jobSkill *models.JobSkill) error {
	if _, err := s.jobRepo.GetByID(ctx, jobSkill.JobID.Hex()); err != nil {
		return fmt.Errorf("job not found")
	}

	if _, err := s.skillRepo.GetByID(ctx, jobSkill.SkillID.Hex()); err != nil {
		return fmt.Errorf("skill not found")
	}

	return s.repo.Create(ctx, jobSkill)
}

// UpdateJobSkillProficiencyLevel updates the proficiency level required for a job skill
func (s *JobSkillService) UpdateJobSkillProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error {
	return s.repo.UpdateProficiencyLevel(ctx, id, proficiencyLevel)
}

// DeleteJobSkill deletes a job skill by ID
func (s *JobSkillService) DeleteJobSkill(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
