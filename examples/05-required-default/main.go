package main

import (
	"encoding/json"
	"net/http"
)

// @title Required By Default API
// @version 1.0
// @description Demonstrates required by default behavior
// @host localhost:8080
// @BasePath /api

// Product demonstrates required fields
type Product struct {
	// Campos normais - tornam-se REQUIRED com flag
	ID    int     `json:"id" example:"1"`
	Name  string  `json:"name" example:"Laptop"`
	Price float64 `json:"price" example:"999.99"`

	// Campo com omitempty - continua OPTIONAL
	Description string `json:"description,omitempty" example:"High-end laptop"`

	// Campo pointer - continua OPTIONAL
	Discount *float64 `json:"discount" example:"10.5"`

	// Campo com binding omitempty - continua OPTIONAL
	Category string `json:"category" binding:"omitempty" example:"Electronics"`
}

// CreateProduct creates a product
// @Summary Create product
// @Description Create a new product (note required fields)
// @Tags products
// @Accept json
// @Produce json
// @Param product body Product true "Product object"
// @Success 201 {object} Product
// @Failure 400 {string} string "Validation error"
// @Router /products [post]
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func main() {
	http.HandleFunc("/api/products", CreateProduct)
	http.ListenAndServe(":8080", nil)
}
