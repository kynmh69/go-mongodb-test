package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-mongodb-test/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Mock UserService for testing
type mockUserService struct {
	createUserFunc      func(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	getUserByIDFunc     func(ctx context.Context, id string) (*models.User, error)
	getUserByUserIDFunc func(ctx context.Context, userID string) (*models.User, error)
	getUserByEmailFunc  func(ctx context.Context, email string) (*models.User, error)
	updateUserFunc      func(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error)
	deleteUserFunc      func(ctx context.Context, id string) error
	listUsersFunc       func(ctx context.Context) ([]*models.User, error)
}

func (m *mockUserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetUserByUserID(ctx context.Context, userID string) (*models.User, error) {
	if m.getUserByUserIDFunc != nil {
		return m.getUserByUserIDFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, email)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) UpdateUser(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error) {
	if m.updateUserFunc != nil {
		return m.updateUserFunc(ctx, id, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) DeleteUser(ctx context.Context, id string) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockUserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	if m.listUsersFunc != nil {
		return m.listUsersFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func TestNewUserHandler(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	if handler == nil {
		t.Fatal("Expected NewUserHandler to return a non-nil UserHandler")
	}

	if handler.userService == nil {
		t.Error("Expected userService to be set")
	}
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	mockService := &mockUserService{
		createUserFunc: func(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
			return &models.User{
				ID:        bson.NewObjectID(),
				UserID:    req.UserID,
				Email:     req.Email,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewUserHandler(mockService)
	e := echo.New()

	reqBody := models.CreateUserRequest{
		UserID:   "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateUser(c)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var user models.User
	if err := json.Unmarshal(rec.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if user.UserID != reqBody.UserID {
		t.Errorf("Expected UserID %s, got %s", reqBody.UserID, user.UserID)
	}
}

func TestUserHandler_CreateUser_InvalidRequest(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)
	e := echo.New()

	// Test invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("invalid json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateUser(c)
	if err != nil {
		t.Fatalf("Expected no error from handler, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUserHandler_CreateUser_MissingFields(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)
	e := echo.New()

	tests := []struct {
		name string
		body models.CreateUserRequest
	}{
		{
			name: "Missing UserID",
			body: models.CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
		},
		{
			name: "Missing Email",
			body: models.CreateUserRequest{
				UserID:   "testuser",
				Password: "password123",
			},
		},
		{
			name: "Missing Password",
			body: models.CreateUserRequest{
				UserID: "testuser",
				Email:  "test@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.CreateUser(c)
			if err != nil {
				t.Fatalf("Expected no error from handler, got %v", err)
			}

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
			}
		})
	}
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	userID := bson.NewObjectID()
	mockService := &mockUserService{
		getUserByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
			return &models.User{
				ID:        userID,
				UserID:    "testuser",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewUserHandler(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.Hex(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.Hex())

	err := handler.GetUser(c)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	mockService := &mockUserService{
		getUserByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
			return nil, errors.New("user not found")
		},
	}

	handler := NewUserHandler(mockService)
	e := echo.New()

	userID := bson.NewObjectID()
	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.Hex(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.Hex())

	err := handler.GetUser(c)
	if err != nil {
		t.Fatalf("Expected no error from handler, got %v", err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUserHandler_GetUserByUserID_Success(t *testing.T) {
	mockService := &mockUserService{
		getUserByUserIDFunc: func(ctx context.Context, userID string) (*models.User, error) {
			return &models.User{
				ID:        bson.NewObjectID(),
				UserID:    userID,
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewUserHandler(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/users/search?user_id=testuser", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetUserByUserID(c)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestUserHandler_GetUserByUserID_MissingParam(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/users/search", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetUserByUserID(c)
	if err != nil {
		t.Fatalf("Expected no error from handler, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUserHandler_ListUsers_Success(t *testing.T) {
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

	mockService := &mockUserService{
		listUsersFunc: func(ctx context.Context) ([]*models.User, error) {
			return users, nil
		},
	}

	handler := NewUserHandler(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.ListUsers(c)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if count, ok := response["count"].(float64); !ok || int(count) != len(users) {
		t.Errorf("Expected count %d, got %v", len(users), response["count"])
	}
}

func TestUserHandler_DeleteUser_Success(t *testing.T) {
	mockService := &mockUserService{
		deleteUserFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}

	handler := NewUserHandler(mockService)
	e := echo.New()

	userID := bson.NewObjectID()
	req := httptest.NewRequest(http.MethodDelete, "/users/"+userID.Hex(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.Hex())

	err := handler.DeleteUser(c)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "User deleted successfully" {
		t.Errorf("Expected success message, got %s", response["message"])
	}
}

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	userID := bson.NewObjectID()
	updatedUser := &models.User{
		ID:        userID,
		UserID:    "updateduser",
		Email:     "updated@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService := &mockUserService{
		updateUserFunc: func(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error) {
			return updatedUser, nil
		},
	}

	handler := NewUserHandler(mockService)
	e := echo.New()

	updateReq := models.UpdateUserRequest{
		UserID: stringPtr("updateduser"),
		Email:  stringPtr("updated@example.com"),
	}

	jsonBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/users/"+userID.Hex(), bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.Hex())

	err := handler.UpdateUser(c)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var user models.User
	if err := json.Unmarshal(rec.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if user.UserID != "updateduser" {
		t.Errorf("Expected UserID %s, got %s", "updateduser", user.UserID)
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}