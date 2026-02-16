package repositories

import (
	"context"
	"go-mongodb-api/helpers"
	"go-mongodb-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CompanyRepository struct {
	collection *mongo.Collection
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository(db *mongo.Database) *CompanyRepository {
	return &CompanyRepository{
		collection: db.Collection("companies"),
	}
}

// GetAll retrieves all companies with pagination and optional filtering
func (r *CompanyRepository) GetAll(ctx context.Context, page, limit int, filters map[string]string, sort, order string) ([]models.Company, int64, error) {
	pagination := helpers.NewPagination(page, limit)

	// Build filter query
	filter := bson.M{}

	// Search by name (partial match)
	if name, exists := filters["name"]; exists && name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// Search by location (searches country, city, postal_code)
	if location, exists := filters["location"]; exists && location != "" {
		filter["$or"] = bson.A{
			bson.M{"country": bson.M{"$regex": location, "$options": "i"}},
			bson.M{"city": bson.M{"$regex": location, "$options": "i"}},
			bson.M{"postal_code": bson.M{"$regex": location, "$options": "i"}},
		}
	}

	// Search by individual location fields
	if country, exists := filters["country"]; exists && country != "" {
		filter["country"] = bson.M{"$regex": country, "$options": "i"}
	}
	if city, exists := filters["city"]; exists && city != "" {
		filter["city"] = bson.M{"$regex": city, "$options": "i"}
	}
	if postalCode, exists := filters["postal_code"]; exists && postalCode != "" {
		filter["postal_code"] = bson.M{"$regex": postalCode, "$options": "i"}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sortableFields := []string{"name", "country", "city", "postal_code", "created_time"}
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

	var companies []models.Company
	if err = cursor.All(ctx, &companies); err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

// GetByID retrieves a company by ID
func (r *CompanyRepository) GetByID(ctx context.Context, id string) (*models.Company, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var company models.Company
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&company)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// Create inserts a new company
func (r *CompanyRepository) Create(ctx context.Context, company *models.Company) error {
	result, err := r.collection.InsertOne(ctx, company)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil
	}
	company.ID = objID
	return nil
}

// Delete removes a company by ID
func (r *CompanyRepository) Delete(ctx context.Context, id string) error {
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
