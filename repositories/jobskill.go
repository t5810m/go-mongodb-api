package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JobSkillRepository struct {
	collection *mongo.Collection
}

// NewJobSkillRepository creates a new job skill repository
func NewJobSkillRepository(db *mongo.Database) *JobSkillRepository {
	return &JobSkillRepository{
		collection: db.Collection("jobskills"),
	}
}

// GetAll retrieves all job skills with pagination and optional filtering
func (r *JobSkillRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobSkill, int64, error) {
	pagination := helpers.NewPagination(page, limit)

	// Build filter query
	filter := bson.M{}

	// Search by job_id (exact match)
	if jobID, exists := filters["job_id"]; exists && jobID != "" {
		objID, err := bson.ObjectIDFromHex(jobID)
		if err == nil {
			filter["job_id"] = objID
		}
	}

	// Search by skill_id (exact match)
	if skillID, exists := filters["skill_id"]; exists && skillID != "" {
		objID, err := bson.ObjectIDFromHex(skillID)
		if err == nil {
			filter["skill_id"] = objID
		}
	}

	// Search by proficiency_level_required (partial match)
	if proficiencyLevel, exists := filters["proficiency_level_required"]; exists && proficiencyLevel != "" {
		filter["proficiency_level_required"] = bson.M{"$regex": proficiencyLevel, "$options": "i"}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sortableFields := []string{"job_id", "skill_id", "proficiency_level_required", "created_time"}
	sortField := "created_time"
	if sort != "" {
		for _, field := range sortableFields {
			if field == sort {
				sortField = sort
				break
			}
		}
	}

	sortOrder := int32(-1)
	if order == "asc" {
		sortOrder = 1
	}

	opts := options.Find().
		SetSkip(int64(pagination.GetSkip())).
		SetLimit(int64(pagination.Limit)).
		SetSort(bson.M{sortField: sortOrder})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var jobSkills []models.JobSkill
	if err = cursor.All(ctx, &jobSkills); err != nil {
		return nil, 0, err
	}

	return jobSkills, total, nil
}

// GetByID retrieves a job skill by ID
func (r *JobSkillRepository) GetByID(ctx context.Context, id string) (*models.JobSkill, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var jobSkill models.JobSkill
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&jobSkill)
	if err != nil {
		return nil, err
	}
	return &jobSkill, nil
}

// GetByJobID retrieves all skills required for a specific job
func (r *JobSkillRepository) GetByJobID(ctx context.Context, jobID string) ([]models.JobSkill, error) {
	objID, err := bson.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"job_id": objID})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var jobSkills []models.JobSkill
	if err = cursor.All(ctx, &jobSkills); err != nil {
		return nil, err
	}

	return jobSkills, nil
}

// Create inserts a new job skill
func (r *JobSkillRepository) Create(ctx context.Context, jobSkill *models.JobSkill) error {
	result, err := r.collection.InsertOne(ctx, jobSkill)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	jobSkill.ID = objID
	return nil
}

// UpdateProficiencyLevel updates the proficiency level required for a job skill
func (r *JobSkillRepository) UpdateProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error {
	if proficiencyLevel == "" {
		return mongo.ErrNoDocuments
	}
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"proficiency_level_required": proficiencyLevel}},
	)
	return err
}

// Delete removes a job skill by ID
func (r *JobSkillRepository) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
