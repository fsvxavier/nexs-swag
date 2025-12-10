package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDependencyPackage(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mock vendor directory
	vendorDir := filepath.Join(tmpDir, "vendor", "github.com", "test", "package")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatalf("Failed to create vendor directory: %v", err)
	}

	// Create a simple Go file
	goFile := filepath.Join(vendorDir, "types.go")
	content := `package testpkg

// User represents a user
type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write Go file: %v", err)
	}

	// Change to tmpDir so vendor is found
	t.Chdir(tmpDir)

	tests := []struct {
		name  string
		level int
	}{
		{"level 1 - models only", 1},
		{"level 2 - operations only", 2},
		{"level 3 - all", 3},
		{"level 0 - none", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.SetParseDependencyLevel(tt.level)

			err := p.parseDependencyPackage("github.com/test/package")
			if err != nil {
				t.Errorf("parseDependencyPackage() error = %v", err)
			}
		})
	}
}

func TestParseDependencyModels(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple Go file with model
	goFile := filepath.Join(tmpDir, "models.go")
	content := `package models

// Product represents a product
type Product struct {
	ID    int     ` + "`json:\"id\"`" + `
	Name  string  ` + "`json:\"name\"`" + `
	Price float64 ` + "`json:\"price\"`" + `
}

// Category represents a category
type Category struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write Go file: %v", err)
	}

	p := New()
	err := p.parseDependencyModels(tmpDir)
	if err != nil {
		t.Errorf("parseDependencyModels() error = %v", err)
	}
}

func TestParseDependencyOperations(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple Go file with operation
	goFile := filepath.Join(tmpDir, "handlers.go")
	content := `package handlers

import "net/http"

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieves user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Implementation
}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write Go file: %v", err)
	}

	p := New()
	err := p.parseDependencyOperations(tmpDir)
	if err != nil {
		t.Errorf("parseDependencyOperations() error = %v", err)
	}
}

func TestParseDependencyPackageNotFound(t *testing.T) {
	p := New()
	p.SetParseDependencyLevel(1)

	// Try to parse a non-existent package
	err := p.parseDependencyPackage("github.com/nonexistent/package")
	if err != nil {
		t.Errorf("parseDependencyPackage() should not error on missing package, got: %v", err)
	}
}

func TestParseDependenciesWithGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mock go.mod file
	goModContent := `module github.com/test/project

go 1.20

require (
	github.com/test/dep1 v1.0.0
	github.com/test/dep2 v2.1.0
)
`
	goModFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModFile, []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	// Change to tmpDir so go.mod is found
	t.Chdir(tmpDir)

	tests := []struct {
		name            string
		parseDependency bool
		dependencyLevel int
		shouldParse     bool
	}{
		{
			name:            "parse dependency disabled",
			parseDependency: false,
			dependencyLevel: 1,
			shouldParse:     false,
		},
		{
			name:            "level 0 - no parsing",
			parseDependency: true,
			dependencyLevel: 0,
			shouldParse:     false,
		},
		{
			name:            "level 1 - models only",
			parseDependency: true,
			dependencyLevel: 1,
			shouldParse:     true,
		},
		{
			name:            "level 2 - operations only",
			parseDependency: true,
			dependencyLevel: 2,
			shouldParse:     true,
		},
		{
			name:            "level 3 - all",
			parseDependency: true,
			dependencyLevel: 3,
			shouldParse:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.SetParseDependency(tt.parseDependency)
			p.SetParseDependencyLevel(tt.dependencyLevel)

			err := p.parseDependencies()
			if err != nil {
				t.Errorf("parseDependencies() error = %v", err)
			}
		})
	}
}

func TestParseDependenciesNoGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	// Change to tmpDir where there's no go.mod
	t.Chdir(tmpDir)

	p := New()
	p.SetParseDependency(true)
	p.SetParseDependencyLevel(1)

	// Should not error when go.mod is missing
	err := p.parseDependencies()
	if err != nil {
		t.Errorf("parseDependencies() should not error when go.mod is missing, got: %v", err)
	}
}
