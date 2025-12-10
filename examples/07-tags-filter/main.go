package main

import (
	"encoding/json"
	"net/http"
)

// @title Tag Filtering API
// @version 1.0
// @description Demonstrates --tags flag for filtering
// @host localhost:8080
// @BasePath /api

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetUsers returns all users
// @Summary List users
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Router /users [get].
func GetUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]User{{ID: 1, Name: "John"}})
}

// CreateUser creates a user
// @Summary Create user
// @Tags users,admin
// @Accept json
// @Success 201 {object} User
// @Router /users [post].
func CreateUser(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(User{ID: 1, Name: "New User"})
}

// DeleteUser deletes a user
// @Summary Delete user (admin only)
// @Tags admin
// @Success 204
// @Router /users/{id} [delete].
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// GetInternalConfig returns internal config
// @Summary Get internal config
// @Tags internal
// @Success 200 {object} map[string]string
// @Router /internal/config [get].
func GetInternalConfig(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"key": "value"})
}

func main() {
	http.HandleFunc("/api/users", GetUsers)
	http.ListenAndServe(":8080", nil)
}
