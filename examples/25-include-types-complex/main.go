package main

import (
	"encoding/json"
	"net/http"

	"github.com/fsvxavier/nexs-swag/examples/25-include-types-complex/models"
	"github.com/fsvxavier/nexs-swag/examples/25-include-types-complex/services"
)

// @title Include Types Complex Example
// @version 2.0
// @description Advanced demonstration of --includeTypes with complex scenarios
// @description Features: nested structs, swaggertype overrides, multiple formats, transitive dependencies
// @host localhost:8080
// @BasePath /api/v2

// CreateOrder creates a new order with payment
// @Summary Create order
// @Description Create a new order with customer, items, and payment information
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.OrderRequest true "Order with payment details"
// @Success 201 {object} models.OrderResponse "Order created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /orders [post]
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.OrderRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Process order using service
	response := services.ProcessOrder(&req)
	json.NewEncoder(w).Encode(response)
}

// GetOrder retrieves an order by ID
// @Summary Get order
// @Description Get order details including all items and status
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID" format:"uuid"
// @Success 200 {object} models.OrderResponse "Order found"
// @Failure 404 {object} models.ErrorResponse "Order not found"
// @Router /orders/{id} [get]
func GetOrder(w http.ResponseWriter, r *http.Request) {
	response := models.OrderResponse{}
	json.NewEncoder(w).Encode(response)
}

// UpdateOrderStatus updates the order status
// @Summary Update order status
// @Description Update the status of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID" format:"uuid"
// @Param status body models.StatusUpdate true "New status"
// @Success 200 {object} models.OrderResponse "Status updated"
// @Failure 404 {object} models.ErrorResponse "Order not found"
// @Router /orders/{id}/status [patch]
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	var req models.StatusUpdate
	json.NewDecoder(r.Body).Decode(&req)
	response := models.OrderResponse{}
	json.NewEncoder(w).Encode(response)
}

// GetCustomerOrders retrieves all orders for a customer
// @Summary Get customer orders
// @Description Get all orders for a specific customer with pagination
// @Tags customers,orders
// @Accept json
// @Produce json
// @Param customerId path int true "Customer ID"
// @Param page query int false "Page number" default:1
// @Param pageSize query int false "Page size" default:10
// @Success 200 {object} models.OrderListResponse "Customer orders"
// @Failure 404 {object} models.ErrorResponse "Customer not found"
// @Router /customers/{customerId}/orders [get]
func GetCustomerOrders(w http.ResponseWriter, r *http.Request) {
	response := models.OrderListResponse{}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/v2/orders", CreateOrder)
	http.HandleFunc("/api/v2/orders/", GetOrder)
	http.ListenAndServe(":8080", nil)
}
