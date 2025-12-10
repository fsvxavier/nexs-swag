package main

import (
	"encoding/json"
	"net/http"
)

// @title Markdown Files API
// @version 1.0
// @description Demonstrates markdown files integration
// @host localhost:8080
// @BasePath /api

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CreateUser creates a new user
// @Summary Create user
// @Description file(create-user.md)
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User data"
// @Success 201 {object} User
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser retrieves a user
// @Summary Get user
// @Description file(get-user.md)
// @Tags users
// @Produce json
// @Success 200 {object} User
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(User{ID: 1, Name: "John"})
}

func main() {
	http.HandleFunc("/api/users", CreateUser)
	http.HandleFunc("/api/users/", GetUser)
	http.ListenAndServe(":8080", nil)
}
