package services

import (
	"context"
	"go-mongodb-api/models"
	"go-mongodb-api/repositories"
)

type CompanyService struct {
	repo *repositories.CompanyRepository
}

// NewCompanyService creates a new company service
func NewCompanyService(repo *repositories.CompanyRepository) *CompanyService {
	return &CompanyService{
		repo: repo,
	}
}

// GetAllCompanies retrieves all companies with pagination, filtering, and sorting
func (s *CompanyService) GetAllCompanies(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Company, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

// GetCompanyByID retrieves a company by ID
func (s *CompanyService) GetCompanyByID(ctx context.Context, id string) (*models.Company, error) {
	return s.repo.GetByID(ctx, id)
}

// CreateCompany creates a new company
func (s *CompanyService) CreateCompany(ctx context.Context, company *models.Company) error {
	return s.repo.Create(ctx, company)
}

// DeleteCompany deletes a company by ID
func (s *CompanyService) DeleteCompany(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
