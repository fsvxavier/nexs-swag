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

// User demonstrates various struct tags.
type User struct {
	// ID with basic tags
	ID int `example:"123" json:"id"`

	// Name with validation
	Name string `example:"John Doe" json:"name" maxLength:"50" minLength:"2"`

	// Email with format and validation
	Email string `example:"john@example.com" format:"email" json:"email"`

	// Age with range validation
	Age int `example:"25" json:"age" maximum:"150" minimum:"0"`

	// Custom type override
	Birthday time.Time `example:"2000-01-01" format:"date" json:"birthday" swaggertype:"string"`

	// Ignored field
	Password string `json:"-" swaggerignore:"true"`

	// Custom extension
	Metadata map[string]interface{} `json:"metadata" x-nullable:"true"`
}

// Config with advanced tags.
type Config struct {
	// Read-only field
	CreatedAt time.Time `json:"created_at" readonly:"true"`

	// Default value
	Timeout int `default:"30" example:"60" json:"timeout"`

	// Enum values
	Status string `enums:"active,inactive,pending" example:"active" json:"status"`

	// Array with validation
	Tags []string `json:"tags" maxItems:"10" minItems:"1"`

	// Nested with swaggertype
	Settings interface{} `json:"settings" swaggertype:"object"`
}

// GetUser returns a user
// @Summary Get user
// @Tags users
// @Produce json
// @Success 200 {object} User
// @Router /users/{id} [get].
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
// @Router /config [get].
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
