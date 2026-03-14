package services

import (
	"context"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type KnowledgeLevelService struct {
	repo interfaces.KnowledgeLevelRepository
}

func NewKnowledgeLevelService(repo interfaces.KnowledgeLevelRepository) *KnowledgeLevelService {
	return &KnowledgeLevelService{repo: repo}
}

func (s *KnowledgeLevelService) GetAllKnowledgeLevels(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.KnowledgeLevel, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

func (s *KnowledgeLevelService) GetKnowledgeLevelByID(ctx context.Context, id string) (*models.KnowledgeLevel, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *KnowledgeLevelService) CreateKnowledgeLevel(ctx context.Context, knowledgeLevel *models.KnowledgeLevel) error {
	return s.repo.Create(ctx, knowledgeLevel)
}

func (s *KnowledgeLevelService) UpdateKnowledgeLevel(ctx context.Context, id string, knowledgeLevel *models.KnowledgeLevel) (*models.KnowledgeLevel, error) {
	return s.repo.Update(ctx, id, knowledgeLevel)
}

func (s *KnowledgeLevelService) DeleteKnowledgeLevel(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
