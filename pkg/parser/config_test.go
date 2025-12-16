package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetGeneralInfoFile(t *testing.T) {
	t.Parallel()
	p := New()
	testPath := "/path/to/main.go"
	p.SetGeneralInfoFile(testPath)

	if p.generalInfoFile != testPath {
		t.Errorf("SetGeneralInfoFile() = %q, want %q", p.generalInfoFile, testPath)
	}
}

func TestSetExcludePatterns(t *testing.T) {
	t.Parallel()
	p := New()
	patterns := []string{"*.pb.go", "mock_*.go", "vendor"}
	p.SetExcludePatterns(patterns)

	if len(p.excludePatterns) != len(patterns) {
		t.Errorf("SetExcludePatterns() length = %d, want %d", len(p.excludePatterns), len(patterns))
	}

	for i, pattern := range patterns {
		if p.excludePatterns[i] != pattern {
			t.Errorf("excludePatterns[%d] = %q, want %q", i, p.excludePatterns[i], pattern)
		}
	}
}

func TestSetPropertyStrategy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		strategy string
	}{
		{"camelcase", "camelcase"},
		{"snakecase", "snakecase"},
		{"pascalcase", "pascalcase"},
		{"custom", "custom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.SetPropertyStrategy(tt.strategy)
			if p.propertyStrategy != tt.strategy {
				t.Errorf("SetPropertyStrategy() = %q, want %q", p.propertyStrategy, tt.strategy)
			}
		})
	}
}

func TestSetRequiredByDefault(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		required bool
	}{
		{"required true", true},
		{"required false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.SetRequiredByDefault(tt.required)
			if p.requiredByDefault != tt.required {
				t.Errorf("SetRequiredByDefault() = %v, want %v", p.requiredByDefault, tt.required)
			}
		})
	}
}

func TestSetParseInternal(t *testing.T) {
	t.Parallel()
	p := New()

	p.SetParseInternal(true)
	if !p.parseInternal {
		t.Error("SetParseInternal(true) failed")
	}

	p.SetParseInternal(false)
	if p.parseInternal {
		t.Error("SetParseInternal(false) failed")
	}
}

func TestSetParseDependency(t *testing.T) {
	t.Parallel()
	p := New()

	p.SetParseDependency(true)
	if !p.parseDependency {
		t.Error("SetParseDependency(true) failed")
	}

	p.SetParseDependency(false)
	if p.parseDependency {
		t.Error("SetParseDependency(false) failed")
	}
}

func TestSetParseDepth(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		depth int
	}{
		{"depth 10", 10},
		{"depth 50", 50},
		{"depth 100", 100},
		{"depth 0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.SetParseDepth(tt.depth)
			if p.parseDepth != tt.depth {
				t.Errorf("SetParseDepth() = %d, want %d", p.parseDepth, tt.depth)
			}
		})
	}
}

func TestSetMarkdownFilesDir(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create a markdown file
	mdFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(mdFile, []byte("# Test Markdown"), 0644); err != nil {
		t.Fatalf("Failed to create markdown file: %v", err)
	}

	p := New()
	p.SetMarkdownFilesDir(tmpDir)

	if p.markdownFilesDir != tmpDir {
		t.Errorf("SetMarkdownFilesDir() = %q, want %q", p.markdownFilesDir, tmpDir)
	}
}

func TestSetMarkdownFilesDirEmpty(t *testing.T) {
	t.Parallel()
	p := New()
	p.SetMarkdownFilesDir("")

	if p.markdownFilesDir != "" {
		t.Errorf("SetMarkdownFilesDir('') should set empty string, got %q", p.markdownFilesDir)
	}
}

func TestSetOverridesFile(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create a valid overrides file
	overridesFile := filepath.Join(tmpDir, "overrides.json")
	content := `{"replace": {"time.Time": "string"}}`
	if err := os.WriteFile(overridesFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create overrides file: %v", err)
	}

	p := New()
	p.SetOverridesFile(overridesFile)

	if p.overridesFile != overridesFile {
		t.Errorf("SetOverridesFile() = %q, want %q", p.overridesFile, overridesFile)
	}
}

func TestSetOverridesFileNonExistent(t *testing.T) {
	t.Parallel()
	p := New()
	p.SetOverridesFile("/nonexistent/overrides.json")

	// Should not panic, just set the path
	if p.overridesFile != "/nonexistent/overrides.json" {
		t.Errorf("SetOverridesFile() with nonexistent file failed")
	}
}

func TestSetTagFilters(t *testing.T) {
	t.Parallel()
	p := New()
	include := []string{"users", "products"}
	exclude := []string{"admin", "internal"}

	p.SetTagFilters(include, exclude)

	if len(p.includeTags) != len(include) {
		t.Errorf("includeTags length = %d, want %d", len(p.includeTags), len(include))
	}
	if len(p.excludeTags) != len(exclude) {
		t.Errorf("excludeTags length = %d, want %d", len(p.excludeTags), len(exclude))
	}
}

func TestSetTagFiltersEmpty(t *testing.T) {
	t.Parallel()
	p := New()
	p.SetTagFilters(nil, nil)

	if p.includeTags != nil {
		t.Error("includeTags should be nil")
	}
	if p.excludeTags != nil {
		t.Error("excludeTags should be nil")
	}
}

func TestSetParseFuncBody(t *testing.T) {
	t.Parallel()
	p := New()

	p.SetParseFuncBody(true)
	if !p.parseFuncBody {
		t.Error("SetParseFuncBody(true) failed")
	}

	p.SetParseFuncBody(false)
	if p.parseFuncBody {
		t.Error("SetParseFuncBody(false) failed")
	}
}

func TestSetParseVendor(t *testing.T) {
	t.Parallel()
	p := New()

	p.SetParseVendor(true)
	if !p.parseVendor {
		t.Error("SetParseVendor(true) failed")
	}

	p.SetParseVendor(false)
	if p.parseVendor {
		t.Error("SetParseVendor(false) failed")
	}
}

func TestLoadTypeOverridesInvalidJSON(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create invalid JSON file
	invalidFile := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(invalidFile, []byte("not valid json"), 0644); err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	p := New()
	p.SetOverridesFile(invalidFile)

	// Should not panic, just skip loading
	if len(p.typeOverrides) != 0 {
		t.Error("typeOverrides should be empty after loading invalid JSON")
	}
}

func TestLoadTypeOverridesValidJSON(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create valid overrides file
	overridesFile := filepath.Join(tmpDir, "overrides.json")
	content := `{
		"replace": {
			"time.Time": "string",
			"uuid.UUID": "string"
		}
	}`
	if err := os.WriteFile(overridesFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create overrides file: %v", err)
	}

	p := New()
	p.SetOverridesFile(overridesFile)

	if len(p.typeOverrides) != 2 {
		t.Errorf("typeOverrides length = %d, want 2", len(p.typeOverrides))
	}
}

func TestMultipleConfigurationMethods(t *testing.T) {
	t.Parallel()
	p := New()

	// Set multiple configurations
	p.SetPropertyStrategy("snakecase")
	p.SetRequiredByDefault(true)
	p.SetParseInternal(true)
	p.SetParseDependency(false)
	p.SetParseDepth(50)
	p.SetParseFuncBody(true)
	p.SetParseVendor(false)

	// Verify all were set correctly
	if p.propertyStrategy != "snakecase" {
		t.Error("propertyStrategy not set correctly")
	}
	if !p.requiredByDefault {
		t.Error("requiredByDefault not set correctly")
	}
	if !p.parseInternal {
		t.Error("parseInternal not set correctly")
	}
	if p.parseDependency {
		t.Error("parseDependency not set correctly")
	}
	if p.parseDepth != 50 {
		t.Error("parseDepth not set correctly")
	}
	if !p.parseFuncBody {
		t.Error("parseFuncBody not set correctly")
	}
	if p.parseVendor {
		t.Error("parseVendor not set correctly")
	}
}

func TestSetExcludePatternsMultipleTimes(t *testing.T) {
	t.Parallel()
	p := New()

	// First set
	p.SetExcludePatterns([]string{"*.pb.go"})
	if len(p.excludePatterns) != 1 {
		t.Errorf("First set: excludePatterns length = %d, want 1", len(p.excludePatterns))
	}

	// Second set (should replace)
	p.SetExcludePatterns([]string{"mock_*.go", "test_*.go"})
	if len(p.excludePatterns) != 2 {
		t.Errorf("Second set: excludePatterns length = %d, want 2", len(p.excludePatterns))
	}
}

func TestSetGeneralInfoFileRelativePath(t *testing.T) {
	t.Parallel()
	p := New()
	relativePath := "./cmd/api/main.go"
	p.SetGeneralInfoFile(relativePath)

	if p.generalInfoFile != relativePath {
		t.Errorf("SetGeneralInfoFile() with relative path = %q, want %q", p.generalInfoFile, relativePath)
	}
}

func TestSetGeneralInfoFileAbsolutePath(t *testing.T) {
	t.Parallel()
	p := New()
	absolutePath := "/home/user/project/main.go"
	p.SetGeneralInfoFile(absolutePath)

	if p.generalInfoFile != absolutePath {
		t.Errorf("SetGeneralInfoFile() with absolute path = %q, want %q", p.generalInfoFile, absolutePath)
	}
}

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

func TestFindModuleInCache(t *testing.T) {
	t.Parallel()
	p := New()
	tmpDir := t.TempDir()

	tests := []struct {
		name           string
		modulePath     string
		setupFunc      func(string) string
		expectedExists bool
	}{
		{
			name:       "Module with uppercase letters",
			modulePath: "github.com/User/Repo",
			setupFunc: func(cacheDir string) string {
				// Create escaped path: github.com/!user/!repo/v1.0.0
				modDir := filepath.Join(cacheDir, "github.com", "!user", "!repo", "v1.0.0")
				os.MkdirAll(modDir, 0755)
				// Create go.mod to make it valid
				os.WriteFile(filepath.Join(modDir, "go.mod"), []byte("module github.com/User/Repo\n"), 0644)
				return modDir
			},
			expectedExists: true,
		},
		{
			name:       "Module with all lowercase",
			modulePath: "github.com/user/repo",
			setupFunc: func(cacheDir string) string {
				modDir := filepath.Join(cacheDir, "github.com", "user", "repo", "v1.2.3")
				os.MkdirAll(modDir, 0755)
				os.WriteFile(filepath.Join(modDir, "go.mod"), []byte("module github.com/user/repo\n"), 0644)
				return modDir
			},
			expectedExists: true,
		},
		{
			name:       "Module not in cache",
			modulePath: "github.com/nonexistent/module",
			setupFunc: func(cacheDir string) string {
				return ""
			},
			expectedExists: false,
		},
		{
			name:       "Module without version subdirectory",
			modulePath: "example.com/simple",
			setupFunc: func(cacheDir string) string {
				modDir := filepath.Join(cacheDir, "example.com", "simple")
				os.MkdirAll(modDir, 0755)
				os.WriteFile(filepath.Join(modDir, "go.mod"), []byte("module example.com/simple\n"), 0644)
				return modDir
			},
			expectedExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheDir := filepath.Join(tmpDir, tt.name)
			os.MkdirAll(cacheDir, 0755)

			expectedPath := tt.setupFunc(cacheDir)
			result := p.findModuleInCache(cacheDir, tt.modulePath)

			if tt.expectedExists {
				if result == "" {
					t.Errorf("findModuleInCache() returned empty string, expected to find module at %s", expectedPath)
				}
			} else {
				if result != "" {
					t.Errorf("findModuleInCache() = %q, expected empty string", result)
				}
			}
		})
	}
}

func TestFindModuleInCacheEscaping(t *testing.T) {
	t.Parallel()
	p := New()
	tmpDir := t.TempDir()

	// Test uppercase letter escaping
	modulePath := "github.com/MyOrg/MyRepo"

	// Create directory with escaped uppercase
	escapedPath := filepath.Join(tmpDir, "github.com", "!my!org", "!my!repo", "v1.0.0")
	if err := os.MkdirAll(escapedPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create go.mod to make it valid
	if err := os.WriteFile(filepath.Join(escapedPath, "go.mod"), []byte("module github.com/MyOrg/MyRepo\n"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	result := p.findModuleInCache(tmpDir, modulePath)
	if result == "" {
		t.Error("findModuleInCache() should find module with escaped uppercase letters")
	}
	if result != escapedPath {
		t.Errorf("findModuleInCache() = %q, want %q", result, escapedPath)
	}
}
