package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-mongodb-test/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseCollectionProvider interface for database operations
type DatabaseCollectionProvider interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(db DatabaseCollectionProvider) *UserService {
	return &UserService{
		collection: db.Collection("users"),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, _ := s.GetUserByUserID(ctx, req.UserID)
	if existingUser != nil {
		return nil, errors.New("user with this user_id already exists")
	}

	existingUser, _ = s.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	user := &models.User{
		UserID:    req.UserID,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = result.InsertedID.(bson.ObjectID)
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var user models.User
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *UserService) GetUserByUserID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	updateFields := bson.M{
		"updated_at": time.Now(),
	}

	if req.UserID != nil {
		// Check if the new user_id is already taken
		existingUser, _ := s.GetUserByUserID(ctx, *req.UserID)
		if existingUser != nil && existingUser.ID != objectID {
			return nil, errors.New("user with this user_id already exists")
		}
		updateFields["user_id"] = *req.UserID
	}

	if req.Email != nil {
		// Check if the new email is already taken
		existingUser, _ := s.GetUserByEmail(ctx, *req.Email)
		if existingUser != nil && existingUser.ID != objectID {
			return nil, errors.New("user with this email already exists")
		}
		updateFields["email"] = *req.Email
	}

	if req.Password != nil {
		user := &models.User{}
		if err := user.HashPassword(*req.Password); err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		updateFields["password"] = user.Password
	}

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.GetUserByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}
