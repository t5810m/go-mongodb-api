package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type KnowledgeLevelRepository struct {
	collection *mongo.Collection
}

func NewKnowledgeLevelRepository(db *mongo.Database) *KnowledgeLevelRepository {
	return &KnowledgeLevelRepository{
		collection: db.Collection("knowledgelevels"),
	}
}

func (r *KnowledgeLevelRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.KnowledgeLevel, int64, error) {
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

	var knowledgeLevels []models.KnowledgeLevel
	if err = cursor.All(ctx, &knowledgeLevels); err != nil {
		return nil, 0, err
	}

	return knowledgeLevels, total, nil
}

func (r *KnowledgeLevelRepository) GetByID(ctx context.Context, id string) (*models.KnowledgeLevel, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var knowledgeLevel models.KnowledgeLevel
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&knowledgeLevel)
	if err != nil {
		return nil, err
	}
	return &knowledgeLevel, nil
}

func (r *KnowledgeLevelRepository) Create(ctx context.Context, knowledgeLevel *models.KnowledgeLevel) error {
	result, err := r.collection.InsertOne(ctx, knowledgeLevel)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	knowledgeLevel.ID = objID
	return nil
}

func (r *KnowledgeLevelRepository) Update(ctx context.Context, id string, knowledgeLevel *models.KnowledgeLevel) (*models.KnowledgeLevel, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"title":        knowledgeLevel.Title,
			"updated_time": knowledgeLevel.UpdatedTime,
			"updated_by":   knowledgeLevel.UpdatedBy,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.KnowledgeLevel
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *KnowledgeLevelRepository) Delete(ctx context.Context, id string) error {
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
