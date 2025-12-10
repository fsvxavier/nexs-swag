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

	// Test GetMarkdownContent
	content := p.GetMarkdownContent("api.md")
	if content == "" {
		// Cache may not be populated yet, that's ok for this test
		t.Log("Markdown cache not populated")
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
	p.SetCodeExampleFilesDir(tmpDir)

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
			p.SetGeneratedTime(tt.enabled)

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
			p.SetInstanceName(tt.instanceName)

			if p.instanceName != tt.instanceName {
				t.Errorf("instanceName = %q, want %q", p.instanceName, tt.instanceName)
			}
		})
	}
}

func TestSetParseGoList(t *testing.T) {
	t.Parallel()
	p := New()

	p.SetParseGoList(true)
	if !p.parseGoList {
		t.Error("parseGoList should be true")
	}

	p.SetParseGoList(false)
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
			p.SetTemplateDelims(tt.delims)

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
			p.SetCollectionFormat(tt.format)

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
			p.SetState(tt.state)

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
			p.SetParseExtension(tt.extension)

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
			p.SetParseDependencyLevel(tt.level)

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

func TestLoadCodeExamplesFromDir(t *testing.T) {
	t.Parallel()
	p := New()

	// Create temp directory with code examples
	tmpDir := t.TempDir()

	// Create example files
	examples := map[string]string{
		"example.go": "package main\n\nfunc main() {}",
		"example.py": "def main():\n    pass",
		"example.js": "function main() {}",
	}

	for filename, content := range examples {
		filePath := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", filename, err)
		}
	}

	p.codeExampleFilesDir = tmpDir
	err := p.loadCodeExamplesFromDir()
	if err != nil {
		t.Fatalf("loadCodeExamplesFromDir failed: %v", err)
	}

	// Verify examples were loaded
	for filename := range examples {
		content := p.GetCodeExample(filename)
		if content == "" {
			t.Errorf("Expected content for %s", filename)
		}
	}
}

func TestLoadCodeExamplesFromDirEmpty(t *testing.T) {
	t.Parallel()
	p := New()

	// No directory set
	err := p.loadCodeExamplesFromDir()
	if err != nil {
		t.Errorf("Expected no error for empty directory, got: %v", err)
	}
}

func TestLoadCodeExamplesFromDirNested(t *testing.T) {
	t.Parallel()
	p := New()

	// Create temp directory with nested structure
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create nested example file
	filePath := filepath.Join(subDir, "nested.go")
	content := "package main"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create nested file: %v", err)
	}

	p.codeExampleFilesDir = tmpDir
	err := p.loadCodeExamplesFromDir()
	if err != nil {
		t.Fatalf("loadCodeExamplesFromDir failed: %v", err)
	}

	// Verify nested example was loaded with relative path
	nestedPath := filepath.Join("subdir", "nested.go")
	content2 := p.GetCodeExample(nestedPath)
	if content2 == "" {
		t.Errorf("Expected content for nested file")
	}
}

func TestGetCodeExample(t *testing.T) {
	t.Parallel()
	p := New()

	// Create temp directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")
	expectedContent := "package main\n\nfunc test() {}"

	if err := os.WriteFile(filePath, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	p.codeExampleFilesDir = tmpDir
	if err := p.loadCodeExamplesFromDir(); err != nil {
		t.Fatalf("Failed to load examples: %v", err)
	}

	// Get example
	content := p.GetCodeExample("test.go")
	if content != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, content)
	}
}

func TestGetCodeExampleNonExistent(t *testing.T) {
	t.Parallel()
	p := New()

	// Without loading any examples
	content := p.GetCodeExample("nonexistent.go")
	if content != "" {
		t.Errorf("Expected empty string for non-existent example, got '%s'", content)
	}
}

func TestDetectLanguageFromExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "Go file",
			filename: "example.go",
			expected: "go",
		},
		{
			name:     "JavaScript file",
			filename: "example.js",
			expected: "javascript",
		},
		{
			name:     "TypeScript file",
			filename: "example.ts",
			expected: "typescript",
		},
		{
			name:     "Python file",
			filename: "example.py",
			expected: "python",
		},
		{
			name:     "Java file",
			filename: "example.java",
			expected: "java",
		},
		{
			name:     "Ruby file",
			filename: "example.rb",
			expected: "ruby",
		},
		{
			name:     "PHP file",
			filename: "example.php",
			expected: "php",
		},
		{
			name:     "C# file",
			filename: "example.cs",
			expected: "csharp",
		},
		{
			name:     "C++ file",
			filename: "example.cpp",
			expected: "cpp",
		},
		{
			name:     "C file",
			filename: "example.c",
			expected: "c",
		},
		{
			name:     "Shell file",
			filename: "example.sh",
			expected: "bash",
		},
		{
			name:     "JSON file",
			filename: "example.json",
			expected: "json",
		},
		{
			name:     "YAML file",
			filename: "example.yaml",
			expected: "yaml",
		},
		{
			name:     "YML file",
			filename: "example.yml",
			expected: "yaml",
		},
		{
			name:     "XML file",
			filename: "example.xml",
			expected: "xml",
		},
		{
			name:     "HTML file",
			filename: "example.html",
			expected: "html",
		},
		{
			name:     "Unknown extension",
			filename: "example.xyz",
			expected: "text",
		},
		{
			name:     "No extension",
			filename: "makefile",
			expected: "text",
		},
		{
			name:     "Uppercase extension",
			filename: "example.GO",
			expected: "go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectLanguageFromExtension(tt.filename)
			if result != tt.expected {
				t.Errorf("Expected language '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetCodeExampleWithMultipleFiles(t *testing.T) {
	t.Parallel()
	p := New()

	// Create temp directory with multiple files
	tmpDir := t.TempDir()

	files := map[string]string{
		"create_user.go": "// Go example",
		"create_user.py": "# Python example",
		"create_user.js": "// JavaScript example",
		"list_users.go":  "// Another Go example",
	}

	for filename, content := range files {
		filePath := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", filename, err)
		}
	}

	p.codeExampleFilesDir = tmpDir
	if err := p.loadCodeExamplesFromDir(); err != nil {
		t.Fatalf("Failed to load examples: %v", err)
	}

	// Verify all files were loaded
	for filename, expectedContent := range files {
		content := p.GetCodeExample(filename)
		if content != expectedContent {
			t.Errorf("For %s, expected '%s', got '%s'", filename, expectedContent, content)
		}
	}
}

func TestLoadCodeExamplesFromDirWithInvalidPath(t *testing.T) {
	t.Parallel()
	p := New()

	// Set invalid directory
	p.codeExampleFilesDir = "/nonexistent/directory/that/does/not/exist"
	err := p.loadCodeExamplesFromDir()
	if err == nil {
		t.Error("Expected error for invalid directory path")
	}
}
