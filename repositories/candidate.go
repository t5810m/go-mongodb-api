package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CandidateRepository struct {
	collection *mongo.Collection
}

// NewCandidateRepository creates a new candidate repository
func NewCandidateRepository(db *mongo.Database) *CandidateRepository {
	return &CandidateRepository{
		collection: db.Collection("candidates"),
	}
}

// GetAll retrieves all candidates with pagination, filtering, and sorting
func (r *CandidateRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Candidate, int64, error) {
	// Create pagination instance with validation
	pagination := helpers.NewPagination(page, limit)

	// Build filter query
	filter := bson.M{}
	searchableFields := []string{"first_name", "last_name", "email", "location"}
	for _, field := range searchableFields {
		if value, exists := filters[field]; exists && value != "" {
			// Case-insensitive partial match using regex
			filter[field] = bson.M{"$regex": value, "$options": "i"}
		}
	}

	// Count total documents matching filter
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build sort options
	sortableFields := []string{"first_name", "last_name", "email", "location", "created_time"}
	sortField := "created_time"
	if sort != "" {
		for _, field := range sortableFields {
			if field == sort {
				sortField = sort
				break
			}
		}
	}

	sortOrder := int32(-1) // desc
	if order == "asc" {
		sortOrder = 1
	}

	// Query with skip, limit, and sort
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

	var candidates []models.Candidate
	if err = cursor.All(ctx, &candidates); err != nil {
		return nil, 0, err
	}

	return candidates, total, nil
}

// GetByID retrieves a candidate by ID
func (r *CandidateRepository) GetByID(ctx context.Context, id string) (*models.Candidate, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var candidate models.Candidate
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&candidate)
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}

// Create inserts a new candidate
func (r *CandidateRepository) Create(ctx context.Context, candidate *models.Candidate) error {
	result, err := r.collection.InsertOne(ctx, candidate)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	candidate.ID = objID
	return nil
}

// Delete removes a candidate by ID
func (r *CandidateRepository) Delete(ctx context.Context, id string) error {
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
