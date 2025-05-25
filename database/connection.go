package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const AuthSource = "admin"

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewConnection() (*Database, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "user_management"
	}

	dbUser := os.Getenv("MONGODB_USER")
	if dbUser == "" {
		dbUser = "admin"
	}

	dbPassword := os.Getenv("MONGODB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}

	credential := options.Credential{
		Username:   dbUser,
		Password:   dbPassword,
		AuthSource: AuthSource,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
  
	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(mongoURI).
			SetAuth(credential),
	)
  
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("Connected to MongoDB at %s", mongoURI)

	return &Database{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

func (d *Database) Close() error {
	if d.Client == nil {
		return errors.New("client is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return d.Client.Disconnect(ctx)
}
