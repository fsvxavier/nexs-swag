// Package parser - Configuration methods for Parser
package parser

import (
	"encoding/json"
	"go/parser"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// markdownCache stores loaded markdown files.
var markdownCache map[string]string
var markdownCacheMutex sync.RWMutex

// SetGeneralInfoFile sets the file path for general API info.
func (p *Parser) SetGeneralInfoFile(path string) {
	p.generalInfoFile = path
}

// SetExcludePatterns sets patterns to exclude from parsing.
func (p *Parser) SetExcludePatterns(patterns []string) {
	p.excludePatterns = patterns
}

// SetPropertyStrategy sets the property naming strategy.
// Valid values: "camelcase", "snakecase", "pascalcase".
func (p *Parser) SetPropertyStrategy(strategy string) {
	p.propertyStrategy = strategy
}

// SetRequiredByDefault sets whether all fields should be required by default.
func (p *Parser) SetRequiredByDefault(required bool) {
	p.requiredByDefault = required
}

// SetParseInternal sets whether to parse internal packages.
func (p *Parser) SetParseInternal(parse bool) {
	p.parseInternal = parse
}

// SetParseDependency sets whether to parse dependencies.
func (p *Parser) SetParseDependency(parse bool) {
	p.parseDependency = parse
}

// SetParseDepth sets the maximum parse depth.
func (p *Parser) SetParseDepth(depth int) {
	p.parseDepth = depth
}

// SetMarkdownFilesDir sets the directory containing markdown files.
func (p *Parser) SetMarkdownFilesDir(dir string) {
	p.markdownFilesDir = dir
	if dir != "" {
		p.loadMarkdownFiles()
	}
}

// SetOverridesFile sets the file path for type overrides.
func (p *Parser) SetOverridesFile(path string) {
	p.overridesFile = path
	p.loadTypeOverrides()
}

// SetTagFilters sets tag filters for API operations.
func (p *Parser) SetTagFilters(include, exclude []string) {
	p.includeTags = include
	p.excludeTags = exclude
}

// SetParseFuncBody sets whether to parse function bodies.
func (p *Parser) SetParseFuncBody(parse bool) {
	p.parseFuncBody = parse
}

// SetParseVendor sets whether to parse vendor directory.
func (p *Parser) SetParseVendor(parse bool) {
	p.parseVendor = parse
}

// loadTypeOverrides loads type overrides from file.
func (p *Parser) loadTypeOverrides() {
	if p.overridesFile == "" {
		return
	}

	data, err := os.ReadFile(p.overridesFile)
	if err != nil {
		// File doesn't exist or can't be read, skip silently
		return
	}

	var config struct {
		Replace map[string]string `json:"replace"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return
	}

	p.typeOverrides = config.Replace
}

// loadMarkdownFiles loads all markdown files from the specified directory.
func (p *Parser) loadMarkdownFiles() {
	if markdownCache == nil {
		markdownCache = make(map[string]string)
	}

	if p.markdownFilesDir == "" {
		return
	}

	_ = filepath.Walk(p.markdownFilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Only process .md files
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Use filename without extension as key
		filename := filepath.Base(path)
		key := strings.TrimSuffix(filename, ".md")

		markdownCacheMutex.Lock()
		markdownCache[key] = string(content)
		markdownCacheMutex.Unlock()

		return nil
	})
}

// GetMarkdownContent retrieves markdown content by filename.
func (p *Parser) GetMarkdownContent(filename string) string {
	markdownCacheMutex.RLock()
	defer markdownCacheMutex.RUnlock()

	if markdownCache == nil {
		return ""
	}
	return markdownCache[filename]
}

// GetParseDepth returns the current parse depth limit.
func (p *Parser) GetParseDepth() int {
	return p.parseDepth
}

// GetTypeOverride returns the override type for a given type name.
// Supports both exact matches and partial matches (e.g., "NullInt64" matches "database/sql.NullInt64").
func (p *Parser) GetTypeOverride(typeName string) (string, bool) {
	if p.typeOverrides == nil {
		return "", false
	}

	// Try exact match first
	if override, exists := p.typeOverrides[typeName]; exists {
		return override, true
	}

	// Try partial match (match suffix after last dot)
	for fullName, override := range p.typeOverrides {
		if strings.HasSuffix(fullName, "."+typeName) {
			return override, true
		}
	}

	return "", false
}

// ShouldIncludeOperation checks if an operation should be included based on tag filters.
func (p *Parser) ShouldIncludeOperation(tags []string) bool {
	// If no filters, include everything
	if len(p.includeTags) == 0 && len(p.excludeTags) == 0 {
		return true
	}

	// Check exclude tags first
	for _, tag := range tags {
		for _, excludeTag := range p.excludeTags {
			if tag == excludeTag {
				return false
			}
		}
	}

	// If include tags specified, operation must have at least one
	if len(p.includeTags) > 0 {
		for _, tag := range tags {
			for _, includeTag := range p.includeTags {
				if tag == includeTag {
					return true
				}
			}
		}
		return false
	}

	return true
}

// SetParseDependencyLevel sets the dependency parse level.
func (p *Parser) SetParseDependencyLevel(level int) {
	p.parseDependencyLevel = level
}

// SetCodeExampleFilesDir sets the directory containing code example files.
func (p *Parser) SetCodeExampleFilesDir(dir string) {
	p.codeExampleFilesDir = dir
	if dir != "" {
		// Ignore error as code examples are optional
		_ = p.loadCodeExamplesFromDir()
	}
}

// SetGeneratedTime sets whether to include generation timestamp.
func (p *Parser) SetGeneratedTime(enabled bool) {
	p.generatedTime = enabled
}

// SetInstanceName sets the swagger instance name.
func (p *Parser) SetInstanceName(name string) {
	p.instanceName = name
}

// SetParseGoList sets whether to use 'go list' for parsing.
func (p *Parser) SetParseGoList(enabled bool) {
	p.parseGoList = enabled
}

// SetTemplateDelims sets custom template delimiters.
func (p *Parser) SetTemplateDelims(delims string) {
	p.templateDelims = delims
}

// SetCollectionFormat sets the default collection format.
func (p *Parser) SetCollectionFormat(format string) {
	p.collectionFormat = TransToValidCollectionFormat(format)
}

// SetState sets the state file path.
func (p *Parser) SetState(path string) {
	p.state = path
}

// SetParseExtension sets the extension filter for operations.
func (p *Parser) SetParseExtension(ext string) {
	p.parseExtension = ext
}

// parseDependencies reads go.mod and parses external dependencies if enabled.
func (p *Parser) parseDependencies() error {
	if !p.parseDependency {
		return nil
	}

	// Read go.mod file
	modData, err := os.ReadFile("go.mod")
	if err != nil {
		// go.mod not found, skip silently
		return nil
	}

	// Parse go.mod for require statements
	lines := strings.Split(string(modData), "\n")
	inRequire := false
	var dependencies []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "require") {
			inRequire = true
			continue
		}

		if inRequire {
			if strings.HasPrefix(line, ")") {
				break
			}

			// Extract module path
			if line != "" && !strings.HasPrefix(line, "//") {
				parts := strings.Fields(line)
				if len(parts) >= 1 {
					modulePath := parts[0]
					dependencies = append(dependencies, modulePath)
				}
			}
		}
	}

	// Parse dependencies based on level
	// Level 0: none (already handled by checking parseDependency)
	// Level 1: models only - parse types/structs from dependencies
	// Level 2: operations only - parse API operations from dependencies
	// Level 3: all - parse everything from dependencies
	if p.parseDependencyLevel > 0 {
		for _, dep := range dependencies {
			if err := p.parseDependencyPackage(dep); err != nil {
				// Log error but continue with other dependencies
				continue
			}
		}
	}

	return nil
}

// parseDependencyPackage parses a single dependency package based on the level.
func (p *Parser) parseDependencyPackage(modulePath string) error {
	// Find dependency in vendor or GOPATH/GOMODCACHE
	var depDir string

	// Check vendor first
	vendorDir := filepath.Join("vendor", modulePath)
	if stat, err := os.Stat(vendorDir); err == nil && stat.IsDir() {
		depDir = vendorDir
	} else {
		// Try GOMODCACHE
		_ = os.Getenv("GOMODCACHE")
		// Note: Module cache lookup would be implemented here in a full solution

		// Find the module in cache (this is simplified - real impl would need version matching)
		// For now, skip if not in vendor
		return nil
	}

	// Parse based on level
	switch p.parseDependencyLevel {
	case 1: // Models only
		return p.parseDependencyModels(depDir)
	case 2: // Operations only
		return p.parseDependencyOperations(depDir)
	case 3: // All
		return p.ParseDir(depDir)
	}

	return nil
}

// parseDependencyModels parses only model definitions from a dependency.
func (p *Parser) parseDependencyModels(dir string) error {
	// Parse Go files and extract only type definitions
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}

		file, err := parser.ParseFile(p.fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil // Skip files with errors
		}

		// Extract only schema/type declarations
		return p.parseSchemas(file)
	})
}

// parseDependencyOperations parses only operations from a dependency.
func (p *Parser) parseDependencyOperations(dir string) error {
	// Parse Go files and extract only operation annotations
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}

		file, err := parser.ParseFile(p.fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil // Skip files with errors
		}

		// Extract only operations
		return p.parseOperations(file)
	})
}
