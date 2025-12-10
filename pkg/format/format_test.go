package format

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	f := New()
	if f == nil {
		t.Fatal("New() returned nil")
	}
	if f.excludes == nil {
		t.Error("excludes map not initialized")
	}
}

func TestExcludeDir(t *testing.T) {
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
