package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// EnsureIndexes creates all required indexes across every collection.
// It is idempotent: running it multiple times does not return an error for
// indexes that already exist.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	specs := []struct {
		collection string
		models     []mongo.IndexModel
	}{
		{
			collection: "users",
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "email", Value: 1}},
					Options: options.Index().SetUnique(true).SetName("email_unique"),
				},
				{
					Keys:    bson.D{{Key: "role", Value: 1}},
					Options: options.Index().SetName("role"),
				},
			},
		},
		{
			collection: "jobs",
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "user_id", Value: 1}},
					Options: options.Index().SetName("user_id"),
				},
				{
					Keys:    bson.D{{Key: "category_id", Value: 1}},
					Options: options.Index().SetName("category_id"),
				},
				{
					Keys:    bson.D{{Key: "status", Value: 1}},
					Options: options.Index().SetName("status"),
				},
				{
					Keys:    bson.D{{Key: "created_time", Value: -1}},
					Options: options.Index().SetName("created_time_desc"),
				},
			},
		},
		{
			collection: "applications",
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "job_id", Value: 1}},
					Options: options.Index().SetName("job_id"),
				},
				{
					Keys:    bson.D{{Key: "user_id", Value: 1}},
					Options: options.Index().SetName("user_id"),
				},
				{
					Keys:    bson.D{{Key: "status", Value: 1}},
					Options: options.Index().SetName("status"),
				},
				{
					Keys: bson.D{
						{Key: "job_id", Value: 1},
						{Key: "user_id", Value: 1},
					},
					Options: options.Index().SetUnique(true).SetName("job_user_unique"),
				},
			},
		},
		{
			collection: "skills",
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "name", Value: 1}},
					Options: options.Index().SetUnique(true).SetName("name_unique"),
				},
			},
		},
		{
			collection: "jobcategories",
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "name", Value: 1}},
					Options: options.Index().SetUnique(true).SetName("name_unique"),
				},
			},
		},
		{
			collection: "candidateskills",
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "user_id", Value: 1}},
					Options: options.Index().SetName("user_id"),
				},
				{
					Keys: bson.D{
						{Key: "user_id", Value: 1},
						{Key: "skill_id", Value: 1},
					},
					Options: options.Index().SetUnique(true).SetName("user_skill_unique"),
				},
			},
		},
		{
			collection: "jobskills",
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "job_id", Value: 1}},
					Options: options.Index().SetName("job_id"),
				},
				{
					Keys: bson.D{
						{Key: "job_id", Value: 1},
						{Key: "skill_id", Value: 1},
					},
					Options: options.Index().SetUnique(true).SetName("job_skill_unique"),
				},
			},
		},
	}

	for _, spec := range specs {
		col := db.Collection(spec.collection)
		if _, err := col.Indexes().CreateMany(ctx, spec.models); err != nil {
			return fmt.Errorf("indexes for %q: %w", spec.collection, err)
		}
	}

	return nil
}
