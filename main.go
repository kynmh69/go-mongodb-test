package main

import (
	"log"
	"os"

	"go-mongodb-test/database"
	"go-mongodb-test/handlers"
	"go-mongodb-test/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize database connection
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize services
	userService := services.NewUserService(db.DB)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Add JSON content type validation middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == "POST" || c.Request().Method == "PUT" {
				contentType := c.Request().Header.Get("Content-Type")
				if contentType != "" && contentType != "application/json" {
					return c.JSON(400, map[string]string{
						"error": "Content-Type must be application/json",
					})
				}
			}
			return next(c)
		}
	})

	// Routes
	api := e.Group("/api/v1")

	// User routes
	users := api.Group("/users")
	users.POST("", userHandler.CreateUser)        // Create user
	users.GET("", userHandler.ListUsers)          // List all users
	users.GET("/search", userHandler.GetUserByUserID) // Search by user_id (query param)
	users.GET("/search/email", userHandler.GetUserByEmail) // Search by email (query param)
	users.GET("/:id", userHandler.GetUser)        // Get user by MongoDB ID
	users.PUT("/:id", userHandler.UpdateUser)     // Update user
	users.DELETE("/:id", userHandler.DeleteUser)  // Delete user

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
			"message": "User management service is running",
		})
	})

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(e.Start(":" + port))
}