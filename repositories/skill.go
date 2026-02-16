package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SkillRepository struct {
	collection *mongo.Collection
}

// NewSkillRepository creates a new skill repository
func NewSkillRepository(db *mongo.Database) *SkillRepository {
	return &SkillRepository{
		collection: db.Collection("skills"),
	}
}

// GetAll retrieves all skills with pagination and optional filtering
func (r *SkillRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Skill, int64, error) {
	pagination := helpers.NewPagination(page, limit)

	// Build filter query
	filter := bson.M{}
	if name, exists := filters["name"]; exists && name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sortableFields := []string{"name", "created_time"}
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

	var skills []models.Skill
	if err = cursor.All(ctx, &skills); err != nil {
		return nil, 0, err
	}

	return skills, total, nil
}

// GetByID retrieves a skill by ID
func (r *SkillRepository) GetByID(ctx context.Context, id string) (*models.Skill, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var skill models.Skill
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&skill)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

// Create inserts a new skill
func (r *SkillRepository) Create(ctx context.Context, skill *models.Skill) error {
	result, err := r.collection.InsertOne(ctx, skill)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	skill.ID = objID
	return nil
}

// Delete removes a skill by ID
func (r *SkillRepository) Delete(ctx context.Context, id string) error {
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
