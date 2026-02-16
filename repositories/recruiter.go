package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RecruiterRepository struct {
	collection *mongo.Collection
}

// NewRecruiterRepository creates a new candidate repository
func NewRecruiterRepository(db *mongo.Database) *RecruiterRepository {
	return &RecruiterRepository{
		collection: db.Collection("recruiters"),
	}
}

// GetAll retrieves all recruiters with pagination and optional filtering
func (r *RecruiterRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Recruiter, int64, error) {
	pagination := helpers.NewPagination(page, limit)

	// Build filter query
	filter := bson.M{}
	for _, field := range []string{"first_name", "last_name", "email"} {
		if value, exists := filters[field]; exists && value != "" {
			// Case-insensitive partial match using regex
			filter[field] = bson.M{"$regex": value, "$options": "i"}
		}
	}

	// Handle company_id as exact match (not regex)
	if companyID, exists := filters["company_id"]; exists && companyID != "" {
		objID, err := bson.ObjectIDFromHex(companyID)
		if err == nil {
			filter["company_id"] = objID
		}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build sort options
	sortableFields := []string{"first_name", "last_name", "email", "company_id", "created_time"}
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

	var recruiters []models.Recruiter
	if err = cursor.All(ctx, &recruiters); err != nil {
		return nil, 0, err
	}

	return recruiters, total, nil
}

// GetByID retrieves a recruiter by ID
func (r *RecruiterRepository) GetByID(ctx context.Context, id string) (*models.Recruiter, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var recruiter models.Recruiter
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&recruiter)
	if err != nil {
		return nil, err
	}
	return &recruiter, nil
}

// Create inserts a new recruiter
func (r *RecruiterRepository) Create(ctx context.Context, recruiter *models.Recruiter) error {
	result, err := r.collection.InsertOne(ctx, recruiter)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	recruiter.ID = objID
	return nil
}

// Delete removes a recruiter by ID
func (r *RecruiterRepository) Delete(ctx context.Context, id string) error {
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
