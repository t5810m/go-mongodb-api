package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ArticleRepository struct {
	collection *mongo.Collection
}

func NewArticleRepository(db *mongo.Database) *ArticleRepository {
	return &ArticleRepository{
		collection: db.Collection("articles"),
	}
}

func (r *ArticleRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Article, int64, error) {
	pagination := helpers.NewPagination(page, limit)

	filter := bson.M{}
	if name, exists := filters["name"]; exists && name != "" {
		filter["title"] = bson.M{"$regex": name, "$options": "i"}
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

	var articles []models.Article
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *ArticleRepository) GetByID(ctx context.Context, id string) (*models.Article, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var article models.Article
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&article)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepository) Create(ctx context.Context, article *models.Article) error {
	result, err := r.collection.InsertOne(ctx, article)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	article.ID = objID
	return nil
}

func (r *ArticleRepository) Update(ctx context.Context, id string, article *models.Article) (*models.Article, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"title":        article.Title,
			"content":      article.Content,
			"slug":         article.Slug,
			"active":       article.Active,
			"updated_time": article.UpdatedTime,
			"updated_by":   article.UpdatedBy,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Article
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *ArticleRepository) Delete(ctx context.Context, id string) error {
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
