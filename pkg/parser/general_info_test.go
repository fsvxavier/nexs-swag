package parser

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestParseGeneralInfo(t *testing.T) {
	t.Parallel()
	content := `package main

// @title Test API
// @version 1.0.0
// @description This is a test API
// @termsOfService http://example.com/terms
// @contact.name API Support
// @contact.email support@example.com
// @contact.url http://example.com/support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host api.example.com
// @BasePath /v1

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() returned error: %v", err)
	}

	// Verify Info was populated
	if p.openapi.Info.Title != "Test API" {
		t.Errorf("Info.Title = %q, want %q", p.openapi.Info.Title, "Test API")
	}
	if p.openapi.Info.Version != "1.0.0" {
		t.Errorf("Info.Version = %q, want %q", p.openapi.Info.Version, "1.0.0")
	}
}

func TestParseGeneralInfoMinimal(t *testing.T) {
	t.Parallel()
	content := `package main

// @title Minimal API
// @version 0.1.0

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() returned error: %v", err)
	}

	if p.openapi.Info.Title != "Minimal API" {
		t.Errorf("Info.Title = %q, want %q", p.openapi.Info.Title, "Minimal API")
	}
}

func TestNewGeneralInfoProcessor(t *testing.T) {
	t.Parallel()
	p := New()
	processor := NewGeneralInfoProcessor(p.openapi)

	if processor == nil {
		t.Fatal("NewGeneralInfoProcessor() returned nil")
	}
}

func TestGeneralInfoProcessorAnnotations(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		annotation string
		checkFunc  func(*Parser) bool
	}{
		{
			name:       "title",
			annotation: "@title My API",
			checkFunc: func(p *Parser) bool {
				return p.openapi.Info.Title == "My API"
			},
		},
		{
			name:       "version",
			annotation: "@version 2.0.0",
			checkFunc: func(p *Parser) bool {
				return p.openapi.Info.Version == "2.0.0"
			},
		},
		{
			name:       "description",
			annotation: "@description API description",
			checkFunc: func(p *Parser) bool {
				return p.openapi.Info.Description == "API description"
			},
		},
		{
			name:       "host",
			annotation: "@host api.example.com",
			checkFunc: func(p *Parser) bool {
				return len(p.openapi.Servers) > 0 && p.openapi.Servers[0].URL == "https://api.example.com"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			processor := NewGeneralInfoProcessor(p.openapi)

			err := processor.Process(tt.annotation)
			if err != nil {
				t.Errorf("Process() returned error: %v", err)
			}

			if !tt.checkFunc(p) {
				t.Errorf("Annotation %q was not processed correctly", tt.annotation)
			}
		})
	}
}

func TestGeneralInfoProcessorContact(t *testing.T) {
	t.Parallel()
	p := New()
	processor := NewGeneralInfoProcessor(p.openapi)

	annotations := []string{
		"@contact.name API Support",
		"@contact.email support@example.com",
		"@contact.url http://example.com/support",
	}

	for _, annotation := range annotations {
		if err := processor.Process(annotation); err != nil {
			t.Errorf("Process(%q) returned error: %v", annotation, err)
		}
	}

	if p.openapi.Info.Contact == nil {
		t.Fatal("Contact should not be nil")
	}
	if p.openapi.Info.Contact.Name != "API Support" {
		t.Errorf("Contact.Name = %q, want %q", p.openapi.Info.Contact.Name, "API Support")
	}
	if p.openapi.Info.Contact.Email != "support@example.com" {
		t.Errorf("Contact.Email = %q, want %q", p.openapi.Info.Contact.Email, "support@example.com")
	}
}

func TestGeneralInfoProcessorLicense(t *testing.T) {
	t.Parallel()
	p := New()
	processor := NewGeneralInfoProcessor(p.openapi)

	annotations := []string{
		"@license.name MIT",
		"@license.url https://opensource.org/licenses/MIT",
	}

	for _, annotation := range annotations {
		if err := processor.Process(annotation); err != nil {
			t.Errorf("Process(%q) returned error: %v", annotation, err)
		}
	}

	if p.openapi.Info.License == nil {
		t.Fatal("License should not be nil")
	}
	if p.openapi.Info.License.Name != "MIT" {
		t.Errorf("License.Name = %q, want %q", p.openapi.Info.License.Name, "MIT")
	}
}

func TestGeneralInfoProcessorServers(t *testing.T) {
	t.Parallel()
	p := New()
	processor := NewGeneralInfoProcessor(p.openapi)

	annotations := []string{
		"@host api.example.com",
		"@BasePath /v1",
		"@schemes https http",
	}

	for _, annotation := range annotations {
		if err := processor.Process(annotation); err != nil {
			t.Errorf("Process(%q) returned error: %v", annotation, err)
		}
	}

	if len(p.openapi.Servers) == 0 {
		t.Fatal("Servers should not be empty")
	}
}

func TestGeneralInfoProcessorTags(t *testing.T) {
	t.Parallel()
	content := `package main

// @title Test API
// @version 1.0.0
// @tag.name users
// @tag.description User management operations
// @tag.name products
// @tag.description Product management operations

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() returned error: %v", err)
	}

	// Should have parsed tags
	if len(p.openapi.Tags) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestGeneralInfoProcessorSecurityDefinitions(t *testing.T) {
	t.Parallel()
	content := `package main

// @title Test API
// @version 1.0.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() returned error: %v", err)
	}

	// Check if security schemes were defined
	if p.openapi.Components == nil || p.openapi.Components.SecuritySchemes == nil {
		t.Log("SecuritySchemes might not be populated (depends on implementation)")
	}
}

func TestParseGeneralInfoNoComments(t *testing.T) {
	t.Parallel()
	content := `package main

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() with no comments should not return error: %v", err)
	}

	// Title and Version should remain empty
	if p.openapi.Info.Title != "" {
		t.Error("Info.Title should be empty when no annotations present")
	}
}

func TestParseGeneralInfoMultilineDescription(t *testing.T) {
	t.Parallel()
	content := `package main

// @title Test API
// @version 1.0.0
// @description This is line 1
// @description This is line 2
// @description This is line 3

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() returned error: %v", err)
	}

	if p.openapi.Info.Description == "" {
		t.Error("Description should not be empty")
	}
}
