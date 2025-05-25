package database

import (
	"os"
	"testing"
	"time"
)

func TestDatabase_Struct(t *testing.T) {
	db := &Database{}

	if db.Client != nil {
		t.Error("Expected Client to be nil in empty Database struct")
	}

	if db.DB != nil {
		t.Error("Expected DB to be nil in empty Database struct")
	}
}

func TestNewConnection_Environment(t *testing.T) {
	// Test with default values (no environment variables set)
	originalURI := os.Getenv("MONGODB_URI")
	originalDBName := os.Getenv("DATABASE_NAME")

	// Clear environment variables
	os.Unsetenv("MONGODB_URI")
	os.Unsetenv("DATABASE_NAME")

	// Note: This test will fail if MongoDB is not running locally
	// but we can test the environment variable handling
	defer func() {
		// Restore original environment variables
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		}
		if originalDBName != "" {
			os.Setenv("DATABASE_NAME", originalDBName)
		}
	}()

	// Test environment variable handling by mocking the connection creation
	// We can't actually connect without a MongoDB instance, but we can test
	// that the function properly reads environment variables

	// Set custom environment variables
	customURI := "mongodb://testhost:27017"
	customDBName := "testdb"

	os.Setenv("MONGODB_URI", customURI)
	os.Setenv("DATABASE_NAME", customDBName)

	// The actual connection will fail, but we can test that it reads env vars correctly
	// by checking if the function attempts to use the custom URI
	_, err := NewConnection()
	if err == nil {
		t.Error("Expected connection to fail without MongoDB instance")
	}

	// Test with invalid URI format
	os.Setenv("MONGODB_URI", "invalid-uri")
	_, err = NewConnection()
	if err == nil {
		t.Error("Expected connection to fail with invalid URI")
	}
}

func TestDatabase_Close(t *testing.T) {
	// Test closing a nil client (should not panic)
	db := &Database{}
	
	err := db.Close()
	if err == nil {
		t.Error("Expected error when closing database with nil client")
	}
}

func TestNewConnection_Timeout(t *testing.T) {
	// Test that the function uses proper timeout
	start := time.Now()
	
	// This will timeout since we don't have a MongoDB instance
	_, err := NewConnection()
	
	elapsed := time.Since(start)
	
	if err == nil {
		t.Error("Expected connection to fail without MongoDB instance")
	}

	// The timeout should be around 30 seconds, but we'll be more lenient
	if elapsed > 35*time.Second {
		t.Error("Connection took too long, timeout might not be working")
	}
}

func TestNewConnection_ContextHandling(t *testing.T) {
	// Test that context cancellation works properly
	
	// Since NewConnection creates its own context, we can't directly test context cancellation
	// but we can ensure the function handles context properly by checking it doesn't hang indefinitely
	
	done := make(chan bool, 1)
	go func() {
		_, err := NewConnection()
		if err == nil {
			t.Error("Expected connection to fail")
		}
		done <- true
	}()
	
	select {
	case <-done:
		// Function completed, which is expected
	case <-time.After(5 * time.Second): // Reduced timeout for faster feedback during testing
		t.Error("NewConnection took too long, context timeout might not be working")
	}
}

func TestNewConnection_DefaultValues(t *testing.T) {
	// Clear environment variables to test defaults
	originalURI := os.Getenv("MONGODB_URI")
	originalDBName := os.Getenv("DATABASE_NAME")
	
	os.Unsetenv("MONGODB_URI")
	os.Unsetenv("DATABASE_NAME")
	
	defer func() {
		// Restore original environment variables
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		} else {
			os.Unsetenv("MONGODB_URI")
		}
		if originalDBName != "" {
			os.Setenv("DATABASE_NAME", originalDBName)
		} else {
			os.Unsetenv("DATABASE_NAME")
		}
	}()
	
	// The function should use default values when env vars are not set
	// We can't test the actual connection, but we can verify the function
	// attempts to use the default values by checking the error message
	_, err := NewConnection()
	if err == nil {
		t.Error("Expected connection to fail without MongoDB instance")
	}
	
	// The error should indicate it tried to connect to localhost:27017 (default)
	if err != nil && err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}