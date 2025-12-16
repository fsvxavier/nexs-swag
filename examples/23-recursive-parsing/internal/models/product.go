package models

// Product represents a product in the catalog
type Product struct {
	ID          int       `json:"id" example:"1"`
	Name        string    `json:"name" example:"Laptop"`
	Description string    `json:"description" example:"High performance laptop"`
	Price       float64   `json:"price" example:"999.99"`
	Stock       int       `json:"stock" example:"50"`
	Category    *Category `json:"category"`
	Supplier    *Supplier `json:"supplier"`
}

// GetProduct retrieves a product
// @Summary Get product (internal model)
// @Description Retrieves product information from internal models
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 404 {string} string "Product not found"
// @Router /internal/products/{id} [get]
func GetProduct() {}

// ListProducts lists all products
// @Summary List all products (internal model)
// @Description Lists all available products from internal catalog
// @Tags products
// @Produce json
// @Success 200 {array} Product
// @Router /internal/products [get]
func ListProducts() {}
