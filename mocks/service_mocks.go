package mocks

import (
	"context"
	"go-mongodb-api/models"

	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock for interfaces.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAllUsers(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.User, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error) {
	args := m.Called(ctx, id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockAuthService is a mock for interfaces.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockApplicationService is a mock for interfaces.ApplicationService
type MockApplicationService struct {
	mock.Mock
}

func (m *MockApplicationService) GetAllApplications(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Application, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.Application), args.Get(1).(int64), args.Error(2)
}

func (m *MockApplicationService) GetApplicationByID(ctx context.Context, id string) (*models.Application, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Application), args.Error(1)
}

func (m *MockApplicationService) GetApplicationsByJobID(ctx context.Context, jobID string) ([]models.Application, error) {
	args := m.Called(ctx, jobID)
	return args.Get(0).([]models.Application), args.Error(1)
}

func (m *MockApplicationService) GetApplicationsByUserID(ctx context.Context, userID string) ([]models.Application, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Application), args.Error(1)
}

func (m *MockApplicationService) CreateApplication(ctx context.Context, application *models.Application) error {
	args := m.Called(ctx, application)
	return args.Error(0)
}

func (m *MockApplicationService) UpdateApplicationStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockApplicationService) DeleteApplication(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockJobService is a mock for interfaces.JobService
type MockJobService struct {
	mock.Mock
}

func (m *MockJobService) GetAllJobs(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Job, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.Job), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobService) GetJobByID(ctx context.Context, id string) (*models.Job, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Job), args.Error(1)
}

func (m *MockJobService) GetJobsByUser(ctx context.Context, userID string) ([]models.Job, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Job), args.Error(1)
}

func (m *MockJobService) CreateJob(ctx context.Context, job *models.Job) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *MockJobService) DeleteJob(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockSkillService is a mock for interfaces.SkillService
type MockSkillService struct {
	mock.Mock
}

func (m *MockSkillService) GetAllSkills(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Skill, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.Skill), args.Get(1).(int64), args.Error(2)
}

func (m *MockSkillService) GetSkillByID(ctx context.Context, id string) (*models.Skill, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Skill), args.Error(1)
}

func (m *MockSkillService) CreateSkill(ctx context.Context, skill *models.Skill) error {
	args := m.Called(ctx, skill)
	return args.Error(0)
}

func (m *MockSkillService) UpdateSkill(ctx context.Context, id string, skill *models.Skill) (*models.Skill, error) {
	args := m.Called(ctx, id, skill)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Skill), args.Error(1)
}

func (m *MockSkillService) DeleteSkill(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockJobCategoryService is a mock for interfaces.JobCategoryService
type MockJobCategoryService struct {
	mock.Mock
}

func (m *MockJobCategoryService) GetAllJobCategories(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobCategory, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.JobCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobCategoryService) GetJobCategoryByID(ctx context.Context, id string) (*models.JobCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobCategory), args.Error(1)
}

func (m *MockJobCategoryService) CreateJobCategory(ctx context.Context, jobCategory *models.JobCategory) error {
	args := m.Called(ctx, jobCategory)
	return args.Error(0)
}

func (m *MockJobCategoryService) UpdateJobCategory(ctx context.Context, id string, jobCategory *models.JobCategory) (*models.JobCategory, error) {
	args := m.Called(ctx, id, jobCategory)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobCategory), args.Error(1)
}

func (m *MockJobCategoryService) DeleteJobCategory(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockArticleService is a mock for interfaces.ArticleService
type MockArticleService struct {
	mock.Mock
}

func (m *MockArticleService) GetAllArticles(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Article, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleService) GetArticleByID(ctx context.Context, id string) (*models.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Article), args.Error(1)
}

func (m *MockArticleService) CreateArticle(ctx context.Context, article *models.Article) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockArticleService) UpdateArticle(ctx context.Context, id string, article *models.Article) (*models.Article, error) {
	args := m.Called(ctx, id, article)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Article), args.Error(1)
}

func (m *MockArticleService) DeleteArticle(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCountryService is a mock for interfaces.CountryService
type MockCountryService struct {
	mock.Mock
}

func (m *MockCountryService) GetAllCountries(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Country, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.Country), args.Get(1).(int64), args.Error(2)
}

func (m *MockCountryService) GetCountryByID(ctx context.Context, id string) (*models.Country, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryService) CreateCountry(ctx context.Context, country *models.Country) error {
	args := m.Called(ctx, country)
	return args.Error(0)
}

func (m *MockCountryService) UpdateCountry(ctx context.Context, id string, country *models.Country) (*models.Country, error) {
	args := m.Called(ctx, id, country)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryService) DeleteCountry(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEducationLevelService is a mock for interfaces.EducationLevelService
type MockEducationLevelService struct {
	mock.Mock
}

func (m *MockEducationLevelService) GetAllEducationLevels(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.EducationLevel, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.EducationLevel), args.Get(1).(int64), args.Error(2)
}

func (m *MockEducationLevelService) GetEducationLevelByID(ctx context.Context, id string) (*models.EducationLevel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EducationLevel), args.Error(1)
}

func (m *MockEducationLevelService) CreateEducationLevel(ctx context.Context, educationLevel *models.EducationLevel) error {
	args := m.Called(ctx, educationLevel)
	return args.Error(0)
}

func (m *MockEducationLevelService) UpdateEducationLevel(ctx context.Context, id string, educationLevel *models.EducationLevel) (*models.EducationLevel, error) {
	args := m.Called(ctx, id, educationLevel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EducationLevel), args.Error(1)
}

func (m *MockEducationLevelService) DeleteEducationLevel(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockJobTypeService is a mock for interfaces.JobTypeService
type MockJobTypeService struct {
	mock.Mock
}

func (m *MockJobTypeService) GetAllJobTypes(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobType, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.JobType), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobTypeService) GetJobTypeByID(ctx context.Context, id string) (*models.JobType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobType), args.Error(1)
}

func (m *MockJobTypeService) CreateJobType(ctx context.Context, jobType *models.JobType) error {
	args := m.Called(ctx, jobType)
	return args.Error(0)
}

func (m *MockJobTypeService) UpdateJobType(ctx context.Context, id string, jobType *models.JobType) (*models.JobType, error) {
	args := m.Called(ctx, id, jobType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobType), args.Error(1)
}

func (m *MockJobTypeService) DeleteJobType(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockKnowledgeLevelService is a mock for interfaces.KnowledgeLevelService
type MockKnowledgeLevelService struct {
	mock.Mock
}

func (m *MockKnowledgeLevelService) GetAllKnowledgeLevels(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.KnowledgeLevel, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.KnowledgeLevel), args.Get(1).(int64), args.Error(2)
}

func (m *MockKnowledgeLevelService) GetKnowledgeLevelByID(ctx context.Context, id string) (*models.KnowledgeLevel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KnowledgeLevel), args.Error(1)
}

func (m *MockKnowledgeLevelService) CreateKnowledgeLevel(ctx context.Context, knowledgeLevel *models.KnowledgeLevel) error {
	args := m.Called(ctx, knowledgeLevel)
	return args.Error(0)
}

func (m *MockKnowledgeLevelService) UpdateKnowledgeLevel(ctx context.Context, id string, knowledgeLevel *models.KnowledgeLevel) (*models.KnowledgeLevel, error) {
	args := m.Called(ctx, id, knowledgeLevel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KnowledgeLevel), args.Error(1)
}

func (m *MockKnowledgeLevelService) DeleteKnowledgeLevel(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockLocationAvailabilityService is a mock for interfaces.LocationAvailabilityService
type MockLocationAvailabilityService struct {
	mock.Mock
}

func (m *MockLocationAvailabilityService) GetAllLocationAvailabilities(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.LocationAvailability, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.LocationAvailability), args.Get(1).(int64), args.Error(2)
}

func (m *MockLocationAvailabilityService) GetLocationAvailabilityByID(ctx context.Context, id string) (*models.LocationAvailability, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LocationAvailability), args.Error(1)
}

func (m *MockLocationAvailabilityService) CreateLocationAvailability(ctx context.Context, locationAvailability *models.LocationAvailability) error {
	args := m.Called(ctx, locationAvailability)
	return args.Error(0)
}

func (m *MockLocationAvailabilityService) UpdateLocationAvailability(ctx context.Context, id string, locationAvailability *models.LocationAvailability) (*models.LocationAvailability, error) {
	args := m.Called(ctx, id, locationAvailability)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LocationAvailability), args.Error(1)
}

func (m *MockLocationAvailabilityService) DeleteLocationAvailability(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCandidateSkillService is a mock for interfaces.CandidateSkillService
type MockCandidateSkillService struct {
	mock.Mock
}

func (m *MockCandidateSkillService) GetAllCandidateSkills(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.CandidateSkill, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.CandidateSkill), args.Get(1).(int64), args.Error(2)
}

func (m *MockCandidateSkillService) GetCandidateSkillByID(ctx context.Context, id string) (*models.CandidateSkill, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CandidateSkill), args.Error(1)
}

func (m *MockCandidateSkillService) GetCandidateSkillsByUserID(ctx context.Context, userID string) ([]models.CandidateSkill, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.CandidateSkill), args.Error(1)
}

func (m *MockCandidateSkillService) CreateCandidateSkill(ctx context.Context, candidateSkill *models.CandidateSkill) error {
	args := m.Called(ctx, candidateSkill)
	return args.Error(0)
}

func (m *MockCandidateSkillService) UpdateCandidateSkillProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error {
	args := m.Called(ctx, id, proficiencyLevel)
	return args.Error(0)
}

func (m *MockCandidateSkillService) DeleteCandidateSkill(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockJobSkillService is a mock for interfaces.JobSkillService
type MockJobSkillService struct {
	mock.Mock
}

func (m *MockJobSkillService) GetAllJobSkills(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobSkill, int64, error) {
	args := m.Called(ctx, page, limit, filters, sort, order)
	return args.Get(0).([]models.JobSkill), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobSkillService) GetJobSkillByID(ctx context.Context, id string) (*models.JobSkill, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobSkill), args.Error(1)
}

func (m *MockJobSkillService) GetJobSkillsByJobID(ctx context.Context, jobID string) ([]models.JobSkill, error) {
	args := m.Called(ctx, jobID)
	return args.Get(0).([]models.JobSkill), args.Error(1)
}

func (m *MockJobSkillService) CreateJobSkill(ctx context.Context, jobSkill *models.JobSkill) error {
	args := m.Called(ctx, jobSkill)
	return args.Error(0)
}

func (m *MockJobSkillService) UpdateJobSkillProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error {
	args := m.Called(ctx, id, proficiencyLevel)
	return args.Error(0)
}

func (m *MockJobSkillService) DeleteJobSkill(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
