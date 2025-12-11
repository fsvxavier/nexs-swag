package parser

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
	"time"

	v3 "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

func TestNewParser(t *testing.T) {
	t.Parallel()
	p := New()

	if p == nil {
		t.Fatal("New() returned nil")
	}
	if p.openapi == nil {
		t.Error("openapi should be initialized")
	}
	if p.openapi.OpenAPI != "3.1.0" {
		t.Errorf("OpenAPI version = %q, want %q", p.openapi.OpenAPI, "3.1.0")
	}
	if p.files == nil {
		t.Error("files map should be initialized")
	}
	if p.typeCache == nil {
		t.Error("typeCache should be initialized")
	}
	if p.propertyStrategy != "camelcase" {
		t.Errorf("default propertyStrategy = %q, want %q", p.propertyStrategy, "camelcase")
	}
	if p.parseDepth != 100 {
		t.Errorf("default parseDepth = %d, want 100", p.parseDepth)
	}
}

func TestGetOpenAPI(t *testing.T) {
	t.Parallel()
	p := New()
	p.openapi.Info.Title = "Test API"
	p.openapi.Info.Version = "1.0.0"

	spec := p.GetOpenAPI()
	if spec == nil {
		t.Fatal("GetOpenAPI() returned nil")
	}
	if spec.Info.Title != "Test API" {
		t.Errorf("Info.Title = %q, want %q", spec.Info.Title, "Test API")
	}
}

func TestGetOpenAPIWithGeneratedTime(t *testing.T) {
	t.Parallel()
	p := New()
	p.openapi.Info.Version = "1.0.0"
	p.generatedTime = true

	spec := p.GetOpenAPI()
	if spec.Info.Version == "1.0.0" {
		t.Error("Version should include generated timestamp when generatedTime is true")
	}
}

func TestParseFile(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create a simple Go file
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

// @title Test API
// @version 1.0

// @Summary Get user
// @Router /users [get]
func GetUser() {}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p := New()
	p.SetGeneralInfoFile(testFile)

	err := p.ParseFile(testFile)
	if err != nil {
		t.Errorf("ParseFile() returned error: %v", err)
	}

	if len(p.files) != 1 {
		t.Errorf("files map length = %d, want 1", len(p.files))
	}
}

func TestParseFileInvalid(t *testing.T) {
	t.Parallel()
	p := New()

	err := p.ParseFile("/nonexistent/file.go")
	if err == nil {
		t.Error("ParseFile() with invalid file should return error")
	}
}

func TestParseFileInvalidSyntax(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "invalid.go")
	content := `package main

func Invalid( {
	// missing closing brace
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p := New()
	err := p.ParseFile(testFile)
	if err == nil {
		t.Error("ParseFile() with invalid Go syntax should return error")
	}
}

func TestParseDirSimple(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create a simple Go file
	testFile := filepath.Join(tmpDir, "main.go")
	content := `package main

// @title Test API
// @version 1.0.0

func main() {}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p := New()
	err := p.ParseDir(tmpDir)
	if err != nil {
		t.Errorf("ParseDir() returned error: %v", err)
	}
}

func TestParseDirInvalidPath(t *testing.T) {
	t.Parallel()
	p := New()

	err := p.ParseDir("/nonexistent/path")
	if err == nil {
		t.Error("ParseDir() with invalid path should return error")
	}
}

func TestParseDirWithVendor(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create vendor directory
	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.Mkdir(vendorDir, 0755); err != nil {
		t.Fatalf("Failed to create vendor dir: %v", err)
	}

	vendorFile := filepath.Join(vendorDir, "vendor.go")
	if err := os.WriteFile(vendorFile, []byte("package vendor"), 0644); err != nil {
		t.Fatalf("Failed to write vendor file: %v", err)
	}

	// Create main file
	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte("package main\nfunc main() {}"), 0644); err != nil {
		t.Fatalf("Failed to write main file: %v", err)
	}

	p := New()
	p.SetParseVendor(false)

	err := p.ParseDir(tmpDir)
	if err != nil {
		t.Errorf("ParseDir() returned error: %v", err)
	}
}

func TestParseDirWithInternal(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create internal directory
	internalDir := filepath.Join(tmpDir, "internal")
	if err := os.Mkdir(internalDir, 0755); err != nil {
		t.Fatalf("Failed to create internal dir: %v", err)
	}

	internalFile := filepath.Join(internalDir, "internal.go")
	if err := os.WriteFile(internalFile, []byte("package internal"), 0644); err != nil {
		t.Fatalf("Failed to write internal file: %v", err)
	}

	// Create main file
	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte("package main\nfunc main() {}"), 0644); err != nil {
		t.Fatalf("Failed to write main file: %v", err)
	}

	p := New()
	p.SetParseInternal(false)

	err := p.ParseDir(tmpDir)
	if err != nil {
		t.Errorf("ParseDir() returned error: %v", err)
	}
}

func TestShouldExclude(t *testing.T) {
	t.Parallel()
	p := New()
	p.SetExcludePatterns([]string{"*.pb.go", "mock_*.go", "vendor"})

	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		filename string
		isDir    bool
		expected bool
	}{
		{"protobuf file", "test.pb.go", false, true},
		{"mock file", "mock_user.go", false, true},
		{"vendor dir", "vendor", true, true},
		{"normal file", "main.go", false, false},
		{"normal dir", "pkg", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.filename)
			info := &mockFileInfo{name: tt.filename, isDir: tt.isDir}

			result := p.shouldExclude(path, info)
			if result != tt.expected {
				t.Errorf("shouldExclude(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		setup   func(*Parser)
		wantErr bool
	}{
		{
			name: "valid spec",
			setup: func(p *Parser) {
				p.openapi.Info.Title = "Test API"
				p.openapi.Info.Version = "1.0.0"
			},
			wantErr: false,
		},
		{
			name: "missing title",
			setup: func(p *Parser) {
				p.openapi.Info.Version = "1.0.0"
			},
			wantErr: true,
		},
		{
			name: "missing version",
			setup: func(p *Parser) {
				p.openapi.Info.Title = "Test API"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			tt.setup(p)

			err := p.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHasGeneralInfo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "has title",
			content: `package main
// @title Test API
func main() {}`,
			expected: true,
		},
		{
			name: "has version",
			content: `package main
// @version 1.0.0
func main() {}`,
			expected: true,
		},
		{
			name: "no general info",
			content: `package main
// @Summary test
func main() {}`,
			expected: false,
		},
		{
			name:     "empty file",
			content:  `package main`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.content, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse file: %v", err)
			}

			p := New()
			result := p.hasGeneralInfo(file)
			if result != tt.expected {
				t.Errorf("hasGeneralInfo() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShouldIncludeOperation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		includeTags []string
		excludeTags []string
		opTags      []string
		expected    bool
	}{
		{
			name:     "no filters",
			opTags:   []string{"users"},
			expected: true,
		},
		{
			name:        "included tag",
			includeTags: []string{"users", "products"},
			opTags:      []string{"users"},
			expected:    true,
		},
		{
			name:        "not in include list",
			includeTags: []string{"users"},
			opTags:      []string{"admin"},
			expected:    false,
		},
		{
			name:        "excluded tag",
			excludeTags: []string{"admin"},
			opTags:      []string{"admin"},
			expected:    false,
		},
		{
			name:        "both include and exclude",
			includeTags: []string{"users"},
			excludeTags: []string{"admin"},
			opTags:      []string{"users"},
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.SetTagFilters(tt.includeTags, tt.excludeTags)

			result := p.ShouldIncludeOperation(tt.opTags)
			if result != tt.expected {
				t.Errorf("ShouldIncludeOperation() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Mock FileInfo for testing.
type mockFileInfo struct {
	name  string
	isDir bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() os.FileMode  { return 0644 }
func (m *mockFileInfo) ModTime() time.Time { return time.Now() }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

func TestParseDirSkipsTestFiles(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create test file (should be skipped)
	testFile := filepath.Join(tmpDir, "main_test.go")
	if err := os.WriteFile(testFile, []byte("package main\nimport \"testing\""), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create normal file
	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte("package main\nfunc main() {}"), 0644); err != nil {
		t.Fatalf("Failed to write main file: %v", err)
	}

	p := New()
	err := p.ParseDir(tmpDir)
	if err != nil {
		t.Errorf("ParseDir() returned error: %v", err)
	}

	// Should only have parsed main.go, not main_test.go
	if len(p.files) != 1 {
		t.Errorf("files count = %d, want 1 (test files should be skipped)", len(p.files))
	}
}

func TestParseDirSkipsHiddenDirs(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create hidden directory
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	if err := os.Mkdir(hiddenDir, 0755); err != nil {
		t.Fatalf("Failed to create hidden dir: %v", err)
	}

	hiddenFile := filepath.Join(hiddenDir, "test.go")
	if err := os.WriteFile(hiddenFile, []byte("package hidden"), 0644); err != nil {
		t.Fatalf("Failed to write hidden file: %v", err)
	}

	// Create normal file
	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to write main file: %v", err)
	}

	p := New()
	err := p.ParseDir(tmpDir)
	if err != nil {
		t.Errorf("ParseDir() returned error: %v", err)
	}
}

func TestValidateSchemaRef(t *testing.T) {
	t.Parallel()
	p := New()
	p.openapi.Components.Schemas["User"] = &v3.Schema{Type: "object"}

	tests := []struct {
		name    string
		ref     string
		wantErr bool
	}{
		{
			name:    "valid reference",
			ref:     "#/components/schemas/User",
			wantErr: false,
		},
		{
			name:    "invalid format",
			ref:     "User",
			wantErr: true,
		},
		{
			name:    "schema not found",
			ref:     "#/components/schemas/NonExistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.validateSchemaRef(tt.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSchemaRef() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
