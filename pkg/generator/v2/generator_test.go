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

func TestGenerateUnsupportedFormat(t *testing.T) {
	t.Parallel()
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
		t.Fatal("Generate() with unsupported format should return error")
	}

	expectedMsg := "unsupported output format: xml"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Error message = %q, want to contain %q", err.Error(), expectedMsg)
	}
}

func TestGenerateInvalidOutputDirectory(t *testing.T) {
	t.Parallel()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	// Use a path that cannot be created (file exists with same name)
	tmpDir := t.TempDir()
	blockingFile := filepath.Join(tmpDir, "blocked")
	if err := os.WriteFile(blockingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create blocking file: %v", err)
	}

	gen := New(spec, blockingFile, []string{"json"})
	err := gen.Generate()
	if err == nil {
		t.Fatal("Generate() with invalid output directory should return error")
	}

	if !strings.Contains(err.Error(), "failed to create output directory") {
		t.Errorf("Error message = %q, want to contain 'failed to create output directory'", err.Error())
	}
}

func TestGenerateJSONWriteError(t *testing.T) {
	t.Parallel()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	// Create a directory with read-only permissions
	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0444); err != nil {
		t.Fatalf("Failed to create readonly directory: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0755) // Restore permissions for cleanup

	gen := New(spec, readOnlyDir, []string{"json"})
	err := gen.Generate()
	if err == nil {
		t.Fatal("Generate() with readonly output directory should return error")
	}

	if !strings.Contains(err.Error(), "failed to write JSON file") {
		t.Errorf("Error message = %q, want to contain 'failed to write JSON file'", err.Error())
	}
}

func TestGenerateYAMLWriteError(t *testing.T) {
	t.Parallel()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	// Create a directory with read-only permissions
	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0444); err != nil {
		t.Fatalf("Failed to create readonly directory: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0755) // Restore permissions for cleanup

	gen := New(spec, readOnlyDir, []string{"yaml"})
	err := gen.Generate()
	if err == nil {
		t.Fatal("Generate() with readonly output directory should return error")
	}

	if !strings.Contains(err.Error(), "failed to write YAML file") {
		t.Errorf("Error message = %q, want to contain 'failed to write YAML file'", err.Error())
	}
}

func TestGenerateGoFileCreationError(t *testing.T) {
	t.Parallel()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	// Create a directory with read-only permissions
	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0444); err != nil {
		t.Fatalf("Failed to create readonly directory: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0755) // Restore permissions for cleanup

	gen := New(spec, readOnlyDir, []string{"go"})
	err := gen.Generate()
	if err == nil {
		t.Fatal("Generate() with readonly output directory should return error")
	}

	if !strings.Contains(err.Error(), "failed to create file") {
		t.Errorf("Error message = %q, want to contain 'failed to create file'", err.Error())
	}
}

func TestGenerateYMLExtension(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"yml"})
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() with 'yml' extension failed: %v", err)
	}

	// Verify YAML file was created
	yamlPath := filepath.Join(tmpDir, "swagger.yaml")
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		t.Error("Expected swagger.yaml to be created with 'yml' output type")
	}
}

func TestGenerateWithCustomInstanceName(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"go"})
	gen.SetInstanceName("customdocs")
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() with custom instance name failed: %v", err)
	}

	// Read generated Go file
	data, err := os.ReadFile(filepath.Join(tmpDir, "docs.go"))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "package customdocs") {
		t.Error("Generated Go file should contain 'package customdocs'")
	}
	if !strings.Contains(content, "Package customdocs") {
		t.Error("Generated Go file should contain package comment with custom name")
	}
}

func TestGenerateGoWithGeneratedTimeEnabled(t *testing.T) {
	t.Parallel()
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
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() with generatedTime enabled failed: %v", err)
	}

	// Read generated Go file
	data, err := os.ReadFile(filepath.Join(tmpDir, "docs.go"))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Generated at:") {
		t.Error("Generated Go file should contain 'Generated at:' timestamp when generatedTime is enabled")
	}
}

func TestGenerateMultipleFormats(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"json", "yaml", "go"})
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() with multiple formats failed: %v", err)
	}

	// Verify all files created
	expectedFiles := []string{"swagger.json", "swagger.yaml", "docs.go"}
	for _, file := range expectedFiles {
		filePath := filepath.Join(tmpDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file not created: %s", file)
		}
	}
}

func TestGenerateGoReadDocFunction(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"go"})
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Read generated Go file
	data, err := os.ReadFile(filepath.Join(tmpDir, "docs.go"))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "func ReadDoc() string {") {
		t.Error("Generated Go file should contain ReadDoc function")
	}
	if !strings.Contains(content, "return SwaggerDoc") {
		t.Error("ReadDoc function should return SwaggerDoc")
	}
}

func TestGenerateFormatCaseInsensitive(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	// Test with uppercase and mixed case formats
	gen := New(spec, tmpDir, []string{"JSON", "YaML", "Go"})
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() with mixed case formats failed: %v", err)
	}

	// Verify all files created
	expectedFiles := []string{"swagger.json", "swagger.yaml", "docs.go"}
	for _, file := range expectedFiles {
		filePath := filepath.Join(tmpDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file not created: %s", file)
		}
	}
}

func TestGenerateJSONIndentation(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: map[string]*swagger.PathItem{
			"/test": {
				Get: &swagger.Operation{
					Summary: "Test endpoint",
				},
			},
		},
	}

	gen := New(spec, tmpDir, []string{"json"})
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Read generated JSON file
	data, err := os.ReadFile(filepath.Join(tmpDir, "swagger.json"))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	content := string(data)
	// Check for proper indentation (2 spaces)
	if !strings.Contains(content, "  \"swagger\"") {
		t.Error("JSON should be indented with 2 spaces")
	}
}

func TestGenerateYAMLStructure(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:       "Test API",
			Version:     "1.0.0",
			Description: "Test description",
		},
		Host:     "api.example.com",
		BasePath: "/v1",
		Schemes:  []string{"https"},
	}

	gen := New(spec, tmpDir, []string{"yaml"})
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Read and parse YAML
	data, err := os.ReadFile(filepath.Join(tmpDir, "swagger.yaml"))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	var result swagger.Swagger
	if err := yaml.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Verify structure
	if result.Host != "api.example.com" {
		t.Errorf("Host = %q, want 'api.example.com'", result.Host)
	}
	if result.BasePath != "/v1" {
		t.Errorf("BasePath = %q, want '/v1'", result.BasePath)
	}
	if len(result.Schemes) != 1 || result.Schemes[0] != "https" {
		t.Errorf("Schemes = %v, want [https]", result.Schemes)
	}
}

func TestFilterSpecByVisibilityV2(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                   string
		visibility             string
		opPublicPath           string
		opPrivatePath          string
		opSharedPath           string
		expectPublicInPublic   bool
		expectPrivateInPublic  bool
		expectSharedInPublic   bool
		expectPublicInPrivate  bool
		expectPrivateInPrivate bool
		expectSharedInPrivate  bool
	}{
		{
			name:                   "public filter includes only public and shared",
			visibility:             "public",
			opPublicPath:           "/public",
			opPrivatePath:          "/private",
			opSharedPath:           "/shared",
			expectPublicInPublic:   true,
			expectPrivateInPublic:  false,
			expectSharedInPublic:   true,
			expectPublicInPrivate:  true,
			expectPrivateInPrivate: true,
			expectSharedInPrivate:  true,
		},
		{
			name:                   "private filter includes all operations",
			visibility:             "private",
			opPublicPath:           "/api/public",
			opPrivatePath:          "/api/private",
			opSharedPath:           "/api/shared",
			expectPublicInPublic:   true,
			expectPrivateInPublic:  false,
			expectSharedInPublic:   true,
			expectPublicInPrivate:  true,
			expectPrivateInPrivate: true,
			expectSharedInPrivate:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()

			spec := &swagger.Swagger{
				Swagger: "2.0",
				Info: swagger.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
				Paths: swagger.Paths{
					tt.opPublicPath: &swagger.PathItem{
						Get: &swagger.Operation{
							Summary: "Public endpoint",
							Responses: map[string]*swagger.Response{
								"200": {Description: "Success"},
							},
							Extensions: map[string]interface{}{
								"x-visibility": "public",
							},
						},
					},
					tt.opPrivatePath: &swagger.PathItem{
						Get: &swagger.Operation{
							Summary: "Private endpoint",
							Responses: map[string]*swagger.Response{
								"200": {Description: "Success"},
							},
							Extensions: map[string]interface{}{
								"x-visibility": "private",
							},
						},
					},
					tt.opSharedPath: &swagger.PathItem{
						Get: &swagger.Operation{
							Summary: "Shared endpoint (no visibility annotation)",
							Responses: map[string]*swagger.Response{
								"200": {Description: "Success"},
							},
						},
					},
				},
				Definitions: make(map[string]*swagger.Schema),
			}

			gen := New(spec, tmpDir, []string{"json"})

			// Test public filter
			publicSpec := gen.filterSpecByVisibility("public")
			_, hasPublic := publicSpec.Paths[tt.opPublicPath]
			_, hasPrivate := publicSpec.Paths[tt.opPrivatePath]
			_, hasShared := publicSpec.Paths[tt.opSharedPath]

			if hasPublic != tt.expectPublicInPublic {
				t.Errorf("Public filter: public path presence = %v, want %v", hasPublic, tt.expectPublicInPublic)
			}
			if hasPrivate != tt.expectPrivateInPublic {
				t.Errorf("Public filter: private path presence = %v, want %v", hasPrivate, tt.expectPrivateInPublic)
			}
			if hasShared != tt.expectSharedInPublic {
				t.Errorf("Public filter: shared path presence = %v, want %v", hasShared, tt.expectSharedInPublic)
			}

			// Test private filter
			privateSpec := gen.filterSpecByVisibility("private")
			hasPublicPrivate, hasPrivatePrivate, hasSharedPrivate := privateSpec.Paths[tt.opPublicPath], privateSpec.Paths[tt.opPrivatePath], privateSpec.Paths[tt.opSharedPath]

			if (hasPublicPrivate != nil) != tt.expectPublicInPrivate {
				t.Errorf("Private filter: public path presence = %v, want %v", hasPublicPrivate != nil, tt.expectPublicInPrivate)
			}
			if (hasPrivatePrivate != nil) != tt.expectPrivateInPrivate {
				t.Errorf("Private filter: private path presence = %v, want %v", hasPrivatePrivate != nil, tt.expectPrivateInPrivate)
			}
			if (hasSharedPrivate != nil) != tt.expectSharedInPrivate {
				t.Errorf("Private filter: shared path presence = %v, want %v", hasSharedPrivate != nil, tt.expectSharedInPrivate)
			}
		})
	}
}

func TestGenerateSeparateSpecsV2(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Visibility Test API",
			Version: "1.0.0",
		},
		Paths: swagger.Paths{
			"/users": &swagger.PathItem{
				Get: &swagger.Operation{
					Summary: "Public get users",
					Responses: map[string]*swagger.Response{
						"200": {
							Description: "Success",
							Schema: &swagger.Schema{
								Ref: "#/definitions/UserPublic",
							},
						},
					},
					Extensions: map[string]interface{}{
						"x-visibility": "public",
					},
				},
			},
			"/admin/users": &swagger.PathItem{
				Get: &swagger.Operation{
					Summary: "Private admin users",
					Responses: map[string]*swagger.Response{
						"200": {
							Description: "Success",
							Schema: &swagger.Schema{
								Ref: "#/definitions/UserPrivate",
							},
						},
					},
					Extensions: map[string]interface{}{
						"x-visibility": "private",
					},
				},
			},
			"/shared/endpoint": &swagger.PathItem{
				Get: &swagger.Operation{
					Summary: "Shared endpoint",
					Responses: map[string]*swagger.Response{
						"200": {Description: "Success"},
					},
				},
			},
		},
		Definitions: map[string]*swagger.Schema{
			"UserPublic": {
				Type: "object",
				Properties: map[string]*swagger.Schema{
					"id":   {Type: "integer"},
					"name": {Type: "string"},
				},
			},
			"UserPrivate": {
				Type: "object",
				Properties: map[string]*swagger.Schema{
					"id":       {Type: "integer"},
					"name":     {Type: "string"},
					"email":    {Type: "string"},
					"password": {Type: "string"},
				},
			},
		},
	}

	gen := New(spec, tmpDir, []string{"json"})
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Verify public spec
	publicData, err := os.ReadFile(filepath.Join(tmpDir, "swagger_public.json"))
	if err != nil {
		t.Fatalf("Failed to read public spec: %v", err)
	}

	var publicSpec swagger.Swagger
	if err := json.Unmarshal(publicData, &publicSpec); err != nil {
		t.Fatalf("Failed to parse public spec: %v", err)
	}

	// Public spec should have /users and /shared/endpoint but NOT /admin/users
	if _, ok := publicSpec.Paths["/users"]; !ok {
		t.Error("Public spec missing /users path")
	}
	if _, ok := publicSpec.Paths["/shared/endpoint"]; !ok {
		t.Error("Public spec missing /shared/endpoint path")
	}
	if _, ok := publicSpec.Paths["/admin/users"]; ok {
		t.Error("Public spec should NOT contain /admin/users path")
	}

	// Verify private spec
	privateData, err := os.ReadFile(filepath.Join(tmpDir, "swagger_private.json"))
	if err != nil {
		t.Fatalf("Failed to read private spec: %v", err)
	}

	var privateSpec swagger.Swagger
	if err := json.Unmarshal(privateData, &privateSpec); err != nil {
		t.Fatalf("Failed to parse private spec: %v", err)
	}

	// Private spec should have ALL paths: /users, /admin/users, and /shared/endpoint
	if _, ok := privateSpec.Paths["/users"]; !ok {
		t.Error("Private spec missing /users path (should include public endpoints)")
	}
	if _, ok := privateSpec.Paths["/admin/users"]; !ok {
		t.Error("Private spec missing /admin/users path")
	}
	if _, ok := privateSpec.Paths["/shared/endpoint"]; !ok {
		t.Error("Private spec missing /shared/endpoint path")
	}

	// Verify schema filtering
	if _, ok := publicSpec.Definitions["UserPublic"]; !ok {
		t.Error("Public spec missing UserPublic definition")
	}
	if _, ok := publicSpec.Definitions["UserPrivate"]; ok {
		t.Error("Public spec should NOT contain UserPrivate definition")
	}

	if _, ok := privateSpec.Definitions["UserPrivate"]; !ok {
		t.Error("Private spec missing UserPrivate definition")
	}
	if _, ok := privateSpec.Definitions["UserPublic"]; !ok {
		t.Error("Private spec missing UserPublic definition (should include public schemas)")
	}
}

func TestHasVisibilityAnnotationsV2(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		spec     *swagger.Swagger
		expected bool
	}{
		{
			name: "spec with x-visibility",
			spec: &swagger.Swagger{
				Paths: swagger.Paths{
					"/test": &swagger.PathItem{
						Get: &swagger.Operation{
							Extensions: map[string]interface{}{
								"x-visibility": "public",
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "spec without x-visibility",
			spec: &swagger.Swagger{
				Paths: swagger.Paths{
					"/test": &swagger.PathItem{
						Get: &swagger.Operation{},
					},
				},
			},
			expected: false,
		},
		{
			name: "spec with empty paths",
			spec: &swagger.Swagger{
				Paths: swagger.Paths{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gen := &Generator{spec: tt.spec}
			result := gen.hasVisibilityAnnotations()
			if result != tt.expected {
				t.Errorf("hasVisibilityAnnotations() = %v, want %v", result, tt.expected)
			}
		})
	}
}
