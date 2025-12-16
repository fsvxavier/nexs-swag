package models

import "time"

// Product represents a product in the catalog
// This struct will be included when using --includeTypes="struct"
type Product struct {
	ID          int       `json:"id" example:"1"`
	Name        string    `json:"name" example:"Laptop"`
	Price       float64   `json:"price" format:"decimal" example:"999.99"`
	CreatedAt   time.Time `json:"created_at" swaggertype:"string" format:"date-time" example:"2025-12-16T10:00:00Z"`
	Category    string    `json:"category" example:"Electronics"`
	IsAvailable bool      `json:"is_available" example:"true"`
}

// ProductSummary is a simplified view of Product
// Demonstrates that only referenced structs are included
type ProductSummary struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Laptop"`
}

// UnusedModel demonstrates selective parsing
// This struct will NOT be included in the output because it's not referenced
type UnusedModel struct {
	Data string `json:"data"`
}
