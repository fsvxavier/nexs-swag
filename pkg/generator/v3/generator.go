// Package v3 implements OpenAPI 3.x specification generation.
package v3

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	openapi "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

// Generator generates OpenAPI 3.x specification files.
type Generator struct {
	spec           *openapi.OpenAPI
	outputDir      string
	outputType     []string
	instanceName   string
	generatedTime  bool
	templateDelims []string // [leftDelim, rightDelim]
	openapiVersion string   // OpenAPI version to use
}

// New creates a new OpenAPI 3.x Generator instance.
func New(spec *openapi.OpenAPI, outputDir string, outputType []string) *Generator {
	return &Generator{
		spec:          spec,
		outputDir:     outputDir,
		outputType:    outputType,
		instanceName:  "docs",
		generatedTime: false,
	}
}

// SetInstanceName sets the package name for generated Go file.
func (g *Generator) SetInstanceName(name string) {
	g.instanceName = name
}

// SetGeneratedTime sets whether to include generation timestamp.
func (g *Generator) SetGeneratedTime(enabled bool) {
	g.generatedTime = enabled
}

// SetTemplateDelims sets custom template delimiters.
// Format: "leftDelim,rightDelim".
func (g *Generator) SetTemplateDelims(delims string) {
	if delims == "" {
		g.templateDelims = []string{"{{", "}}"}
		return
	}
	parts := strings.Split(delims, ",")
	if len(parts) == 2 {
		g.templateDelims = []string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])}
	} else {
		g.templateDelims = []string{"{{", "}}"}
	}
}

// SetOpenAPIVersion sets the OpenAPI version for the specification.
// Supported versions: 3.0.0, 3.0.1, 3.0.2, 3.0.3, 3.0.4, 3.1.0, 3.1.1, 3.1.2, 3.2.0
func (g *Generator) SetOpenAPIVersion(version string) {
	g.openapiVersion = version
}

// Generate generates the OpenAPI specification files.
func (g *Generator) Generate() error {
	// Update spec OpenAPI version if set
	if g.openapiVersion != "" {
		g.spec.OpenAPI = g.openapiVersion
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Check if we need to generate separate public/private specs
	hasVisibility := g.hasVisibilityAnnotations()

	if hasVisibility {
		// Generate separate specs for public and private
		if err := g.generateSeparateSpecs(); err != nil {
			return err
		}
	} else {
		// Generate single spec as before
		for _, outputType := range g.outputType {
			switch outputType {
			case "json":
				if err := g.generateJSON(); err != nil {
					return fmt.Errorf("failed to generate JSON: %w", err)
				}
			case "yaml", "yml":
				if err := g.generateYAML(); err != nil {
					return fmt.Errorf("failed to generate YAML: %w", err)
				}
			case "go":
				if err := g.generateGo(); err != nil {
					return fmt.Errorf("failed to generate Go file: %w", err)
				}
			default:
				return fmt.Errorf("unsupported output type: %s", outputType)
			}
		}
	}

	return nil
}

// generateJSON generates the OpenAPI specification in JSON format.
func (g *Generator) generateJSON() error {
	filePath := filepath.Join(g.outputDir, "openapi.json")

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(g.spec); err != nil {
		return err
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}

// generateYAML generates the OpenAPI specification in YAML format.
func (g *Generator) generateYAML() error {
	filePath := filepath.Join(g.outputDir, "openapi.yaml")

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	if err := encoder.Encode(g.spec); err != nil {
		return err
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}

// hasVisibilityAnnotations checks if any operation has x-visibility extension.
func (g *Generator) hasVisibilityAnnotations() bool {
	for _, pathItem := range g.spec.Paths {
		for _, op := range []*openapi.Operation{
			pathItem.Get, pathItem.Post, pathItem.Put, pathItem.Delete,
			pathItem.Patch, pathItem.Options, pathItem.Head, pathItem.Trace,
		} {
			if op != nil && op.Extensions != nil {
				if _, ok := op.Extensions["x-visibility"]; ok {
					return true
				}
			}
		}
	}
	return false
}

// generateSeparateSpecs generates separate public and private OpenAPI specs.
func (g *Generator) generateSeparateSpecs() error {
	// Create public spec
	publicSpec := g.filterSpecByVisibility("public")
	// Create private spec
	privateSpec := g.filterSpecByVisibility("private")

	// Generate files for each visibility level
	for _, outputType := range g.outputType {
		switch outputType {
		case "json":
			if err := g.generateJSONWithSuffix(publicSpec, "_public"); err != nil {
				return fmt.Errorf("failed to generate public JSON: %w", err)
			}
			if err := g.generateJSONWithSuffix(privateSpec, "_private"); err != nil {
				return fmt.Errorf("failed to generate private JSON: %w", err)
			}
		case "yaml", "yml":
			if err := g.generateYAMLWithSuffix(publicSpec, "_public"); err != nil {
				return fmt.Errorf("failed to generate public YAML: %w", err)
			}
			if err := g.generateYAMLWithSuffix(privateSpec, "_private"); err != nil {
				return fmt.Errorf("failed to generate private YAML: %w", err)
			}
		case "go":
			if err := g.generateGoWithSuffix(publicSpec, "_public"); err != nil {
				return fmt.Errorf("failed to generate public Go file: %w", err)
			}
			if err := g.generateGoWithSuffix(privateSpec, "_private"); err != nil {
				return fmt.Errorf("failed to generate private Go file: %w", err)
			}
		default:
			return fmt.Errorf("unsupported output type: %s", outputType)
		}
	}

	return nil
}

// filterSpecByVisibility creates a new spec containing only operations with the specified visibility.
func (g *Generator) filterSpecByVisibility(visibility string) *openapi.OpenAPI {
	filteredSpec := &openapi.OpenAPI{
		OpenAPI:           g.spec.OpenAPI,
		JSONSchemaDialect: g.spec.JSONSchemaDialect,
		Info:              g.spec.Info,
		Servers:           g.spec.Servers,
		Paths:             make(map[string]*openapi.PathItem),
		Components:        &openapi.Components{Schemas: make(map[string]*openapi.Schema)},
		Security:          g.spec.Security,
		Tags:              g.spec.Tags,
		ExternalDocs:      g.spec.ExternalDocs,
	}

	usedSchemas := make(map[string]bool)

	// Filter paths based on visibility
	for path, pathItem := range g.spec.Paths {
		filteredPathItem := &openapi.PathItem{}
		hasOperations := false

		for method, op := range map[string]*openapi.Operation{
			"get":     pathItem.Get,
			"post":    pathItem.Post,
			"put":     pathItem.Put,
			"delete":  pathItem.Delete,
			"patch":   pathItem.Patch,
			"options": pathItem.Options,
			"head":    pathItem.Head,
			"trace":   pathItem.Trace,
		} {
			if op == nil {
				continue
			}

			// Check visibility - if not set or empty, include in both
			opVisibility := ""
			if op.Extensions != nil {
				if vis, ok := op.Extensions["x-visibility"].(string); ok {
					opVisibility = vis
				}
			}

			// Include operation if:
			// 1. No x-visibility set (empty) - include in both specs
			// 2. x-visibility matches the current visibility filter
			if opVisibility == "" || opVisibility == visibility {
				hasOperations = true
				switch method {
				case "get":
					filteredPathItem.Get = op
				case "post":
					filteredPathItem.Post = op
				case "put":
					filteredPathItem.Put = op
				case "delete":
					filteredPathItem.Delete = op
				case "patch":
					filteredPathItem.Patch = op
				case "options":
					filteredPathItem.Options = op
				case "head":
					filteredPathItem.Head = op
				case "trace":
					filteredPathItem.Trace = op
				}

				// Collect schemas used in this operation
				g.collectSchemasFromOperation(op, usedSchemas)
			}
		}

		if hasOperations {
			filteredSpec.Paths[path] = filteredPathItem
		}
	}

	// Copy only used schemas
	for schemaName := range usedSchemas {
		if schema, ok := g.spec.Components.Schemas[schemaName]; ok {
			filteredSpec.Components.Schemas[schemaName] = schema
		}
	}

	return filteredSpec
}

// collectSchemasFromOperation collects all schema names referenced in an operation.
func (g *Generator) collectSchemasFromOperation(op *openapi.Operation, usedSchemas map[string]bool) {
	// Collect from request body
	if op.RequestBody != nil && op.RequestBody.Content != nil {
		for _, mediaType := range op.RequestBody.Content {
			if mediaType.Schema != nil {
				g.collectSchemaRefs(mediaType.Schema, usedSchemas)
			}
		}
	}

	// Collect from responses
	for _, response := range op.Responses {
		if response.Content != nil {
			for _, mediaType := range response.Content {
				if mediaType.Schema != nil {
					g.collectSchemaRefs(mediaType.Schema, usedSchemas)
				}
			}
		}
	}

	// Collect from parameters
	for _, param := range op.Parameters {
		if param.Schema != nil {
			g.collectSchemaRefs(param.Schema, usedSchemas)
		}
	}
}

// collectSchemaRefs recursively collects schema references.
func (g *Generator) collectSchemaRefs(schema *openapi.Schema, usedSchemas map[string]bool) {
	if schema == nil {
		return
	}

	// Check for $ref
	if schema.Ref != "" {
		// Extract schema name from #/components/schemas/SchemaName
		parts := strings.Split(schema.Ref, "/")
		if len(parts) > 0 {
			schemaName := parts[len(parts)-1]
			if !usedSchemas[schemaName] {
				usedSchemas[schemaName] = true
				// Recursively collect schemas from referenced schema
				if refSchema, ok := g.spec.Components.Schemas[schemaName]; ok {
					g.collectSchemaRefs(refSchema, usedSchemas)
				}
			}
		}
	}

	// Check properties
	for _, propSchema := range schema.Properties {
		g.collectSchemaRefs(propSchema, usedSchemas)
	}

	// Check items (for arrays)
	if schema.Items != nil {
		g.collectSchemaRefs(schema.Items, usedSchemas)
	}

	// Check allOf, anyOf, oneOf
	for _, s := range schema.AllOf {
		g.collectSchemaRefs(&s, usedSchemas)
	}
	for _, s := range schema.AnyOf {
		g.collectSchemaRefs(&s, usedSchemas)
	}
	for _, s := range schema.OneOf {
		g.collectSchemaRefs(&s, usedSchemas)
	}
}

// generateJSONWithSuffix generates JSON with a filename suffix.
func (g *Generator) generateJSONWithSuffix(spec *openapi.OpenAPI, suffix string) error {
	filePath := filepath.Join(g.outputDir, fmt.Sprintf("openapi%s.json", suffix))

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(spec); err != nil {
		return err
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}

// generateYAMLWithSuffix generates YAML with a filename suffix.
func (g *Generator) generateYAMLWithSuffix(spec *openapi.OpenAPI, suffix string) error {
	filePath := filepath.Join(g.outputDir, fmt.Sprintf("openapi%s.yaml", suffix))

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	if err := encoder.Encode(spec); err != nil {
		return err
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}

// generateGoWithSuffix generates Go file with a suffix.
func (g *Generator) generateGoWithSuffix(spec *openapi.OpenAPI, suffix string) error {
	filePath := filepath.Join(g.outputDir, fmt.Sprintf("docs%s.go", suffix))

	// Marshal to JSON for embedding
	jsonData, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write package declaration
	varName := "SwaggerDoc"
	if suffix == "_public" {
		varName = "SwaggerDocPublic"
	} else if suffix == "_private" {
		varName = "SwaggerDocPrivate"
	}

	if _, err := fmt.Fprintf(file, "// Package %s Code generated by nexs-swag. DO NOT EDIT\n", g.instanceName); err != nil {
		return err
	}

	if g.generatedTime {
		if _, err := fmt.Fprintf(file, "// Generated at: %s\n", time.Now().Format("2006-01-02 15:04:05")); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(file, "package %s\n\n", g.instanceName); err != nil {
		return err
	}

	// Write SwaggerDoc variable
	if _, err := fmt.Fprintf(file, "// %s is the OpenAPI v3 specification in JSON format\n", varName); err != nil {
		return err
	}

	// Apply template delimiters if configured
	leftDelim := "{{"
	rightDelim := "}}"
	if len(g.templateDelims) == 2 {
		leftDelim = g.templateDelims[0]
		rightDelim = g.templateDelims[1]
	}

	docStr := string(jsonData)
	docStr = strings.ReplaceAll(docStr, "{{", leftDelim)
	docStr = strings.ReplaceAll(docStr, "}}", rightDelim)

	if _, err := fmt.Fprintf(file, "var %s = `%s`\n", varName, docStr); err != nil {
		return err
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}

// generateGo generates a Go file with embedded OpenAPI specification.
func (g *Generator) generateGo() error {
	filePath := filepath.Join(g.outputDir, "docs.go")

	// Marshal to JSON for embedding
	jsonData, err := json.MarshalIndent(g.spec, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Write package declaration
	if _, err := fmt.Fprintf(file, "// Package %s Code generated by nexs-swag. DO NOT EDIT\n", g.instanceName); err != nil {
		return err
	}

	if g.generatedTime {
		if _, err := fmt.Fprintf(file, "// Generated at: %s\n", time.Now().Format("2006-01-02 15:04:05")); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(file, "package %s\n\n", g.instanceName); err != nil {
		return err
	}

	// Write SwaggerDoc variable
	if _, err := fmt.Fprintln(file, "// SwaggerDoc is the OpenAPI v3 specification in JSON format"); err != nil {
		return err
	}

	// Apply template delimiters if configured
	leftDelim := "{{"
	rightDelim := "}}"
	if len(g.templateDelims) == 2 {
		leftDelim = g.templateDelims[0]
		rightDelim = g.templateDelims[1]
	}

	docStr := string(jsonData)
	docStr = strings.ReplaceAll(docStr, "{{", leftDelim)
	docStr = strings.ReplaceAll(docStr, "}}", rightDelim)

	if _, err := fmt.Fprintf(file, "var SwaggerDoc = `%s`\n\n", docStr); err != nil {
		return err
	}

	// Write ReadDoc function
	if _, err := fmt.Fprintln(file, "// ReadDoc returns the OpenAPI specification"); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(file, "func ReadDoc() string {"); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(file, "\treturn SwaggerDoc"); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(file, "}"); err != nil {
		return err
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}
