package main

import (
	"encoding/json"
	"net/http"
)

// @title Generated Time API
// @version 1.0
// @host localhost:8080
// @BasePath /api

type Status struct {
	Message string `json:"message"`
}

// GetStatus returns status
// @Summary Get status
// @Tags system
// @Produce json
// @Success 200 {object} Status
// @Router /status [get].
func GetStatus(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Status{Message: "OK"})
}

func main() {
	http.HandleFunc("/api/status", GetStatus)
	http.ListenAndServe(":8080", nil)
}
