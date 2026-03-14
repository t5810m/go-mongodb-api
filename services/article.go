package services

import (
	"context"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type ArticleService struct {
	repo interfaces.ArticleRepository
}

func NewArticleService(repo interfaces.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) GetAllArticles(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Article, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

func (s *ArticleService) GetArticleByID(ctx context.Context, id string) (*models.Article, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ArticleService) CreateArticle(ctx context.Context, article *models.Article) error {
	return s.repo.Create(ctx, article)
}

func (s *ArticleService) UpdateArticle(ctx context.Context, id string, article *models.Article) (*models.Article, error) {
	return s.repo.Update(ctx, id, article)
}

func (s *ArticleService) DeleteArticle(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
