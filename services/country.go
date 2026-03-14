package services

import (
	"context"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type CountryService struct {
	repo interfaces.CountryRepository
}

func NewCountryService(repo interfaces.CountryRepository) *CountryService {
	return &CountryService{repo: repo}
}

func (s *CountryService) GetAllCountries(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Country, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

func (s *CountryService) GetCountryByID(ctx context.Context, id string) (*models.Country, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CountryService) CreateCountry(ctx context.Context, country *models.Country) error {
	return s.repo.Create(ctx, country)
}

func (s *CountryService) UpdateCountry(ctx context.Context, id string, country *models.Country) (*models.Country, error) {
	return s.repo.Update(ctx, id, country)
}

func (s *CountryService) DeleteCountry(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
