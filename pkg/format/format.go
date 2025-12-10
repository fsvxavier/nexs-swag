package format

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsvxavier/nexs-swag/pkg/parser"
	"golang.org/x/tools/imports"
)

// Config holds format configuration
type Config struct {
	SearchDir string
	Excludes  string
	MainFile  string // Deprecated: no longer needed
}

// Format handles swagger comment formatting
type Format struct {
	excludes map[string]struct{}
}

// New creates a new Format instance
func New() *Format {
	return &Format{
		excludes: make(map[string]struct{}),
	}
}

// Build formats swagger comments in all Go files in the search directory
func (f *Format) Build(config *Config) error {
	if config.SearchDir == "" {
		config.SearchDir = "./"
	}

	// Setup default excludes
	defaultExcludes := []string{"docs", "vendor"}
	for _, exclude := range defaultExcludes {
		f.excludes[exclude] = struct{}{}
	}

	// Add configured excludes
	if config.Excludes != "" {
		for _, exclude := range strings.Split(config.Excludes, ",") {
			exclude = strings.TrimSpace(exclude)
			if exclude != "" {
				f.excludes[exclude] = struct{}{}
			}
		}
	}

	// Walk through directory
	return filepath.WalkDir(config.SearchDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		if d.IsDir() {
			if f.excludeDir(path) {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .go files
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Format the file
		return f.formatFile(path)
	})
}

// excludeDir checks if a directory should be excluded
func (f *Format) excludeDir(path string) bool {
	// Skip hidden directories
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") {
		return true
	}

	// Check if in exclude list
	_, excluded := f.excludes[base]
	return excluded
}

// formatFile formats swagger comments in a single Go file
func (f *Format) formatFile(path string) error {
	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	// Format swagger comments
	formatter := parser.NewFormatter()
	formatted, err := formatter.Format(path, content)
	if err != nil {
		return fmt.Errorf("failed to format swagger comments in %s: %w", path, err)
	}

	// Format imports and general Go code
	finalFormatted, err := imports.Process(path, formatted, nil)
	if err != nil {
		// If imports.Process fails, use the swagger-formatted version
		finalFormatted = formatted
	}

	// Only write if content changed
	if string(content) != string(finalFormatted) {
		if err := os.WriteFile(path, finalFormatted, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}
		fmt.Printf("  Formatted: %s\n", path)
	}

	return nil
}
