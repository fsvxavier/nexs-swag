package v3

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	v3 "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

func TestNew(t *testing.T) {
	spec := &v3.OpenAPI{
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
		spec: &v3.OpenAPI{},
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
		spec: &v3.OpenAPI{},
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
		spec: &v3.OpenAPI{},
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

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
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

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
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

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
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

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
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
	if !strings.Contains(contentStr, "Generated at:") {
		t.Error("Generated file doesn't contain timestamp")
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
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

	spec := &v3.OpenAPI{
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

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
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
	spec := &v3.OpenAPI{
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
		spec: &v3.OpenAPI{},
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

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:       "Complex API",
			Version:     "2.0.0",
			Description: "A complex API specification",
		},
		Servers: []v3.Server{
			{
				URL:         "https://api.example.com",
				Description: "Production server",
			},
		},
		Paths: v3.Paths{
			"/users": &v3.PathItem{
				Get: &v3.Operation{
					Summary: "List users",
					Responses: v3.Responses{
						"200": &v3.Response{
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

func TestSetOpenAPIVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		initialVersion  string
		setVersion      string
		expectedVersion string
	}{
		{
			name:            "Set version 3.0.0",
			initialVersion:  "3.1.0",
			setVersion:      "3.0.0",
			expectedVersion: "3.0.0",
		},
		{
			name:            "Set version 3.0.1",
			initialVersion:  "3.1.0",
			setVersion:      "3.0.1",
			expectedVersion: "3.0.1",
		},
		{
			name:            "Set version 3.0.3",
			initialVersion:  "3.1.0",
			setVersion:      "3.0.3",
			expectedVersion: "3.0.3",
		},
		{
			name:            "Set version 3.1.0",
			initialVersion:  "3.0.0",
			setVersion:      "3.1.0",
			expectedVersion: "3.1.0",
		},
		{
			name:            "Set version 3.2.0",
			initialVersion:  "3.1.0",
			setVersion:      "3.2.0",
			expectedVersion: "3.2.0",
		},
		{
			name:            "Empty version does not update",
			initialVersion:  "3.1.0",
			setVersion:      "",
			expectedVersion: "3.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()

			spec := &v3.OpenAPI{
				OpenAPI: tt.initialVersion,
				Info: v3.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
				Paths: v3.Paths{},
			}

			gen := New(spec, tmpDir, []string{"json"})
			gen.SetOpenAPIVersion(tt.setVersion)
			if err := gen.Generate(); err != nil {
				t.Fatalf("Generate() failed: %v", err)
			}

			// Read generated JSON to verify version
			data, err := os.ReadFile(filepath.Join(tmpDir, "openapi.json"))
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			var result v3.OpenAPI
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if result.OpenAPI != tt.expectedVersion {
				t.Errorf("OpenAPI version = %q, want %q", result.OpenAPI, tt.expectedVersion)
			}
		})
	}
}

func TestGenerateUnsupportedOutputType(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: v3.Paths{},
	}

	gen := New(spec, tmpDir, []string{"xml"})
	err := gen.Generate()
	if err == nil {
		t.Fatal("Generate() with unsupported output type should return error")
	}

	expectedMsg := "unsupported output type: xml"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Error message = %q, want to contain %q", err.Error(), expectedMsg)
	}
}

func TestGenerateInvalidOutputDir(t *testing.T) {
	t.Parallel()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: v3.Paths{},
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

func TestGenerateJSONFileCreationError(t *testing.T) {
	t.Parallel()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: v3.Paths{},
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

	if !strings.Contains(err.Error(), "failed to generate JSON") {
		t.Errorf("Error message = %q, want to contain 'failed to generate JSON'", err.Error())
	}
}

func TestGenerateYAMLFileCreationError(t *testing.T) {
	t.Parallel()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: v3.Paths{},
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

	if !strings.Contains(err.Error(), "failed to generate YAML") {
		t.Errorf("Error message = %q, want to contain 'failed to generate YAML'", err.Error())
	}
}

func TestGenerateGoFileCreationError(t *testing.T) {
	t.Parallel()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: v3.Paths{},
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

	if !strings.Contains(err.Error(), "failed to generate Go file") {
		t.Errorf("Error message = %q, want to contain 'failed to generate Go file'", err.Error())
	}
}

func TestGenerateYMLExtension(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: v3.Paths{},
	}

	gen := New(spec, tmpDir, []string{"yml"})
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() with 'yml' extension failed: %v", err)
	}

	// Verify YAML file was created
	yamlPath := filepath.Join(tmpDir, "openapi.yaml")
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		t.Error("Expected openapi.yaml to be created with 'yml' output type")
	}
}

func TestGenerateGoWithGeneratedTime(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: v3.Paths{},
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

func TestGenerateGoWithCustomInstanceName(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: v3.Paths{},
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

func TestGenerateGoWithTemplateDelims(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:       "Test API",
			Version:     "1.0.0",
			Description: "API with {{variable}} placeholder",
		},
		Paths: v3.Paths{},
	}

	gen := New(spec, tmpDir, []string{"go"})
	gen.SetTemplateDelims("[[,]]")
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() with template delims failed: %v", err)
	}

	// Read generated Go file
	data, err := os.ReadFile(filepath.Join(tmpDir, "docs.go"))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "[[variable]]") {
		t.Error("Generated Go file should replace {{ }} with custom delimiters")
	}
	if strings.Contains(content, "{{variable}}") {
		t.Error("Generated Go file should not contain original {{ }} delimiters")
	}
}

func TestGenerateMultipleOutputTypes(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: v3.Paths{},
	}

	gen := New(spec, tmpDir, []string{"json", "yaml", "go"})
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() with multiple output types failed: %v", err)
	}

	// Verify all files created
	expectedFiles := []string{"openapi.json", "openapi.yaml", "docs.go"}
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

	spec := &v3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: v3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: v3.Paths{},
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
