package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LocationAvailabilityRepository struct {
	collection *mongo.Collection
}

func NewLocationAvailabilityRepository(db *mongo.Database) *LocationAvailabilityRepository {
	return &LocationAvailabilityRepository{
		collection: db.Collection("locationavailabilities"),
	}
}

func (r *LocationAvailabilityRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.LocationAvailability, int64, error) {
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

	var locationAvailabilities []models.LocationAvailability
	if err = cursor.All(ctx, &locationAvailabilities); err != nil {
		return nil, 0, err
	}

	return locationAvailabilities, total, nil
}

func (r *LocationAvailabilityRepository) GetByID(ctx context.Context, id string) (*models.LocationAvailability, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var locationAvailability models.LocationAvailability
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&locationAvailability)
	if err != nil {
		return nil, err
	}
	return &locationAvailability, nil
}

func (r *LocationAvailabilityRepository) Create(ctx context.Context, locationAvailability *models.LocationAvailability) error {
	result, err := r.collection.InsertOne(ctx, locationAvailability)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	locationAvailability.ID = objID
	return nil
}

func (r *LocationAvailabilityRepository) Update(ctx context.Context, id string, locationAvailability *models.LocationAvailability) (*models.LocationAvailability, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"title":        locationAvailability.Title,
			"updated_time": locationAvailability.UpdatedTime,
			"updated_by":   locationAvailability.UpdatedBy,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.LocationAvailability
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *LocationAvailabilityRepository) Delete(ctx context.Context, id string) error {
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
