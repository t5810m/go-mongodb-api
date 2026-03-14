package interfaces

import (
	"context"
	"go-mongodb-api/models"
)

type UserService interface {
	GetAllUsers(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.User, int64, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, user *models.User) error
}

type ApplicationService interface {
	GetAllApplications(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Application, int64, error)
	GetApplicationByID(ctx context.Context, id string) (*models.Application, error)
	GetApplicationsByJobID(ctx context.Context, jobID string) ([]models.Application, error)
	GetApplicationsByUserID(ctx context.Context, userID string) ([]models.Application, error)
	CreateApplication(ctx context.Context, application *models.Application) error
	UpdateApplicationStatus(ctx context.Context, id string, status string) error
	DeleteApplication(ctx context.Context, id string) error
}

type JobService interface {
	GetAllJobs(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Job, int64, error)
	GetJobByID(ctx context.Context, id string) (*models.Job, error)
	GetJobsByUser(ctx context.Context, userID string) ([]models.Job, error)
	CreateJob(ctx context.Context, job *models.Job) error
	DeleteJob(ctx context.Context, id string) error
}

type SkillService interface {
	GetAllSkills(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Skill, int64, error)
	GetSkillByID(ctx context.Context, id string) (*models.Skill, error)
	CreateSkill(ctx context.Context, skill *models.Skill) error
	UpdateSkill(ctx context.Context, id string, skill *models.Skill) (*models.Skill, error)
	DeleteSkill(ctx context.Context, id string) error
}

type JobCategoryService interface {
	GetAllJobCategories(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobCategory, int64, error)
	GetJobCategoryByID(ctx context.Context, id string) (*models.JobCategory, error)
	CreateJobCategory(ctx context.Context, jobCategory *models.JobCategory) error
	UpdateJobCategory(ctx context.Context, id string, jobCategory *models.JobCategory) (*models.JobCategory, error)
	DeleteJobCategory(ctx context.Context, id string) error
}

type ArticleService interface {
	GetAllArticles(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Article, int64, error)
	GetArticleByID(ctx context.Context, id string) (*models.Article, error)
	CreateArticle(ctx context.Context, article *models.Article) error
	UpdateArticle(ctx context.Context, id string, article *models.Article) (*models.Article, error)
	DeleteArticle(ctx context.Context, id string) error
}

type CountryService interface {
	GetAllCountries(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Country, int64, error)
	GetCountryByID(ctx context.Context, id string) (*models.Country, error)
	CreateCountry(ctx context.Context, country *models.Country) error
	UpdateCountry(ctx context.Context, id string, country *models.Country) (*models.Country, error)
	DeleteCountry(ctx context.Context, id string) error
}

type EducationLevelService interface {
	GetAllEducationLevels(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.EducationLevel, int64, error)
	GetEducationLevelByID(ctx context.Context, id string) (*models.EducationLevel, error)
	CreateEducationLevel(ctx context.Context, educationLevel *models.EducationLevel) error
	UpdateEducationLevel(ctx context.Context, id string, educationLevel *models.EducationLevel) (*models.EducationLevel, error)
	DeleteEducationLevel(ctx context.Context, id string) error
}

type JobTypeService interface {
	GetAllJobTypes(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobType, int64, error)
	GetJobTypeByID(ctx context.Context, id string) (*models.JobType, error)
	CreateJobType(ctx context.Context, jobType *models.JobType) error
	UpdateJobType(ctx context.Context, id string, jobType *models.JobType) (*models.JobType, error)
	DeleteJobType(ctx context.Context, id string) error
}

type KnowledgeLevelService interface {
	GetAllKnowledgeLevels(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.KnowledgeLevel, int64, error)
	GetKnowledgeLevelByID(ctx context.Context, id string) (*models.KnowledgeLevel, error)
	CreateKnowledgeLevel(ctx context.Context, knowledgeLevel *models.KnowledgeLevel) error
	UpdateKnowledgeLevel(ctx context.Context, id string, knowledgeLevel *models.KnowledgeLevel) (*models.KnowledgeLevel, error)
	DeleteKnowledgeLevel(ctx context.Context, id string) error
}

type LocationAvailabilityService interface {
	GetAllLocationAvailabilities(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.LocationAvailability, int64, error)
	GetLocationAvailabilityByID(ctx context.Context, id string) (*models.LocationAvailability, error)
	CreateLocationAvailability(ctx context.Context, locationAvailability *models.LocationAvailability) error
	UpdateLocationAvailability(ctx context.Context, id string, locationAvailability *models.LocationAvailability) (*models.LocationAvailability, error)
	DeleteLocationAvailability(ctx context.Context, id string) error
}

type CandidateSkillService interface {
	GetAllCandidateSkills(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.CandidateSkill, int64, error)
	GetCandidateSkillByID(ctx context.Context, id string) (*models.CandidateSkill, error)
	GetCandidateSkillsByUserID(ctx context.Context, userID string) ([]models.CandidateSkill, error)
	CreateCandidateSkill(ctx context.Context, candidateSkill *models.CandidateSkill) error
	UpdateCandidateSkillProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error
	DeleteCandidateSkill(ctx context.Context, id string) error
}

type JobSkillService interface {
	GetAllJobSkills(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobSkill, int64, error)
	GetJobSkillByID(ctx context.Context, id string) (*models.JobSkill, error)
	GetJobSkillsByJobID(ctx context.Context, jobID string) ([]models.JobSkill, error)
	CreateJobSkill(ctx context.Context, jobSkill *models.JobSkill) error
	UpdateJobSkillProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error
	DeleteJobSkill(ctx context.Context, id string) error
}
