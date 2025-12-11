// Package generator defines the common interface for OpenAPI/Swagger generators.
package generator

// Generator is the common interface for all specification generators.
type Generator interface {
	// Generate generates all requested output formats.
	Generate() error

	// SetInstanceName sets the package name for generated Go files.
	SetInstanceName(name string)

	// SetGeneratedTime sets whether to include generation timestamp.
	SetGeneratedTime(enabled bool)
}
