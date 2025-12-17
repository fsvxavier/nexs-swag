// Package v2 implements Swagger 2.0 specification generation.
package v2

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	swagger "github.com/fsvxavier/nexs-swag/pkg/openapi/v2"
)

// Generator generates Swagger 2.0 specification files.
type Generator struct {
	spec          *swagger.Swagger
	outputDir     string
	outputType    []string
	instanceName  string
	generatedTime bool
}

// New creates a new Swagger 2.0 Generator instance.
func New(spec *swagger.Swagger, outputDir string, outputType []string) *Generator {
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

// Generate generates all requested output formats.
func (g *Generator) Generate() error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Check if any operation has x-visibility annotation
	hasVisibility := g.hasVisibilityAnnotations()

	if hasVisibility {
		// Generate separate specs for public and private
		return g.generateSeparateSpecs()
	}

	// Generate normal spec
	for _, format := range g.outputType {
		format = strings.ToLower(strings.TrimSpace(format))
		switch format {
		case "json":
			if err := g.generateJSON(); err != nil {
				return err
			}
		case "yaml", "yml":
			if err := g.generateYAML(); err != nil {
				return err
			}
		case "go":
			if err := g.generateGo(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported output format: %s", format)
		}
	}

	return nil
}

// hasVisibilityAnnotations checks if any operation has x-visibility extension.
func (g *Generator) hasVisibilityAnnotations() bool {
	for _, pathItem := range g.spec.Paths {
		operations := []*swagger.Operation{
			pathItem.Get, pathItem.Post, pathItem.Put, pathItem.Delete,
			pathItem.Patch, pathItem.Options, pathItem.Head,
		}
		for _, op := range operations {
			if op != nil && op.Extensions != nil {
				if _, ok := op.Extensions["x-visibility"]; ok {
					return true
				}
			}
		}
	}
	return false
}

// generateSeparateSpecs generates separate public and private specifications.
func (g *Generator) generateSeparateSpecs() error {
	publicSpec := g.filterSpecByVisibility("public")
	privateSpec := g.filterSpecByVisibility("private")

	for _, format := range g.outputType {
		format = strings.ToLower(strings.TrimSpace(format))
		switch format {
		case "json":
			if err := g.generateJSONWithSuffix(publicSpec, "_public"); err != nil {
				return err
			}
			if err := g.generateJSONWithSuffix(privateSpec, "_private"); err != nil {
				return err
			}
		case "yaml", "yml":
			if err := g.generateYAMLWithSuffix(publicSpec, "_public"); err != nil {
				return err
			}
			if err := g.generateYAMLWithSuffix(privateSpec, "_private"); err != nil {
				return err
			}
		case "go":
			if err := g.generateGoWithSuffix(publicSpec, "_public"); err != nil {
				return err
			}
			if err := g.generateGoWithSuffix(privateSpec, "_private"); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported output format: %s", format)
		}
	}

	return nil
}

// filterSpecByVisibility creates a new spec containing only operations with the specified visibility.
func (g *Generator) filterSpecByVisibility(visibility string) *swagger.Swagger {
	filteredSpec := &swagger.Swagger{
		Swagger:             g.spec.Swagger,
		Info:                g.spec.Info,
		Host:                g.spec.Host,
		BasePath:            g.spec.BasePath,
		Schemes:             g.spec.Schemes,
		Consumes:            g.spec.Consumes,
		Produces:            g.spec.Produces,
		Paths:               make(swagger.Paths),
		Definitions:         make(map[string]*swagger.Schema),
		Parameters:          g.spec.Parameters,
		Responses:           g.spec.Responses,
		SecurityDefinitions: g.spec.SecurityDefinitions,
		Security:            g.spec.Security,
		Tags:                g.spec.Tags,
		ExternalDocs:        g.spec.ExternalDocs,
	}

	usedSchemas := make(map[string]bool)

	// Filter paths based on visibility
	for path, pathItem := range g.spec.Paths {
		filteredPathItem := &swagger.PathItem{}
		hasOperations := false

		for method, op := range map[string]*swagger.Operation{
			"get":     pathItem.Get,
			"post":    pathItem.Post,
			"put":     pathItem.Put,
			"delete":  pathItem.Delete,
			"patch":   pathItem.Patch,
			"options": pathItem.Options,
			"head":    pathItem.Head,
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
				}
				hasOperations = true

				// Collect schemas used by this operation
				g.collectSchemasFromOperation(op, usedSchemas)
			}
		}

		if hasOperations {
			filteredSpec.Paths[path] = filteredPathItem
		}
	}

	// Copy only used schemas
	for schemaName := range usedSchemas {
		if schema, ok := g.spec.Definitions[schemaName]; ok {
			filteredSpec.Definitions[schemaName] = schema
		}
	}

	return filteredSpec
}

// collectSchemasFromOperation collects all schemas referenced by an operation.
func (g *Generator) collectSchemasFromOperation(op *swagger.Operation, usedSchemas map[string]bool) {
	// Check parameters
	for _, param := range op.Parameters {
		if param.Schema != nil {
			g.collectSchemaRefs(param.Schema, usedSchemas)
		}
	}

	// Check responses
	for _, resp := range op.Responses {
		if resp.Schema != nil {
			g.collectSchemaRefs(resp.Schema, usedSchemas)
		}
	}
}

// collectSchemaRefs recursively collects all schema references.
func (g *Generator) collectSchemaRefs(schema *swagger.Schema, usedSchemas map[string]bool) {
	if schema == nil {
		return
	}

	// Check $ref
	if schema.Ref != "" {
		// Extract schema name from #/definitions/SchemaName
		schemaName := strings.TrimPrefix(schema.Ref, "#/definitions/")
		if schemaName != "" && !usedSchemas[schemaName] {
			usedSchemas[schemaName] = true
			// Recursively collect schemas referenced by this schema
			if refSchema, ok := g.spec.Definitions[schemaName]; ok {
				g.collectSchemaRefs(refSchema, usedSchemas)
			}
		}
		return
	}

	// Check properties
	for _, prop := range schema.Properties {
		g.collectSchemaRefs(prop, usedSchemas)
	}

	// Check items (for arrays)
	if schema.Items != nil {
		g.collectSchemaRefs(schema.Items, usedSchemas)
	}

	// Check allOf
	for _, s := range schema.AllOf {
		g.collectSchemaRefs(s, usedSchemas)
	}
}

// generateJSONWithSuffix generates JSON with a filename suffix.
func (g *Generator) generateJSONWithSuffix(spec *swagger.Swagger, suffix string) error {
	filePath := filepath.Join(g.outputDir, "swagger"+suffix+".json")
	jsonData, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}

// generateYAMLWithSuffix generates YAML with a filename suffix.
func (g *Generator) generateYAMLWithSuffix(spec *swagger.Swagger, suffix string) error {
	filePath := filepath.Join(g.outputDir, "swagger"+suffix+".yaml")
	yamlData, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}

// generateGoWithSuffix generates Go file with a filename suffix.
func (g *Generator) generateGoWithSuffix(spec *swagger.Swagger, suffix string) error {
	filePath := filepath.Join(g.outputDir, "docs"+suffix+".go")

	// Marshal JSON for embedding
	jsonData, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

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

	// Determine variable name based on suffix
	varName := "SwaggerDoc"
	if suffix == "_public" {
		varName = "SwaggerDocPublic"
	} else if suffix == "_private" {
		varName = "SwaggerDocPrivate"
	}

	// Write SwaggerDoc
	if _, err := fmt.Fprintf(file, "// %s is the Swagger 2.0 specification in JSON format\n", varName); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(file, "var %s = `%s`\n\n", varName, string(jsonData)); err != nil {
		return err
	}

	// Write ReadDoc function
	funcName := "ReadDoc"
	if suffix == "_public" {
		funcName = "ReadDocPublic"
	} else if suffix == "_private" {
		funcName = "ReadDocPrivate"
	}

	if _, err := fmt.Fprintf(file, "// %s returns the Swagger specification\n", funcName); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(file, "func %s() string {\n", funcName); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(file, "\treturn %s\n", varName); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(file, "}"); err != nil {
		return err
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}

// generateJSON generates swagger.json file.
func (g *Generator) generateJSON() error {
	filePath := filepath.Join(g.outputDir, "swagger.json")
	jsonData, err := json.MarshalIndent(g.spec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}

// generateYAML generates swagger.yaml file.
func (g *Generator) generateYAML() error {
	filePath := filepath.Join(g.outputDir, "swagger.yaml")
	yamlData, err := yaml.Marshal(g.spec)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	fmt.Printf("Generated: %s\n", filePath)
	return nil
}

// generateGo generates docs.go file with embedded specification.
func (g *Generator) generateGo() error {
	filePath := filepath.Join(g.outputDir, "docs.go")

	// Marshal JSON for embedding
	jsonData, err := json.Marshal(g.spec)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

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

	// Write SwaggerDoc
	if _, err := fmt.Fprintln(file, "// SwaggerDoc is the Swagger 2.0 specification in JSON format"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(file, "var SwaggerDoc = `%s`\n\n", string(jsonData)); err != nil {
		return err
	}

	// Write ReadDoc function
	if _, err := fmt.Fprintln(file, "// ReadDoc returns the Swagger specification"); err != nil {
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
