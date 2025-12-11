// Package parser implements the comment parser for OpenAPI documentation generation.
package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"

	openapi "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

// Parser parses Go source files and extracts OpenAPI documentation from comments.
type Parser struct {
	openapi         *openapi.OpenAPI
	files           map[string]*ast.File
	fset            *token.FileSet
	generalInfoFile string
	typeCache       map[string]*TypeInfo

	// Configuration options
	excludePatterns      []string
	propertyStrategy     string
	requiredByDefault    bool
	parseInternal        bool
	parseDependency      bool
	parseDepth           int
	markdownFilesDir     string
	overridesFile        string
	includeTags          []string
	excludeTags          []string
	parseFuncBody        bool
	parseVendor          bool
	typeOverrides        map[string]string
	parseDependencyLevel int
	codeExampleFilesDir  string
	generatedTime        bool
	instanceName         string
	parseGoList          bool
	templateDelims       string
	collectionFormat     string
	state                string
	parseExtension       string
	openapiVersion       string // Target OpenAPI version: "2.0", "3.0.0", "3.1.0"
}

// TypeInfo stores information about a parsed type.
type TypeInfo struct {
	Name    string
	Package string
	Schema  *openapi.Schema
	ASTNode ast.Node
}

// New creates a new Parser instance.
func New() *Parser {
	return &Parser{
		openapi: &openapi.OpenAPI{
			OpenAPI:           "3.1.0",
			JSONSchemaDialect: "https://spec.openapis.org/oas/3.1/dialect/base",
			Info:              openapi.Info{},
			Paths:             make(openapi.Paths),
			Components: &openapi.Components{
				Schemas:         make(map[string]*openapi.Schema),
				Responses:       make(map[string]*openapi.Response),
				Parameters:      make(map[string]*openapi.Parameter),
				Examples:        make(map[string]*openapi.Example),
				RequestBodies:   make(map[string]*openapi.RequestBody),
				Headers:         make(map[string]*openapi.Header),
				SecuritySchemes: make(map[string]*openapi.SecurityScheme),
				Links:           make(map[string]*openapi.Link),
				Callbacks:       make(map[string]*openapi.Callback),
				PathItems:       make(map[string]*openapi.PathItem),
			},
		},
		files:                make(map[string]*ast.File),
		fset:                 token.NewFileSet(),
		typeCache:            make(map[string]*TypeInfo),
		propertyStrategy:     "camelcase",
		parseDepth:           100,
		typeOverrides:        make(map[string]string),
		instanceName:         "swagger",
		collectionFormat:     "csv",
		parseDependencyLevel: 0,
		openapiVersion:       "3.1.0", // Default to latest
	}
}

// SetOpenAPIVersion sets the target OpenAPI version.
func (p *Parser) SetOpenAPIVersion(version string) {
	p.openapiVersion = version
}

// GetOpenAPIVersion returns the target OpenAPI version.
func (p *Parser) GetOpenAPIVersion() string {
	return p.openapiVersion
}

// ParseDir parses all Go files in the specified directory recursively.
func (p *Parser) ParseDir(dir string) error {
	// Parse dependencies from go.mod if enabled
	if err := p.parseDependencies(); err != nil {
		return fmt.Errorf("failed to parse dependencies: %w", err)
	}

	// Use go list if enabled
	if p.parseGoList {
		if err := p.parseWithGoList(dir); err != nil {
			return fmt.Errorf("failed to parse with go list: %w", err)
		}
		return nil
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if path matches exclude patterns
		if p.shouldExclude(path, info) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip vendor (unless parseVendor is true), testdata, and hidden directories
		if info.IsDir() {
			name := info.Name()
			if (!p.parseVendor && name == "vendor") ||
				name == "testdata" ||
				name == "docs" ||
				(strings.HasPrefix(name, ".") && name != ".") {
				return filepath.SkipDir
			}

			// Skip internal packages unless parseInternal is true
			if !p.parseInternal && name == "internal" {
				return filepath.SkipDir
			}

			return nil
		}

		// Parse only Go files (excluding test files)
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		return p.ParseFile(path)
	})
}

// shouldExclude checks if a path matches any exclude pattern.
func (p *Parser) shouldExclude(path string, info os.FileInfo) bool {
	if len(p.excludePatterns) == 0 {
		return false
	}

	name := info.Name()
	for _, pattern := range p.excludePatterns {
		// Simple pattern matching
		pattern = strings.TrimSpace(pattern)

		// Check if it's a wildcard pattern
		if strings.Contains(pattern, "*") {
			if matched, err := filepath.Match(pattern, name); err == nil && matched {
				return true
			}
		} else if strings.Contains(path, pattern) || name == pattern {
			return true
		}
	}

	return false
}

// ParseFile parses a single Go file and extracts documentation.
func (p *Parser) ParseFile(path string) error {
	file, err := parser.ParseFile(p.fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", path, err)
	}

	p.files[path] = file

	// Check if this file should be used for general API info
	shouldParseGeneralInfo := false

	if p.generalInfoFile != "" {
		// If generalInfo flag is set, only parse that specific file
		if absGeneralInfo, err := filepath.Abs(p.generalInfoFile); err == nil {
			if absPath, err := filepath.Abs(path); err == nil {
				shouldParseGeneralInfo = (absGeneralInfo == absPath)
			}
		}
	} else {
		// Auto-detect general info file
		if p.hasGeneralInfo(file) {
			p.generalInfoFile = path
			shouldParseGeneralInfo = true
		}
	}

	if shouldParseGeneralInfo {
		if err := p.parseGeneralInfo(file); err != nil {
			return fmt.Errorf("failed to parse general info from %s: %w", path, err)
		}
	}

	// Parse operations from function comments
	if err := p.parseOperations(file); err != nil {
		return fmt.Errorf("failed to parse operations from %s: %w", path, err)
	}

	// Parse schemas from type definitions
	if err := p.parseSchemas(file); err != nil {
		return fmt.Errorf("failed to parse schemas from %s: %w", path, err)
	}

	return nil
}

// GetOpenAPI returns the parsed OpenAPI specification.
func (p *Parser) GetOpenAPI() *openapi.OpenAPI {
	// Add generated timestamp if enabled
	if p.generatedTime && p.openapi.Info.Version != "" {
		p.openapi.Info.Version += " (generated at " + time.Now().Format("2006-01-02 15:04:05") + ")"
	}
	return p.openapi
}

// hasGeneralInfo checks if the file contains general API information.
func (p *Parser) hasGeneralInfo(file *ast.File) bool {
	for _, comment := range file.Comments {
		for _, line := range comment.List {
			text := strings.TrimSpace(strings.TrimPrefix(line.Text, "//"))
			if strings.HasPrefix(text, "@title") ||
				strings.HasPrefix(text, "@version") {
				return true
			}
		}
	}
	return false
}

// parseGeneralInfo extracts general API information from file comments.
func (p *Parser) parseGeneralInfo(file *ast.File) error {
	processor := NewGeneralInfoProcessor(p.openapi)

	for _, comment := range file.Comments {
		for _, line := range comment.List {
			text := strings.TrimSpace(strings.TrimPrefix(line.Text, "//"))
			if text == "" || !strings.HasPrefix(text, "@") {
				continue
			}

			if err := processor.Process(text); err != nil {
				return err
			}
		}
	}

	return nil
}

// parseOperations extracts operation information from function comments.
func (p *Parser) parseOperations(file *ast.File) error {
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Doc == nil {
			return true
		}

		processor := NewOperationProcessor(p, p.openapi, p.typeCache)
		op := processor.Process(funcDecl.Doc)
		if op == nil {
			return true
		}

		// Process function body comments if parseFuncBody is enabled
		if p.parseFuncBody && funcDecl.Body != nil {
			// Iterate through all comments in the file
			for _, commentGroup := range file.Comments {
				// Check if comment is within the function body range
				if commentGroup.Pos() > funcDecl.Body.Lbrace &&
					commentGroup.End() < funcDecl.Body.Rbrace {
					processor.Process(commentGroup)
				}
			}
		}

		// Check if operation should be included based on tag filters
		if !p.ShouldIncludeOperation(op.Tags) {
			return true
		}

		// Check if operation should be included based on extension filters
		if p.parseExtension != "" && !p.hasExtension(op) {
			return true
		}

		// Extract path and method from @Router annotation
		routeInfo := processor.GetRouteInfo(funcDecl.Doc)
		if routeInfo.Path == "" || routeInfo.Method == "" {
			return true
		}

		// Add operation to path
		pathItem := p.openapi.Paths[routeInfo.Path]
		if pathItem == nil {
			pathItem = &openapi.PathItem{}
			p.openapi.Paths[routeInfo.Path] = pathItem
		}

		// Set operation based on HTTP method
		switch strings.ToLower(routeInfo.Method) {
		case "get":
			pathItem.Get = op
		case "post":
			pathItem.Post = op
		case "put":
			pathItem.Put = op
		case "delete":
			pathItem.Delete = op
		case "patch":
			pathItem.Patch = op
		case "options":
			pathItem.Options = op
		case "head":
			pathItem.Head = op
		case "trace":
			pathItem.Trace = op
		}

		return true
	})

	return nil
}

// parseSchemas extracts schema definitions from type declarations.
func (p *Parser) parseSchemas(file *ast.File) error {
	processor := NewSchemaProcessor(p, p.openapi, p.typeCache)

	ast.Inspect(file, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		schema := processor.ProcessStruct(structType, typeSpec.Doc, typeSpec.Name.Name)
		if schema != nil {
			p.openapi.Components.Schemas[typeSpec.Name.Name] = schema

			// Cache the type info
			p.typeCache[typeSpec.Name.Name] = &TypeInfo{
				Name:    typeSpec.Name.Name,
				Package: file.Name.Name,
				Schema:  schema,
				ASTNode: typeSpec,
			}
		}

		return true
	})

	return nil
}

// Validate performs validation on the parsed OpenAPI specification.
func (p *Parser) Validate() error {
	if p.openapi.Info.Title == "" {
		return errors.New("API title is required (@title)")
	}

	if p.openapi.Info.Version == "" {
		return errors.New("API version is required (@version)")
	}

	// Validate that all schema references exist
	for path, pathItem := range p.openapi.Paths {
		operations := []*openapi.Operation{
			pathItem.Get, pathItem.Post, pathItem.Put, pathItem.Delete,
			pathItem.Patch, pathItem.Options, pathItem.Head, pathItem.Trace,
		}

		for _, op := range operations {
			if op == nil {
				continue
			}

			if err := p.validateOperation(op, path); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateOperation validates a single operation.
func (p *Parser) validateOperation(op *openapi.Operation, path string) error {
	// Validate parameters
	for _, param := range op.Parameters {
		if param.Schema != nil && param.Schema.Ref != "" {
			if err := p.validateSchemaRef(param.Schema.Ref); err != nil {
				return fmt.Errorf("invalid parameter schema reference in %s: %w", path, err)
			}
		}
	}

	// Validate request body
	if op.RequestBody != nil {
		for contentType, media := range op.RequestBody.Content {
			if media.Schema != nil && media.Schema.Ref != "" {
				if err := p.validateSchemaRef(media.Schema.Ref); err != nil {
					return fmt.Errorf("invalid request body schema reference (%s) in %s: %w", contentType, path, err)
				}
			}
		}
	}

	// Validate responses
	for statusCode, response := range op.Responses {
		if response.Content != nil {
			for contentType, media := range response.Content {
				if media.Schema != nil && media.Schema.Ref != "" {
					if err := p.validateSchemaRef(media.Schema.Ref); err != nil {
						return fmt.Errorf("invalid response schema reference (status %s, %s) in %s: %w",
							statusCode, contentType, path, err)
					}
				}
			}
		}
	}

	return nil
}

// validateSchemaRef validates that a schema reference exists.
func (p *Parser) validateSchemaRef(ref string) error {
	// Extract schema name from reference
	const prefix = "#/components/schemas/"
	if !strings.HasPrefix(ref, prefix) {
		return fmt.Errorf("invalid schema reference format: %s", ref)
	}

	schemaName := strings.TrimPrefix(ref, prefix)
	if _, exists := p.openapi.Components.Schemas[schemaName]; !exists {
		return fmt.Errorf("schema '%s' not found", schemaName)
	}

	return nil
}

// hasExtension checks if operation has the required extension.
func (p *Parser) hasExtension(op *openapi.Operation) bool {
	if op.Extensions == nil {
		return false
	}

	// Check if any extension key matches the filter
	for key := range op.Extensions {
		if strings.HasPrefix(key, p.parseExtension) {
			return true
		}
	}

	return false
}
