package main

import (
	"encoding/json"
	"net/http"
)

// @title Code Examples API
// @version 1.0
// @host localhost:8080
// @BasePath /api

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CreateUser creates a user
// @Summary Create user
// @Description Create a new user with code examples
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User data"
// @Success 201 {object} User
// @Router /users [post]
// @x-codeSamples file(create_user).
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func main() {
	http.HandleFunc("/api/users", CreateUser)
	http.ListenAndServe(":8080", nil)
}
