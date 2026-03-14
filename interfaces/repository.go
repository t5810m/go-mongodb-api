package interfaces

import (
	"context"
	"go-mongodb-api/models"
)

type UserRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.User, int64, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, id string, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id string) error
}

type ApplicationRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Application, int64, error)
	GetByID(ctx context.Context, id string) (*models.Application, error)
	GetByJobID(ctx context.Context, jobID string) ([]models.Application, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Application, error)
	Create(ctx context.Context, application *models.Application) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
}

type JobRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Job, int64, error)
	GetByID(ctx context.Context, id string) (*models.Job, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Job, error)
	Create(ctx context.Context, job *models.Job) error
	Delete(ctx context.Context, id string) error
}

type SkillRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Skill, int64, error)
	GetByID(ctx context.Context, id string) (*models.Skill, error)
	Create(ctx context.Context, skill *models.Skill) error
	Update(ctx context.Context, id string, skill *models.Skill) (*models.Skill, error)
	Delete(ctx context.Context, id string) error
}

type JobCategoryRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobCategory, int64, error)
	GetByID(ctx context.Context, id string) (*models.JobCategory, error)
	Create(ctx context.Context, jobCategory *models.JobCategory) error
	Update(ctx context.Context, id string, jobCategory *models.JobCategory) (*models.JobCategory, error)
	Delete(ctx context.Context, id string) error
}

type ArticleRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Article, int64, error)
	GetByID(ctx context.Context, id string) (*models.Article, error)
	Create(ctx context.Context, article *models.Article) error
	Update(ctx context.Context, id string, article *models.Article) (*models.Article, error)
	Delete(ctx context.Context, id string) error
}

type CountryRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Country, int64, error)
	GetByID(ctx context.Context, id string) (*models.Country, error)
	Create(ctx context.Context, country *models.Country) error
	Update(ctx context.Context, id string, country *models.Country) (*models.Country, error)
	Delete(ctx context.Context, id string) error
}

type EducationLevelRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.EducationLevel, int64, error)
	GetByID(ctx context.Context, id string) (*models.EducationLevel, error)
	Create(ctx context.Context, educationLevel *models.EducationLevel) error
	Update(ctx context.Context, id string, educationLevel *models.EducationLevel) (*models.EducationLevel, error)
	Delete(ctx context.Context, id string) error
}

type JobTypeRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobType, int64, error)
	GetByID(ctx context.Context, id string) (*models.JobType, error)
	Create(ctx context.Context, jobType *models.JobType) error
	Update(ctx context.Context, id string, jobType *models.JobType) (*models.JobType, error)
	Delete(ctx context.Context, id string) error
}

type KnowledgeLevelRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.KnowledgeLevel, int64, error)
	GetByID(ctx context.Context, id string) (*models.KnowledgeLevel, error)
	Create(ctx context.Context, knowledgeLevel *models.KnowledgeLevel) error
	Update(ctx context.Context, id string, knowledgeLevel *models.KnowledgeLevel) (*models.KnowledgeLevel, error)
	Delete(ctx context.Context, id string) error
}

type LocationAvailabilityRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.LocationAvailability, int64, error)
	GetByID(ctx context.Context, id string) (*models.LocationAvailability, error)
	Create(ctx context.Context, locationAvailability *models.LocationAvailability) error
	Update(ctx context.Context, id string, locationAvailability *models.LocationAvailability) (*models.LocationAvailability, error)
	Delete(ctx context.Context, id string) error
}

type CandidateSkillRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.CandidateSkill, int64, error)
	GetByID(ctx context.Context, id string) (*models.CandidateSkill, error)
	GetByUserID(ctx context.Context, userID string) ([]models.CandidateSkill, error)
	Create(ctx context.Context, candidateSkill *models.CandidateSkill) error
	UpdateProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error
	Delete(ctx context.Context, id string) error
}

type JobSkillRepository interface {
	GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobSkill, int64, error)
	GetByID(ctx context.Context, id string) (*models.JobSkill, error)
	GetByJobID(ctx context.Context, jobID string) ([]models.JobSkill, error)
	Create(ctx context.Context, jobSkill *models.JobSkill) error
	UpdateProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error
	Delete(ctx context.Context, id string) error
}
