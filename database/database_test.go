package database

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestConnectAndDisconnect(t *testing.T) {
	// Only run if MONGODB_URI is set
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// Set environment variables for testing
	os.Setenv("MONGODB_URI", mongoURI)
	os.Setenv("MONGODB_DB_NAME", "testdb")

	// Test successful connection
	db, err := Connect()
	if err != nil {
		t.Skipf("Skipping test due to MongoDB connection error: %v", err)
		return
	}

	// Check that the database is not nil
	if db == nil {
		t.Error("Expected database to be non-nil")
	}

	// Test the DB name
	if db.Name() != "testdb" {
		t.Errorf("Expected database name to be 'testdb', got '%s'", db.Name())
	}

	// Test collection method
	coll := db.Collection("users")
	if coll == nil {
		t.Error("Expected collection to be non-nil")
	}

	// Test close function
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err = Close(ctx)
	if err != nil {
		t.Errorf("Expected no error when closing connection, got %v", err)
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	// Test with existing environment variable
	testVarName := "TEST_ENV_VAR"
	os.Setenv(testVarName, "test_value")
	
	value := getEnvWithDefault(testVarName, "default_value")
	if value != "test_value" {
		t.Errorf("Expected '%s', got '%s'", "test_value", value)
	}
	
	// Test with non-existing environment variable
	nonExistingVarName := "NON_EXISTING_TEST_ENV_VAR"
	os.Unsetenv(nonExistingVarName)
	
	value = getEnvWithDefault(nonExistingVarName, "default_value")
	if value != "default_value" {
		t.Errorf("Expected '%s', got '%s'", "default_value", value)
	}
}

func TestGetMongoURI(t *testing.T) {
	// Test with MONGODB_URI set
	os.Setenv("MONGODB_URI", "mongodb://customhost:27017")
	
	uri := getMongoURI()
	if uri != "mongodb://customhost:27017" {
		t.Errorf("Expected '%s', got '%s'", "mongodb://customhost:27017", uri)
	}
	
	// Test with MONGODB_URI unset
	os.Unsetenv("MONGODB_URI")
	
	uri = getMongoURI()
	if uri != "mongodb://localhost:27017" {
		t.Errorf("Expected '%s', got '%s'", "mongodb://localhost:27017", uri)
	}
}

func TestGetDatabaseName(t *testing.T) {
	// Test with MONGODB_DB_NAME set
	os.Setenv("MONGODB_DB_NAME", "custom_db")
	
	name := getDatabaseName()
	if name != "custom_db" {
		t.Errorf("Expected '%s', got '%s'", "custom_db", name)
	}
	
	// Test with MONGODB_DB_NAME unset
	os.Unsetenv("MONGODB_DB_NAME")
	
	name = getDatabaseName()
	if name != "godb" {
		t.Errorf("Expected '%s', got '%s'", "godb", name)
	}
}