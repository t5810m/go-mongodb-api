package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JobCategoryRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

// NewJobCategoryRepository creates a new job category repository
func NewJobCategoryRepository(db *mongo.Database) *JobCategoryRepository {
	return &JobCategoryRepository{
		collection: db.Collection("jobcategories"),
		db:         db,
	}
}

// GetAll retrieves all JobCategories with pagination and optional filtering
func (r *JobCategoryRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobCategory, int64, error) {
	pagination := helpers.NewPagination(page, limit)

	// Build filter query
	filter := bson.M{}

	// Search by name (partial match)
	if name, exists := filters["name"]; exists && name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// Search by description (partial match)
	if description, exists := filters["description"]; exists && description != "" {
		filter["description"] = bson.M{"$regex": description, "$options": "i"}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sortableFields := []string{"name", "description", "created_time"}
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

	var jobCategories []models.JobCategory
	if err = cursor.All(ctx, &jobCategories); err != nil {
		return nil, 0, err
	}

	return jobCategories, total, nil
}

// GetByID retrieves a job category by ID
func (r *JobCategoryRepository) GetByID(ctx context.Context, id string) (*models.JobCategory, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var jobCategory models.JobCategory
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&jobCategory)
	if err != nil {
		return nil, err
	}
	return &jobCategory, nil
}

// Create inserts a new job category
func (r *JobCategoryRepository) Create(ctx context.Context, jobCategory *models.JobCategory) error {
	result, err := r.collection.InsertOne(ctx, jobCategory)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	jobCategory.ID = objID
	return nil
}

// Delete removes a job category by ID only if no jobs use it
func (r *JobCategoryRepository) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	jobsCollection := r.db.Collection("jobs")
	jobCount, err := jobsCollection.CountDocuments(ctx, bson.M{"category_id": objID})
	if err != nil {
		return err
	}

	if jobCount > 0 {
		return mongo.ErrNoDocuments
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
