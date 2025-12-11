// Package openapi provides common interfaces for OpenAPI/Swagger specifications.
package openapi

// Specification is a common interface for all OpenAPI/Swagger specification versions.
type Specification interface {
	// GetVersion returns the specification version (e.g., "2.0", "3.0.0", "3.1.0")
	GetVersion() string

	// GetTitle returns the API title
	GetTitle() string

	// GetInfo returns the Info object
	GetInfo() interface{}

	// Validate performs validation of the specification
	Validate() error

	// MarshalJSON serializes the specification to JSON
	MarshalJSON() ([]byte, error)
}
