package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"go-mongodb-test/handlers"
)

// MockUserService is a mock implementation of the UserServiceInterface
type MockUserService struct{}

// MockUserHandler is a simplified handler for testing routes
type MockUserHandler struct {
	handlers.UserHandler
}

func (m *MockUserHandler) CreateUser(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{"status": "created"})
}

func (m *MockUserHandler) GetUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "user found"})
}

func (m *MockUserHandler) GetUserByUserID(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "user found by user_id"})
}

func (m *MockUserHandler) GetUserByEmail(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "user found by email"})
}

func (m *MockUserHandler) UpdateUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "user updated"})
}

func (m *MockUserHandler) DeleteUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "user deleted"})
}

func (m *MockUserHandler) ListUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "users listed", "count": 0, "users": []string{}})
}

func TestSetupRoutes(t *testing.T) {
	// Create echo instance
	e := echo.New()
	
	// Create a mock handler
	mockHandler := &MockUserHandler{}
	
	// Setup routes with mock handler
	SetupRoutes(e, mockHandler)
	
	// Test all routes
	testRoutes := []struct {
		name       string
		method     string
		path       string
		statusCode int
	}{
		{"CreateUser", http.MethodPost, "/api/users", http.StatusCreated},
		{"GetUser", http.MethodGet, "/api/users/123", http.StatusOK},
		{"UpdateUser", http.MethodPut, "/api/users/123", http.StatusOK},
		{"DeleteUser", http.MethodDelete, "/api/users/123", http.StatusOK},
		{"ListUsers", http.MethodGet, "/api/users", http.StatusOK},
		{"GetUserByUserID", http.MethodGet, "/api/users/search?user_id=testuser", http.StatusOK},
		{"GetUserByEmail", http.MethodGet, "/api/users/search?email=test@example.com", http.StatusOK},
	}
	
	for _, tc := range testRoutes {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			
			if rec.Code != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, rec.Code)
			}
		})
	}
}

func TestGetUserSearchHandler(t *testing.T) {
	// Create echo instance
	e := echo.New()
	
	// Create a mock handler
	mockHandler := &MockUserHandler{}
	
	// Test search handler with user_id
	t.Run("Search by user_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users/search?user_id=testuser", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Call the search handler
		err := getUserSearchHandler(c, mockHandler)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})
	
	// Test search handler with email
	t.Run("Search by email", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users/search?email=test@example.com", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Call the search handler
		err := getUserSearchHandler(c, mockHandler)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})
	
	// Test search handler with no parameters
	t.Run("Search with no parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users/search", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Call the search handler
		err := getUserSearchHandler(c, mockHandler)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
	
	// Test search handler with both parameters
	t.Run("Search with both parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users/search?user_id=testuser&email=test@example.com", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Call the search handler
		err := getUserSearchHandler(c, mockHandler)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// First parameter takes precedence (user_id)
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})
}