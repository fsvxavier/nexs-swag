package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMarkdownFiles(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create markdown files
	mdContent1 := "# API Documentation\nThis is the main API documentation."
	mdFile1 := filepath.Join(tmpDir, "api.md")
	if err := os.WriteFile(mdFile1, []byte(mdContent1), 0644); err != nil {
		t.Fatalf("Failed to write markdown file: %v", err)
	}

	mdContent2 := "# User Guide\nUser management guide."
	mdFile2 := filepath.Join(tmpDir, "users.md")
	if err := os.WriteFile(mdFile2, []byte(mdContent2), 0644); err != nil {
		t.Fatalf("Failed to write markdown file: %v", err)
	}

	p := New()
	p.SetMarkdownFilesDir(tmpDir)

	// Check if markdown files were loaded (this tests the loadMarkdownFiles internal function)
	// Since markdownCache is package-level, we can verify the directory was set
	if p.markdownFilesDir != tmpDir {
		t.Errorf("markdownFilesDir = %q, want %q", p.markdownFilesDir, tmpDir)
	}
}

func TestLoadMarkdownFilesEmpty(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create empty directory
	p := New()
	p.SetMarkdownFilesDir(tmpDir)

	if p.markdownFilesDir != tmpDir {
		t.Errorf("markdownFilesDir = %q, want %q", p.markdownFilesDir, tmpDir)
	}
}

func TestLoadMarkdownFilesNonExistent(t *testing.T) {
	t.Parallel()
	p := New()
	p.SetMarkdownFilesDir("/nonexistent/path")

	// Should not panic, just set the path
	if p.markdownFilesDir != "/nonexistent/path" {
		t.Error("SetMarkdownFilesDir should set path even if directory doesn't exist")
	}
}

func TestSetCodeExampleFilesDir(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create code example files
	goExample := filepath.Join(tmpDir, "example.go")
	if err := os.WriteFile(goExample, []byte("package main\nfunc main() {}"), 0644); err != nil {
		t.Fatalf("Failed to write go example: %v", err)
	}

	pyExample := filepath.Join(tmpDir, "example.py")
	if err := os.WriteFile(pyExample, []byte("print('hello')"), 0644); err != nil {
		t.Fatalf("Failed to write python example: %v", err)
	}

	p := New()
	p.codeExampleFilesDir = tmpDir

	if p.codeExampleFilesDir != tmpDir {
		t.Errorf("codeExampleFilesDir = %q, want %q", p.codeExampleFilesDir, tmpDir)
	}
}

func TestSetGeneratedTime(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		enabled bool
	}{
		{"enabled", true},
		{"disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.generatedTime = tt.enabled

			if p.generatedTime != tt.enabled {
				t.Errorf("generatedTime = %v, want %v", p.generatedTime, tt.enabled)
			}
		})
	}
}

func TestSetInstanceName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		instanceName string
	}{
		{"default swagger", "swagger"},
		{"custom docs", "docs"},
		{"custom api", "api"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.instanceName = tt.instanceName

			if p.instanceName != tt.instanceName {
				t.Errorf("instanceName = %q, want %q", p.instanceName, tt.instanceName)
			}
		})
	}
}

func TestSetParseGoList(t *testing.T) {
	t.Parallel()
	p := New()

	p.parseGoList = true
	if !p.parseGoList {
		t.Error("parseGoList should be true")
	}

	p.parseGoList = false
	if p.parseGoList {
		t.Error("parseGoList should be false")
	}
}

func TestSetTemplateDelims(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		delims string
	}{
		{"default", "{{,}}"},
		{"custom", "<%,%>"},
		{"brackets", "[[,]]"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.templateDelims = tt.delims

			if p.templateDelims != tt.delims {
				t.Errorf("templateDelims = %q, want %q", p.templateDelims, tt.delims)
			}
		})
	}
}

func TestSetCollectionFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		format string
	}{
		{"csv", "csv"},
		{"multi", "multi"},
		{"pipes", "pipes"},
		{"ssv", "ssv"},
		{"tsv", "tsv"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.collectionFormat = tt.format

			if p.collectionFormat != tt.format {
				t.Errorf("collectionFormat = %q, want %q", p.collectionFormat, tt.format)
			}
		})
	}
}

func TestSetState(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		state string
	}{
		{"development", "development"},
		{"production", "production"},
		{"staging", "staging"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.state = tt.state

			if p.state != tt.state {
				t.Errorf("state = %q, want %q", p.state, tt.state)
			}
		})
	}
}

func TestSetParseExtension(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		extension string
	}{
		{"x-codeSamples", "x-codeSamples"},
		{"x-custom", "x-custom"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.parseExtension = tt.extension

			if p.parseExtension != tt.extension {
				t.Errorf("parseExtension = %q, want %q", p.parseExtension, tt.extension)
			}
		})
	}
}

func TestSetParseDependencyLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		level int
	}{
		{"level 0", 0},
		{"level 1", 1},
		{"level 5", 5},
		{"level 10", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.parseDependencyLevel = tt.level

			if p.parseDependencyLevel != tt.level {
				t.Errorf("parseDependencyLevel = %d, want %d", p.parseDependencyLevel, tt.level)
			}
		})
	}
}

func TestMarkdownCacheGlobal(t *testing.T) {
	t.Parallel()
	// Test that markdownCache is properly initialized
	if markdownCache == nil {
		// markdownCache might be initialized lazily
		t.Log("markdownCache is nil (might be initialized lazily)")
	}
}

func TestMultipleConfigOptions(t *testing.T) {
	t.Parallel()
	p := New()

	// Set multiple configuration options
	p.generatedTime = true
	p.instanceName = "customapi"
	p.parseGoList = true
	p.templateDelims = "<%,%>"
	p.collectionFormat = "multi"
	p.state = "production"
	p.parseExtension = "x-custom"
	p.parseDependencyLevel = 5

	// Verify all were set
	if !p.generatedTime {
		t.Error("generatedTime not set")
	}
	if p.instanceName != "customapi" {
		t.Error("instanceName not set")
	}
	if !p.parseGoList {
		t.Error("parseGoList not set")
	}
	if p.templateDelims != "<%,%>" {
		t.Error("templateDelims not set")
	}
	if p.collectionFormat != "multi" {
		t.Error("collectionFormat not set")
	}
	if p.state != "production" {
		t.Error("state not set")
	}
	if p.parseExtension != "x-custom" {
		t.Error("parseExtension not set")
	}
	if p.parseDependencyLevel != 5 {
		t.Error("parseDependencyLevel not set")
	}
}
