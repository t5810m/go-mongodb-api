package services

import (
	"context"
	"go-mongodb-api/interfaces"
	"go-mongodb-api/models"
)

type LocationAvailabilityService struct {
	repo interfaces.LocationAvailabilityRepository
}

func NewLocationAvailabilityService(repo interfaces.LocationAvailabilityRepository) *LocationAvailabilityService {
	return &LocationAvailabilityService{repo: repo}
}

func (s *LocationAvailabilityService) GetAllLocationAvailabilities(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.LocationAvailability, int64, error) {
	return s.repo.GetAll(ctx, page, limit, filters, sort, order)
}

func (s *LocationAvailabilityService) GetLocationAvailabilityByID(ctx context.Context, id string) (*models.LocationAvailability, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *LocationAvailabilityService) CreateLocationAvailability(ctx context.Context, locationAvailability *models.LocationAvailability) error {
	return s.repo.Create(ctx, locationAvailability)
}

func (s *LocationAvailabilityService) UpdateLocationAvailability(ctx context.Context, id string, locationAvailability *models.LocationAvailability) (*models.LocationAvailability, error) {
	return s.repo.Update(ctx, id, locationAvailability)
}

func (s *LocationAvailabilityService) DeleteLocationAvailability(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
