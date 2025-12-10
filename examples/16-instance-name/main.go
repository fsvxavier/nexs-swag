package main

import (
	"encoding/json"
	"net/http"
)

// @title Instance Name API
// @version 1.0
// @host localhost:8080
// @BasePath /api

type Response struct {
	Data string `json:"data"`
}

// GetData returns data
// @Summary Get data
// @Tags data
// @Produce json
// @Success 200 {object} Response
// @Router /data [get]
func GetData(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{Data: "Hello"})
}

func main() {
	http.HandleFunc("/api/data", GetData)
	http.ListenAndServe(":8080", nil)
}
