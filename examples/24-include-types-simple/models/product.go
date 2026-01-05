package models

import "time"

// Product represents a product in the catalog
// This struct will be included when using --includeTypes="struct".
type Product struct {
	ID          int       `example:"1"                    json:"id"`
	Name        string    `example:"Laptop"               json:"name"`
	Price       float64   `example:"999.99"               format:"decimal"    json:"price"`
	CreatedAt   time.Time `example:"2025-12-16T10:00:00Z" format:"date-time"  json:"created_at" swaggertype:"string"`
	Category    string    `example:"Electronics"          json:"category"`
	IsAvailable bool      `example:"true"                 json:"is_available"`
}

// ProductSummary is a simplified view of Product
// Demonstrates that only referenced structs are included.
type ProductSummary struct {
	ID   int    `example:"1"      json:"id"`
	Name string `example:"Laptop" json:"name"`
}

// UnusedModel demonstrates selective parsing
// This struct will NOT be included in the output because it's not referenced.
type UnusedModel struct {
	Data string `json:"data"`
}
