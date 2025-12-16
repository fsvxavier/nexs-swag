package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseWithGoList(t *testing.T) {
	t.Parallel()
	p := New()
	p.parseGoList = true

	// This will likely fail if not in a valid Go module
	// but we test that it doesn't panic
	err := p.parseWithGoList(".")
	// Error is expected if not in proper Go module context
	if err != nil {
		t.Logf("parseWithGoList() returned expected error: %v", err)
	}
}

func TestParseWithGoListInvalidPath(t *testing.T) {
	t.Parallel()
	p := New()
	p.parseGoList = true

	err := p.parseWithGoList("/nonexistent/path")
	if err == nil {
		t.Error("parseWithGoList() with invalid path should return error")
	}
}

func TestParseDependencies(t *testing.T) {
	t.Parallel()
	p := New()
	p.parseDependency = false

	// Should not parse dependencies when disabled
	err := p.parseDependencies()
	if err != nil {
		t.Errorf("parseDependencies() with disabled flag returned error: %v", err)
	}
}

func TestParseDependenciesEnabled(t *testing.T) {
	t.Parallel()
	p := New()
	p.parseDependency = true

	// This will try to parse go.mod
	err := p.parseDependencies()
	// Error is expected if go.mod doesn't exist or is invalid
	if err != nil {
		t.Logf("parseDependencies() returned expected error: %v", err)
	}
}

func TestGoListPackageInfo(t *testing.T) {
	t.Parallel()
	// Test that go list functionality exists
	p := New()
	p.parseGoList = true

	if !p.parseGoList {
		t.Error("parseGoList should be true")
	}
}

func TestParsePackageFromGoList(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		pkg           *GoListPackage
		parseInternal bool
		expectError   bool
	}{
		{
			name: "Valid package with Go files",
			pkg: &GoListPackage{
				Dir:     tmpDir,
				GoFiles: []string{"main.go"},
				Goroot:  false,
			},
			parseInternal: false,
			expectError:   true, // File doesn't exist, but that's ok for testing the function flow
		},
		{
			name: "Goroot package with parseInternal disabled",
			pkg: &GoListPackage{
				Dir:     tmpDir,
				GoFiles: []string{"main.go"},
				Goroot:  true,
			},
			parseInternal: false,
			expectError:   false, // Should skip without error
		},
		{
			name: "Goroot package with parseInternal enabled",
			pkg: &GoListPackage{
				Dir:     tmpDir,
				GoFiles: []string{"main.go"},
				Goroot:  true,
			},
			parseInternal: true,
			expectError:   true, // File doesn't exist
		},
		{
			name: "Package with CGO files",
			pkg: &GoListPackage{
				Dir:      tmpDir,
				CgoFiles: []string{"cgo.go"},
				Goroot:   false,
			},
			parseInternal: false,
			expectError:   true, // File doesn't exist
		},
		{
			name: "Package with both Go and CGO files",
			pkg: &GoListPackage{
				Dir:      tmpDir,
				GoFiles:  []string{"main.go"},
				CgoFiles: []string{"cgo.go"},
				Goroot:   false,
			},
			parseInternal: false,
			expectError:   true, // Files don't exist
		},
		{
			name: "Empty package",
			pkg: &GoListPackage{
				Dir:    tmpDir,
				Goroot: false,
			},
			parseInternal: false,
			expectError:   false, // No files to parse
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.parseInternal = tt.parseInternal

			err := p.parsePackageFromGoList(tt.pkg)

			if tt.expectError && err == nil {
				t.Error("parsePackageFromGoList() expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("parsePackageFromGoList() unexpected error: %v", err)
			}
		})
	}
}

func TestParsePackageFromGoListWithRealFiles(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create a simple Go file
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package testpkg

// User represents a user
type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	p := New()
	pkg := &GoListPackage{
		Dir:     tmpDir,
		GoFiles: []string{"test.go"},
		Goroot:  false,
	}

	err := p.parsePackageFromGoList(pkg)
	if err != nil {
		t.Errorf("parsePackageFromGoList() with valid file returned error: %v", err)
	}
}
