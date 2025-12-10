package main

import (
	"encoding/json"
	"net/http"
)

// @title Parse Func Body API
// @version 1.0
// @host localhost:8080
// @BasePath /api

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CreateItem creates an item
// @Summary Create item
// @Tags items
// @Accept json
// @Param item body Item true "Item data"
// @Success 201 {object} Item
// @Failure 400 {object} map[string]string
// @Router /items [post].
func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}

	// Validação interna que será detectada com --parseFuncBody
	if item.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "name is required"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func main() {
	http.HandleFunc("/api/items", CreateItem)
	http.ListenAndServe(":8080", nil)
}
