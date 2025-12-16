package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderRequest represents a request to create an order
// Demonstrates nested struct dependencies and swaggertype usage
type OrderRequest struct {
	CustomerID   int         `json:"customer_id" example:"123" minimum:"1"`
	Items        []OrderItem `json:"items" minItems:"1"`
	Payment      PaymentInfo `json:"payment"`
	ShippingAddr Address     `json:"shipping_address"`
	BillingAddr  Address     `json:"billing_address"`
	Notes        string      `json:"notes,omitempty" maxLength:"500"`
}

// OrderResponse represents the response after creating/retrieving an order
type OrderResponse struct {
	OrderID      uuid.UUID   `json:"order_id" swaggertype:"string" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	CustomerID   int         `json:"customer_id" example:"123"`
	Items        []OrderItem `json:"items"`
	Payment      PaymentInfo `json:"payment"`
	ShippingAddr Address     `json:"shipping_address"`
	Total        Money       `json:"total"`
	Status       OrderStatus `json:"status"`
	CreatedAt    time.Time   `json:"created_at" swaggertype:"string" format:"date-time" example:"2025-12-16T10:30:00Z"`
	UpdatedAt    time.Time   `json:"updated_at" swaggertype:"string" format:"date-time" example:"2025-12-16T11:00:00Z"`
	Metadata     interface{} `json:"metadata,omitempty" swaggertype:"object"`
}

// OrderItem represents an item in an order
// This will be transitively included because OrderRequest references it
type OrderItem struct {
	ProductID   int    `json:"product_id" example:"456"`
	ProductName string `json:"product_name" example:"Laptop"`
	Quantity    int    `json:"quantity" example:"2" minimum:"1"`
	UnitPrice   Money  `json:"unit_price"`
	Subtotal    Money  `json:"subtotal"`
	Discount    *Money `json:"discount,omitempty"` // Pointer demonstrates optional nested struct
}

// PaymentInfo contains payment details
// Demonstrates swaggertype for sensitive data
type PaymentInfo struct {
	Method      PaymentMethod `json:"method"`
	CardLast4   string        `json:"card_last4,omitempty" pattern:"^[0-9]{4}$" example:"1234"`
	CardBrand   string        `json:"card_brand,omitempty" example:"Visa"`
	ProcessorID string        `json:"processor_id" swaggertype:"string" format:"uuid"`
	ProcessedAt time.Time     `json:"processed_at" swaggertype:"string" format:"date-time"`
}

// Address represents a physical address
// Reused for both shipping and billing
type Address struct {
	Street     string  `json:"street" example:"123 Main St" maxLength:"100"`
	City       string  `json:"city" example:"New York" maxLength:"50"`
	State      string  `json:"state" example:"NY" minLength:"2" maxLength:"2"`
	PostalCode string  `json:"postal_code" example:"10001" pattern:"^[0-9]{5}$"`
	Country    string  `json:"country" example:"US" minLength:"2" maxLength:"2"`
	Latitude   float64 `json:"latitude,omitempty" format:"double" minimum:"-90" maximum:"90"`
	Longitude  float64 `json:"longitude,omitempty" format:"double" minimum:"-180" maximum:"180"`
}

// Money represents a monetary amount with currency
// Uses swaggertype to convert to string for precision
type Money struct {
	Amount   int64  `json:"amount" swaggertype:"string" format:"int64" example:"99999"` // Amount in cents
	Currency string `json:"currency" example:"USD" minLength:"3" maxLength:"3"`
}

// OrderStatus represents the current status of an order
type OrderStatus struct {
	Code        StatusCode `json:"code"`
	Description string     `json:"description" example:"Order is being processed"`
	UpdatedBy   string     `json:"updated_by,omitempty" example:"system"`
}

// StatusCode is a string enum for order status
// Demonstrates enum with swaggertype
type StatusCode string

const (
	StatusPending    StatusCode = "pending"
	StatusProcessing StatusCode = "processing"
	StatusShipped    StatusCode = "shipped"
	StatusDelivered  StatusCode = "delivered"
	StatusCancelled  StatusCode = "cancelled"
)

// PaymentMethod is a string enum for payment methods
type PaymentMethod string

const (
	PaymentCard   PaymentMethod = "card"
	PaymentPaypal PaymentMethod = "paypal"
	PaymentCrypto PaymentMethod = "crypto"
)

// StatusUpdate is used to update order status
type StatusUpdate struct {
	Status StatusCode `json:"status" example:"shipped"`
	Notes  string     `json:"notes,omitempty" maxLength:"200"`
}

// OrderListResponse contains paginated list of orders
type OrderListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	TotalCount int             `json:"total_count" example:"50"`
	Page       int             `json:"page" example:"1" minimum:"1"`
	PageSize   int             `json:"page_size" example:"10" minimum:"1" maximum:"100"`
	HasMore    bool            `json:"has_more" example:"true"`
}

// ErrorResponse represents an error message
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid request"`
	Details string `json:"details,omitempty"`
}

// UnusedComplexModel demonstrates that complex unused types are excluded
type UnusedComplexModel struct {
	ID        uuid.UUID              `json:"id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}
