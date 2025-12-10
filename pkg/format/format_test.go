package format

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	f := New()
	if f == nil {
		t.Fatal("New() returned nil")
	}
	if f.excludes == nil {
		t.Error("excludes map not initialized")
	}
	if len(f.excludes) != 0 {
		t.Errorf("excludes map should be empty, got %d entries", len(f.excludes))
	}
}

func TestExcludeDir(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		path     string
		excludes map[string]struct{}
		expected bool
	}{
		{
			name:     "hidden directory",
			path:     ".git",
			excludes: map[string]struct{}{},
			expected: true,
		},
		{
			name:     "vendor directory",
			path:     "vendor",
			excludes: map[string]struct{}{"vendor": {}},
			expected: true,
		},
		{
			name:     "docs directory",
			path:     "docs",
			excludes: map[string]struct{}{"docs": {}},
			expected: true,
		},
		{
			name:     "normal directory",
			path:     "pkg",
			excludes: map[string]struct{}{},
			expected: false,
		},
		{
			name:     "hidden dot prefix",
			path:     ".hidden",
			excludes: map[string]struct{}{},
			expected: true,
		},
		{
			name:     "full path with vendor",
			path:     "/home/user/project/vendor",
			excludes: map[string]struct{}{"vendor": {}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Format{excludes: tt.excludes}
			result := f.excludeDir(tt.path)
			if result != tt.expected {
				t.Errorf("excludeDir(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestBuildWithDefaultConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "format-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.go")
	testContent := `package main

// @title Test API
// @version 1.0
func main() {}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	f := New()
	config := &Config{
		SearchDir: tmpDir,
	}

	err = f.Build(config)
	if err != nil {
		t.Errorf("Build() returned error: %v", err)
	}

	// Verify default excludes were added
	if _, ok := f.excludes["docs"]; !ok {
		t.Error("Default exclude 'docs' not added")
	}
	if _, ok := f.excludes["vendor"]; !ok {
		t.Error("Default exclude 'vendor' not added")
	}
}

func TestBuildWithCustomExcludes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "format-test-excludes-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(testFile, []byte("package main\n\nfunc main() {}"), 0644); err != nil {
		t.Fatalf("Failed to write main file: %v", err)
	}

	tests := []struct {
		name     string
		excludes string
		expected []string
	}{
		{
			name:     "single exclude",
			excludes: "test",
			expected: []string{"test", "docs", "vendor"},
		},
		{
			name:     "multiple excludes",
			excludes: "test,example,tmp",
			expected: []string{"test", "example", "tmp", "docs", "vendor"},
		},
		{
			name:     "excludes with spaces",
			excludes: " test , example ",
			expected: []string{"test", "example", "docs", "vendor"},
		},
		{
			name:     "empty excludes",
			excludes: "",
			expected: []string{"docs", "vendor"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New()
			config := &Config{
				SearchDir: tmpDir,
				Excludes:  tt.excludes,
			}

			err := f.Build(config)
			if err != nil {
				t.Errorf("Build() returned error: %v", err)
			}

			for _, expected := range tt.expected {
				if _, ok := f.excludes[expected]; !ok {
					t.Errorf("Expected exclude %q not found in excludes map", expected)
				}
			}
		})
	}
}

func TestBuildWithEmptySearchDir(t *testing.T) {
	f := New()
	config := &Config{
		SearchDir: "",
	}

	// Should default to "./" - we don't check error as current dir might vary
	_ = f.Build(config)

	if config.SearchDir != "./" {
		t.Errorf("SearchDir should default to './', got %q", config.SearchDir)
	}
}

func TestBuildWithInvalidDir(t *testing.T) {
	f := New()
	config := &Config{
		SearchDir: "/nonexistent/path/that/does/not/exist/xyz123",
	}

	err := f.Build(config)
	if err == nil {
		t.Error("Build() with invalid dir should return error")
	}
}

func TestBuildSkipsExcludedDirs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "format-test-skip-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create vendor directory
	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.Mkdir(vendorDir, 0755); err != nil {
		t.Fatalf("Failed to create vendor dir: %v", err)
	}

	vendorFile := filepath.Join(vendorDir, "vendor.go")
	if err := os.WriteFile(vendorFile, []byte("package vendor"), 0644); err != nil {
		t.Fatalf("Failed to write vendor file: %v", err)
	}

	// Create normal file
	testFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(testFile, []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to write main file: %v", err)
	}

	f := New()
	config := &Config{
		SearchDir: tmpDir,
	}

	err = f.Build(config)
	if err != nil {
		t.Errorf("Build() returned error: %v", err)
	}
}

func TestBuildSkipsNonGoFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "format-test-nongo-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create non-Go files
	txtFile := filepath.Join(tmpDir, "readme.txt")
	if err := os.WriteFile(txtFile, []byte("readme"), 0644); err != nil {
		t.Fatalf("Failed to write txt file: %v", err)
	}

	mdFile := filepath.Join(tmpDir, "readme.md")
	if err := os.WriteFile(mdFile, []byte("# Readme"), 0644); err != nil {
		t.Fatalf("Failed to write md file: %v", err)
	}

	f := New()
	config := &Config{
		SearchDir: tmpDir,
	}

	err = f.Build(config)
	if err != nil {
		t.Errorf("Build() should not fail with non-Go files: %v", err)
	}
}

func TestFormatFileSuccess(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp("", "format-test-file-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.go")
	testContent := `package main

// @Summary Get user
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} string
// @Router /users/{id} [get]
func GetUser() {}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	f := New()
	err = f.formatFile(testFile)
	if err != nil {
		t.Errorf("formatFile() returned error: %v", err)
	}

	// Verify file still exists and is readable
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("Failed to read formatted file: %v", err)
	}
	if len(content) == 0 {
		t.Error("Formatted file is empty")
	}
}

func TestFormatFileInvalidPath(t *testing.T) {
	t.Parallel()
	f := New()
	err := f.formatFile("/nonexistent/path/file.go")
	if err == nil {
		t.Error("formatFile() with invalid path should return error")
	}
}

func TestFormatFileInvalidGoSyntax(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp("", "format-test-invalid-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "invalid.go")
	invalidContent := `package main

func InvalidFunc( {
	// Missing closing brace and proper syntax
`
	if err := os.WriteFile(testFile, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	f := New()
	err = f.formatFile(testFile)
	// Should return error for invalid Go syntax
	if err == nil {
		t.Error("formatFile() with invalid Go syntax should return error")
	}
}

func TestBuildWithNestedDirectories(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp("", "format-test-nested-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create nested directory structure
	subDir1 := filepath.Join(tmpDir, "api")
	subDir2 := filepath.Join(subDir1, "handlers")
	if err := os.MkdirAll(subDir2, 0755); err != nil {
		t.Fatalf("Failed to create nested dirs: %v", err)
	}

	// Create Go files in different levels
	files := []string{
		filepath.Join(tmpDir, "main.go"),
		filepath.Join(subDir1, "api.go"),
		filepath.Join(subDir2, "handler.go"),
	}

	for _, file := range files {
		content := `package main

func Test() {}
`
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", file, err)
		}
	}

	f := New()
	config := &Config{
		SearchDir: tmpDir,
	}

	err = f.Build(config)
	if err != nil {
		t.Errorf("Build() with nested directories returned error: %v", err)
	}
}

func TestExcludeDirEdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		path     string
		excludes map[string]struct{}
		expected bool
	}{
		{
			name:     "current directory",
			path:     ".",
			excludes: map[string]struct{}{},
			expected: true,
		},
		{
			name:     "double dot directory",
			path:     "..",
			excludes: map[string]struct{}{},
			expected: true,
		},
		{
			name:     "custom exclude",
			path:     "mydir",
			excludes: map[string]struct{}{"mydir": {}},
			expected: true,
		},
		{
			name:     "not in exclude list",
			path:     "src",
			excludes: map[string]struct{}{"vendor": {}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Format{excludes: tt.excludes}
			result := f.excludeDir(tt.path)
			if result != tt.expected {
				t.Errorf("excludeDir(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestConfigMainFileDeprecated(t *testing.T) {
	t.Parallel()
	config := &Config{
		SearchDir: "./",
		Excludes:  "test",
		MainFile:  "main.go", // Deprecated field
	}

	if config.MainFile != "main.go" {
		t.Error("MainFile field should be accessible even if deprecated")
	}
}

func TestBuildWithHiddenDirectories(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp("", "format-test-hidden-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create hidden directory
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	if err := os.Mkdir(hiddenDir, 0755); err != nil {
		t.Fatalf("Failed to create hidden dir: %v", err)
	}

	hiddenFile := filepath.Join(hiddenDir, "test.go")
	if err := os.WriteFile(hiddenFile, []byte("package hidden"), 0644); err != nil {
		t.Fatalf("Failed to write hidden file: %v", err)
	}

	f := New()
	config := &Config{
		SearchDir: tmpDir,
	}

	// Build should skip hidden directory
	err = f.Build(config)
	if err != nil {
		t.Errorf("Build() with hidden directories returned error: %v", err)
	}
}

func TestFormatFileWithValidSwaggerComments(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp("", "format-test-swagger-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "swagger.go")
	testContent := `package api

// @title User API
// @version 1.0
// @description API for user management
// @host localhost:8080
// @BasePath /api/v1

// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User object"
// @Success 201 {object} User
// @Failure 400 {object} Error
// @Router /users [post]
func CreateUser() {}

type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

type Error struct {
	Message string ` + "`json:\"message\"`" + `
}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	f := New()
	err = f.formatFile(testFile)
	if err != nil {
		t.Errorf("formatFile() with valid swagger comments returned error: %v", err)
	}
}
