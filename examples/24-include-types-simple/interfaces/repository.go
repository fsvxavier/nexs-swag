package interfaces

import "github.com/fsvxavier/nexs-swag/examples/24-include-types-simple/models"

// ProductRepository defines the interface for product data access
// This interface demonstrates interface type category filtering
type ProductRepository interface {
	// GetByID retrieves a product by its ID
	GetByID(id int) (*models.Product, error)

	// List returns all products as summaries
	List() ([]models.ProductSummary, error)

	// Create adds a new product
	Create(product *models.Product) error
}

// UnusedInterface demonstrates selective parsing
// This interface will NOT be included because it's not referenced
type UnusedInterface interface {
	DoSomething() string
}
