package models

import "time"

// ComplexType is a complex type that should be converted to primitive
type ComplexType struct {
	InternalData string
	Secret       string
}

// Address is a nested struct that should be included
type Address struct {
	Street  string `json:"street" example:"123 Main St"`
	City    string `json:"city" example:"San Francisco"`
	ZipCode string `json:"zipCode" example:"94102"`
}

// ContactInfo has complex types with swaggertype override
type ContactInfo struct {
	Email       string `json:"email" format:"email"`
	PhoneNumber string `json:"phoneNumber" format:"phone"`
}

// OrderStatus is used in Order (transitive dependency)
type OrderStatus struct {
	Code        string    `json:"code" example:"PENDING"`
	Description string    `json:"description" example:"Order is being processed"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Order has complex type conversions
type Order struct {
	ID         string       `json:"id" example:"ORD-123"`
	ProductID  int          `json:"productId" example:"1"`
	Quantity   int          `json:"quantity" example:"2"`
	TotalPrice float64      `json:"totalPrice" example:"1999.98" format:"decimal"`
	Status     *OrderStatus `json:"status"`
	// ComplexType converted to string primitive
	Metadata ComplexType `json:"metadata" swaggertype:"string" format:"json" example:"{\"key\":\"value\"}"`
	// Time with custom format
	CreatedAt time.Time `json:"createdAt"`
	// Address is a normal nested struct (should be included)
	ShippingAddress *Address `json:"shippingAddress"`
	// ContactInfo should be included
	Contact ContactInfo `json:"contact"`
	// Array of primitive with swaggertype override
	Tags []string `json:"tags" swaggertype:"array,string"`
}

// Payment has swaggertype converting struct to primitive
type Payment struct {
	ID      string  `json:"id" example:"PAY-456"`
	OrderID string  `json:"orderId" example:"ORD-123"`
	Amount  float64 `json:"amount" example:"1999.98"`
	// ComplexType converted to object primitive
	PaymentDetails ComplexType `json:"paymentDetails" swaggertype:"object"`
	ProcessedAt    time.Time   `json:"processedAt" format:"date-time"`
}

// UnusedComplex should NOT be included (not referenced)
type UnusedComplex struct {
	Field1 ComplexType `json:"field1"`
	Field2 string      `json:"field2"`
}

// GetOrder retrieves an order with complex types
// @Summary Get order with complex types
// @Description Retrieves order with nested structs and swaggertype conversions
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} Order
// @Failure 404 {string} string "Order not found"
// @Router /orders/{id} [get]
func GetOrder() {}

// CreatePayment creates a payment
// @Summary Create payment
// @Description Creates a new payment with type conversions
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body Payment true "Payment data"
// @Success 201 {object} Payment
// @Failure 400 {string} string "Bad request"
// @Router /payments [post]
func CreatePayment() {}
