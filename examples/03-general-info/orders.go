package main

import (
	"encoding/json"
	"net/http"
)

// This file contains orders endpoints
// Notice: NO general API info here

type Order struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
}

// GetOrders returns all orders
// @Summary List orders
// @Tags orders
// @Produce json
// @Success 200 {array} Order
// @Router /orders [get].
func GetOrders(w http.ResponseWriter, r *http.Request) {
	orders := []Order{{ID: 1, UserID: 1}, {ID: 2, UserID: 2}}
	json.NewEncoder(w).Encode(orders)
}
