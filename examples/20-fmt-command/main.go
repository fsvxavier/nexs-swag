package main

import (
	"encoding/json"
	"net/http"
)

// @title Fmt Command Demo API
// @version 1.0
// @host localhost:8080
// @BasePath /api

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// GetProduct returns a product
// @Summary      Get product
// @Description  Get product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  Product
// @Failure      404  {object}  map[string]string
// @Router       /products/{id} [get]
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
