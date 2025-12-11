package v2

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	swagger "github.com/fsvxavier/nexs-swag/pkg/openapi/v2"
)

func TestNew(t *testing.T) {
	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, "./docs", []string{"json", "yaml"})

	if gen == nil {
		t.Fatal("New() returned nil")
	}

	if gen.spec != spec {
		t.Error("spec not set correctly")
	}

	if gen.outputDir != "./docs" {
		t.Errorf("outputDir = %v, want ./docs", gen.outputDir)
	}

	if len(gen.outputType) != 2 {
		t.Errorf("len(outputType) = %v, want 2", len(gen.outputType))
	}

	if gen.instanceName != "docs" {
		t.Errorf("instanceName = %v, want docs", gen.instanceName)
	}

	if gen.generatedTime != false {
		t.Error("generatedTime should default to false")
	}
}

func TestSetInstanceName(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		newName  string
		expected string
	}{
		{
			name:     "default_swagger",
			initial:  "docs",
			newName:  "swagger",
			expected: "swagger",
		},
		{
			name:     "custom_name",
			initial:  "docs",
			newName:  "myapi",
			expected: "myapi",
		},
		{
			name:     "empty_name",
			initial:  "docs",
			newName:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := &Generator{instanceName: tt.initial}
			gen.SetInstanceName(tt.newName)

			if gen.instanceName != tt.expected {
				t.Errorf("SetInstanceName() instanceName = %v, want %v", gen.instanceName, tt.expected)
			}
		})
	}
}

func TestSetGeneratedTime(t *testing.T) {
	gen := &Generator{}

	if gen.generatedTime != false {
		t.Error("generatedTime should default to false")
	}

	gen.SetGeneratedTime(true)

	if gen.generatedTime != true {
		t.Error("SetGeneratedTime(true) failed")
	}

	gen.SetGeneratedTime(false)

	if gen.generatedTime != false {
		t.Error("SetGeneratedTime(false) failed")
	}
}

func TestGenerateJSON(t *testing.T) {
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"json"})
	err := gen.generateJSON()
	if err != nil {
		t.Fatalf("generateJSON() returned error: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "swagger.json")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// Verify JSON content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result swagger.Swagger
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if result.Swagger != "2.0" {
		t.Errorf("swagger version = %v, want 2.0", result.Swagger)
	}

	if result.Info.Title != "Test API" {
		t.Errorf("title = %v, want Test API", result.Info.Title)
	}
}

func TestGenerateYAML(t *testing.T) {
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"yaml"})
	err := gen.generateYAML()
	if err != nil {
		t.Fatalf("generateYAML() returned error: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "swagger.yaml")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// Verify YAML content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result swagger.Swagger
	if err := yaml.Unmarshal(content, &result); err != nil {
		t.Fatalf("Output is not valid YAML: %v", err)
	}

	if result.Swagger != "2.0" {
		t.Errorf("swagger version = %v, want 2.0", result.Swagger)
	}
}

func TestGenerateGo(t *testing.T) {
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"go"})
	err := gen.generateGo()
	if err != nil {
		t.Fatalf("generateGo() returned error: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "docs.go")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// Verify content structure
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Check for required elements
	requiredStrings := []string{
		"package docs",
		"SwaggerDoc",
		"ReadDoc",
		"2.0",
	}

	for _, str := range requiredStrings {
		if !strings.Contains(contentStr, str) {
			t.Errorf("Generated file missing required string: %s", str)
		}
	}
}

func TestGenerateGoWithTimestamp(t *testing.T) {
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"go"})
	gen.SetGeneratedTime(true)
	err := gen.generateGo()
	if err != nil {
		t.Fatalf("generateGo() with timestamp returned error: %v", err)
	}

	// Read and verify content includes timestamp
	content, err := os.ReadFile(filepath.Join(tmpDir, "docs.go"))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Generated at:") {
		t.Error("Generated file doesn't contain timestamp")
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		outputTypes []string
		wantFiles   []string
	}{
		{
			name:        "json only",
			outputTypes: []string{"json"},
			wantFiles:   []string{"swagger.json"},
		},
		{
			name:        "yaml only",
			outputTypes: []string{"yaml"},
			wantFiles:   []string{"swagger.yaml"},
		},
		{
			name:        "yml alias",
			outputTypes: []string{"yml"},
			wantFiles:   []string{"swagger.yaml"},
		},
		{
			name:        "go only",
			outputTypes: []string{"go"},
			wantFiles:   []string{"docs.go"},
		},
		{
			name:        "multiple formats",
			outputTypes: []string{"json", "yaml", "go"},
			wantFiles:   []string{"swagger.json", "swagger.yaml", "docs.go"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := filepath.Join(t.TempDir(), tt.name)

			spec := &swagger.Swagger{
				Swagger: "2.0",
				Info: swagger.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
			}

			gen := New(spec, tmpDir, tt.outputTypes)
			err := gen.Generate()
			if err != nil {
				t.Fatalf("Generate() returned error: %v", err)
			}

			// Verify all expected files were created
			for _, file := range tt.wantFiles {
				fullPath := filepath.Join(tmpDir, file)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("Expected file not created: %s", file)
				}
			}
		})
	}
}

func TestGenerateUnsupportedType(t *testing.T) {
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"xml"})
	err := gen.Generate()
	if err == nil {
		t.Error("Generate() should return error for unsupported format")
	}

	if !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("Error message should mention unsupported format, got: %v", err)
	}
}

func TestGenerateCreatesDirIfNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentDir := filepath.Join(tmpDir, "newdir")

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, nonExistentDir, []string{"json"})
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(nonExistentDir, "swagger.json")); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestNewDefaults(t *testing.T) {
	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, "./docs", []string{"json"})

	if gen.instanceName != "docs" {
		t.Errorf("Default instanceName = %v, want docs", gen.instanceName)
	}

	if gen.generatedTime != false {
		t.Error("Default generatedTime should be false")
	}
}

func TestGenerateComplexSpec(t *testing.T) {
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:       "Complex API",
			Description: "A complex API with multiple endpoints",
			Version:     "2.0.0",
			Contact: &swagger.Contact{
				Name:  "API Support",
				Email: "support@example.com",
			},
		},
		Host:     "api.example.com",
		BasePath: "/v2",
		Schemes:  []string{"https"},
		Paths: map[string]*swagger.PathItem{
			"/users": {
				Get: &swagger.Operation{
					Summary:     "List users",
					Description: "Returns a list of users",
					OperationID: "listUsers",
					Responses: map[string]*swagger.Response{
						"200": {
							Description: "Successful response",
						},
					},
				},
			},
		},
	}

	gen := New(spec, tmpDir, []string{"json", "yaml", "go"})
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	expectedFiles := []string{"swagger.json", "swagger.yaml", "docs.go"}
	for _, file := range expectedFiles {
		fullPath := filepath.Join(tmpDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file not created: %s", file)
		}
	}

	// Verify JSON content has complex structure
	jsonContent, err := os.ReadFile(filepath.Join(tmpDir, "swagger.json"))
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	if !strings.Contains(string(jsonContent), "api.example.com") {
		t.Error("JSON output missing host information")
	}

	if !strings.Contains(string(jsonContent), "/users") {
		t.Error("JSON output missing paths information")
	}
}
