package parser

import (
	"go/ast"
	"regexp"
	"strconv"
	"strings"

	"github.com/fsvxavier/nexs-swag/pkg/openapi"
)

// SchemaProcessor processes struct type definitions to generate OpenAPI schemas.
type SchemaProcessor struct {
	parser    *Parser
	openapi   *openapi.OpenAPI
	typeCache map[string]*TypeInfo
	depth     int // Current parsing depth for nested structures
}

// NewSchemaProcessor creates a new schema processor.
func NewSchemaProcessor(p *Parser, spec *openapi.OpenAPI, typeCache map[string]*TypeInfo) *SchemaProcessor {
	return &SchemaProcessor{
		parser:    p,
		openapi:   spec,
		typeCache: typeCache,
		depth:     0,
	}
}

// ProcessStruct processes a struct type and returns an OpenAPI schema.
func (s *SchemaProcessor) ProcessStruct(structType *ast.StructType, doc *ast.CommentGroup, typeName string) *openapi.Schema {
	schema := &openapi.Schema{
		Type:       "object",
		Properties: make(map[string]*openapi.Schema),
		Required:   []string{},
	}

	// Parse struct-level documentation
	if doc != nil {
		s.parseStructDoc(doc, schema)
	}

	// Parse struct fields
	for _, field := range structType.Fields.List {
		s.processField(field, schema)
	}

	return schema
}

// parseStructDoc parses struct-level documentation comments.
func (s *SchemaProcessor) parseStructDoc(doc *ast.CommentGroup, schema *openapi.Schema) {
	for _, comment := range doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))

		// @Description annotation
		if strings.HasPrefix(text, "@Description ") {
			desc := strings.TrimPrefix(text, "@Description ")
			if schema.Description == "" {
				schema.Description = desc
			} else {
				schema.Description += "\n" + desc
			}
		}

		// @Title annotation
		if strings.HasPrefix(text, "@Title ") {
			schema.Title = strings.TrimPrefix(text, "@Title ")
		}

		// @Deprecated annotation
		if strings.TrimSpace(text) == "@Deprecated" {
			schema.Deprecated = true
		}

		// @Example annotation
		if strings.HasPrefix(text, "@Example ") {
			// Example would need JSON parsing
			// For now, store as string
			schema.Example = strings.TrimPrefix(text, "@Example ")
		}
	}
}

// processField processes a single struct field.
func (s *SchemaProcessor) processField(field *ast.Field, schema *openapi.Schema) {
	if len(field.Names) == 0 {
		// Handle embedded fields
		s.processEmbeddedField(field, schema)
		return
	}

	fieldName := field.Names[0].Name

	// Skip unexported fields
	if !ast.IsExported(fieldName) {
		return
	}

	// Check swaggerignore tag first (before parsing all tags)
	if field.Tag != nil {
		tagStr := strings.Trim(field.Tag.Value, "`")
		if ignore := extractTag(tagStr, "swaggerignore"); strings.EqualFold(ignore, "true") {
			return // Skip this field completely
		}
	}

	// Parse struct tags
	tags := s.parseStructTags(field)

	// Skip if json tag is "-"
	if tags.JSON == "-" {
		return
	}

	// Get JSON name using property naming strategy
	jsonName := s.applyPropertyNaming(fieldName, tags.JSON)

	// Check if field type is a pointer
	isPointer := false
	fieldType := field.Type
	if starExpr, ok := fieldType.(*ast.StarExpr); ok {
		isPointer = true
		fieldType = starExpr.X
	}

	// Create field schema
	fieldSchema := s.processFieldType(field.Type)

	// Add field documentation
	if field.Doc != nil {
		s.parseFieldDoc(field.Doc, fieldSchema)
	}

	// Apply struct tag validations and attributes
	s.applyStructTagAttributes(tags, fieldSchema)

	// Add to properties
	schema.Properties[jsonName] = fieldSchema

	// Determine if field should be required
	if s.shouldBeRequired(tags, isPointer) {
		schema.Required = append(schema.Required, jsonName)
	}
}

// processEmbeddedField processes an embedded/anonymous field.
func (s *SchemaProcessor) processEmbeddedField(field *ast.Field, schema *openapi.Schema) {
	// Get the type name of the embedded field
	var typeName string

	switch t := field.Type.(type) {
	case *ast.Ident:
		typeName = t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			typeName = ident.Name
		}
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			typeName = ident.Name + "." + t.Sel.Name
		}
	}

	// Use allOf composition for embedded fields
	if typeName != "" {
		ref := openapi.Schema{
			Ref: "#/components/schemas/" + typeName,
		}
		schema.AllOf = append(schema.AllOf, ref)
	}
}

// parseFieldDoc parses field-level documentation.
func (s *SchemaProcessor) parseFieldDoc(doc *ast.CommentGroup, schema *openapi.Schema) {
	for _, comment := range doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))

		// Simple description (no annotation prefix)
		if !strings.HasPrefix(text, "@") && text != "" {
			if schema.Description == "" {
				schema.Description = text
			} else {
				schema.Description += " " + text
			}
		}

		// @Description annotation
		if strings.HasPrefix(text, "@Description ") {
			desc := strings.TrimPrefix(text, "@Description ")
			if schema.Description == "" {
				schema.Description = desc
			} else {
				schema.Description += " " + desc
			}
		}

		// @Example annotation
		if strings.HasPrefix(text, "@Example ") {
			schema.Example = strings.TrimPrefix(text, "@Example ")
		}

		// @Deprecated annotation
		if text == "@Deprecated" {
			schema.Deprecated = true
		}
	}
}

// StructTags represents parsed struct tag information.
type StructTags struct {
	JSON        string
	Binding     string
	Validate    string
	SwaggerType string // swaggertype tag for type override
	Extensions  string // extensions tag for x-* custom extensions
	Example     string
	Format      string
	Default     string
	Enum        string
	Minimum     string
	Maximum     string
	MinLength   string
	MaxLength   string
	Pattern     string
	Required    bool
	OmitEmpty   bool
	ReadOnly    bool
	WriteOnly   bool
}

// parseStructTags parses struct tags.
func (s *SchemaProcessor) parseStructTags(field *ast.Field) StructTags {
	tags := StructTags{}

	if field.Tag == nil {
		return tags
	}

	tagStr := strings.Trim(field.Tag.Value, "`")

	// Parse JSON tag
	if jsonTag := extractTag(tagStr, "json"); jsonTag != "" {
		tags.JSON = jsonTag
		if strings.Contains(jsonTag, "omitempty") {
			tags.OmitEmpty = true
		}
	}

	// Parse binding tag (gin framework)
	if bindingTag := extractTag(tagStr, "binding"); bindingTag != "" {
		tags.Binding = bindingTag
		if strings.Contains(bindingTag, "required") {
			tags.Required = true
		}
	}

	// Parse validate tag
	if validateTag := extractTag(tagStr, "validate"); validateTag != "" {
		tags.Validate = validateTag
		if strings.Contains(validateTag, "required") {
			tags.Required = true
		}
	}

	// Parse custom OpenAPI tags
	tags.SwaggerType = extractTag(tagStr, "swaggertype")
	tags.Extensions = extractTag(tagStr, "extensions")
	tags.Example = extractTag(tagStr, "example")
	tags.Format = extractTag(tagStr, "format")
	tags.Default = extractTag(tagStr, "default")
	tags.Enum = extractTag(tagStr, "enum")
	tags.Minimum = extractTag(tagStr, "minimum")
	tags.Maximum = extractTag(tagStr, "maximum")
	tags.MinLength = extractTag(tagStr, "minLength")
	tags.MaxLength = extractTag(tagStr, "maxLength")
	tags.Pattern = extractTag(tagStr, "pattern")

	// Check for readonly/writeonly
	if extractTag(tagStr, "readonly") == "true" {
		tags.ReadOnly = true
	}
	if extractTag(tagStr, "writeonly") == "true" {
		tags.WriteOnly = true
	}

	return tags
}

// applyStructTagAttributes applies struct tag attributes to schema.
func (s *SchemaProcessor) applyStructTagAttributes(tags StructTags, schema *openapi.Schema) {
	// Apply swaggertype override first (highest priority)
	if tags.SwaggerType != "" {
		s.applySwaggerType(tags.SwaggerType, schema)
	}

	// Apply extensions
	if tags.Extensions != "" {
		s.applyExtensions(tags.Extensions, schema)
	}

	// Apply explicit tag attributes
	if tags.Example != "" {
		schema.Example = tags.Example
	}

	if tags.Format != "" {
		schema.Format = tags.Format
	}

	if tags.Default != "" {
		schema.Default = tags.Default
	}

	if tags.Enum != "" {
		enumValues := strings.Split(tags.Enum, ",")
		for _, v := range enumValues {
			schema.Enum = append(schema.Enum, strings.TrimSpace(v))
		}
	}

	if tags.Minimum != "" {
		if min, err := strconv.ParseFloat(tags.Minimum, 64); err == nil {
			schema.Minimum = min
		}
	}

	if tags.Maximum != "" {
		if max, err := strconv.ParseFloat(tags.Maximum, 64); err == nil {
			schema.Maximum = max
		}
	}

	if tags.MinLength != "" {
		if minLen, err := strconv.Atoi(tags.MinLength); err == nil {
			schema.MinLength = minLen
		}
	}

	if tags.MaxLength != "" {
		if maxLen, err := strconv.Atoi(tags.MaxLength); err == nil {
			schema.MaxLength = maxLen
		}
	}

	if tags.Pattern != "" {
		schema.Pattern = tags.Pattern
	}

	if tags.ReadOnly {
		schema.ReadOnly = true
	}

	if tags.WriteOnly {
		schema.WriteOnly = true
	}

	// Process binding tag validations (Gin framework)
	if tags.Binding != "" {
		s.applyBindingValidations(tags.Binding, schema)
	}

	// Process validate tag validations (go-playground/validator)
	if tags.Validate != "" {
		s.applyValidateRules(tags.Validate, schema)
	}
}

// applyBindingValidations parses binding tag and applies validations.
// Common binding rules: required, email, min, max, len, etc.
func (s *SchemaProcessor) applyBindingValidations(binding string, schema *openapi.Schema) {
	rules := strings.Split(binding, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		// Skip 'required' as it's handled separately
		if rule == "required" {
			continue
		}

		// Handle email validation
		if rule == "email" {
			schema.Format = "email"
		}

		// Handle url validation
		if rule == "url" {
			schema.Format = "uri"
		}

		// Handle min/max with values: min=1, max=100
		if strings.HasPrefix(rule, "min=") {
			value := strings.TrimPrefix(rule, "min=")
			if schema.Type == "string" || schema.Type == "array" {
				if minLen, err := strconv.Atoi(value); err == nil {
					if schema.Type == "string" {
						schema.MinLength = minLen
					} else if schema.Type == "array" {
						schema.MinItems = minLen
					}
				}
			} else if schema.Type == "integer" || schema.Type == "number" {
				if min, err := strconv.ParseFloat(value, 64); err == nil {
					schema.Minimum = min
				}
			}
		}

		if strings.HasPrefix(rule, "max=") {
			value := strings.TrimPrefix(rule, "max=")
			if schema.Type == "string" || schema.Type == "array" {
				if maxLen, err := strconv.Atoi(value); err == nil {
					if schema.Type == "string" {
						schema.MaxLength = maxLen
					} else if schema.Type == "array" {
						schema.MaxItems = maxLen
					}
				}
			} else if schema.Type == "integer" || schema.Type == "number" {
				if max, err := strconv.ParseFloat(value, 64); err == nil {
					schema.Maximum = max
				}
			}
		}

		// Handle len (exact length for strings/arrays)
		if strings.HasPrefix(rule, "len=") {
			value := strings.TrimPrefix(rule, "len=")
			if length, err := strconv.Atoi(value); err == nil {
				schema.MinLength = length
				schema.MaxLength = length
			}
		}

		// Handle gte (greater than or equal)
		if strings.HasPrefix(rule, "gte=") {
			value := strings.TrimPrefix(rule, "gte=")
			if min, err := strconv.ParseFloat(value, 64); err == nil {
				schema.Minimum = min
			}
		}

		// Handle lte (less than or equal)
		if strings.HasPrefix(rule, "lte=") {
			value := strings.TrimPrefix(rule, "lte=")
			if max, err := strconv.ParseFloat(value, 64); err == nil {
				schema.Maximum = max
			}
		}

		// Handle gt (greater than)
		if strings.HasPrefix(rule, "gt=") {
			value := strings.TrimPrefix(rule, "gt=")
			if min, err := strconv.ParseFloat(value, 64); err == nil {
				schema.ExclusiveMinimum = min
			}
		}

		// Handle lt (less than)
		if strings.HasPrefix(rule, "lt=") {
			value := strings.TrimPrefix(rule, "lt=")
			if max, err := strconv.ParseFloat(value, 64); err == nil {
				schema.ExclusiveMaximum = max
			}
		}

		// Handle oneof (enum)
		if strings.HasPrefix(rule, "oneof=") {
			value := strings.TrimPrefix(rule, "oneof=")
			enumValues := strings.Split(value, " ")
			for _, v := range enumValues {
				schema.Enum = append(schema.Enum, v)
			}
		}
	}
}

// applyValidateRules parses validate tag and applies validations.
// Supports go-playground/validator rules.
func (s *SchemaProcessor) applyValidateRules(validate string, schema *openapi.Schema) {
	// Similar to binding, but more comprehensive
	// For now, reuse the binding logic as they share many rules
	s.applyBindingValidations(validate, schema)

	// Additional validator-specific rules can be added here
	rules := strings.Split(validate, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		// Handle uuid validation
		if rule == "uuid" || rule == "uuid4" {
			schema.Format = "uuid"
		}

		// Handle datetime validation
		if rule == "datetime" {
			schema.Format = "date-time"
		}

		// Handle date validation
		if rule == "date" {
			schema.Format = "date"
		}

		// Handle numeric validation
		if rule == "numeric" {
			schema.Pattern = "^[0-9]+$"
		}

		// Handle alpha validation
		if rule == "alpha" {
			schema.Pattern = "^[a-zA-Z]+$"
		}

		// Handle alphanum validation
		if rule == "alphanum" {
			schema.Pattern = "^[a-zA-Z0-9]+$"
		}
	}
}

// processFieldType processes a field type and returns a schema.
func (s *SchemaProcessor) processFieldType(expr ast.Expr) *openapi.Schema {
	// Check depth limit if parseDepth is set
	maxDepth := s.parser.GetParseDepth()
	if maxDepth > 0 && s.depth >= maxDepth {
		// Return empty schema when depth limit is reached
		return &openapi.Schema{}
	}

	schema := &openapi.Schema{}

	switch t := expr.(type) {
	case *ast.Ident:
		// Simple type (built-in or defined in same package)
		return s.identToSchema(t.Name)

	case *ast.ArrayType:
		// Array type
		schema.Type = "array"
		s.depth++
		schema.Items = s.processFieldType(t.Elt)
		s.depth--
		return schema

	case *ast.MapType:
		// Map type
		schema.Type = "object"
		s.depth++
		valueSchema := s.processFieldType(t.Value)
		s.depth--
		schema.AdditionalProperties = valueSchema
		return schema

	case *ast.StarExpr:
		// Pointer type - process the underlying type
		return s.processFieldType(t.X)

	case *ast.SelectorExpr:
		// External package type
		if ident, ok := t.X.(*ast.Ident); ok {
			typeName := ident.Name + "." + t.Sel.Name
			fullTypeName := typeName

			// Check for type override from .swaggo file
			if override, exists := s.parser.GetTypeOverride(fullTypeName); exists {
				return s.parseOverrideType(override)
			}

			schema.Ref = "#/components/schemas/" + typeName
		}
		return schema

	case *ast.InterfaceType:
		// Interface type - use empty schema (any type)
		return schema

	case *ast.StructType:
		// Inline struct - process it with depth tracking
		s.depth++
		result := s.ProcessStruct(t, nil, "")
		s.depth--
		return result
	}

	return schema
}

// identToSchema converts an identifier to a schema.
func (s *SchemaProcessor) identToSchema(name string) *openapi.Schema {
	schema := &openapi.Schema{}

	// Check for primitive types
	switch name {
	case "string":
		schema.Type = "string"
	case "int", "int8", "int16", "int32":
		schema.Type = "integer"
		schema.Format = "int32"
	case "int64":
		schema.Type = "integer"
		schema.Format = "int64"
	case "uint", "uint8", "uint16", "uint32":
		schema.Type = "integer"
		schema.Format = "int32"
	case "uint64":
		schema.Type = "integer"
		schema.Format = "int64"
	case "float32":
		schema.Type = "number"
		schema.Format = "float"
	case "float64":
		schema.Type = "number"
		schema.Format = "double"
	case "bool":
		schema.Type = "boolean"
	case "byte":
		schema.Type = "string"
		schema.Format = "byte"
	case "rune":
		schema.Type = "integer"
		schema.Format = "int32"
	default:
		// Reference to another schema
		schema.Ref = "#/components/schemas/" + name
	}

	return schema
}

// applySwaggerType applies swaggertype tag override to schema.
// Supports formats:
// - "integer", "string", "number", "boolean", "object", "array"
// - "primitive,integer" - convert struct to primitive type
// - "array,number" - convert to array of numbers
func (s *SchemaProcessor) applySwaggerType(swaggerType string, schema *openapi.Schema) {
	if swaggerType == "" {
		return
	}

	parts := strings.Split(swaggerType, ",")

	if len(parts) == 1 {
		// Simple type override: swaggertype:"integer"
		typeStr := strings.TrimSpace(parts[0])
		switch typeStr {
		case "string", "integer", "number", "boolean", "object", "array":
			schema.Type = typeStr
			// Clear Ref when overriding type
			schema.Ref = ""
		}
	} else if len(parts) == 2 {
		modifier := strings.TrimSpace(parts[0])
		typeStr := strings.TrimSpace(parts[1])

		if modifier == "primitive" {
			// Primitive type: swaggertype:"primitive,integer"
			switch typeStr {
			case "string", "integer", "number", "boolean":
				schema.Type = typeStr
				schema.Ref = ""
			}
		} else if modifier == "array" {
			// Array type: swaggertype:"array,number"
			schema.Type = "array"
			schema.Ref = ""
			switch typeStr {
			case "string", "integer", "number", "boolean", "object":
				schema.Items = &openapi.Schema{
					Type: typeStr,
				}
			}
		}
	}
}

// applyExtensions applies extensions tag to schema.
// Supports formats:
// - "x-nullable" - boolean true
// - "x-abc=def" - string value
// - "!x-omitempty" - boolean false (negation)
// - "x-nullable,x-abc=def,!x-omitempty" - multiple extensions
func (s *SchemaProcessor) applyExtensions(extensionsTag string, schema *openapi.Schema) {
	if extensionsTag == "" {
		return
	}

	if schema.Extensions == nil {
		schema.Extensions = make(map[string]interface{})
	}

	// Parse: "x-nullable,x-abc=def,!x-omitempty"
	extensions := strings.Split(extensionsTag, ",")
	for _, ext := range extensions {
		ext = strings.TrimSpace(ext)
		if ext == "" {
			continue
		}

		if strings.HasPrefix(ext, "!") {
			// Negation: !x-omitempty → x-omitempty: false
			key := strings.TrimPrefix(ext, "!")
			if strings.HasPrefix(key, "x-") {
				schema.Extensions[key] = false
			}
		} else if strings.Contains(ext, "=") {
			// With value: x-abc=def
			parts := strings.SplitN(ext, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if strings.HasPrefix(key, "x-") {
				// Try to parse as number, otherwise string
				if numValue, err := strconv.ParseFloat(value, 64); err == nil {
					schema.Extensions[key] = numValue
				} else if value == "true" {
					schema.Extensions[key] = true
				} else if value == "false" {
					schema.Extensions[key] = false
				} else {
					schema.Extensions[key] = value
				}
			}
		} else {
			// Boolean true: x-nullable
			if strings.HasPrefix(ext, "x-") {
				schema.Extensions[ext] = true
			}
		}
	}
}

// extractTag extracts a specific tag value from a struct tag string.
func extractTag(tagStr, key string) string {
	// Find the tag using regex to handle quoted values
	pattern := key + `:"([^"]*?)"`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(tagStr)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// parseOverrideType parses a type override string into a schema.
// Supports: "string", "number", "integer", "boolean", "time.Time", etc.
func (s *SchemaProcessor) parseOverrideType(override string) *openapi.Schema {
	schema := &openapi.Schema{}

	switch override {
	case "string":
		schema.Type = "string"
	case "number":
		schema.Type = "number"
	case "integer":
		schema.Type = "integer"
	case "boolean":
		schema.Type = "boolean"
	case "time.Time":
		schema.Type = "string"
		schema.Format = "date-time"
	default:
		// If not a primitive, treat as a reference
		schema.Ref = "#/components/schemas/" + override
	}

	return schema
}
