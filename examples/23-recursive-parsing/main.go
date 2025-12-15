package main

// @title Recursive Parsing API
// @version 1.0
// @description Demonstrates recursive parsing with --parseInternal and --exclude flags
// @host localhost:8080
// @BasePath /api

// User represents a user in the main package
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetUser retrieves a user
// @Summary Get user from main package
// @Description Retrieves user information from the main package
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /users/{id} [get]
func GetUser() {}

func main() {}
