package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CandidateSkillRepository struct {
	collection *mongo.Collection
}

// NewCandidateSkillRepository creates a new candidate skill repository
func NewCandidateSkillRepository(db *mongo.Database) *CandidateSkillRepository {
	return &CandidateSkillRepository{
		collection: db.Collection("candidateskills"),
	}
}

// GetAll retrieves all candidate skills with pagination and optional filtering
func (r *CandidateSkillRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.CandidateSkill, int64, error) {
	pagination := helpers.NewPagination(page, limit)

	// Build filter query
	filter := bson.M{}

	// Search by candidate_id (exact match)
	if candidateID, exists := filters["candidate_id"]; exists && candidateID != "" {
		objID, err := bson.ObjectIDFromHex(candidateID)
		if err == nil {
			filter["candidate_id"] = objID
		}
	}

	// Search by skill_id (exact match)
	if skillID, exists := filters["skill_id"]; exists && skillID != "" {
		objID, err := bson.ObjectIDFromHex(skillID)
		if err == nil {
			filter["skill_id"] = objID
		}
	}

	// Search by proficiency_level (partial match)
	if proficiencyLevel, exists := filters["proficiency_level"]; exists && proficiencyLevel != "" {
		filter["proficiency_level"] = bson.M{"$regex": proficiencyLevel, "$options": "i"}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sortableFields := []string{"candidate_id", "skill_id", "proficiency_level", "created_time"}
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

	var candidateSkills []models.CandidateSkill
	if err = cursor.All(ctx, &candidateSkills); err != nil {
		return nil, 0, err
	}

	return candidateSkills, total, nil
}

// GetByID retrieves a candidate skill by ID
func (r *CandidateSkillRepository) GetByID(ctx context.Context, id string) (*models.CandidateSkill, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var candidateSkill models.CandidateSkill
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&candidateSkill)
	if err != nil {
		return nil, err
	}
	return &candidateSkill, nil
}

// GetByCandidateID retrieves all skills for a specific candidate
func (r *CandidateSkillRepository) GetByCandidateID(ctx context.Context, candidateID string) ([]models.CandidateSkill, error) {
	objID, err := bson.ObjectIDFromHex(candidateID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"candidate_id": objID})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var candidateSkills []models.CandidateSkill
	if err = cursor.All(ctx, &candidateSkills); err != nil {
		return nil, err
	}

	return candidateSkills, nil
}

// Create inserts a new candidate skill
func (r *CandidateSkillRepository) Create(ctx context.Context, candidateSkill *models.CandidateSkill) error {
	result, err := r.collection.InsertOne(ctx, candidateSkill)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	candidateSkill.ID = objID
	return nil
}

// UpdateProficiencyLevel updates the proficiency level of a candidate skill
func (r *CandidateSkillRepository) UpdateProficiencyLevel(ctx context.Context, id string, proficiencyLevel string) error {
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
		bson.M{"$set": bson.M{"proficiency_level": proficiencyLevel}},
	)
	return err
}

// Delete removes a candidate skill by ID
func (r *CandidateSkillRepository) Delete(ctx context.Context, id string) error {
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
