package models

// Category is referenced by Product (transitive dependency)
type Category struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Electronics"`
}

// Supplier is also referenced by Product (transitive dependency)
type Supplier struct {
	ID      int    `json:"id" example:"100"`
	Name    string `json:"name" example:"Tech Supplies Inc"`
	Contact string `json:"contact" example:"contact@techsupplies.com"`
}

// UnrelatedModel is NOT referenced anywhere
type UnrelatedModel struct {
	Data string `json:"data"`
}
