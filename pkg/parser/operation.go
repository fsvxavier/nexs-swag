package parser

import (
	"go/ast"
	"regexp"
	"strconv"
	"strings"

	v3 "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

// OperationProcessor processes operation annotations from function comments.
type OperationProcessor struct {
	parser    *Parser
	openapi   *v3.OpenAPI
	typeCache map[string]*TypeInfo
}

// RouteInfo contains routing information for an operation.
type RouteInfo struct {
	Path   string
	Method string
}

// NewOperationProcessor creates a new operation processor.
func NewOperationProcessor(p *Parser, spec *v3.OpenAPI, typeCache map[string]*TypeInfo) *OperationProcessor {
	return &OperationProcessor{
		parser:    p,
		openapi:   spec,
		typeCache: typeCache,
	}
}

var (
	// Operation-level annotations.
	summaryOpRegex     = regexp.MustCompile(`^@Summary\s+(.+)$`)
	descriptionOpRegex = regexp.MustCompile(`^@Description\s+(.+)$`)
	idRegex            = regexp.MustCompile(`^@ID\s+(.+)$`)
	tagsOpRegex        = regexp.MustCompile(`^@Tags\s+(.+)$`)
	acceptRegex        = regexp.MustCompile(`^@Accept\s+(.+)$`)
	produceRegex       = regexp.MustCompile(`^@Produce\s+(.+)$`)
	deprecatedOpRegex  = regexp.MustCompile(`^@Deprecated\s*$`)
	stateRegex         = regexp.MustCompile(`^@State\s+(.+)$`)

	// Router annotation.
	routerRegex = regexp.MustCompile(`^@Router\s+(\S+)\s+\[(\w+)\]`)

	// Parameter annotations.
	paramRegex = regexp.MustCompile(`^@Param\s+(\S+)\s+(\w+)\s+(\S+)\s+(true|false)\s+"([^"]*)"(?:\s+(.+))?`)

	// Response annotations.
	successRegex  = regexp.MustCompile(`^@Success\s+(\d+)\s+\{(\w+)\}\s+(\S+)(?:\s+"([^"]*)")?`)
	failureRegex  = regexp.MustCompile(`^@Failure\s+(\d+)\s+\{(\w+)\}\s+(\S+)(?:\s+"([^"]*)")?`)
	responseRegex = regexp.MustCompile(`^@Response\s+(\d+)\s+\{(\w+)\}\s+(\S+)(?:\s+"([^"]*)")?`)

	// Header annotation.
	headerRegex = regexp.MustCompile(`^@Header\s+(\d+)\s+\{(\w+)\}\s+(\S+)\s+"([^"]*)"`)

	// Security annotation - capture name before [ and scopes inside [].
	securityOpRegex = regexp.MustCompile(`^@Security\s+([^\[\s]+)(?:\[([^\]]+)\])?`)

	// Callback annotation.
	callbackRegex = regexp.MustCompile(`^@Callback\s+(\S+)\s+(\S+)\s+\[(\w+)\]`)

	// Extension annotations.
	xCodeSamplesRegex = regexp.MustCompile(`^@x-codeSamples\s+(.+)$`)
)

// Process processes function documentation and returns an Operation.
func (o *OperationProcessor) Process(doc *ast.CommentGroup) *v3.Operation {
	op := &v3.Operation{
		Responses: make(v3.Responses),
	}

	hasAnnotations := false

	for _, comment := range doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if text == "" || !strings.HasPrefix(text, "@") {
			continue
		}

		hasAnnotations = true

		switch {
		case summaryOpRegex.MatchString(text):
			o.processSummary(text, op)

		case descriptionOpRegex.MatchString(text):
			o.processDescription(text, op)

		case idRegex.MatchString(text):
			o.processID(text, op)

		case tagsOpRegex.MatchString(text):
			o.processTags(text, op)

		case acceptRegex.MatchString(text):
			o.processAccept(text, op)

		case produceRegex.MatchString(text):
			o.processProduce(text, op)

		case paramRegex.MatchString(text):
			o.processParameter(text, op)

		case successRegex.MatchString(text):
			o.processResponse(text, successRegex, op)

		case failureRegex.MatchString(text):
			o.processResponse(text, failureRegex, op)

		case responseRegex.MatchString(text):
			o.processResponse(text, responseRegex, op)

		case headerRegex.MatchString(text):
			o.processHeader(text, op)

		case securityOpRegex.MatchString(text):
			o.processSecurity(text, op)

		case deprecatedOpRegex.MatchString(text):
			op.Deprecated = true

		case stateRegex.MatchString(text):
			o.processState(text, op)

		case xCodeSamplesRegex.MatchString(text):
			o.processCodeSamples(text, op)
		}
	}

	if !hasAnnotations {
		return nil
	}

	return op
}

// GetRouteInfo extracts routing information from function documentation.
func (o *OperationProcessor) GetRouteInfo(doc *ast.CommentGroup) RouteInfo {
	for _, comment := range doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if routerRegex.MatchString(text) {
			matches := routerRegex.FindStringSubmatch(text)
			return RouteInfo{
				Path:   matches[1],
				Method: matches[2],
			}
		}
	}
	return RouteInfo{}
}

// processSummary processes @Summary annotation.
func (o *OperationProcessor) processSummary(text string, op *v3.Operation) {
	matches := summaryOpRegex.FindStringSubmatch(text)
	op.Summary = matches[1]
}

// processDescription processes @Description annotation.
// Supports markdown file substitution: @Description file(docs/endpoint.md).
func (o *OperationProcessor) processDescription(text string, op *v3.Operation) {
	matches := descriptionOpRegex.FindStringSubmatch(text)
	description := matches[1]

	// Check for file() substitution pattern
	if strings.HasPrefix(description, "file(") && strings.HasSuffix(description, ")") {
		filename := description[5 : len(description)-1]
		if content := o.parser.GetMarkdownContent(filename); content != "" {
			description = content
		}
	}

	if op.Description == "" {
		op.Description = description
	} else {
		op.Description += "\n" + description
	}
}

// processID processes @ID annotation.
func (o *OperationProcessor) processID(text string, op *v3.Operation) {
	matches := idRegex.FindStringSubmatch(text)
	op.OperationID = matches[1]
}

// processTags processes @Tags annotation.
func (o *OperationProcessor) processTags(text string, op *v3.Operation) {
	matches := tagsOpRegex.FindStringSubmatch(text)
	tags := strings.Split(matches[1], ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	op.Tags = tags
}

// processAccept processes @Accept annotation (consumes).
// Supports multiple MIME types: @Accept json,xml,plain.
func (o *OperationProcessor) processAccept(text string, op *v3.Operation) {
	matches := acceptRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return
	}

	// Store accepted content types for later use in RequestBody
	contentTypes := o.parseMimeTypes(matches[1])

	// Store in operation metadata for use when creating RequestBody
	if op.RequestBody != nil && len(contentTypes) > 0 {
		// Update existing request body with multiple content types
		if len(op.RequestBody.Content) > 0 {
			// Get the existing schema
			var existingSchema *v3.Schema
			for _, mt := range op.RequestBody.Content {
				existingSchema = mt.Schema
				break
			}

			// Apply to all accepted content types
			newContent := make(map[string]*v3.MediaType)
			for _, ct := range contentTypes {
				newContent[ct] = &v3.MediaType{
					Schema: existingSchema,
				}
			}
			op.RequestBody.Content = newContent
		}
	}
}

// processProduce processes @Produce annotation (produces).
// Supports multiple MIME types: @Produce json,xml,plain.
func (o *OperationProcessor) processProduce(text string, op *v3.Operation) {
	matches := produceRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return
	}

	// Store produced content types for later use in Responses
	contentTypes := o.parseMimeTypes(matches[1])

	// Store in operation for use when creating responses
	// We'll use this when processing @Success/@Failure annotations
	if len(contentTypes) > 0 {
		// Update existing responses with multiple content types
		for _, response := range op.Responses {
			if len(response.Content) > 0 {
				// Get the existing schema
				var existingSchema *v3.Schema
				for _, mt := range response.Content {
					existingSchema = mt.Schema
					break
				}

				// Apply to all produced content types
				newContent := make(map[string]*v3.MediaType)
				for _, ct := range contentTypes {
					newContent[ct] = &v3.MediaType{
						Schema: existingSchema,
					}
				}
				response.Content = newContent
			}
		}
	}
}

// parseMimeTypes converts short MIME type names to full MIME types.
// Supports: json, xml, plain, html, form, mpfd, etc.
func (o *OperationProcessor) parseMimeTypes(mimeList string) []string {
	var contentTypes []string

	// MIME type shortcuts
	mimeMap := map[string]string{
		"json":                  "application/json",
		"xml":                   "text/xml",
		"plain":                 "text/plain",
		"html":                  "text/html",
		"form":                  "application/x-www-form-urlencoded",
		"mpfd":                  "multipart/form-data",
		"multipart":             "multipart/form-data",
		"x-www-form-urlencoded": "application/x-www-form-urlencoded",
		"json-api":              "application/vnd.api+json",
		"json-stream":           "application/x-json-stream",
		"octet-stream":          "application/octet-stream",
		"png":                   "image/png",
		"jpeg":                  "image/jpeg",
		"jpg":                   "image/jpeg",
		"gif":                   "image/gif",
		"svg":                   "image/svg+xml",
		"pdf":                   "application/pdf",
		"zip":                   "application/zip",
		"csv":                   "text/csv",
	}

	types := strings.Split(mimeList, ",")
	for _, t := range types {
		t = strings.TrimSpace(t)

		// Check if it's a shortcut
		if fullType, ok := mimeMap[t]; ok {
			contentTypes = append(contentTypes, fullType)
		} else if strings.Contains(t, "/") {
			// Already a full MIME type
			contentTypes = append(contentTypes, t)
		} else {
			// Unknown shortcut, use application/* as fallback
			contentTypes = append(contentTypes, "application/"+t)
		}
	}

	return contentTypes
}

// processParameter processes @Param annotation.
func (o *OperationProcessor) processParameter(text string, op *v3.Operation) {
	matches := paramRegex.FindStringSubmatch(text)
	if len(matches) < 6 {
		return
	}

	name := matches[1]
	in := matches[2]
	schemaType := matches[3]
	required := matches[4] == valueTrue
	description := matches[5]

	// Parse additional attributes if present
	var attributes map[string]string
	if len(matches) > 6 && matches[6] != "" {
		attributes = o.parseAttributes(matches[6])
	}

	// Handle body parameter (request body in OpenAPI 3.x)
	if in == "body" {
		o.processRequestBody(schemaType, required, description, op)
		return
	}

	// Create parameter
	param := v3.Parameter{
		Name:        name,
		In:          in,
		Required:    required,
		Description: description,
		Schema:      o.parseSchemaType(schemaType),
	}

	// Note: collectionFormat is handled via Schema properties in OpenAPI 3.x
	// The parser.collectionFormat field can be used for default array serialization
	// but requires extending the Parameter struct to support Style/Explode fields

	// Apply additional attributes
	if attributes != nil {
		o.applyParameterAttributes(&param, attributes)
	}

	op.Parameters = append(op.Parameters, param)
}

// processRequestBody processes body parameter as RequestBody.
func (o *OperationProcessor) processRequestBody(schemaType string, required bool, description string, op *v3.Operation) {
	if op.RequestBody == nil {
		op.RequestBody = &v3.RequestBody{
			Required:    required,
			Description: description,
			Content:     make(map[string]*v3.MediaType),
		}
	}

	mediaType := &v3.MediaType{
		Schema: o.parseSchemaType(schemaType),
	}

	// Default to application/json
	op.RequestBody.Content["application/json"] = mediaType
}

// processResponse processes @Success, @Failure, and @Response annotations.
func (o *OperationProcessor) processResponse(text string, regex *regexp.Regexp, op *v3.Operation) {
	matches := regex.FindStringSubmatch(text)
	if len(matches) < 4 {
		return
	}

	statusCode := matches[1]
	responseType := matches[2]
	schemaRef := matches[3]
	description := "Success"
	if len(matches) > 4 && matches[4] != "" {
		description = matches[4]
	}

	response := &v3.Response{
		Description: description,
	}

	// Add content if schema is specified
	if responseType == typeObject || responseType == typeArray {
		content := make(map[string]*v3.MediaType)

		schema := o.parseSchemaType(schemaRef)
		if responseType == typeArray {
			schema = &v3.Schema{
				Type:  typeArray,
				Items: schema,
			}
		}

		content["application/json"] = &v3.MediaType{
			Schema: schema,
		}
		response.Content = content
	}

	op.Responses[statusCode] = response
}

// processHeader processes @Header annotation.
func (o *OperationProcessor) processHeader(text string, op *v3.Operation) {
	matches := headerRegex.FindStringSubmatch(text)
	if len(matches) < 5 {
		return
	}

	statusCode := matches[1]
	headerType := matches[2]
	headerName := matches[3]
	description := matches[4]

	response := op.Responses[statusCode]
	if response == nil {
		response = &v3.Response{
			Description: "Response " + statusCode,
			Headers:     make(map[string]*v3.Header),
		}
		op.Responses[statusCode] = response
	}

	if response.Headers == nil {
		response.Headers = make(map[string]*v3.Header)
	}

	response.Headers[headerName] = &v3.Header{
		Description: description,
		Schema:      o.parseSchemaType(headerType),
	}
}

// processSecurity processes @Security annotation.
func (o *OperationProcessor) processSecurity(text string, op *v3.Operation) {
	matches := securityOpRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return
	}

	schemeName := matches[1]
	var scopes []string

	if len(matches) > 2 && matches[2] != "" {
		scopeList := strings.Split(matches[2], ",")
		for _, scope := range scopeList {
			scopes = append(scopes, strings.TrimSpace(scope))
		}
	}

	security := v3.SecurityRequirement{
		schemeName: scopes,
	}

	op.Security = append(op.Security, security)
}

// parseAttributes parses additional parameter attributes.
// Supports: minimum(10), maximum(100), minLength(1), maxLength(255), pattern(^[a-z]+$),
// enum(A,B,C), default(value), example(value), format(email), collectionFormat(multi).
func (o *OperationProcessor) parseAttributes(attrStr string) map[string]string {
	attrs := make(map[string]string)

	// Enhanced regex to support nested parentheses and complex values
	attrRegex := regexp.MustCompile(`(\w+)\(([^)]+)\)`)
	matches := attrRegex.FindAllStringSubmatch(attrStr, -1)

	for _, match := range matches {
		if len(match) == 3 {
			attrs[strings.ToLower(match[1])] = match[2]
		}
	}

	return attrs
}

// applyParameterAttributes applies parsed attributes to a parameter.
func (o *OperationProcessor) applyParameterAttributes(param *v3.Parameter, attrs map[string]string) {
	if param.Schema == nil {
		param.Schema = &v3.Schema{}
	}

	for key, value := range attrs {
		switch key {
		case "minimum", "min":
			if minVal, err := strconv.ParseFloat(value, 64); err == nil {
				param.Schema.Minimum = minVal
			}

		case "maximum", "max":
			if maxVal, err := strconv.ParseFloat(value, 64); err == nil {
				param.Schema.Maximum = maxVal
			}

		case "exclusiveminimum":
			// In JSON Schema 2020-12, exclusiveMinimum is a number
			if minVal, err := strconv.ParseFloat(value, 64); err == nil {
				param.Schema.ExclusiveMinimum = minVal
			}

		case "exclusivemaximum":
			// In JSON Schema 2020-12, exclusiveMaximum is a number
			if maxVal, err := strconv.ParseFloat(value, 64); err == nil {
				param.Schema.ExclusiveMaximum = maxVal
			}

		case "minlength":
			if minLen, err := strconv.Atoi(value); err == nil {
				param.Schema.MinLength = minLen
			}

		case "maxlength":
			if maxLen, err := strconv.Atoi(value); err == nil {
				param.Schema.MaxLength = maxLen
			}

		case "pattern":
			param.Schema.Pattern = value

		case "multipleof":
			if mult, err := strconv.ParseFloat(value, 64); err == nil {
				param.Schema.MultipleOf = mult
			}

		case "minitems":
			if minItems, err := strconv.Atoi(value); err == nil {
				param.Schema.MinItems = minItems
			}

		case "maxitems":
			if maxItems, err := strconv.Atoi(value); err == nil {
				param.Schema.MaxItems = maxItems
			}

		case "uniqueitems":
			param.Schema.UniqueItems = value == valueTrue

		case "enum", "enums":
			// Parse enum values, handling different types
			schemaType := o.getSchemaTypeString(param.Schema)
			enumValues := o.parseEnumValues(value, schemaType)
			param.Schema.Enum = enumValues

		case "default":
			// Parse default value based on schema type
			schemaType := o.getSchemaTypeString(param.Schema)
			param.Schema.Default = o.parseValue(value, schemaType)

		case "example":
			// Parse example value based on schema type
			schemaType := o.getSchemaTypeString(param.Schema)
			param.Example = o.parseValue(value, schemaType)

		case "format":
			param.Schema.Format = value

		case "collectionformat":
			// For array parameters, this affects serialization
			// OpenAPI 3.0 uses 'style' and 'explode' instead
			// We'll store it as a comment for now
			// TODO: Convert to proper style/explode based on format

		case "readonly":
			param.Schema.ReadOnly = value == valueTrue

		case "writeonly":
			param.Schema.WriteOnly = value == valueTrue

		case "nullable":
			param.Schema.Nullable = value == valueTrue

		case "deprecated":
			param.Deprecated = value == valueTrue

		case "allowemptyvalue":
			param.AllowEmptyValue = value == valueTrue
		}
	}
}

// getSchemaTypeString extracts the type as string from a Schema.
// Schema.Type can be string or []string in JSON Schema 2020-12.
func (o *OperationProcessor) getSchemaTypeString(schema *v3.Schema) string {
	if schema == nil || schema.Type == nil {
		return typeString
	}

	switch v := schema.Type.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return s
			}
		}
	case []string:
		if len(v) > 0 {
			return v[0]
		}
	}

	return typeString
}

// parseEnumValues parses enum values based on the schema type.
func (o *OperationProcessor) parseEnumValues(enumStr string, schemaType string) []interface{} {
	var enums []interface{}
	values := strings.Split(enumStr, ",")

	for _, v := range values {
		v = strings.TrimSpace(v)
		enums = append(enums, o.parseValue(v, schemaType))
	}

	return enums
}

// parseValue parses a string value to the appropriate type.
func (o *OperationProcessor) parseValue(value string, schemaType string) interface{} {
	value = strings.TrimSpace(value)

	switch schemaType {
	case typeInteger:
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}

	case typeNumber:
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}

	case typeBoolean:
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}

	case typeArray:
		// For arrays, split by comma
		parts := strings.Split(value, ",")
		var arr []interface{}
		for _, p := range parts {
			arr = append(arr, strings.TrimSpace(p))
		}
		return arr
	}

	// Default: return as string
	return value
}

// parseSchemaType converts a type string to an OpenAPI schema.
func (o *OperationProcessor) parseSchemaType(typeName string) *v3.Schema {
	schema := &v3.Schema{}

	// Handle array types
	if strings.HasPrefix(typeName, "[]") {
		schema.Type = typeArray
		itemType := strings.TrimPrefix(typeName, "[]")
		schema.Items = o.parseSchemaType(itemType)
		return schema
	}

	// Handle map types: map[string]Type
	if strings.HasPrefix(typeName, "map[") {
		schema.Type = typeObject
		// Extract value type from map[keyType]valueType
		re := regexp.MustCompile(`map\[[^\]]+\](.+)`)
		if matches := re.FindStringSubmatch(typeName); len(matches) > 1 {
			valueType := matches[1]
			schema.AdditionalProperties = o.parseSchemaType(valueType)
		} else {
			// If can't parse, use generic object
			schema.AdditionalProperties = true
		}
		return schema
	}

	// Handle primitive types
	switch typeName {
	case typeString:
		schema.Type = typeString
	case formatInt, formatInt8, formatInt16, formatInt32, formatInteger:
		schema.Type = typeInteger
		schema.Format = formatInt32
	case formatInt64:
		schema.Type = typeInteger
		schema.Format = formatInt64
	case formatUInt, formatUInt8, formatUInt16, formatUInt32:
		schema.Type = typeInteger
		schema.Format = formatInt32
	case formatUInt64:
		schema.Type = typeInteger
		schema.Format = formatInt64
	case formatFloat32:
		schema.Type = typeNumber
		schema.Format = formatFloat32
	case formatFloat64, formatFloat, typeNumber, formatDouble:
		schema.Type = typeNumber
		schema.Format = formatDouble
	case formatBoolean, formatBool:
		schema.Type = typeBoolean
	case formatByte:
		schema.Type = typeString
		schema.Format = formatByte
	case formatDate:
		schema.Type = typeString
		schema.Format = formatDate
	case formatDateTime, "time.Time":
		schema.Type = typeString
		schema.Format = formatDateTime
	case "file":
		schema.Type = typeString
		schema.Format = formatBinary
	default:
		// Assume it's a reference to a schema
		schema.Ref = "#/components/schemas/" + typeName
	}

	return schema
}

// processState processes @State annotation.
// @State adds x-state extension to operation.
func (o *OperationProcessor) processState(text string, op *v3.Operation) {
	matches := stateRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return
	}

	state := strings.TrimSpace(matches[1])
	if op.Extensions == nil {
		op.Extensions = make(map[string]interface{})
	}
	op.Extensions["x-state"] = state
}

// processCodeSamples processes @x-codeSamples annotation.
// Format: @x-codeSamples lang:filename
// Example: @x-codeSamples go:examples/create_user.go.
func (o *OperationProcessor) processCodeSamples(text string, op *v3.Operation) {
	matches := xCodeSamplesRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return
	}

	// Parse format: lang:filename
	parts := strings.SplitN(matches[1], ":", 2)
	if len(parts) != 2 {
		return
	}

	lang := strings.TrimSpace(parts[0])
	filename := strings.TrimSpace(parts[1])

	// Get code example content
	source := o.parser.GetCodeExample(filename)
	if source == "" {
		return
	}

	// Use provided language or detect from file
	if lang == "" {
		lang = detectLanguageFromExtension(filename)
	}

	// Initialize extensions and x-codeSamples
	if op.Extensions == nil {
		op.Extensions = make(map[string]interface{})
	}

	var samples []map[string]interface{}
	if existing, ok := op.Extensions["x-codeSamples"].([]map[string]interface{}); ok {
		samples = existing
	}

	samples = append(samples, map[string]interface{}{
		"lang":   lang,
		"source": source,
	})

	op.Extensions["x-codeSamples"] = samples
}

// TransToValidCollectionFormat validates and normalizes collection format.
// Valid formats: csv, multi, pipes, tsv, ssv.
func TransToValidCollectionFormat(format string) string {
	format = strings.ToLower(strings.TrimSpace(format))
	validFormats := map[string]bool{
		"csv":   true,
		"multi": true,
		"pipes": true,
		"tsv":   true,
		"ssv":   true,
	}
	if validFormats[format] {
		return format
	}
	return "csv" // Default
}
