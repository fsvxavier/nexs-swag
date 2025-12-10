package generator

import (
	"os"
	"strings"
	"testing"

	"github.com/fsvxavier/nexs-swag/pkg/openapi"
)

func TestNew(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
	}
	outputDir := "/tmp/test"
	outputTypes := []string{"json", "yaml"}

	gen := New(spec, outputDir, outputTypes)
	if gen == nil {
		t.Fatal("New() returned nil")
	}

	if gen.spec != spec {
		t.Error("Generator spec not set correctly")
	}
	if gen.outputDir != outputDir {
		t.Error("Generator outputDir not set correctly")
	}
	if len(gen.outputType) != len(outputTypes) {
		t.Errorf("Generator outputType length = %d, want %d", len(gen.outputType), len(outputTypes))
	}
}

func TestSetInstanceName(t *testing.T) {
	gen := &Generator{
		spec: &openapi.OpenAPI{},
	}

	tests := []struct {
		name         string
		instanceName string
	}{
		{"default swagger", "swagger"},
		{"custom name", "api"},
		{"empty name", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen.SetInstanceName(tt.instanceName)
			if gen.instanceName != tt.instanceName {
				t.Errorf("SetInstanceName() = %q, want %q", gen.instanceName, tt.instanceName)
			}
		})
	}
}

func TestSetGeneratedTime(t *testing.T) {
	gen := &Generator{
		spec: &openapi.OpenAPI{},
	}

	gen.SetGeneratedTime(true)
	if !gen.generatedTime {
		t.Error("SetGeneratedTime(true) failed")
	}

	gen.SetGeneratedTime(false)
	if gen.generatedTime {
		t.Error("SetGeneratedTime(false) failed")
	}
}

func TestSetTemplateDelims(t *testing.T) {
	gen := &Generator{
		spec: &openapi.OpenAPI{},
	}

	tests := []struct {
		name      string
		delims    string
		wantLeft  string
		wantRight string
	}{
		{"default", "{{,}}", "{{", "}}"},
		{"custom", "<%,%>", "<%", "%>"},
		{"brackets", "[[,]]", "[[", "]]"},
		{"empty uses default", "", "{{", "}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen.SetTemplateDelims(tt.delims)
			if len(gen.templateDelims) != 2 {
				t.Fatalf("templateDelims length = %d, want 2", len(gen.templateDelims))
			}
			if gen.templateDelims[0] != tt.wantLeft {
				t.Errorf("leftDelim = %q, want %q", gen.templateDelims[0], tt.wantLeft)
			}
			if gen.templateDelims[1] != tt.wantRight {
				t.Errorf("rightDelim = %q, want %q", gen.templateDelims[1], tt.wantRight)
			}
		})
	}
}

func TestGenerateJSON(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"json"})
	err := gen.generateJSON()
	if err != nil {
		t.Fatalf("generateJSON() returned error: %v", err)
	}

	// Verify file was created
	outputFile := tmpDir + "/openapi.json"
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %s", outputFile)
	}
}

func TestGenerateYAML(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"yaml"})
	err := gen.generateYAML()
	if err != nil {
		t.Fatalf("generateYAML() returned error: %v", err)
	}

	// Verify file was created
	outputFile := tmpDir + "/openapi.yaml"
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %s", outputFile)
	}
}

func TestGenerateGo(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, tmpDir, []string{"go"})
	gen.SetInstanceName("testdocs")
	err := gen.generateGo()
	if err != nil {
		t.Fatalf("generateGo() returned error: %v", err)
	}

	// Verify file was created
	outputFile := tmpDir + "/docs.go"
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %s", outputFile)
	}

	// Read and verify content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "package testdocs") {
		t.Error("Generated file doesn't contain correct package name")
	}
	if !strings.Contains(contentStr, "SwaggerDoc") {
		t.Error("Generated file doesn't contain SwaggerDoc variable")
	}
	if !strings.Contains(contentStr, "ReadDoc") {
		t.Error("Generated file doesn't contain ReadDoc function")
	}
}

func TestGenerateGoWithTimestamp(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
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
	content, err := os.ReadFile(tmpDir + "/docs.go")
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "generated by nexs-swag at") {
		t.Error("Generated file doesn't contain timestamp")
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	tests := []struct {
		name        string
		outputTypes []string
		wantFiles   []string
	}{
		{
			name:        "json only",
			outputTypes: []string{"json"},
			wantFiles:   []string{"openapi.json"},
		},
		{
			name:        "yaml only",
			outputTypes: []string{"yaml"},
			wantFiles:   []string{"openapi.yaml"},
		},
		{
			name:        "yml alias",
			outputTypes: []string{"yml"},
			wantFiles:   []string{"openapi.yaml"},
		},
		{
			name:        "go only",
			outputTypes: []string{"go"},
			wantFiles:   []string{"docs.go"},
		},
		{
			name:        "multiple formats",
			outputTypes: []string{"json", "yaml", "go"},
			wantFiles:   []string{"openapi.json", "openapi.yaml", "docs.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subDir := tmpDir + "/" + tt.name
			gen := New(spec, subDir, tt.outputTypes)
			err := gen.Generate()
			if err != nil {
				t.Fatalf("Generate() returned error: %v", err)
			}

			// Verify all expected files were created
			for _, wantFile := range tt.wantFiles {
				filePath := subDir + "/" + wantFile
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file was not created: %s", wantFile)
				}
			}
		})
	}
}

func TestGenerateUnsupportedType(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
	}

	gen := New(spec, tmpDir, []string{"unsupported"})
	err := gen.Generate()
	if err == nil {
		t.Error("Generate() with unsupported type should return error")
	}
	if !strings.Contains(err.Error(), "unsupported output type") {
		t.Errorf("Error message should mention unsupported type, got: %v", err)
	}
}

func TestGenerateCreatesDirIfNotExists(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	nonExistentDir := tmpDir + "/newdir/subdir"

	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	gen := New(spec, nonExistentDir, []string{"json"})
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(nonExistentDir); os.IsNotExist(err) {
		t.Error("Output directory was not created")
	}

	// Verify file was created
	if _, err := os.Stat(nonExistentDir + "/openapi.json"); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestNewDefaults(t *testing.T) {
	t.Parallel()
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
	}

	gen := New(spec, "/tmp/test", []string{"json"})

	if gen.instanceName != "docs" {
		t.Errorf("Default instanceName = %q, want %q", gen.instanceName, "docs")
	}
	if gen.generatedTime != false {
		t.Error("Default generatedTime should be false")
	}
}

func TestSetTemplateDelimsInvalidFormat(t *testing.T) {
	t.Parallel()
	gen := &Generator{
		spec: &openapi.OpenAPI{},
	}

	// Test with single value (invalid format)
	gen.SetTemplateDelims("{{")
	if len(gen.templateDelims) != 2 {
		t.Error("Invalid delims format should default to {{,}}")
	}
	if gen.templateDelims[0] != "{{" || gen.templateDelims[1] != "}}" {
		t.Error("Invalid delims format should default to {{,}}")
	}

	// Test with three values (invalid format)
	gen.SetTemplateDelims("{{,}},extra")
	if len(gen.templateDelims) != 2 {
		t.Error("Invalid delims format should default to {{,}}")
	}
}

func TestGenerateComplexSpec(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:       "Complex API",
			Version:     "2.0.0",
			Description: "A complex API specification",
		},
		Servers: []openapi.Server{
			{
				URL:         "https://api.example.com",
				Description: "Production server",
			},
		},
		Paths: openapi.Paths{
			"/users": &openapi.PathItem{
				Get: &openapi.Operation{
					Summary: "List users",
					Responses: openapi.Responses{
						"200": &openapi.Response{
							Description: "Success",
						},
					},
				},
			},
		},
	}

	gen := New(spec, tmpDir, []string{"json", "yaml", "go"})
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() with complex spec returned error: %v", err)
	}

	// Verify all files created
	expectedFiles := []string{"openapi.json", "openapi.yaml", "docs.go"}
	for _, file := range expectedFiles {
		if _, err := os.Stat(tmpDir + "/" + file); os.IsNotExist(err) {
			t.Errorf("Expected file not created: %s", file)
		}
	}
}
