package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client

// InitMongo initializes MongoDB connection using global config
func InitMongo() (*mongo.Client, error) {
	cfg, err := Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// Set up MongoDB client options with timeout from config
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI).SetTimeout(cfg.Timeout))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	MongoClient = client
	log.Printf("successfully connected to MongoDB")
	return client, nil
}

// GetDatabase returns MongoDB database instance, returns error if not initialized or invalid dbName
func GetDatabase(dbName string) (*mongo.Database, error) {
	if MongoClient == nil {
		return nil, fmt.Errorf("MongoDB client not initialized: call InitMongo() first")
	}
	if dbName == "" {
		return nil, fmt.Errorf("database name cannot be empty")
	}
	return MongoClient.Database(dbName), nil
}

// DisconnectMongo closes MongoDB connection
func DisconnectMongo() error {
	if MongoClient == nil {
		log.Printf("info: DisconnectMongo called but client not initialized")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := MongoClient.Disconnect(ctx); err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %w", err)
	}
	log.Printf("successfully disconnected from MongoDB")
	return nil
}
