package main

import (
	"encoding/json"
	"net/http"
)

// @title Multiple Formats API
// @version 1.0
// @description Example demonstrating multiple output formats
// @host localhost:8080
// @BasePath /api/v1

type Product struct {
	ID    int     `json:"id" example:"1"`
	Name  string  `json:"name" example:"Laptop"`
	Price float64 `json:"price" example:"999.99"`
}

// GetProduct returns a product
// @Summary Get product
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Router /products/{id} [get]
func GetProduct(w http.ResponseWriter, r *http.Request) {
	product := Product{ID: 1, Name: "Laptop", Price: 999.99}
	json.NewEncoder(w).Encode(product)
}

func main() {
	http.HandleFunc("/api/v1/products/", GetProduct)
	http.ListenAndServe(":8080", nil)
}
