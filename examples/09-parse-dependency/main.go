package main

import (
	"encoding/json"
	"net/http"
)

// Product represents a product in models subpackage
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// @title Parse Dependency API
// @version 1.0
// @host localhost:8080
// @BasePath /api

// GetProduct returns a product
// @Summary Get product
// @Tags products
// @Produce json
// @Success 200 {object} Product
// @Router /products/{id} [get]
func GetProduct(w http.ResponseWriter, r *http.Request) {
	product := Product{
		ID:    1,
		Name:  "Laptop",
		Price: 999.99,
	}
	json.NewEncoder(w).Encode(product)
}

func main() {
	http.HandleFunc("/api/products/", GetProduct)
	http.ListenAndServe(":8080", nil)
}
