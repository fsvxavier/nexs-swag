package main

import (
	"encoding/json"
	"net/http"
)

// @title Exclude Patterns API
// @version 1.0
// @description Demonstrates --exclude flag
// @host localhost:8080
// @BasePath /api

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetUser in main API
// @Summary Get user
// @Tags users
// @Success 200 {object} User
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(User{ID: 1, Name: "John"})
}

func main() {
	http.HandleFunc("/api/users/", GetUser)
	http.ListenAndServe(":8080", nil)
}
