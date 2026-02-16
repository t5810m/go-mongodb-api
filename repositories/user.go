package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// GetAll retrieves all users with pagination and optional filtering
func (r *UserRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.User, int64, error) {
	pagination := helpers.NewPagination(page, limit)

	// Build filter query
	filter := bson.M{}
	for _, field := range []string{"first_name", "last_name", "email"} {
		if value, exists := filters[field]; exists && value != "" {
			filter[field] = bson.M{"$regex": value, "$options": "i"}
		}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sortableFields := []string{"first_name", "last_name", "email", "verified", "active", "created_time"}
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

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create inserts a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	user.ID = objID
	return nil
}

// Delete removes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
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
