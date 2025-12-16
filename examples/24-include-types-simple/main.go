package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fsvxavier/nexs-swag/examples/24-include-types-simple/models"
)

// @title Include Types Simple Example
// @version 1.0
// @description Demonstrates --includeTypes flag with basic struct filtering
// @description Shows how swaggertype and format tags work with type filtering
// @description Structs are in separate files: models/product.go, interfaces/repository.go
// @host localhost:8080
// @BasePath /api/v1

// GetProduct returns a product by ID
// @Summary Get product
// @Description Retrieve product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product "Product found"
// @Failure 404 {string} string "Product not found"
// @Router /products/{id} [get]
func GetProduct(w http.ResponseWriter, r *http.Request) {
	product := models.Product{
		ID:          1,
		Name:        "Laptop",
		Price:       999.99,
		CreatedAt:   time.Now(),
		Category:    "Electronics",
		IsAvailable: true,
	}
	json.NewEncoder(w).Encode(product)
}

// ListProducts returns product summaries
// @Summary List products
// @Description Get list of all products (summary view)
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} models.ProductSummary "List of products"
// @Router /products [get]
func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := []models.ProductSummary{
		{ID: 1, Name: "Laptop"},
		{ID: 2, Name: "Mouse"},
	}
	json.NewEncoder(w).Encode(products)
}

func main() {
	http.HandleFunc("/api/v1/products", ListProducts)
	http.HandleFunc("/api/v1/products/", GetProduct)
	http.ListenAndServe(":8080", nil)
}
