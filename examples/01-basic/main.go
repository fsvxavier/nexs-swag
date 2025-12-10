package main

import (
	"encoding/json"
	"net/http"
)

// @title Simple API
// @version 1.0
// @description This is a basic API example
// @host localhost:8080
// @BasePath /api/v1

// User represents a user in the system.
type User struct {
	ID   int    `example:"1"        json:"id"`
	Name string `example:"John Doe" json:"name"`
	Age  int    `example:"30"       json:"age"`
}

// GetUser returns a user by ID
// @Summary Get user
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get].
func GetUser(w http.ResponseWriter, r *http.Request) {
	user := User{ID: 1, Name: "John Doe", Age: 30}
	json.NewEncoder(w).Encode(user)
}

// CreateUser creates a new user
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User object"
// @Success 201 {object} User
// @Failure 400 {string} string "Invalid input"
// @Router /users [post].
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	json.NewEncoder(w).Encode(user)
}

func main() {
	http.HandleFunc("/api/v1/users/", GetUser)
	http.HandleFunc("/api/v1/users", CreateUser)
	http.ListenAndServe(":8080", nil)
}
