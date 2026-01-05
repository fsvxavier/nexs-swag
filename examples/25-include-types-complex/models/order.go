package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderRequest represents a request to create an order
// Demonstrates nested struct dependencies and swaggertype usage.
type OrderRequest struct {
	CustomerID   int         `example:"123"           json:"customer_id" minimum:"1"`
	Items        []OrderItem `json:"items"            minItems:"1"`
	Payment      PaymentInfo `json:"payment"`
	ShippingAddr Address     `json:"shipping_address"`
	BillingAddr  Address     `json:"billing_address"`
	Notes        string      `json:"notes,omitempty"  maxLength:"500"`
}

// OrderResponse represents the response after creating/retrieving an order.
type OrderResponse struct {
	OrderID      uuid.UUID   `example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"        json:"order_id"   swaggertype:"string"`
	CustomerID   int         `example:"123"                                  json:"customer_id"`
	Items        []OrderItem `json:"items"`
	Payment      PaymentInfo `json:"payment"`
	ShippingAddr Address     `json:"shipping_address"`
	Total        Money       `json:"total"`
	Status       OrderStatus `json:"status"`
	CreatedAt    time.Time   `example:"2025-12-16T10:30:00Z"                 format:"date-time"   json:"created_at" swaggertype:"string"`
	UpdatedAt    time.Time   `example:"2025-12-16T11:00:00Z"                 format:"date-time"   json:"updated_at" swaggertype:"string"`
	Metadata     interface{} `json:"metadata,omitempty"                      swaggertype:"object"`
}

// OrderItem represents an item in an order
// This will be transitively included because OrderRequest references it.
type OrderItem struct {
	ProductID   int    `example:"456"             json:"product_id"`
	ProductName string `example:"Laptop"          json:"product_name"`
	Quantity    int    `example:"2"               json:"quantity"     minimum:"1"`
	UnitPrice   Money  `json:"unit_price"`
	Subtotal    Money  `json:"subtotal"`
	Discount    *Money `json:"discount,omitempty"` // Pointer demonstrates optional nested struct
}

// PaymentInfo contains payment details
// Demonstrates swaggertype for sensitive data.
type PaymentInfo struct {
	Method      PaymentMethod `json:"method"`
	CardLast4   string        `example:"1234"     json:"card_last4,omitempty" pattern:"^[0-9]{4}$"`
	CardBrand   string        `example:"Visa"     json:"card_brand,omitempty"`
	ProcessorID string        `format:"uuid"      json:"processor_id"         swaggertype:"string"`
	ProcessedAt time.Time     `format:"date-time" json:"processed_at"         swaggertype:"string"`
}

// Address represents a physical address
// Reused for both shipping and billing.
type Address struct {
	Street     string  `example:"123 Main St" json:"street"              maxLength:"100"`
	City       string  `example:"New York"    json:"city"                maxLength:"50"`
	State      string  `example:"NY"          json:"state"               maxLength:"2"        minLength:"2"`
	PostalCode string  `example:"10001"       json:"postal_code"         pattern:"^[0-9]{5}$"`
	Country    string  `example:"US"          json:"country"             maxLength:"2"        minLength:"2"`
	Latitude   float64 `format:"double"       json:"latitude,omitempty"  maximum:"90"         minimum:"-90"`
	Longitude  float64 `format:"double"       json:"longitude,omitempty" maximum:"180"        minimum:"-180"`
}

// Money represents a monetary amount with currency
// Uses swaggertype to convert to string for precision.
type Money struct {
	Amount   int64  `example:"99999" format:"int64"  json:"amount" swaggertype:"string"` // Amount in cents
	Currency string `example:"USD"   json:"currency" maxLength:"3" minLength:"3"`
}

// OrderStatus represents the current status of an order.
type OrderStatus struct {
	Code        StatusCode `json:"code"`
	Description string     `example:"Order is being processed" json:"description"`
	UpdatedBy   string     `example:"system"                   json:"updated_by,omitempty"`
}

// StatusCode is a string enum for order status
// Demonstrates enum with swaggertype.
type StatusCode string

const (
	StatusPending    StatusCode = "pending"
	StatusProcessing StatusCode = "processing"
	StatusShipped    StatusCode = "shipped"
	StatusDelivered  StatusCode = "delivered"
	StatusCancelled  StatusCode = "cancelled"
)

// PaymentMethod is a string enum for payment methods.
type PaymentMethod string

const (
	PaymentCard   PaymentMethod = "card"
	PaymentPaypal PaymentMethod = "paypal"
	PaymentCrypto PaymentMethod = "crypto"
)

// StatusUpdate is used to update order status.
type StatusUpdate struct {
	Status StatusCode `example:"shipped"      json:"status"`
	Notes  string     `json:"notes,omitempty" maxLength:"200"`
}

// OrderListResponse contains paginated list of orders.
type OrderListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	TotalCount int             `example:"50"   json:"total_count"`
	Page       int             `example:"1"    json:"page"        minimum:"1"`
	PageSize   int             `example:"10"   json:"page_size"   maximum:"100" minimum:"1"`
	HasMore    bool            `example:"true" json:"has_more"`
}

// ErrorResponse represents an error message.
type ErrorResponse struct {
	Code    int    `example:"400"             json:"code"`
	Message string `example:"Invalid request" json:"message"`
	Details string `json:"details,omitempty"`
}

// UnusedComplexModel demonstrates that complex unused types are excluded.
type UnusedComplexModel struct {
	ID        uuid.UUID              `json:"id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}
