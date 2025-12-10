package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// @title Struct Tags API
// @version 1.0
// @host localhost:8080
// @BasePath /api

// User demonstrates various struct tags
type User struct {
	// ID with basic tags
	ID int `json:"id" example:"123"`

	// Name with validation
	Name string `json:"name" example:"John Doe" minLength:"2" maxLength:"50"`

	// Email with format and validation
	Email string `json:"email" example:"john@example.com" format:"email"`

	// Age with range validation
	Age int `json:"age" example:"25" minimum:"0" maximum:"150"`

	// Custom type override
	Birthday time.Time `json:"birthday" swaggertype:"string" format:"date" example:"2000-01-01"`

	// Ignored field
	Password string `json:"-" swaggerignore:"true"`

	// Custom extension
	Metadata map[string]interface{} `json:"metadata" x-nullable:"true"`
}

// Config with advanced tags
type Config struct {
	// Read-only field
	CreatedAt time.Time `json:"created_at" readonly:"true"`

	// Default value
	Timeout int `json:"timeout" default:"30" example:"60"`

	// Enum values
	Status string `json:"status" enums:"active,inactive,pending" example:"active"`

	// Array with validation
	Tags []string `json:"tags" minItems:"1" maxItems:"10"`

	// Nested with swaggertype
	Settings interface{} `json:"settings" swaggertype:"object"`
}

// GetUser returns a user
// @Summary Get user
// @Tags users
// @Produce json
// @Success 200 {object} User
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	user := User{
		ID:    123,
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}
	json.NewEncoder(w).Encode(user)
}

// GetConfig returns config
// @Summary Get config
// @Tags config
// @Produce json
// @Success 200 {object} Config
// @Router /config [get]
func GetConfig(w http.ResponseWriter, r *http.Request) {
	config := Config{
		CreatedAt: time.Now(),
		Timeout:   30,
		Status:    "active",
		Tags:      []string{"api", "web"},
	}
	json.NewEncoder(w).Encode(config)
}

func main() {
	http.HandleFunc("/api/users/", GetUser)
	http.HandleFunc("/api/config", GetConfig)
	http.ListenAndServe(":8080", nil)
}
