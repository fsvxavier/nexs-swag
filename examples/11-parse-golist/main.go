package main

import (
	"encoding/json"
	"net/http"
)

// @title Parse GoList API
// @version 1.0
// @host localhost:8080
// @BasePath /api

type Data struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

// GetData returns data
// @Summary Get data
// @Tags data
// @Produce json
// @Success 200 {object} Data
// @Router /data [get].
func GetData(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Data{ID: 1, Value: "test"})
}

func main() {
	http.HandleFunc("/api/data", GetData)
	http.ListenAndServe(":8080", nil)
}
