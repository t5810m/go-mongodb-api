package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ApplicationRepository struct {
	collection *mongo.Collection
}

// NewApplicationRepository creates a new application repository
func NewApplicationRepository(db *mongo.Database) *ApplicationRepository {
	return &ApplicationRepository{
		collection: db.Collection("applications"),
	}
}

// GetAll retrieves all applications with pagination and optional filtering
func (r *ApplicationRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Application, int64, error) {
	pagination := helpers.NewPagination(page, limit)

	// Build filter query
	filter := bson.M{}

	// Search by status (partial match)
	if status, exists := filters["status"]; exists && status != "" {
		filter["status"] = bson.M{"$regex": status, "$options": "i"}
	}

	// Search by job_id (exact match)
	if jobID, exists := filters["job_id"]; exists && jobID != "" {
		objID, err := bson.ObjectIDFromHex(jobID)
		if err == nil {
			filter["job_id"] = objID
		}
	}

	// Search by candidate_id (exact match)
	if candidateID, exists := filters["candidate_id"]; exists && candidateID != "" {
		objID, err := bson.ObjectIDFromHex(candidateID)
		if err == nil {
			filter["candidate_id"] = objID
		}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sortableFields := []string{"status", "job_id", "candidate_id", "applied_time"}
	sortField := "applied_time"
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

	var applications []models.Application
	if err = cursor.All(ctx, &applications); err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

// GetByID retrieves an application by ID
func (r *ApplicationRepository) GetByID(ctx context.Context, id string) (*models.Application, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var application models.Application
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&application)
	if err != nil {
		return nil, err
	}
	return &application, nil
}

// GetByJobID retrieves all applications for a specific job
func (r *ApplicationRepository) GetByJobID(ctx context.Context, jobID string) ([]models.Application, error) {
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

	var applications []models.Application
	if err = cursor.All(ctx, &applications); err != nil {
		return nil, err
	}

	return applications, nil
}

// GetByCandidateID retrieves all applications from a specific candidate
func (r *ApplicationRepository) GetByCandidateID(ctx context.Context, candidateID string) ([]models.Application, error) {
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

	var applications []models.Application
	if err = cursor.All(ctx, &applications); err != nil {
		return nil, err
	}

	return applications, nil
}

// Create inserts a new application
func (r *ApplicationRepository) Create(ctx context.Context, application *models.Application) error {
	result, err := r.collection.InsertOne(ctx, application)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	application.ID = objID
	return nil
}

// UpdateStatus updates the status of an application
func (r *ApplicationRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if status == "" {
		return mongo.ErrNoDocuments
	}
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

// Delete removes an application by ID
func (r *ApplicationRepository) Delete(ctx context.Context, id string) error {
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
