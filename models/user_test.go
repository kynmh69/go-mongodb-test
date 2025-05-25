package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func TestUser_HashPassword(t *testing.T) {
	user := &User{}
	password := "testpassword123"

	err := user.HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Password == "" {
		t.Fatal("Expected password to be hashed, got empty string")
	}

	if user.Password == password {
		t.Fatal("Expected password to be hashed, got original password")
	}

	// Verify the password is properly hashed
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		t.Fatalf("Expected hashed password to match original, got error: %v", err)
	}
}

func TestUser_CheckPassword(t *testing.T) {
	user := &User{}
	password := "testpassword123"

	// Hash the password first
	err := user.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	if !user.CheckPassword(password) {
		t.Fatal("Expected password check to return true for correct password")
	}

	// Test incorrect password
	if user.CheckPassword("wrongpassword") {
		t.Fatal("Expected password check to return false for incorrect password")
	}

	// Test empty password
	if user.CheckPassword("") {
		t.Fatal("Expected password check to return false for empty password")
	}
}

func TestCreateUserRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateUserRequest
		valid   bool
	}{
		{
			name: "Valid request",
			request: CreateUserRequest{
				UserID:   "test123",
				Email:    "test@example.com",
				Password: "password123",
			},
			valid: true,
		},
		{
			name: "Missing UserID",
			request: CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "Missing Email",
			request: CreateUserRequest{
				UserID:   "test123",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "Missing Password",
			request: CreateUserRequest{
				UserID: "test123",
				Email:  "test@example.com",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEmpty := tt.request.UserID == "" || tt.request.Email == "" || tt.request.Password == ""
			if tt.valid && isEmpty {
				t.Error("Expected valid request but required fields are empty")
			}
			if !tt.valid && !isEmpty {
				t.Error("Expected invalid request but all required fields are present")
			}
		})
	}
}

func TestUpdateUserRequest_FieldUpdates(t *testing.T) {
	userID := "newuserid"
	email := "newemail@example.com"
	password := "newpassword"

	req := UpdateUserRequest{
		UserID:   &userID,
		Email:    &email,
		Password: &password,
	}

	if req.UserID == nil || *req.UserID != userID {
		t.Error("Expected UserID to be set")
	}

	if req.Email == nil || *req.Email != email {
		t.Error("Expected Email to be set")
	}

	if req.Password == nil || *req.Password != password {
		t.Error("Expected Password to be set")
	}

	// Test partial update
	partialReq := UpdateUserRequest{
		Email: &email,
	}

	if partialReq.UserID != nil {
		t.Error("Expected UserID to be nil in partial update")
	}

	if partialReq.Email == nil || *partialReq.Email != email {
		t.Error("Expected Email to be set in partial update")
	}

	if partialReq.Password != nil {
		t.Error("Expected Password to be nil in partial update")
	}
}

func TestUser_StructFields(t *testing.T) {
	id := primitive.NewObjectID()
	userID := "test123"
	email := "test@example.com"
	password := "hashedpassword"
	now := time.Now()

	user := User{
		ID:        id,
		UserID:    userID,
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if user.ID != id {
		t.Errorf("Expected ID %v, got %v", id, user.ID)
	}

	if user.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, user.UserID)
	}

	if user.Email != email {
		t.Errorf("Expected Email %s, got %s", email, user.Email)
	}

	if user.Password != password {
		t.Errorf("Expected Password %s, got %s", password, user.Password)
	}

	if !user.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt %v, got %v", now, user.CreatedAt)
	}

	if !user.UpdatedAt.Equal(now) {
		t.Errorf("Expected UpdatedAt %v, got %v", now, user.UpdatedAt)
	}
}