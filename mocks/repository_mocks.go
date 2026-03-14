package mocks

import (
	"context"
	"go-mongodb-api/models"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock for interfaces.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.User, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, id string, user *models.User) (*models.User, error) {
	args := m.Called(ctx, id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockApplicationRepository is a mock for interfaces.ApplicationRepository
type MockApplicationRepository struct {
	mock.Mock
}

func (m *MockApplicationRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Application, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.Application), args.Get(1).(int64), args.Error(2)
}

func (m *MockApplicationRepository) GetByID(ctx context.Context, id string) (*models.Application, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Application), args.Error(1)
}

func (m *MockApplicationRepository) GetByJobID(ctx context.Context, jobID string) ([]models.Application, error) {
	args := m.Called(ctx, jobID)
	return args.Get(0).([]models.Application), args.Error(1)
}

func (m *MockApplicationRepository) GetByUserID(ctx context.Context, userID string) ([]models.Application, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Application), args.Error(1)
}

func (m *MockApplicationRepository) Create(ctx context.Context, application *models.Application) error {
	args := m.Called(ctx, application)
	return args.Error(0)
}

func (m *MockApplicationRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockApplicationRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockJobRepository is a mock for interfaces.JobRepository
type MockJobRepository struct {
	mock.Mock
}

func (m *MockJobRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Job, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.Job), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobRepository) GetByID(ctx context.Context, id string) (*models.Job, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Job), args.Error(1)
}

func (m *MockJobRepository) GetByUserID(ctx context.Context, userID string) ([]models.Job, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Job), args.Error(1)
}

func (m *MockJobRepository) Create(ctx context.Context, job *models.Job) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *MockJobRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockSkillRepository is a mock for interfaces.SkillRepository
type MockSkillRepository struct {
	mock.Mock
}

func (m *MockSkillRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Skill, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.Skill), args.Get(1).(int64), args.Error(2)
}

func (m *MockSkillRepository) GetByID(ctx context.Context, id string) (*models.Skill, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Skill), args.Error(1)
}

func (m *MockSkillRepository) Create(ctx context.Context, skill *models.Skill) error {
	args := m.Called(ctx, skill)
	return args.Error(0)
}

func (m *MockSkillRepository) Update(ctx context.Context, id string, skill *models.Skill) (*models.Skill, error) {
	args := m.Called(ctx, id, skill)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Skill), args.Error(1)
}

func (m *MockSkillRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockJobCategoryRepository is a mock for interfaces.JobCategoryRepository
type MockJobCategoryRepository struct {
	mock.Mock
}

func (m *MockJobCategoryRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobCategory, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.JobCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobCategoryRepository) GetByID(ctx context.Context, id string) (*models.JobCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobCategory), args.Error(1)
}

func (m *MockJobCategoryRepository) Create(ctx context.Context, jobCategory *models.JobCategory) error {
	args := m.Called(ctx, jobCategory)
	return args.Error(0)
}

func (m *MockJobCategoryRepository) Update(ctx context.Context, id string, jobCategory *models.JobCategory) (*models.JobCategory, error) {
	args := m.Called(ctx, id, jobCategory)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobCategory), args.Error(1)
}

func (m *MockJobCategoryRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockArticleRepository is a mock for interfaces.ArticleRepository
type MockArticleRepository struct {
	mock.Mock
}

func (m *MockArticleRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Article, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) GetByID(ctx context.Context, id string) (*models.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Article), args.Error(1)
}

func (m *MockArticleRepository) Create(ctx context.Context, article *models.Article) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockArticleRepository) Update(ctx context.Context, id string, article *models.Article) (*models.Article, error) {
	args := m.Called(ctx, id, article)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Article), args.Error(1)
}

func (m *MockArticleRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCountryRepository is a mock for interfaces.CountryRepository
type MockCountryRepository struct {
	mock.Mock
}

func (m *MockCountryRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Country, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.Country), args.Get(1).(int64), args.Error(2)
}

func (m *MockCountryRepository) GetByID(ctx context.Context, id string) (*models.Country, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryRepository) Create(ctx context.Context, country *models.Country) error {
	args := m.Called(ctx, country)
	return args.Error(0)
}

func (m *MockCountryRepository) Update(ctx context.Context, id string, country *models.Country) (*models.Country, error) {
	args := m.Called(ctx, id, country)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEducationLevelRepository is a mock for interfaces.EducationLevelRepository
type MockEducationLevelRepository struct {
	mock.Mock
}

func (m *MockEducationLevelRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.EducationLevel, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.EducationLevel), args.Get(1).(int64), args.Error(2)
}

func (m *MockEducationLevelRepository) GetByID(ctx context.Context, id string) (*models.EducationLevel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EducationLevel), args.Error(1)
}

func (m *MockEducationLevelRepository) Create(ctx context.Context, educationLevel *models.EducationLevel) error {
	args := m.Called(ctx, educationLevel)
	return args.Error(0)
}

func (m *MockEducationLevelRepository) Update(ctx context.Context, id string, educationLevel *models.EducationLevel) (*models.EducationLevel, error) {
	args := m.Called(ctx, id, educationLevel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EducationLevel), args.Error(1)
}

func (m *MockEducationLevelRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockJobTypeRepository is a mock for interfaces.JobTypeRepository
type MockJobTypeRepository struct {
	mock.Mock
}

func (m *MockJobTypeRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobType, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.JobType), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobTypeRepository) GetByID(ctx context.Context, id string) (*models.JobType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobType), args.Error(1)
}

func (m *MockJobTypeRepository) Create(ctx context.Context, jobType *models.JobType) error {
	args := m.Called(ctx, jobType)
	return args.Error(0)
}

func (m *MockJobTypeRepository) Update(ctx context.Context, id string, jobType *models.JobType) (*models.JobType, error) {
	args := m.Called(ctx, id, jobType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobType), args.Error(1)
}

func (m *MockJobTypeRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockKnowledgeLevelRepository is a mock for interfaces.KnowledgeLevelRepository
type MockKnowledgeLevelRepository struct {
	mock.Mock
}

func (m *MockKnowledgeLevelRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.KnowledgeLevel, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.KnowledgeLevel), args.Get(1).(int64), args.Error(2)
}

func (m *MockKnowledgeLevelRepository) GetByID(ctx context.Context, id string) (*models.KnowledgeLevel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KnowledgeLevel), args.Error(1)
}

func (m *MockKnowledgeLevelRepository) Create(ctx context.Context, knowledgeLevel *models.KnowledgeLevel) error {
	args := m.Called(ctx, knowledgeLevel)
	return args.Error(0)
}

func (m *MockKnowledgeLevelRepository) Update(ctx context.Context, id string, knowledgeLevel *models.KnowledgeLevel) (*models.KnowledgeLevel, error) {
	args := m.Called(ctx, id, knowledgeLevel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KnowledgeLevel), args.Error(1)
}

func (m *MockKnowledgeLevelRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockLocationAvailabilityRepository is a mock for interfaces.LocationAvailabilityRepository
type MockLocationAvailabilityRepository struct {
	mock.Mock
}

func (m *MockLocationAvailabilityRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.LocationAvailability, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.LocationAvailability), args.Get(1).(int64), args.Error(2)
}

func (m *MockLocationAvailabilityRepository) GetByID(ctx context.Context, id string) (*models.LocationAvailability, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LocationAvailability), args.Error(1)
}

func (m *MockLocationAvailabilityRepository) Create(ctx context.Context, locationAvailability *models.LocationAvailability) error {
	args := m.Called(ctx, locationAvailability)
	return args.Error(0)
}

func (m *MockLocationAvailabilityRepository) Update(ctx context.Context, id string, locationAvailability *models.LocationAvailability) (*models.LocationAvailability, error) {
	args := m.Called(ctx, id, locationAvailability)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LocationAvailability), args.Error(1)
}

func (m *MockLocationAvailabilityRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCandidateSkillRepository is a mock for interfaces.CandidateSkillRepository
type MockCandidateSkillRepository struct {
	mock.Mock
}

func (m *MockCandidateSkillRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.CandidateSkill, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.CandidateSkill), args.Get(1).(int64), args.Error(2)
}

func (m *MockCandidateSkillRepository) GetByID(ctx context.Context, id string) (*models.CandidateSkill, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CandidateSkill), args.Error(1)
}

func (m *MockCandidateSkillRepository) GetByUserID(ctx context.Context, userID string) ([]models.CandidateSkill, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.CandidateSkill), args.Error(1)
}

func (m *MockCandidateSkillRepository) Create(ctx context.Context, candidateSkill *models.CandidateSkill) error {
	args := m.Called(ctx, candidateSkill)
	return args.Error(0)
}

func (m *MockCandidateSkillRepository) UpdateProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error {
	args := m.Called(ctx, id, proficiencyLevel)
	return args.Error(0)
}

func (m *MockCandidateSkillRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockJobSkillRepository is a mock for interfaces.JobSkillRepository
type MockJobSkillRepository struct {
	mock.Mock
}

func (m *MockJobSkillRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobSkill, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.JobSkill), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobSkillRepository) GetByID(ctx context.Context, id string) (*models.JobSkill, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobSkill), args.Error(1)
}

func (m *MockJobSkillRepository) GetByJobID(ctx context.Context, jobID string) ([]models.JobSkill, error) {
	args := m.Called(ctx, jobID)
	return args.Get(0).([]models.JobSkill), args.Error(1)
}

func (m *MockJobSkillRepository) Create(ctx context.Context, jobSkill *models.JobSkill) error {
	args := m.Called(ctx, jobSkill)
	return args.Error(0)
}

func (m *MockJobSkillRepository) UpdateProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error {
	args := m.Called(ctx, id, proficiencyLevel)
	return args.Error(0)
}

func (m *MockJobSkillRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
