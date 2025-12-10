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
