package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// UserHandlerInterface defines the methods that need to be implemented by a handler
type UserHandlerInterface interface {
	CreateUser(c echo.Context) error
	GetUser(c echo.Context) error
	GetUserByUserID(c echo.Context) error
	GetUserByEmail(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	ListUsers(c echo.Context) error
}

// SetupRoutes configures all the routes for the API
func SetupRoutes(e *echo.Echo, handler UserHandlerInterface) {
	// Create API group
	api := e.Group("/api")

	// User routes
	users := api.Group("/users")
	users.POST("", handler.CreateUser)
	users.GET("", handler.ListUsers)
	users.GET("/:id", handler.GetUser)
	users.PUT("/:id", handler.UpdateUser)
	users.DELETE("/:id", handler.DeleteUser)

	// Search routes
	users.GET("/search", func(c echo.Context) error {
		return getUserSearchHandler(c, handler)
	})
}

// getUserSearchHandler handles requests to search for users by user_id or email
func getUserSearchHandler(c echo.Context, handler UserHandlerInterface) error {
	// Check if user_id parameter is present
	userID := c.QueryParam("user_id")
	if userID != "" {
		return handler.GetUserByUserID(c)
	}

	// Check if email parameter is present
	email := c.QueryParam("email")
	if email != "" {
		return handler.GetUserByEmail(c)
	}

	// If neither parameter is present, return bad request
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "Missing search parameter: user_id or email is required",
	})
}