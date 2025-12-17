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
	parsedModules   map[string]bool              // Track parsed modules to avoid infinite recursion
	referencedTypes map[string]bool              // Track types referenced in operations (selective parsing)
	importMap       map[string]map[string]string // Map of file path -> (package alias -> import path)
	parsingExternal bool                         // Flag to indicate we're parsing external packages

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
	includeTypes         []string // Filter by Go type categories: struct, interface, func, const, type, all
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
		parsedModules:        make(map[string]bool),
		referencedTypes:      make(map[string]bool),
		importMap:            make(map[string]map[string]string),
		includeTypes:         []string{"all"}, // Default: include all referenced types
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
	} // Use go list if enabled
	if p.parseGoList {
		if err := p.parseWithGoList(dir); err != nil {
			return fmt.Errorf("failed to parse with go list: %w", err)
		}
		return nil
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

	if err != nil {
		return err
	}

	// After parsing all files and operations, resolve type dependencies
	p.ResolveTypeDependencies()

	// Now parse schemas for only the referenced types
	for _, file := range p.files {
		if err := p.parseSchemas(file); err != nil {
			return fmt.Errorf("failed to parse schemas: %w", err)
		}
	}

	return nil
} // shouldExclude checks if a path matches any exclude pattern.
func (p *Parser) shouldExclude(path string, info os.FileInfo) bool {
	if len(p.excludePatterns) == 0 {
		return false
	}

	name := info.Name()

	// Clean path for comparison (remove leading ./)
	cleanPath := strings.TrimPrefix(path, "./")

	for _, pattern := range p.excludePatterns {
		// Simple pattern matching
		pattern = strings.TrimSpace(pattern)

		// Remove leading ./ from pattern
		pattern = strings.TrimPrefix(pattern, "./")

		// Check if it's a wildcard pattern
		if strings.Contains(pattern, "*") {
			// Match against both name and full path
			if matched, err := filepath.Match(pattern, name); err == nil && matched {
				return true
			}
			if matched, err := filepath.Match(pattern, cleanPath); err == nil && matched {
				return true
			}
		} else {
			// Exact match on directory/file name
			if name == pattern {
				return true
			}
			// Match if pattern is found in path
			if strings.Contains(cleanPath, pattern) {
				return true
			}
			// Match if path starts with pattern (for directory exclusion)
			if strings.HasPrefix(cleanPath, pattern+"/") || cleanPath == pattern {
				return true
			}
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

	// Collect imports from this file (only for main project files, not external dependencies)
	if !p.parsingExternal {
		p.collectImports(path, file)
	}

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

	// Parse operations from function comments (this populates referencedTypes)
	if err := p.parseOperations(file); err != nil {
		return fmt.Errorf("failed to parse operations from %s: %w", path, err)
	}

	// Note: parseSchemas will be called after all files are parsed
	// to ensure we have all referenced types from all operations

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
// Respects includeTypes filter for func category.
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
		case "query":
			// QUERY method is new in OpenAPI 3.2.0
			pathItem.Query = op
		}

		return true
	})

	return nil
}

// parseSchemas extracts schema definitions from type declarations.
// Only processes structs that are referenced in operations or their dependencies.
// Respects includeTypes filter for type categories.
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

		// Check if struct category should be included
		if !p.ShouldIncludeTypeCategory("struct") {
			return true
		}

		// Use simple name for types in current package
		schemaName := typeSpec.Name.Name

		// For external packages, use qualified name (package.Type)
		packageName := file.Name.Name
		qualifiedName := schemaName
		if packageName != "main" && packageName != "" {
			qualifiedName = packageName + "." + typeSpec.Name.Name
		}

		// Check if this type is referenced (skip if not)
		if !p.IsTypeReferenced(schemaName) && !p.IsTypeReferenced(qualifiedName) {
			return true
		}

		schema := processor.ProcessStruct(structType, typeSpec.Doc, typeSpec.Name.Name)
		if schema != nil {
			if packageName != "main" && packageName != "" {
				// Register both names to support both reference styles
				p.openapi.Components.Schemas[qualifiedName] = schema
				p.typeCache[qualifiedName] = &TypeInfo{
					Name:    qualifiedName,
					Package: packageName,
					Schema:  schema,
					ASTNode: typeSpec,
				}
			}

			// Always register with simple name too
			p.openapi.Components.Schemas[schemaName] = schema
			p.typeCache[schemaName] = &TypeInfo{
				Name:    schemaName,
				Package: packageName,
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
			pathItem.Query, // QUERY method (OpenAPI 3.2.0)
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

// AddReferencedType registers a type as referenced by an operation.
// This method extracts type names from schema references and marks them for processing.
func (p *Parser) AddReferencedType(typeName string) {
	if typeName == "" {
		return
	}

	// Remove #/components/schemas/ prefix if present
	const prefix = "#/components/schemas/"
	if strings.HasPrefix(typeName, prefix) {
		typeName = strings.TrimPrefix(typeName, prefix)
	}

	// Remove []prefix for array types
	typeName = strings.TrimPrefix(typeName, "[]")

	// Handle map types: map[key]value -> extract value type
	if strings.HasPrefix(typeName, "map[") {
		// This will be handled during schema processing
		return
	}

	// Skip primitive types
	primitives := map[string]bool{
		"string": true, "int": true, "int8": true, "int16": true, "int32": true, "int64": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
		"float32": true, "float64": true, "bool": true, "byte": true, "rune": true,
		"time.Time": true, "interface{}": true, "any": true,
	}

	if primitives[typeName] {
		return
	}

	p.referencedTypes[typeName] = true
}

// IsTypeReferenced checks if a type is referenced by operations.
func (p *Parser) IsTypeReferenced(typeName string) bool {
	return p.referencedTypes[typeName]
}

// GetReferencedTypes returns all referenced types.
func (p *Parser) GetReferencedTypes() map[string]bool {
	return p.referencedTypes
}

// ResolveTypeDependencies recursively resolves all type dependencies.
// This method processes structs referenced in operations and extracts their field types,
// including nested structs, to ensure all related types are included in the schema.
func (p *Parser) ResolveTypeDependencies() {
	// Keep processing until no new types are discovered
	maxIterations := 10
	iterations := 0
	for {
		iterations++
		if iterations > maxIterations {
			break
		}
		newTypesFound := false
		currentTypes := make([]string, 0, len(p.referencedTypes))

		// Copy current referenced types
		for typeName := range p.referencedTypes {
			currentTypes = append(currentTypes, typeName)
		}

		// Process each referenced type
		for _, typeName := range currentTypes {
			found := false
			// Find the type in parsed files
			for _, file := range p.files {
				ast.Inspect(file, func(n ast.Node) bool {
					typeSpec, ok := n.(*ast.TypeSpec)
					if !ok {
						return true
					}

					// Check if this is the type we're looking for
					if typeSpec.Name.Name != typeName {
						// Also check package.Type format
						packageName := file.Name.Name
						qualifiedName := packageName + "." + typeSpec.Name.Name
						if qualifiedName != typeName {
							return true
						}
					}

					// Process struct fields to find dependencies
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						return true
					}

					// Extract field types
					for _, field := range structType.Fields.List {
						// Check for swaggertype override first
						if field.Tag != nil {
							tagStr := strings.Trim(field.Tag.Value, "`")
							swaggerType := extractTag(tagStr, "swaggertype")

							// If swaggertype is primitive, skip dependency tracking
							if swaggerType != "" && !p.isPrimitiveSwaggerType(swaggerType) {
								// Extract the actual type from swaggertype
								actualType := p.extractTypeFromSwaggerType(swaggerType)
								if actualType != "" && !p.referencedTypes[actualType] {
									p.referencedTypes[actualType] = true
									newTypesFound = true
								}
								continue
							} else if swaggerType != "" {
								// Primitive swaggertype override, no need to track original type
								continue
							}
						}

						// Extract dependencies from field type
						fieldType := p.extractFieldTypeName(field.Type)
						if fieldType != "" && !p.referencedTypes[fieldType] {
							p.referencedTypes[fieldType] = true
							newTypesFound = true
						}
					}

					return true
				})
			}

			// If type not found in parsed files, try to parse it from external dependencies
			if !found {
				if err := p.parseExternalType(typeName); err == nil {
					newTypesFound = true
				}
			}
		}

		// Stop if no new types were discovered
		if !newTypesFound {
			break
		}
	}
}

// isPrimitiveSwaggerType checks if a swaggertype value represents a primitive type.
func (p *Parser) isPrimitiveSwaggerType(swaggerType string) bool {
	// Handle "primitive,type" format
	parts := strings.Split(swaggerType, ",")
	if len(parts) > 0 && parts[0] == "primitive" {
		return true
	}

	// Handle direct primitive types
	primitives := map[string]bool{
		"string": true, "integer": true, "number": true, "boolean": true,
		"object": true, "array": true,
	}

	return primitives[swaggerType]
}

// extractTypeFromSwaggerType extracts the actual type name from swaggertype tag.
// Example: "array,Product" -> "Product"
func (p *Parser) extractTypeFromSwaggerType(swaggerType string) string {
	parts := strings.Split(swaggerType, ",")
	if len(parts) > 1 {
		typeName := strings.TrimSpace(parts[1])
		// Remove []prefix for array types
		typeName = strings.TrimPrefix(typeName, "[]")
		return typeName
	}
	return ""
}

// extractFieldTypeName extracts the type name from an AST expression.
func (p *Parser) extractFieldTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// Simple type in same package
		return t.Name

	case *ast.StarExpr:
		// Pointer type - recurse
		return p.extractFieldTypeName(t.X)

	case *ast.ArrayType:
		// Array type - get element type
		return p.extractFieldTypeName(t.Elt)

	case *ast.MapType:
		// Map type - get value type
		return p.extractFieldTypeName(t.Value)

	case *ast.SelectorExpr:
		// External package type: package.Type
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name + "." + t.Sel.Name
		}

	case *ast.StructType:
		// Inline struct - no external dependency
		return ""

	case *ast.InterfaceType:
		// Interface - no concrete type dependency
		return ""
	}

	return ""
}

// collectImports collects import statements from a file and stores them in importMap.
// This allows resolving qualified type names (e.g., errors.BadRequest) to their full import paths.
func (p *Parser) collectImports(filePath string, file *ast.File) {
	if p.importMap[filePath] == nil {
		p.importMap[filePath] = make(map[string]string)
	}

	for _, imp := range file.Imports {
		if imp.Path == nil {
			continue
		}

		// Remove quotes from import path
		importPath := strings.Trim(imp.Path.Value, "`\"")

		// Determine package alias
		var packageAlias string
		if imp.Name != nil {
			// Explicit alias (e.g., import foo "github.com/x/y")
			packageAlias = imp.Name.Name
			if packageAlias == "_" || packageAlias == "." {
				// Skip blank imports and dot imports
				continue
			}
		} else {
			// No alias, use last component of import path
			// e.g., "github.com/user/repo/errors" -> "errors"
			parts := strings.Split(importPath, "/")
			packageAlias = parts[len(parts)-1]
		}

		p.importMap[filePath][packageAlias] = importPath
	}
}

// resolveQualifiedType resolves a qualified type name (e.g., "errors.BadRequest") to its full path.
// Returns the full type path (e.g., "github.com/user/repo/errors.BadRequest") or empty string if not found.
func (p *Parser) resolveQualifiedType(filePath, typeName string) string {
	// Check if type is qualified (contains ".")
	if !strings.Contains(typeName, ".") {
		return typeName
	}

	// Split into package alias and type name
	parts := strings.SplitN(typeName, ".", 2)
	if len(parts) != 2 {
		return typeName
	}

	packageAlias := parts[0]
	typeNameOnly := parts[1]

	// Look up import path for this package alias
	fileImports, ok := p.importMap[filePath]
	if !ok {
		return packageAlias + "." + typeNameOnly
	}

	importPath, ok := fileImports[packageAlias]
	if !ok {
		return packageAlias + "." + typeNameOnly
	}

	// Extract package name from import path (last component)
	// e.g., "github.com/dock-tech/isis-golang-lib/domainerrors" -> "domainerrors"
	importParts := strings.Split(importPath, "/")
	packageName := importParts[len(importParts)-1]

	return packageName + "." + typeNameOnly
}

// parseExternalType attempts to parse an external type from a dependency.
// Given a qualified type name (e.g., "domainerrors.BadRequest"), it tries to:
// 1. Find which file has the import for this package
// 2. Resolve the full import path
// 3. Parse the external package to get the type definition
func (p *Parser) parseExternalType(qualifiedType string) error {
	// Check if type is qualified
	if !strings.Contains(qualifiedType, ".") {
		return nil
	}

	parts := strings.SplitN(qualifiedType, ".", 2)
	if len(parts) != 2 {
		return nil
	}

	packageAlias := parts[0]

	// Search all files for an import matching this package alias
	var importPath string
	for _, imports := range p.importMap {
		if path, ok := imports[packageAlias]; ok {
			importPath = path
			break
		}
	}

	if importPath == "" {
		// No import found, type might be in same package
		return nil
	}

	// Try to find and parse this package
	return p.parseExternalPackage(importPath)
}

// parseExternalPackage parses an external package from dependencies.
func (p *Parser) parseExternalPackage(importPath string) error {
	if !p.parseDependency {
		return nil
	}

	// Check if we've already tried to parse this module
	if p.parsedModules[importPath] {
		return nil
	}

	// Mark as parsed to avoid repeated attempts
	p.parsedModules[importPath] = true

	// Find the package directory
	var pkgDir string

	// Check vendor first
	vendorDir := filepath.Join("vendor", importPath)
	if stat, err := os.Stat(vendorDir); err == nil && stat.IsDir() {
		pkgDir = vendorDir
	} else {
		// Try GOMODCACHE
		goModCache := os.Getenv("GOMODCACHE")
		if goModCache == "" {
			goPath := os.Getenv("GOPATH")
			if goPath == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return nil
				}
				goPath = filepath.Join(homeDir, "go")
			}
			goModCache = filepath.Join(goPath, "pkg", "mod")
		}

		// Try to find the package in GOMODCACHE
		// The import path might be a subpackage, so we need to search for it
		pkgDir = p.findPackageInCache(goModCache, importPath)
		if pkgDir == "" {
			return nil
		}
	}

	// Set flag to indicate we're parsing external packages
	// This prevents collecting imports from these files which could cause infinite recursion
	wasParsingExternal := p.parsingExternal
	p.parsingExternal = true
	defer func() {
		p.parsingExternal = wasParsingExternal
	}()

	// Parse all .go files in the package directory (non-recursively)
	return filepath.Walk(pkgDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip subdirectories
		if info.IsDir() && path != pkgDir {
			return filepath.SkipDir
		}

		// Parse only Go files (excluding test files)
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		return p.ParseFile(path)
	})
}

// findPackageInCache searches for a package in the GOMODCACHE.
// It handles both full module paths and subpackages within modules.
func (p *Parser) findPackageInCache(cacheDir, importPath string) string {
	// First, try to find the package directly
	parts := strings.Split(importPath, "/")

	// Try increasingly shorter paths to find the module root
	for i := len(parts); i >= 3; i-- { // Minimum 3 parts for domain/org/repo
		modulePath := strings.Join(parts[:i], "/")
		subPath := ""
		if i < len(parts) {
			subPath = strings.Join(parts[i:], "/")
		}

		// Build cache path with proper escaping for uppercase letters
		var cachePath strings.Builder
		for _, c := range modulePath {
			if c >= 'A' && c <= 'Z' {
				cachePath.WriteByte('!')
				cachePath.WriteRune(c + 32) // convert to lowercase
			} else {
				cachePath.WriteRune(c)
			}
		}

		// Try to find any version of this module
		pattern := filepath.Join(cacheDir, cachePath.String()) + "@*"
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			// Use the first (alphabetically last) version found
			moduleDir := matches[len(matches)-1]
			if subPath != "" {
				return filepath.Join(moduleDir, subPath)
			}
			return moduleDir
		}
	}

	return ""
}
