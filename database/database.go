package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	db     *mongo.Database
)

// Connect establishes a connection to MongoDB
func Connect() (*mongo.Database, error) {
	mongoURI := getMongoURI()
	dbName := getDatabaseName()

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Get database
	db = client.Database(dbName)
	return db, nil
}

// Close disconnects from MongoDB
func Close(ctx context.Context) error {
	if client != nil {
		return client.Disconnect(ctx)
	}
	return nil
}

// GetDB returns the MongoDB database instance
func GetDB() *mongo.Database {
	return db
}

// Helper function to get environment variables with default values
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getMongoURI returns the MongoDB URI from environment variables
func getMongoURI() string {
	return getEnvWithDefault("MONGODB_URI", "mongodb://localhost:27017")
}

// getDatabaseName returns the database name from environment variables
func getDatabaseName() string {
	return getEnvWithDefault("MONGODB_DB_NAME", "godb")
}