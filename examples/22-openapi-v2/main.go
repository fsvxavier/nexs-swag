// Package main demonstrates OpenAPI 2.0 (Swagger 2.0) generation.
//
// This example shows how to generate both OpenAPI 3.1.0 and Swagger 2.0
// specifications from the same Go code annotations.
//
// @title           Product API
// @version         1.0
// @description     A simple product management API demonstrating Swagger 2.0 generation
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8080
// @BasePath  /api/v1
// @schemes   http https
//
// @tag.name         Products
// @tag.description  Operations about products
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Product represents a product in the system.
type Product struct {
	ID          int     `json:"id" example:"1"`
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"High-performance laptop"`
	Price       float64 `json:"price" example:"999.99"`
	Stock       int     `json:"stock" example:"50"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

// products is an in-memory store.
var products = map[int]Product{
	1: {ID: 1, Name: "Laptop", Description: "High-performance laptop", Price: 999.99, Stock: 50},
	2: {ID: 2, Name: "Mouse", Description: "Wireless mouse", Price: 29.99, Stock: 200},
	3: {ID: 3, Name: "Keyboard", Description: "Mechanical keyboard", Price: 79.99, Stock: 100},
}

var nextID = 4

func main() {
	http.HandleFunc("/api/v1/products", handleProducts)
	http.HandleFunc("/api/v1/products/", handleProduct)

	fmt.Println("Server starting on :8080")
	fmt.Println("API documentation:")
	fmt.Println("  - OpenAPI 3.1.0: http://localhost:8080/docs/openapi.json")
	fmt.Println("  - Swagger 2.0:   http://localhost:8080/docs/swagger.json")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ListProducts godoc
//
// @Summary      List all products
// @Description  Get a list of all products in the system
// @Tags         Products
// @Accept       json
// @Produce      json
// @Success      200  {array}   Product
// @Failure      500  {object}  ErrorResponse
// @Router       /products [get]
func handleProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		listProducts(w, r)
	} else if r.Method == http.MethodPost {
		createProduct(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	list := make([]Product, 0, len(products))
	for _, p := range products {
		list = append(list, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// CreateProduct godoc
//
// @Summary      Create a product
// @Description  Add a new product to the system
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      Product  true  "Product object"
// @Success      201      {object}  Product
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /products [post]
func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
		})
		return
	}

	product.ID = nextID
	nextID++
	products[product.ID] = product

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func handleProduct(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	idStr := r.URL.Path[len("/api/v1/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getProduct(w, r, id)
	case http.MethodPut:
		updateProduct(w, r, id)
	case http.MethodDelete:
		deleteProduct(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetProduct godoc
//
// @Summary      Get a product
// @Description  Get a product by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  Product
// @Failure      404  {object}  ErrorResponse
// @Router       /products/{id} [get]
func getProduct(w http.ResponseWriter, r *http.Request, id int) {
	product, exists := products[id]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    404,
			Message: "Product not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
//
// @Summary      Update a product
// @Description  Update an existing product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id       path      int      true  "Product ID"
// @Param        product  body      Product  true  "Updated product object"
// @Success      200      {object}  Product
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Router       /products/{id} [put]
func updateProduct(w http.ResponseWriter, r *http.Request, id int) {
	if _, exists := products[id]; !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    404,
			Message: "Product not found",
		})
		return
	}

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
		})
		return
	}

	product.ID = id
	products[id] = product

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct godoc
//
// @Summary      Delete a product
// @Description  Remove a product from the system
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      204  "No Content"
// @Failure      404  {object}  ErrorResponse
// @Router       /products/{id} [delete]
func deleteProduct(w http.ResponseWriter, r *http.Request, id int) {
	if _, exists := products[id]; !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    404,
			Message: "Product not found",
		})
		return
	}

	delete(products, id)
	w.WriteHeader(http.StatusNoContent)
}
