package main

import (
	"encoding/json"
	"net/http"
)

// @title Property Strategy API
// @version 1.0
// @description Demonstrates different naming strategies
// @host localhost:8080
// @BasePath /api

// UserProfile demonstrates field naming.
type UserProfile struct {
	// Campo COM tag json explícita - SEMPRE usa a tag
	UserID int `example:"123" json:"user_id"`

	// Campos SEM tag json - usa propertyStrategy
	FirstName string `example:"John"` // snake_case: first_name | camelCase: firstName | PascalCase: FirstName
	LastName  string `example:"Doe"`  // snake_case: last_name  | camelCase: lastName  | PascalCase: LastName
	IsActive  bool   `example:"true"` // snake_case: is_active  | camelCase: isActive  | PascalCase: IsActive

	// Campo com omitempty - mantém comportamento
	MiddleName string `example:"Smith" json:",omitempty"`
}

// GetProfile returns user profile
// @Summary Get user profile
// @Description Returns profile demonstrating naming strategy
// @Tags users
// @Produce json
// @Success 200 {object} UserProfile
// @Router /profile [get].
func GetProfile(w http.ResponseWriter, r *http.Request) {
	profile := UserProfile{
		UserID:    123,
		FirstName: "John",
		LastName:  "Doe",
		IsActive:  true,
	}
	json.NewEncoder(w).Encode(profile)
}

func main() {
	http.HandleFunc("/api/profile", GetProfile)
	http.ListenAndServe(":8080", nil)
}
