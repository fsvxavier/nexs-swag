package services

import (
	"time"

	"github.com/fsvxavier/nexs-swag/examples/25-include-types-complex/models"
	"github.com/google/uuid"
)

// OrderProcessor defines the interface for order processing
// This interface won't be included unless --includeTypes="interface" is used
type OrderProcessor interface {
	Process(order *models.OrderRequest) (*models.OrderResponse, error)
	Validate(order *models.OrderRequest) error
	CalculateTotal(items []models.OrderItem) models.Money
}

// ProcessOrder is a simple function to process orders
// This demonstrates that the function itself isn't included in schemas
func ProcessOrder(req *models.OrderRequest) *models.OrderResponse {
	total := models.Money{
		Amount:   calculateTotalAmount(req.Items),
		Currency: "USD",
	}

	return &models.OrderResponse{
		OrderID:      uuid.New(),
		CustomerID:   req.CustomerID,
		Items:        req.Items,
		Payment:      req.Payment,
		ShippingAddr: req.ShippingAddr,
		Total:        total,
		Status: models.OrderStatus{
			Code:        models.StatusPending,
			Description: "Order received and pending processing",
			UpdatedBy:   "system",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func calculateTotalAmount(items []models.OrderItem) int64 {
	var total int64
	for _, item := range items {
		total += item.Subtotal.Amount
	}
	return total
}

// UnusedService demonstrates that unused types are excluded
type UnusedService struct {
	Name string
	ID   int
}
