package main

import (
	"encoding/json"
	"net/http"
)

// @title Dependency Level API
// @version 1.0
// @host localhost:8080
// @BasePath /api

// Meta represents metadata (Level 3)
type Meta struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Item represents an item (Level 2)
type Item struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Metadata Meta   `json:"metadata"`
}

// Order represents an order (Level 1)
type Order struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Items  []Item `json:"items"`
}

// GetOrder returns an order
// @Summary Get order
// @Tags orders
// @Produce json
// @Success 200 {object} Order
// @Router /orders/{id} [get]
func GetOrder(w http.ResponseWriter, r *http.Request) {
	order := Order{
		ID:     1,
		Status: "pending",
		Items: []Item{
			{
				ID:   1,
				Name: "Product 1",
				Metadata: Meta{
					CreatedAt: "2025-01-01",
					UpdatedAt: "2025-01-02",
				},
			},
		},
	}
	json.NewEncoder(w).Encode(order)
}

func main() {
	http.HandleFunc("/api/orders/", GetOrder)
	http.ListenAndServe(":8080", nil)
}
