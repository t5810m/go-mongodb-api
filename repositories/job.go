package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JobRepository struct {
	collection *mongo.Collection
}

// NewJobRepository creates a new job repository
func NewJobRepository(db *mongo.Database) *JobRepository {
	return &JobRepository{
		collection: db.Collection("jobs"),
	}
}

// GetAll retrieves all jobs with pagination, filtering, and sorting
func (r *JobRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Job, int64, error) {
	// Create pagination instance with validation
	pagination := helpers.NewPagination(page, limit)

	// Build filter query
	filter := bson.M{}
	searchableFields := []string{"title", "description", "location", "job_type", "status"}
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
	sortableFields := []string{"title", "description", "location", "job_type", "status", "created_time"}
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

	var jobs []models.Job
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

// GetByID retrieves a job by ID
func (r *JobRepository) GetByID(ctx context.Context, id string) (*models.Job, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var job models.Job
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetByCompanyID retrieves jobs by company ID
func (r *JobRepository) GetByCompanyID(ctx context.Context, companyID string) ([]models.Job, error) {
	objID, err := bson.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"company_id": objID})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var jobs []models.Job
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// Create inserts a new job
func (r *JobRepository) Create(ctx context.Context, job *models.Job) error {
	result, err := r.collection.InsertOne(ctx, job)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	job.ID = objID
	return nil
}

// Delete removes a job by ID
func (r *JobRepository) Delete(ctx context.Context, id string) error {
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
