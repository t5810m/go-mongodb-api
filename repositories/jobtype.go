package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JobTypeRepository struct {
	collection *mongo.Collection
}

func NewJobTypeRepository(db *mongo.Database) *JobTypeRepository {
	return &JobTypeRepository{
		collection: db.Collection("jobtypes"),
	}
}

func (r *JobTypeRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.JobType, int64, error) {
	pagination := helpers.NewPagination(page, limit)

	filter := bson.M{}
	if title, exists := filters["title"]; exists && title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sortableFields := []string{"title", "created_time"}
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

	var jobTypes []models.JobType
	if err = cursor.All(ctx, &jobTypes); err != nil {
		return nil, 0, err
	}

	return jobTypes, total, nil
}

func (r *JobTypeRepository) GetByID(ctx context.Context, id string) (*models.JobType, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var jobType models.JobType
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&jobType)
	if err != nil {
		return nil, err
	}
	return &jobType, nil
}

func (r *JobTypeRepository) Create(ctx context.Context, jobType *models.JobType) error {
	result, err := r.collection.InsertOne(ctx, jobType)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	jobType.ID = objID
	return nil
}

func (r *JobTypeRepository) Update(ctx context.Context, id string, jobType *models.JobType) (*models.JobType, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"title":        jobType.Title,
			"updated_time": jobType.UpdatedTime,
			"updated_by":   jobType.UpdatedBy,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.JobType
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *JobTypeRepository) Delete(ctx context.Context, id string) error {
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
