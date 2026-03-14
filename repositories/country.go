package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CountryRepository struct {
	collection *mongo.Collection
}

func NewCountryRepository(db *mongo.Database) *CountryRepository {
	return &CountryRepository{
		collection: db.Collection("countries"),
	}
}

func (r *CountryRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Country, int64, error) {
	pagination := helpers.NewPagination(page, limit)

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

	var countries []models.Country
	if err = cursor.All(ctx, &countries); err != nil {
		return nil, 0, err
	}

	return countries, total, nil
}

func (r *CountryRepository) GetByID(ctx context.Context, id string) (*models.Country, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var country models.Country
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&country)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *CountryRepository) Create(ctx context.Context, country *models.Country) error {
	result, err := r.collection.InsertOne(ctx, country)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	country.ID = objID
	return nil
}

func (r *CountryRepository) Update(ctx context.Context, id string, country *models.Country) (*models.Country, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"name":         country.Name,
			"updated_time": country.UpdatedTime,
			"updated_by":   country.UpdatedBy,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Country
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *CountryRepository) Delete(ctx context.Context, id string) error {
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
