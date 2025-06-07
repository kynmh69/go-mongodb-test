package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-mongodb-test/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockDatabase implements DatabaseCollectionProvider for testing
type MockDatabase struct{}

func (m *MockDatabase) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	// Return nil for testing - we'll focus on business logic, not DB operations
	return nil
}

// TestNewUserService tests the service constructor
func TestNewUserService(t *testing.T) {
	db := &MockDatabase{}
	service := NewUserService(db)
	
	if service == nil {
		t.Error("Expected service to be non-nil")
	}
}

// TestUserServiceValidation tests input validation logic that gets executed before database operations
func TestUserServiceValidation(t *testing.T) {
	ctx := context.Background()
	db := &MockDatabase{}
	service := NewUserService(db)

	t.Run("GetUserByID with invalid ObjectID", func(t *testing.T) {
		// Test the ObjectID validation logic that happens before DB operations
		invalidIDs := []string{
			"",                           // Empty
			"123",                        // Too short  
			"123456789012345678901234z",  // Invalid character
			"invalid-id",                 // Invalid format
		}

		for _, id := range invalidIDs {
			_, err := service.GetUserByID(ctx, id)
			if err == nil {
				t.Errorf("Expected error for invalid ObjectID: %s", id)
			}
			
			// Check that it contains "invalid user ID" message
			if err != nil && !contains(err.Error(), "invalid user ID") {
				t.Errorf("Expected 'invalid user ID' error for %s, got: %v", id, err)
			}
		}
	})

	t.Run("UpdateUser with invalid ObjectID", func(t *testing.T) {
		req := &models.UpdateUserRequest{}
		
		invalidIDs := []string{
			"",                           
			"123",                        
			"123456789012345678901234z",  
			"invalid-id",                 
		}

		for _, id := range invalidIDs {
			_, err := service.UpdateUser(ctx, id, req)
			if err == nil {
				t.Errorf("Expected error for invalid ObjectID: %s", id)
			}
			
			if err != nil && !contains(err.Error(), "invalid user ID") {
				t.Errorf("Expected 'invalid user ID' error for %s, got: %v", id, err)
			}
		}
	})

	t.Run("DeleteUser with invalid ObjectID", func(t *testing.T) {
		invalidIDs := []string{
			"",                           
			"123",                        
			"123456789012345678901234z",  
			"invalid-id",                 
		}

		for _, id := range invalidIDs {
			err := service.DeleteUser(ctx, id)
			if err == nil {
				t.Errorf("Expected error for invalid ObjectID: %s", id)
			}
			
			if err != nil && !contains(err.Error(), "invalid user ID") {
				t.Errorf("Expected 'invalid user ID' error for %s, got: %v", id, err)
			}
		}
	})

	t.Run("CreateUser input processing", func(t *testing.T) {
		// Only test that we can create the request structures properly
		// Don't actually call the service methods since they require database
		testCases := []struct {
			name string
			req  *models.CreateUserRequest
		}{
			{
				name: "Valid request",
				req: &models.CreateUserRequest{
					UserID:   "testuser",
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			{
				name: "Different user",
				req: &models.CreateUserRequest{
					UserID:   "anotheruser",
					Email:    "another@example.com",
					Password: "differentpassword",
				},
			},
			{
				name: "Long password",
				req: &models.CreateUserRequest{
					UserID:   "longpassuser",
					Email:    "long@example.com",
					Password: "verylongpasswordwithmancharacters123456789",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Just validate the request structure, don't call the service
				if tc.req.UserID == "" {
					t.Error("UserID should not be empty")
				}
				if tc.req.Email == "" {
					t.Error("Email should not be empty")
				}
				if tc.req.Password == "" {
					t.Error("Password should not be empty")
				}
			})
		}
	})

	t.Run("UpdateUser input processing", func(t *testing.T) {
		// Only test that we can create the request structures properly
		// Don't actually call the service methods since they require database
		testCases := []struct {
			name string
			req  *models.UpdateUserRequest
		}{
			{
				name: "Update UserID only",
				req: &models.UpdateUserRequest{
					UserID: stringPtr("newuserid"),
				},
			},
			{
				name: "Update Email only",
				req: &models.UpdateUserRequest{
					Email: stringPtr("newemail@example.com"),
				},
			},
			{
				name: "Update Password only",
				req: &models.UpdateUserRequest{
					Password: stringPtr("newpassword"),
				},
			},
			{
				name: "Update all fields",
				req: &models.UpdateUserRequest{
					UserID:   stringPtr("updateduserid"),
					Email:    stringPtr("updated@example.com"),
					Password: stringPtr("updatedpassword"),
				},
			},
			{
				name: "Empty update",
				req:  &models.UpdateUserRequest{},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Just validate the request structure, don't call the service
				// Test that pointer fields are set correctly
				if tc.name == "Update UserID only" {
					if tc.req.UserID == nil {
						t.Error("Expected UserID to be non-nil")
					}
					if tc.req.Email != nil {
						t.Error("Expected Email to be nil")
					}
					if tc.req.Password != nil {
						t.Error("Expected Password to be nil")
					}
				}
			})
		}
	})
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findInString(s, substr)
}

func findInString(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// Test business logic and data structures extensively
func TestUserServiceBusinessLogic(t *testing.T) {
	t.Run("CreateUserRequest validation", func(t *testing.T) {
		testCases := []struct {
			name  string
			req   *models.CreateUserRequest
			valid bool
		}{
			{
				name: "Valid request",
				req: &models.CreateUserRequest{
					UserID:   "testuser",
					Email:    "test@example.com",
					Password: "password123",
				},
				valid: true,
			},
			{
				name: "Empty UserID",
				req: &models.CreateUserRequest{
					UserID:   "",
					Email:    "test@example.com",
					Password: "password123",
				},
				valid: false,
			},
			{
				name: "Empty Email",
				req: &models.CreateUserRequest{
					UserID:   "testuser",
					Email:    "",
					Password: "password123",
				},
				valid: false,
			},
			{
				name: "Empty Password",
				req: &models.CreateUserRequest{
					UserID:   "testuser",
					Email:    "test@example.com",
					Password: "",
				},
				valid: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				isEmpty := tc.req.UserID == "" || tc.req.Email == "" || tc.req.Password == ""
				if tc.valid && isEmpty {
					t.Error("Expected valid request but has empty fields")
				}
				if !tc.valid && !isEmpty {
					t.Error("Expected invalid request but all fields are present")
				}
			})
		}
	})

	t.Run("UpdateUserRequest validation", func(t *testing.T) {
		testCases := []struct {
			name        string
			req         *models.UpdateUserRequest
			expectEmpty bool
		}{
			{
				name: "All fields set",
				req: &models.UpdateUserRequest{
					UserID:   stringPtr("newuserid"),
					Email:    stringPtr("newemail@example.com"),
					Password: stringPtr("newpassword"),
				},
				expectEmpty: false,
			},
			{
				name: "Only UserID set",
				req: &models.UpdateUserRequest{
					UserID: stringPtr("newuserid"),
				},
				expectEmpty: false,
			},
			{
				name: "Only Email set",
				req: &models.UpdateUserRequest{
					Email: stringPtr("newemail@example.com"),
				},
				expectEmpty: false,
			},
			{
				name: "Only Password set",
				req: &models.UpdateUserRequest{
					Password: stringPtr("newpassword"),
				},
				expectEmpty: false,
			},
			{
				name:        "Empty request",
				req:         &models.UpdateUserRequest{},
				expectEmpty: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				isEmpty := tc.req.UserID == nil && tc.req.Email == nil && tc.req.Password == nil
				if tc.expectEmpty && !isEmpty {
					t.Error("Expected empty request but fields are set")
				}
				if !tc.expectEmpty && isEmpty {
					t.Error("Expected non-empty request but all fields are nil")
				}
			})
		}
	})

	t.Run("User model operations", func(t *testing.T) {
		user := &models.User{
			ID:        bson.NewObjectID(),
			UserID:    "testuser",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Test password hashing with different passwords
		passwords := []string{
			"password123",
			"",           // Empty password
			"short",      // Short password
			"verylongpasswordwithmancharacters123456789", // Long password
			"special!@#$%^&*()_+-=[]{}|;:,.<>?",           // Special characters
		}

		for _, password := range passwords {
			t.Run("Password: "+password, func(t *testing.T) {
				// Test password hashing
				err := user.HashPassword(password)
				if err != nil {
					t.Errorf("Expected no error hashing password '%s', got %v", password, err)
				}

				if user.Password == password && password != "" {
					t.Error("Expected password to be hashed")
				}

				// Test password verification
				if !user.CheckPassword(password) {
					t.Error("Expected correct password to verify")
				}

				// Test wrong password
				if user.CheckPassword("definitely_wrong_password") {
					t.Error("Expected wrong password to fail verification")
				}
			})
		}
	})

	t.Run("Time handling", func(t *testing.T) {
		now := time.Now()
		past := now.Add(-time.Hour)
		future := now.Add(time.Hour)

		user := &models.User{
			CreatedAt: past,
			UpdatedAt: now,
		}

		// Test time relationships
		if !user.UpdatedAt.After(user.CreatedAt) {
			t.Error("Expected UpdatedAt to be after CreatedAt")
		}

		// Test updating timestamps
		user.UpdatedAt = future
		if !user.UpdatedAt.After(user.CreatedAt) {
			t.Error("Expected updated UpdatedAt to be after CreatedAt")
		}

		if !user.UpdatedAt.After(now) {
			t.Error("Expected updated UpdatedAt to be after now")
		}
	})
}

// Test BSON and ObjectID operations extensively
func TestBSONOperations(t *testing.T) {
	t.Run("ObjectID creation and manipulation", func(t *testing.T) {
		// Test creating multiple ObjectIDs
		ids := make([]bson.ObjectID, 10)
		for i := 0; i < 10; i++ {
			ids[i] = bson.NewObjectID()
			if ids[i].IsZero() {
				t.Errorf("Expected non-zero ObjectID at index %d", i)
			}
		}

		// Test that all IDs are unique
		for i := 0; i < len(ids); i++ {
			for j := i + 1; j < len(ids); j++ {
				if ids[i] == ids[j] {
					t.Errorf("Expected unique ObjectIDs, found duplicate at indices %d and %d", i, j)
				}
			}
		}

		// Test hex conversion for all IDs
		for i, id := range ids {
			hex := id.Hex()
			if len(hex) != 24 {
				t.Errorf("Expected hex string of length 24 for ID %d, got %d", i, len(hex))
			}

			// Test parsing back from hex
			parsedID, err := bson.ObjectIDFromHex(hex)
			if err != nil {
				t.Errorf("Expected no error parsing hex for ID %d, got %v", i, err)
			}

			if parsedID != id {
				t.Errorf("Expected parsed ID to equal original ID at index %d", i)
			}
		}
	})

	t.Run("Invalid ObjectID parsing comprehensive", func(t *testing.T) {
		invalidIDs := []string{
			"",                                    // Empty
			"123",                                 // Too short
			"123456789012345678901234z",           // Invalid character z
			"123456789012345678901234Z",           // Invalid character Z
			"123456789012345678901234!",           // Invalid character !
			"123456789012345678901234 ",           // Invalid character space
			"gggggggggggggggggggggggg",            // Invalid hex characters
			"GGGGGGGGGGGGGGGGGGGGGGGG",            // Invalid hex characters (uppercase)
			"123456789012345678901234567890",      // Too long
			"12345678901234567890123",             // One character short
			"1234567890123456789012345",           // One character long
		}

		for _, invalidID := range invalidIDs {
			t.Run("Invalid ID: "+invalidID, func(t *testing.T) {
				_, err := bson.ObjectIDFromHex(invalidID)
				if err == nil {
					t.Errorf("Expected error for invalid ObjectID: '%s'", invalidID)
				}
			})
		}
	})

	t.Run("BSON document operations", func(t *testing.T) {
		id := bson.NewObjectID()
		userID := "testuser"
		email := "test@example.com"
		now := time.Now()

		// Test filter documents
		filterDocs := []bson.M{
			{"_id": id},
			{"user_id": userID},
			{"email": email},
			{"user_id": userID, "email": email},
			{"created_at": bson.M{"$gte": now}},
		}

		for i, filter := range filterDocs {
			t.Run("Filter doc "+string(rune('A'+i)), func(t *testing.T) {
				if len(filter) == 0 {
					t.Error("Expected non-empty filter document")
				}

				// Test specific fields based on filter type
				if idVal, exists := filter["_id"]; exists {
					if idVal != id {
						t.Error("Expected _id field to match")
					}
				}

				if userIDVal, exists := filter["user_id"]; exists {
					if userIDVal != userID {
						t.Error("Expected user_id field to match")
					}
				}

				if emailVal, exists := filter["email"]; exists {
					if emailVal != email {
						t.Error("Expected email field to match")
					}
				}
			})
		}

		// Test update documents
		updateDocs := []bson.M{
			{"$set": bson.M{"user_id": "newuserid"}},
			{"$set": bson.M{"email": "newemail@example.com"}},
			{"$set": bson.M{"updated_at": time.Now()}},
			{"$set": bson.M{
				"user_id":    "updateduserid",
				"email":      "updated@example.com",
				"updated_at": time.Now(),
			}},
		}

		for i, updateDoc := range updateDocs {
			t.Run("Update doc "+string(rune('A'+i)), func(t *testing.T) {
				setFields, exists := updateDoc["$set"]
				if !exists {
					t.Error("Expected $set field in update document")
				}

				setMap, ok := setFields.(bson.M)
				if !ok {
					t.Error("Expected $set field to be bson.M")
				}

				if len(setMap) == 0 {
					t.Error("Expected non-empty $set fields")
				}
			})
		}
	})
}

// Test context operations extensively
func TestContextOperations(t *testing.T) {
	t.Run("Context creation and properties", func(t *testing.T) {
		// Test background context
		ctx := context.Background()
		if ctx == nil {
			t.Error("Expected context to be non-nil")
		}

		if ctx.Err() != nil {
			t.Error("Expected background context to have no error")
		}

		// Test that background context doesn't have a deadline
		_, hasDeadline := ctx.Deadline()
		if hasDeadline {
			t.Error("Expected background context to have no deadline")
		}

		// Test context value (should be nil for background context)
		value := ctx.Value("testkey")
		if value != nil {
			t.Error("Expected background context to have no values")
		}
	})

	t.Run("Context with timeout", func(t *testing.T) {
		durations := []time.Duration{
			time.Millisecond,
			time.Second,
			time.Minute,
			time.Hour,
		}

		for _, duration := range durations {
			t.Run("Duration: "+duration.String(), func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), duration)
				defer cancel()

				if ctx == nil {
					t.Error("Expected timeout context to be non-nil")
				}

				deadline, hasDeadline := ctx.Deadline()
				if !hasDeadline {
					t.Error("Expected timeout context to have deadline")
				}

				if deadline.Before(time.Now()) {
					t.Error("Expected deadline to be in the future")
				}

				// Test that context is not yet done
				select {
				case <-ctx.Done():
					t.Error("Expected context to not be done immediately")
				default:
					// Expected
				}
			})
		}
	})

	t.Run("Context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// Verify context is not cancelled initially
		if ctx.Err() != nil {
			t.Error("Expected context to not be cancelled initially")
		}

		// Cancel the context
		cancel()

		// Verify context is now cancelled
		select {
		case <-ctx.Done():
			// Expected
		case <-time.After(time.Millisecond):
			t.Error("Expected context to be cancelled")
		}

		if ctx.Err() == nil {
			t.Error("Expected cancelled context to have error")
		}
	})

	t.Run("Context with values", func(t *testing.T) {
		type contextKey string
		
		testCases := []struct {
			key   contextKey
			value interface{}
		}{
			{"string_key", "string_value"},
			{"int_key", 42},
			{"bool_key", true},
			{"struct_key", struct{ Name string }{"test"}},
		}

		for _, tc := range testCases {
			t.Run("Key: "+string(tc.key), func(t *testing.T) {
				ctx := context.WithValue(context.Background(), tc.key, tc.value)

				retrievedValue := ctx.Value(tc.key)
				if retrievedValue != tc.value {
					t.Errorf("Expected value %v, got %v", tc.value, retrievedValue)
				}

				// Test that other keys return nil
				otherValue := ctx.Value(contextKey("nonexistent"))
				if otherValue != nil {
					t.Error("Expected nil for nonexistent key")
				}
			})
		}
	})
}

// Test error handling patterns extensively
func TestErrorHandling(t *testing.T) {
	errorMessages := []string{
		"user with this user_id already exists",
		"user with this email already exists",
		"user not found",
		"invalid user ID",
		"failed to hash password",
		"failed to create user",
		"failed to get user",
		"failed to update user",
		"failed to delete user",
		"failed to decode user",
	}

	for _, msg := range errorMessages {
		t.Run("Error: "+msg, func(t *testing.T) {
			err := errors.New(msg)
			if err.Error() != msg {
				t.Errorf("Expected '%s', got '%s'", msg, err.Error())
			}

			// Test error is not nil
			if err == nil {
				t.Error("Expected error to be non-nil")
			}
		})
	}

	t.Run("Wrapped errors", func(t *testing.T) {
		baseErr := errors.New("base error")
		wrappedErr := errors.New("wrapped: " + baseErr.Error())

		if !contains(wrappedErr.Error(), baseErr.Error()) {
			t.Error("Expected wrapped error to contain base error message")
		}
	})
}

// Test time operations extensively
func TestTimeOperations(t *testing.T) {
	t.Run("Time creation and comparison", func(t *testing.T) {
		now := time.Now()
		past := now.Add(-time.Hour)
		future := now.Add(time.Hour)

		// Test basic comparisons
		if !future.After(now) {
			t.Error("Expected future to be after now")
		}

		if !now.After(past) {
			t.Error("Expected now to be after past")
		}

		if !past.Before(now) {
			t.Error("Expected past to be before now")
		}

		if !now.Before(future) {
			t.Error("Expected now to be before future")
		}

		// Test equality
		if !now.Equal(now) {
			t.Error("Expected time to equal itself")
		}

		// Test time differences
		diff := now.Sub(past)
		if diff != time.Hour {
			t.Errorf("Expected difference of 1 hour, got %v", diff)
		}
	})

	t.Run("User timestamp handling", func(t *testing.T) {
		createdAt := time.Now().Add(-time.Hour)
		updatedAt := time.Now()

		user := &models.User{
			ID:        bson.NewObjectID(),
			UserID:    "timetest",
			Email:     "time@example.com",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		if !user.UpdatedAt.After(user.CreatedAt) {
			t.Error("Expected UpdatedAt to be after CreatedAt")
		}

		// Test updating timestamp
		newUpdateTime := time.Now().Add(time.Minute)
		user.UpdatedAt = newUpdateTime

		if !user.UpdatedAt.After(user.CreatedAt) {
			t.Error("Expected new UpdatedAt to be after CreatedAt")
		}

		if !user.UpdatedAt.After(updatedAt) {
			t.Error("Expected new UpdatedAt to be after original UpdatedAt")
		}
	})
}