package main

import (
	"encoding/json"
	"net/http"
)

// This file contains products endpoints
// Notice: NO general API info here (@title, @version, etc)

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetProducts returns all products
// @Summary List products
// @Tags products
// @Produce json
// @Success 200 {array} Product
// @Router /products [get]
func GetProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{{ID: 1, Name: "Book"}, {ID: 2, Name: "Pen"}}
	json.NewEncoder(w).Encode(products)
}
