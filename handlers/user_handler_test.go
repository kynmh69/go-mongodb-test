package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-mongodb-test/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock UserService for testing
type mockUserService struct {
	createUserFunc     func(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	getUserByIDFunc    func(ctx context.Context, id string) (*models.User, error)
	getUserByUserIDFunc func(ctx context.Context, userID string) (*models.User, error)
	getUserByEmailFunc func(ctx context.Context, email string) (*models.User, error)
	updateUserFunc     func(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error)
	deleteUserFunc     func(ctx context.Context, id string) error
	listUsersFunc      func(ctx context.Context) ([]*models.User, error)
}

func (m *mockUserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, req)
	}
	return nil, errors.New("CreateUser not implemented")
}

func (m *mockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, errors.New("GetUserByID not implemented")
}

func (m *mockUserService) GetUserByUserID(ctx context.Context, userID string) (*models.User, error) {
	if m.getUserByUserIDFunc != nil {
		return m.getUserByUserIDFunc(ctx, userID)
	}
	return nil, errors.New("GetUserByUserID not implemented")
}

func (m *mockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, email)
	}
	return nil, errors.New("GetUserByEmail not implemented")
}

func (m *mockUserService) UpdateUser(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error) {
	if m.updateUserFunc != nil {
		return m.updateUserFunc(ctx, id, req)
	}
	return nil, errors.New("UpdateUser not implemented")
}

func (m *mockUserService) DeleteUser(ctx context.Context, id string) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(ctx, id)
	}
	return errors.New("DeleteUser not implemented")
}

func (m *mockUserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	if m.listUsersFunc != nil {
		return m.listUsersFunc(ctx)
	}
	return nil, errors.New("ListUsers not implemented")
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	// Mock service with successful creation
	mockService := &mockUserService{
		createUserFunc: func(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
			return &models.User{
				ID:        primitive.NewObjectID(),
				UserID:    req.UserID,
				Email:     req.Email,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	// Create handler with mock service
	handler := NewUserHandler(mockService)

	// Create Echo instance
	e := echo.New()

	// Create request
	reqBody := `{"user_id":"test123","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test handler
	if err := handler.CreateUser(c); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check status code
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	// Parse response
	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check response fields
	if resp["user_id"] != "test123" {
		t.Errorf("Expected user_id test123, got %v", resp["user_id"])
	}
	if resp["email"] != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", resp["email"])
	}
}

func TestUserHandler_CreateUser_ValidationError(t *testing.T) {
	// Create handler with mock service
	handler := NewUserHandler(&mockUserService{})

	// Create Echo instance
	e := echo.New()

	// Create request with missing fields
	reqBody := `{"user_id":"","email":"test@example.com","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test handler
	err := handler.CreateUser(c)
	if err != nil {
		t.Fatalf("Expected no error to be returned, got %v", err)
	}

	// Check status code
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}