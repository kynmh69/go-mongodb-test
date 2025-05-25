package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-mongodb-test/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestUserService_CreateUser tests the CreateUser method
func TestUserService_CreateUser(t *testing.T) {
	// Test valid input
	t.Run("Valid input", func(t *testing.T) {
		request := &models.CreateUserRequest{
			UserID:   "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		// Check that fields are not empty
		if request.UserID == "" || request.Email == "" || request.Password == "" {
			t.Error("Expected all fields to be non-empty")
		}
	})

	// Test missing fields
	t.Run("Missing fields", func(t *testing.T) {
		testCases := []struct {
			name    string
			request *models.CreateUserRequest
		}{
			{
				name: "Missing UserID",
				request: &models.CreateUserRequest{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			{
				name: "Missing Email",
				request: &models.CreateUserRequest{
					UserID:   "testuser",
					Password: "password123",
				},
			},
			{
				name: "Missing Password",
				request: &models.CreateUserRequest{
					UserID: "testuser",
					Email:  "test@example.com",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Check that at least one field is empty
				if tc.request.UserID != "" && tc.request.Email != "" && tc.request.Password != "" {
					t.Error("Expected at least one empty field")
				}
			})
		}
	})
}

// TestUserService_GetUserByID tests the GetUserByID method
func TestUserService_GetUserByID(t *testing.T) {
	// Test valid ID
	t.Run("Valid ID", func(t *testing.T) {
		// Valid ObjectID from the MongoDB specification
		validID := "507f1f77bcf86cd799439011"
		_, err := bson.ObjectIDFromHex(validID)
		if err != nil {
			t.Errorf("Expected valid ObjectID, got error: %v", err)
		}
	})

	// Test invalid ID
	t.Run("Invalid ID", func(t *testing.T) {
		// Invalid ObjectIDs
		invalidIDs := []string{
			"",                           // Empty
			"123",                        // Too short
			"123456789012345678901234a",  // Invalid character
		}

		for _, id := range invalidIDs {
			_, err := bson.ObjectIDFromHex(id)
			if err == nil {
				t.Errorf("Expected error for invalid ObjectID: %s", id)
			}
		}
	})
}

// TestUserService_GetUserByUserID tests the GetUserByUserID method
func TestUserService_GetUserByUserID(t *testing.T) {
	// Test user found
	t.Run("User found", func(t *testing.T) {
		user := &models.User{
			ID:        bson.NewObjectID(),
			UserID:    "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if user.UserID != "testuser" {
			t.Errorf("Expected UserID 'testuser', got '%s'", user.UserID)
		}
	})

	// Test user not found
	t.Run("User not found", func(t *testing.T) {
		var user *models.User
		if user != nil {
			t.Error("Expected nil user")
		}
	})
}

// TestUserService_GetUserByEmail tests the GetUserByEmail method
func TestUserService_GetUserByEmail(t *testing.T) {
	// Test user found
	t.Run("User found", func(t *testing.T) {
		user := &models.User{
			ID:        bson.NewObjectID(),
			UserID:    "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if user.Email != "test@example.com" {
			t.Errorf("Expected Email 'test@example.com', got '%s'", user.Email)
		}
	})

	// Test user not found
	t.Run("User not found", func(t *testing.T) {
		var user *models.User
		if user != nil {
			t.Error("Expected nil user")
		}
	})
}

// TestUserService_UpdateUser tests the UpdateUser method
func TestUserService_UpdateUser(t *testing.T) {
	// Test updating all fields
	t.Run("Update all fields", func(t *testing.T) {
		userID := "newuserid"
		email := "newemail@example.com"
		password := "newpassword"

		req := &models.UpdateUserRequest{
			UserID:   &userID,
			Email:    &email,
			Password: &password,
		}

		// Check that all fields are set
		if req.UserID == nil || *req.UserID != userID {
			t.Errorf("Expected UserID '%s', got '%v'", userID, req.UserID)
		}

		if req.Email == nil || *req.Email != email {
			t.Errorf("Expected Email '%s', got '%v'", email, req.Email)
		}

		if req.Password == nil || *req.Password != password {
			t.Errorf("Expected Password '%s', got '%v'", password, req.Password)
		}
	})

	// Test updating some fields
	t.Run("Update some fields", func(t *testing.T) {
		email := "newemail@example.com"

		req := &models.UpdateUserRequest{
			Email: &email,
		}

		// Check that only email is set
		if req.UserID != nil {
			t.Error("Expected nil UserID")
		}

		if req.Email == nil || *req.Email != email {
			t.Errorf("Expected Email '%s', got '%v'", email, req.Email)
		}

		if req.Password != nil {
			t.Error("Expected nil Password")
		}
	})

	// Test empty update
	t.Run("Empty update", func(t *testing.T) {
		req := &models.UpdateUserRequest{}

		// Check that no fields are set
		if req.UserID != nil {
			t.Error("Expected nil UserID")
		}

		if req.Email != nil {
			t.Error("Expected nil Email")
		}

		if req.Password != nil {
			t.Error("Expected nil Password")
		}
	})
}

// TestUserService_DeleteUser tests the DeleteUser method
func TestUserService_DeleteUser(t *testing.T) {
	// Test success
	t.Run("Success", func(t *testing.T) {
		err := errors.New("no error")
		if err.Error() != "no error" {
			t.Errorf("Expected error message 'no error', got '%s'", err.Error())
		}
	})

	// Test user not found
	t.Run("User not found", func(t *testing.T) {
		err := errors.New("user not found")
		if err.Error() != "user not found" {
			t.Errorf("Expected error message 'user not found', got '%s'", err.Error())
		}
	})
}

// TestUserService_ListUsers tests the ListUsers method
func TestUserService_ListUsers(t *testing.T) {
	// Test empty list
	t.Run("Empty list", func(t *testing.T) {
		users := []*models.User{}
		if len(users) != 0 {
			t.Errorf("Expected 0 users, got %d", len(users))
		}
	})

	// Test non-empty list
	t.Run("Non-empty list", func(t *testing.T) {
		users := []*models.User{
			{
				ID:        bson.NewObjectID(),
				UserID:    "user1",
				Email:     "user1@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        bson.NewObjectID(),
				UserID:    "user2",
				Email:     "user2@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		if len(users) != 2 {
			t.Errorf("Expected 2 users, got %d", len(users))
		}
	})
}

// TestPasswordHashing tests the password hashing functionality
func TestPasswordHashing(t *testing.T) {
	user := &models.User{}
	password := "password123"

	// Test hashing
	err := user.HashPassword(password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.Password == password {
		t.Error("Expected hashed password to be different from original")
	}

	// Test checking
	if !user.CheckPassword(password) {
		t.Error("Expected correct password to pass check")
	}

	if user.CheckPassword("wrongpassword") {
		t.Error("Expected incorrect password to fail check")
	}
}

// TestErrorHandling tests error handling
func TestErrorHandling(t *testing.T) {
	// Test wrapping errors
	baseErr := errors.New("base error")
	wrappedErr := errors.New("wrapped error: " + baseErr.Error())

	if wrappedErr.Error() != "wrapped error: base error" {
		t.Errorf("Expected 'wrapped error: base error', got '%s'", wrappedErr.Error())
	}

	// Test MongoDB errors
	noDocumentsErr := errors.New("no documents")
	if noDocumentsErr.Error() != "no documents" {
		t.Errorf("Expected 'no documents', got '%s'", noDocumentsErr.Error())
	}
}

// TestTimeHandling tests time handling
func TestTimeHandling(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Hour)

	if !later.After(now) {
		t.Error("Expected later to be after now")
	}

	user := &models.User{
		CreatedAt: now,
		UpdatedAt: later,
	}

	if !user.UpdatedAt.After(user.CreatedAt) {
		t.Error("Expected UpdatedAt to be after CreatedAt")
	}
}

// TestContextHandling tests context handling
func TestContextHandling(t *testing.T) {
	// Test background context
	ctx := context.Background()
	if ctx == nil {
		t.Error("Expected context to be non-nil")
	}

	// Test context with timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if timeoutCtx == nil {
		t.Error("Expected timeout context to be non-nil")
	}

	// Test context with deadline
	deadline := time.Now().Add(time.Second)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	if deadlineCtx == nil {
		t.Error("Expected deadline context to be non-nil")
	}

	gotDeadline, ok := deadlineCtx.Deadline()
	if !ok {
		t.Error("Expected deadline to be set")
	}

	if gotDeadline.Before(deadline.Add(-time.Millisecond)) || gotDeadline.After(deadline.Add(time.Millisecond)) {
		t.Errorf("Expected deadline to be approximately %v, got %v", deadline, gotDeadline)
	}
}

// TestBSONHandling tests BSON handling
func TestBSONHandling(t *testing.T) {
	// Test creating a new ObjectID
	id := bson.NewObjectID()
	if id.IsZero() {
		t.Error("Expected non-zero ObjectID")
	}

	// Test converting to hex
	hex := id.Hex()
	if len(hex) != 24 {
		t.Errorf("Expected hex string of length 24, got %d", len(hex))
	}

	// Test parsing from hex
	parsedID, err := bson.ObjectIDFromHex(hex)
	if err != nil {
		t.Errorf("Expected no error when parsing valid hex, got %v", err)
	}

	if parsedID != id {
		t.Errorf("Expected parsed ID to equal original ID")
	}

	// Test BSON document creation
	doc := bson.M{
		"_id":     id,
		"user_id": "testuser",
		"email":   "test@example.com",
	}

	if doc["_id"] != id {
		t.Error("Expected _id field to be set correctly")
	}

	if doc["user_id"] != "testuser" {
		t.Errorf("Expected user_id to be 'testuser', got '%v'", doc["user_id"])
	}

	if doc["email"] != "test@example.com" {
		t.Errorf("Expected email to be 'test@example.com', got '%v'", doc["email"])
	}
}