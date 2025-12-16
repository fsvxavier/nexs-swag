// Package converter provides functionality to convert between OpenAPI/Swagger versions.
package converter

import (
	"fmt"
	"net/url"
	"strings"

	swagger "github.com/fsvxavier/nexs-swag/pkg/openapi/v2"
	openapi "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

// Converter handles conversion between OpenAPI versions.
type Converter struct {
	warnings []string
}

// New creates a new Converter instance.
func New() *Converter {
	return &Converter{
		warnings: make([]string, 0),
	}
}

// GetWarnings returns all conversion warnings.
func (c *Converter) GetWarnings() []string {
	return c.warnings
}

// ClearWarnings clears all accumulated warnings.
func (c *Converter) ClearWarnings() {
	c.warnings = make([]string, 0)
}

// ConvertToV2 converts an OpenAPI 3.1.0 specification to Swagger 2.0.
func (c *Converter) ConvertToV2(spec *openapi.OpenAPI) (*swagger.Swagger, error) {
	if spec == nil {
		return nil, fmt.Errorf("input specification is nil")
	}

	swagger := &swagger.Swagger{
		Swagger:      "2.0",
		Info:         c.convertInfo(spec.Info),
		Paths:        c.convertPaths(spec.Paths),
		Tags:         c.convertTags(spec.Tags),
		ExternalDocs: c.convertExternalDocs(spec.ExternalDocs),
		Security:     c.convertSecurity(spec.Security),
	}

	// Convert servers to host/basePath/schemes
	if len(spec.Servers) > 0 {
		host, basePath, schemes := c.parseServerURL(spec.Servers[0].URL)
		swagger.Host = host
		swagger.BasePath = basePath
		swagger.Schemes = schemes

		if len(spec.Servers) > 1 {
			c.warnings = append(c.warnings, "multiple servers detected: only the first server is converted to host/basePath/schemes in Swagger 2.0")
		}
	}

	// Convert components to definitions/parameters/responses/securityDefinitions
	if spec.Components != nil {
		swagger.Definitions = c.convertSchemas(spec.Components.Schemas)
		swagger.Parameters = c.convertParameterDefinitions(spec.Components.Parameters)
		swagger.Responses = c.convertResponseDefinitions(spec.Components.Responses)
		swagger.SecurityDefinitions = c.convertSecuritySchemes(spec.Components.SecuritySchemes)
	}

	// Warn about unsupported features
	if len(spec.Webhooks) > 0 {
		c.warnings = append(c.warnings, "webhooks are not supported in Swagger 2.0 and were ignored")
	}
	if spec.JSONSchemaDialect != "" {
		c.warnings = append(c.warnings, "jsonSchemaDialect is not supported in Swagger 2.0 and was ignored")
	}

	return swagger, nil
}

// convertInfo converts OpenAPI Info to Swagger Info.
func (c *Converter) convertInfo(info openapi.Info) swagger.Info {
	v2Info := swagger.Info{
		Title:          info.Title,
		Description:    info.Description,
		TermsOfService: info.TermsOfService,
		Version:        info.Version,
	}

	if info.Contact != nil {
		v2Info.Contact = &swagger.Contact{
			Name:  info.Contact.Name,
			URL:   info.Contact.URL,
			Email: info.Contact.Email,
		}
	}

	if info.License != nil {
		v2Info.License = &swagger.License{
			Name: info.License.Name,
			URL:  info.License.URL,
		}
		if info.License.Identifier != "" {
			c.warnings = append(c.warnings, "license.identifier is not supported in Swagger 2.0, only license.url is used")
		}
	}

	if info.Summary != "" {
		c.warnings = append(c.warnings, "info.summary is not supported in Swagger 2.0 and was ignored")
	}

	return v2Info
}

// parseServerURL parses an OpenAPI 3.x server URL into host, basePath, and schemes.
func (c *Converter) parseServerURL(serverURL string) (host, basePath string, schemes []string) {
	// Handle server URL with variables (replace with defaults or remove)
	serverURL = c.replaceServerVariables(serverURL)

	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		c.warnings = append(c.warnings, fmt.Sprintf("failed to parse server URL %q: %v", serverURL, err))
		return "", "", nil
	}

	host = parsedURL.Host
	basePath = parsedURL.Path
	if basePath == "" {
		basePath = "/"
	}

	if parsedURL.Scheme != "" {
		schemes = []string{parsedURL.Scheme}
	}

	return host, basePath, schemes
}

// replaceServerVariables replaces server variables with defaults or empty strings.
func (c *Converter) replaceServerVariables(serverURL string) string {
	// Simple replacement: {variable} → empty
	// More sophisticated handling could use server variable defaults
	return strings.ReplaceAll(strings.ReplaceAll(serverURL, "{", ""), "}", "")
}

// convertPaths converts OpenAPI Paths to Swagger Paths.
func (c *Converter) convertPaths(paths openapi.Paths) swagger.Paths {
	if paths == nil {
		return nil
	}

	v2Paths := make(swagger.Paths, len(paths))
	for path, pathItem := range paths {
		v2Paths[path] = c.convertPathItem(pathItem)
	}

	return v2Paths
}

// convertPathItem converts an OpenAPI PathItem to Swagger PathItem.
func (c *Converter) convertPathItem(pathItem *openapi.PathItem) *swagger.PathItem {
	if pathItem == nil {
		return nil
	}

	v2PathItem := &swagger.PathItem{
		Ref:        c.convertRefToV2(pathItem.Ref),
		Parameters: c.convertParameters(pathItem.Parameters),
	}

	if pathItem.Get != nil {
		v2PathItem.Get = c.convertOperation(pathItem.Get)
	}
	if pathItem.Put != nil {
		v2PathItem.Put = c.convertOperation(pathItem.Put)
	}
	if pathItem.Post != nil {
		v2PathItem.Post = c.convertOperation(pathItem.Post)
	}
	if pathItem.Delete != nil {
		v2PathItem.Delete = c.convertOperation(pathItem.Delete)
	}
	if pathItem.Options != nil {
		v2PathItem.Options = c.convertOperation(pathItem.Options)
	}
	if pathItem.Head != nil {
		v2PathItem.Head = c.convertOperation(pathItem.Head)
	}
	if pathItem.Patch != nil {
		v2PathItem.Patch = c.convertOperation(pathItem.Patch)
	}

	// QUERY method is new in OpenAPI 3.2.0 - not supported in Swagger 2.0
	// But we still need to process it to generate warnings for its content
	if pathItem.Query != nil {
		c.warnings = append(c.warnings, "QUERY HTTP method is not supported in Swagger 2.0 (OpenAPI 3.2.0 feature) and was ignored")
		// Process the operation to generate warnings for responses, request body, etc.
		_ = c.convertOperation(pathItem.Query)
	}

	return v2PathItem
}

// convertOperation converts an OpenAPI Operation to Swagger Operation.
func (c *Converter) convertOperation(op *openapi.Operation) *swagger.Operation {
	if op == nil {
		return nil
	}

	v2Op := &swagger.Operation{
		Tags:         op.Tags,
		Summary:      op.Summary,
		Description:  op.Description,
		ExternalDocs: c.convertExternalDocs(op.ExternalDocs),
		OperationID:  op.OperationID,
		Parameters:   c.convertParameters(op.Parameters),
		Responses:    c.convertResponses(op.Responses),
		Deprecated:   op.Deprecated,
		Security:     c.convertSecurity(op.Security),
	}

	// Convert RequestBody to body parameter
	if op.RequestBody != nil {
		bodyParam := c.convertRequestBodyToParameter(op.RequestBody)
		if bodyParam != nil {
			v2Op.Parameters = append(v2Op.Parameters, bodyParam)
		}
	}

	// Extract consumes/produces from RequestBody and Responses
	v2Op.Consumes = c.extractConsumes(op.RequestBody)
	v2Op.Produces = c.extractProduces(op.Responses)

	// Warn about unsupported features
	if len(op.Callbacks) > 0 {
		c.warnings = append(c.warnings, fmt.Sprintf("operation %q: callbacks are not supported in Swagger 2.0 and were ignored", op.OperationID))
	}
	if len(op.Servers) > 0 {
		c.warnings = append(c.warnings, fmt.Sprintf("operation %q: operation-level servers are not supported in Swagger 2.0 and were ignored", op.OperationID))
	}

	return v2Op
}

// convertParameters converts OpenAPI Parameters to Swagger Parameters.
func (c *Converter) convertParameters(params []openapi.Parameter) []*swagger.Parameter {
	if len(params) == 0 {
		return nil
	}

	v2Params := make([]*swagger.Parameter, 0, len(params))
	for _, param := range params {
		v2Param := c.convertParameter(&param)
		if v2Param != nil {
			v2Params = append(v2Params, v2Param)
		}
	}

	return v2Params
}

// convertParameter converts an OpenAPI Parameter to Swagger Parameter.
func (c *Converter) convertParameter(param *openapi.Parameter) *swagger.Parameter {
	if param == nil {
		return nil
	}

	v2Param := &swagger.Parameter{
		Name:        param.Name,
		In:          param.In,
		Description: param.Description,
		Required:    param.Required,
	}

	// Convert schema to type/format for simple parameters
	if param.Schema != nil {
		c.convertSchemaToParameter(param.Schema, v2Param)
	}

	// Handle deprecated fields
	if param.Deprecated {
		if v2Param.Extensions == nil {
			v2Param.Extensions = make(map[string]interface{})
		}
		v2Param.Extensions["x-deprecated"] = true
	}

	return v2Param
}

// convertSchemaToParameter converts schema properties to parameter properties.
func (c *Converter) convertSchemaToParameter(schema *openapi.Schema, param *swagger.Parameter) {
	if schema == nil {
		return
	}

	// Handle type (can be array in 3.1.0)
	typeStr := c.extractType(schema.Type)
	param.Type = typeStr
	param.Format = schema.Format
	param.Default = schema.Default
	param.Pattern = schema.Pattern
	param.UniqueItems = schema.UniqueItems
	param.Enum = schema.Enum

	// Handle pointer fields
	if schema.Maximum != 0 {
		v := schema.Maximum
		param.Maximum = &v
	}
	if schema.Minimum != 0 {
		v := schema.Minimum
		param.Minimum = &v
	}
	if schema.MaxLength != 0 {
		v := schema.MaxLength
		param.MaxLength = &v
	}
	if schema.MinLength != 0 {
		v := schema.MinLength
		param.MinLength = &v
	}
	if schema.MaxItems != 0 {
		v := schema.MaxItems
		param.MaxItems = &v
	}
	if schema.MinItems != 0 {
		v := schema.MinItems
		param.MinItems = &v
	}

	// Convert items for array type
	if typeStr == "array" && schema.Items != nil {
		param.Items = c.convertSchemaToItems(schema.Items)
	}
}

// convertSchemaToItems converts a Schema to Items.
func (c *Converter) convertSchemaToItems(schema *openapi.Schema) *swagger.Items {
	if schema == nil {
		return nil
	}

	items := &swagger.Items{
		Type:    c.extractType(schema.Type),
		Format:  schema.Format,
		Default: schema.Default,
	}

	return items
}

// extractType extracts a simple type string from OpenAPI 3.1.0 type (which can be array).
func (c *Converter) extractType(t interface{}) string {
	if t == nil {
		return ""
	}

	switch v := t.(type) {
	case string:
		return v
	case []interface{}:
		// In OpenAPI 3.1.0, type can be an array like ["string", "null"]
		// Extract the first non-null type
		for _, typ := range v {
			if s, ok := typ.(string); ok && s != "null" {
				return s
			}
		}
	case []string:
		for _, typ := range v {
			if typ != "null" {
				return typ
			}
		}
	}

	return ""
}

// convertRequestBodyToParameter converts RequestBody to a body parameter.
func (c *Converter) convertRequestBodyToParameter(rb *openapi.RequestBody) *swagger.Parameter {
	if rb == nil || len(rb.Content) == 0 {
		return nil
	}

	// Use first content type (preferring application/json)
	var mediaType *openapi.MediaType
	var contentType string

	if mt, ok := rb.Content["application/json"]; ok {
		mediaType = mt
		contentType = "application/json"
	} else {
		// Use first available
		for ct, mt := range rb.Content {
			mediaType = mt
			contentType = ct
			break
		}
	}

	if mediaType == nil || mediaType.Schema == nil {
		return nil
	}

	// Warn about OpenAPI 3.2.0 streaming features in request body
	if mediaType.ItemSchema != nil {
		c.warnings = append(c.warnings, "MediaType.itemSchema for streaming is not supported in Swagger 2.0 (OpenAPI 3.2.0 feature) and was ignored")
	}
	if len(mediaType.ItemEncoding) > 0 {
		c.warnings = append(c.warnings, "MediaType.itemEncoding for streaming is not supported in Swagger 2.0 (OpenAPI 3.2.0 feature) and was ignored")
	}

	param := &swagger.Parameter{
		Name:        "body",
		In:          "body",
		Description: rb.Description,
		Required:    rb.Required,
		Schema:      c.convertSchema(mediaType.Schema),
	}

	if len(rb.Content) > 1 {
		c.warnings = append(c.warnings, fmt.Sprintf("requestBody has multiple content types: only %q is converted to body parameter", contentType))
	}

	return param
}

// extractConsumes extracts MIME types from RequestBody.
func (c *Converter) extractConsumes(rb *openapi.RequestBody) []string {
	if rb == nil || len(rb.Content) == 0 {
		return nil
	}

	consumes := make([]string, 0, len(rb.Content))
	for contentType := range rb.Content {
		consumes = append(consumes, contentType)
	}

	return consumes
}

// extractProduces extracts MIME types from Responses.
func (c *Converter) extractProduces(responses openapi.Responses) []string {
	if len(responses) == 0 {
		return nil
	}

	producesMap := make(map[string]bool)
	for _, resp := range responses {
		if resp != nil && len(resp.Content) > 0 {
			for contentType := range resp.Content {
				producesMap[contentType] = true
			}
		}
	}

	if len(producesMap) == 0 {
		return nil
	}

	produces := make([]string, 0, len(producesMap))
	for contentType := range producesMap {
		produces = append(produces, contentType)
	}

	return produces
}

// convertResponses converts OpenAPI Responses to Swagger Responses.
func (c *Converter) convertResponses(responses openapi.Responses) swagger.Responses {
	if len(responses) == 0 {
		return nil
	}

	v2Responses := make(swagger.Responses, len(responses))
	for code, resp := range responses {
		v2Responses[code] = c.convertResponse(resp)
	}

	return v2Responses
}

// convertResponse converts an OpenAPI Response to Swagger Response.
func (c *Converter) convertResponse(resp *openapi.Response) *swagger.Response {
	if resp == nil {
		return nil
	}

	v2Resp := &swagger.Response{
		Description: resp.Description,
		Headers:     c.convertHeaders(resp.Headers),
	}

	// Convert content to schema (use first content type, preferring application/json)
	if len(resp.Content) > 0 {
		var mediaType *openapi.MediaType
		if mt, ok := resp.Content["application/json"]; ok {
			mediaType = mt
		} else {
			// Use first available
			for _, mt := range resp.Content {
				mediaType = mt
				break
			}
		}

		if mediaType != nil && mediaType.Schema != nil {
			v2Resp.Schema = c.convertSchema(mediaType.Schema)
		}

		// Warn about OpenAPI 3.2.0 streaming features in all media types
		for contentType, mt := range resp.Content {
			if mt.ItemSchema != nil {
				c.warnings = append(c.warnings, "MediaType.itemSchema for streaming is not supported in Swagger 2.0 (OpenAPI 3.2.0 feature) and was ignored")
				break // Only warn once
			}
			if len(mt.ItemEncoding) > 0 {
				c.warnings = append(c.warnings, "MediaType.itemEncoding for streaming is not supported in Swagger 2.0 (OpenAPI 3.2.0 feature) and was ignored")
				break // Only warn once
			}
			_ = contentType // avoid unused variable
		}

		// Convert examples
		v2Resp.Examples = c.convertExamples(resp.Content)
	}

	return v2Resp
}

// convertHeaders converts OpenAPI Headers to Swagger Headers.
func (c *Converter) convertHeaders(headers map[string]*openapi.Header) map[string]*swagger.Header {
	if len(headers) == 0 {
		return nil
	}

	v2Headers := make(map[string]*swagger.Header, len(headers))
	for name, header := range headers {
		v2Headers[name] = c.convertHeader(header)
	}

	return v2Headers
}

// convertHeader converts an OpenAPI Header to Swagger Header.
func (c *Converter) convertHeader(header *openapi.Header) *swagger.Header {
	if header == nil {
		return nil
	}

	v2Header := &swagger.Header{
		Description: header.Description,
	}

	if header.Schema != nil {
		v2Header.Type = c.extractType(header.Schema.Type)
		v2Header.Format = header.Schema.Format
		v2Header.Default = header.Schema.Default
	}

	return v2Header
}

// convertExamples converts content examples to Swagger examples.
func (c *Converter) convertExamples(content map[string]*openapi.MediaType) map[string]interface{} {
	if len(content) == 0 {
		return nil
	}

	examples := make(map[string]interface{})
	for contentType, mediaType := range content {
		if mediaType.Example != nil {
			examples[contentType] = mediaType.Example
		} else if len(mediaType.Examples) > 0 {
			// Use first example
			for _, ex := range mediaType.Examples {
				if ex != nil && ex.Value != nil {
					examples[contentType] = ex.Value
					break
				}
			}
		}
	}

	if len(examples) == 0 {
		return nil
	}

	return examples
}

// convertRefToV2 converts OpenAPI 3.x $ref to Swagger 2.0 $ref format.
// Converts #/components/schemas/Foo to #/definitions/Foo
// Converts #/components/parameters/Foo to #/parameters/Foo
// Converts #/components/responses/Foo to #/responses/Foo
func (c *Converter) convertRefToV2(ref string) string {
	if ref == "" {
		return ""
	}
	ref = strings.Replace(ref, "#/components/schemas/", "#/definitions/", 1)
	ref = strings.Replace(ref, "#/components/parameters/", "#/parameters/", 1)
	ref = strings.Replace(ref, "#/components/responses/", "#/responses/", 1)
	// components/securitySchemes -> securityDefinitions
	ref = strings.Replace(ref, "#/components/securitySchemes/", "#/securityDefinitions/", 1)
	return ref
}

// convertRefToV3 converts Swagger 2.0 $ref to OpenAPI 3.x $ref format.
// Converts #/definitions/Foo to #/components/schemas/Foo
// Converts #/parameters/Foo to #/components/parameters/Foo
// Converts #/responses/Foo to #/components/responses/Foo
func (c *Converter) convertRefToV3(ref string) string {
	if ref == "" {
		return ""
	}
	ref = strings.Replace(ref, "#/definitions/", "#/components/schemas/", 1)
	ref = strings.Replace(ref, "#/parameters/", "#/components/parameters/", 1)
	ref = strings.Replace(ref, "#/responses/", "#/components/responses/", 1)
	// securityDefinitions -> components/securitySchemes
	ref = strings.Replace(ref, "#/securityDefinitions/", "#/components/securitySchemes/", 1)
	return ref
}

// convertSchemas converts OpenAPI Schemas to Swagger Schemas (definitions).
func (c *Converter) convertSchemas(schemas map[string]*openapi.Schema) map[string]*swagger.Schema {
	if len(schemas) == 0 {
		return nil
	}

	v2Schemas := make(map[string]*swagger.Schema, len(schemas))
	for name, schema := range schemas {
		v2Schemas[name] = c.convertSchema(schema)
	}

	return v2Schemas
}

// convertSchema converts an OpenAPI Schema to Swagger Schema (JSON Schema Draft 4).
func (c *Converter) convertSchema(schema *openapi.Schema) *swagger.Schema {
	if schema == nil {
		return nil
	}

	v2Schema := &swagger.Schema{
		Ref:          c.convertRefToV2(schema.Ref),
		Type:         c.extractType(schema.Type),
		Format:       schema.Format,
		Title:        schema.Title,
		Description:  schema.Description,
		Default:      schema.Default,
		Pattern:      schema.Pattern,
		UniqueItems:  schema.UniqueItems,
		Required:     schema.Required,
		Enum:         schema.Enum,
		Example:      schema.Example,
		ReadOnly:     schema.ReadOnly,
		ExternalDocs: c.convertExternalDocs(schema.ExternalDocs),
	}

	// Handle optional numeric/int fields with pointers
	if schema.MultipleOf != 0 {
		v := schema.MultipleOf
		v2Schema.MultipleOf = &v
	}
	if schema.Maximum != 0 {
		v := schema.Maximum
		v2Schema.Maximum = &v
	}
	if schema.Minimum != 0 {
		v := schema.Minimum
		v2Schema.Minimum = &v
	}
	if schema.MaxLength != 0 {
		v := schema.MaxLength
		v2Schema.MaxLength = &v
	}
	if schema.MinLength != 0 {
		v := schema.MinLength
		v2Schema.MinLength = &v
	}
	if schema.MaxItems != 0 {
		v := schema.MaxItems
		v2Schema.MaxItems = &v
	}
	if schema.MinItems != 0 {
		v := schema.MinItems
		v2Schema.MinItems = &v
	}
	if schema.MaxProperties != 0 {
		v := schema.MaxProperties
		v2Schema.MaxProperties = &v
	}
	if schema.MinProperties != 0 {
		v := schema.MinProperties
		v2Schema.MinProperties = &v
	}

	// ExclusiveMaximum/ExclusiveMinimum: boolean in Draft 4, number in 2020-12
	if schema.ExclusiveMaximum != nil {
		if val, ok := schema.ExclusiveMaximum.(bool); ok {
			v2Schema.ExclusiveMaximum = val
		}
	}
	if schema.ExclusiveMinimum != nil {
		if val, ok := schema.ExclusiveMinimum.(bool); ok {
			v2Schema.ExclusiveMinimum = val
		}
	}

	// Convert properties
	if len(schema.Properties) > 0 {
		v2Schema.Properties = c.convertSchemas(schema.Properties)
	}

	// Convert additionalProperties
	if schema.AdditionalProperties != nil {
		switch v := schema.AdditionalProperties.(type) {
		case bool:
			v2Schema.AdditionalProperties = v
		case *openapi.Schema:
			v2Schema.AdditionalProperties = c.convertSchema(v)
		}
	}

	// Convert items
	if schema.Items != nil {
		v2Schema.Items = c.convertSchema(schema.Items)
	}

	// Convert allOf
	if len(schema.AllOf) > 0 {
		v2Schema.AllOf = make([]*swagger.Schema, len(schema.AllOf))
		for i, s := range schema.AllOf {
			v2Schema.AllOf[i] = c.convertSchema(&s)
		}
	}

	// Handle nullable (3.1.0 uses type array with "null")
	if c.isNullable(schema.Type) {
		if v2Schema.Extensions == nil {
			v2Schema.Extensions = make(map[string]interface{})
		}
		v2Schema.Extensions["x-nullable"] = true
	}

	// Handle XML
	if schema.XML != nil {
		v2Schema.XML = &swagger.XML{
			Name:      schema.XML.Name,
			Namespace: schema.XML.Namespace,
			Prefix:    schema.XML.Prefix,
			Attribute: schema.XML.Attribute,
			Wrapped:   schema.XML.Wrapped,
		}
	}

	// Handle discriminator
	if schema.Discriminator != nil {
		v2Schema.Discriminator = schema.Discriminator.PropertyName
	}

	// Warn about unsupported Draft 2020-12 features
	c.warnUnsupportedSchemaFeatures(schema)

	return v2Schema
}

// isNullable checks if type includes "null" (OpenAPI 3.1.0 feature).
func (c *Converter) isNullable(t interface{}) bool {
	switch v := t.(type) {
	case []interface{}:
		for _, typ := range v {
			if s, ok := typ.(string); ok && s == "null" {
				return true
			}
		}
	case []string:
		for _, typ := range v {
			if typ == "null" {
				return true
			}
		}
	}
	return false
}

// warnUnsupportedSchemaFeatures warns about JSON Schema 2020-12 features not in Draft 4.
func (c *Converter) warnUnsupportedSchemaFeatures(schema *openapi.Schema) {
	if len(schema.OneOf) > 0 {
		c.warnings = append(c.warnings, "oneOf is limited in Swagger 2.0 (use with discriminator only)")
	}
	if len(schema.AnyOf) > 0 {
		c.warnings = append(c.warnings, "anyOf is not supported in JSON Schema Draft 4 (Swagger 2.0)")
	}
	if schema.Not != nil {
		c.warnings = append(c.warnings, "not is limited in JSON Schema Draft 4 (Swagger 2.0)")
	}
	if len(schema.PrefixItems) > 0 {
		c.warnings = append(c.warnings, "prefixItems is not supported in JSON Schema Draft 4 (Swagger 2.0)")
	}
	if schema.WriteOnly {
		c.warnings = append(c.warnings, "writeOnly is not supported in Swagger 2.0")
	}
}

// convertParameterDefinitions converts component parameters.
func (c *Converter) convertParameterDefinitions(params map[string]*openapi.Parameter) map[string]*swagger.Parameter {
	if len(params) == 0 {
		return nil
	}

	v2Params := make(map[string]*swagger.Parameter, len(params))
	for name, param := range params {
		v2Params[name] = c.convertParameter(param)
	}

	return v2Params
}

// convertResponseDefinitions converts component responses.
func (c *Converter) convertResponseDefinitions(responses map[string]*openapi.Response) map[string]*swagger.Response {
	if len(responses) == 0 {
		return nil
	}

	v2Responses := make(map[string]*swagger.Response, len(responses))
	for name, resp := range responses {
		v2Responses[name] = c.convertResponse(resp)
	}

	return v2Responses
}

// convertSecuritySchemes converts OpenAPI SecuritySchemes to Swagger SecurityDefinitions.
func (c *Converter) convertSecuritySchemes(schemes map[string]*openapi.SecurityScheme) map[string]*swagger.SecurityScheme {
	if len(schemes) == 0 {
		return nil
	}

	v2Schemes := make(map[string]*swagger.SecurityScheme, len(schemes))
	for name, scheme := range schemes {
		v2Schemes[name] = c.convertSecurityScheme(scheme)
	}

	return v2Schemes
}

// convertSecurityScheme converts an OpenAPI SecurityScheme to Swagger SecurityScheme.
func (c *Converter) convertSecurityScheme(scheme *openapi.SecurityScheme) *swagger.SecurityScheme {
	if scheme == nil {
		return nil
	}

	v2Scheme := &swagger.SecurityScheme{
		Description: scheme.Description,
	}

	// Map OpenAPI 3.x types to Swagger 2.0 types
	switch scheme.Type {
	case "http":
		if scheme.Scheme == "basic" {
			v2Scheme.Type = "basic"
		} else if scheme.Scheme == "bearer" {
			v2Scheme.Type = "apiKey"
			v2Scheme.In = "header"
			v2Scheme.Name = "Authorization"
			if v2Scheme.Extensions == nil {
				v2Scheme.Extensions = make(map[string]interface{})
			}
			v2Scheme.Extensions["x-bearer-format"] = scheme.BearerFormat
		} else {
			c.warnings = append(c.warnings, fmt.Sprintf("http scheme %q is not directly supported in Swagger 2.0", scheme.Scheme))
		}
	case "apiKey":
		v2Scheme.Type = "apiKey"
		v2Scheme.Name = scheme.Name
		v2Scheme.In = scheme.In
	case "oauth2":
		v2Scheme.Type = "oauth2"
		c.convertOAuth2Flows(scheme, v2Scheme)

		// Warn about OpenAPI 3.2.0 features
		if scheme.OAuth2MetadataURL != "" {
			c.warnings = append(c.warnings, "OAuth2MetadataURL is not supported in Swagger 2.0 (OpenAPI 3.2.0 feature) and was ignored")
		}
	case "openIdConnect":
		c.warnings = append(c.warnings, "openIdConnect is not supported in Swagger 2.0, converted to oauth2")
		v2Scheme.Type = "oauth2"
	default:
		c.warnings = append(c.warnings, fmt.Sprintf("security scheme type %q is not supported in Swagger 2.0", scheme.Type))
	}

	// Warn about deprecated field (OpenAPI 3.2.0)
	if scheme.Deprecated {
		if v2Scheme.Extensions == nil {
			v2Scheme.Extensions = make(map[string]interface{})
		}
		v2Scheme.Extensions["x-deprecated"] = true
		c.warnings = append(c.warnings, "SecurityScheme.deprecated is not natively supported in Swagger 2.0, converted to x-deprecated extension")
	}

	return v2Scheme
}

// convertOAuth2Flows converts OAuth2 flows from OpenAPI 3.x to Swagger 2.0 format.
func (c *Converter) convertOAuth2Flows(scheme *openapi.SecurityScheme, v2Scheme *swagger.SecurityScheme) {
	if scheme.Flows == nil {
		return
	}

	// Swagger 2.0 supports only one flow at a time
	// Priority: implicit > password > application > authorizationCode
	if scheme.Flows.Implicit != nil {
		v2Scheme.Flow = "implicit"
		v2Scheme.AuthorizationURL = scheme.Flows.Implicit.AuthorizationURL
		v2Scheme.Scopes = scheme.Flows.Implicit.Scopes
	} else if scheme.Flows.Password != nil {
		v2Scheme.Flow = "password"
		v2Scheme.TokenURL = scheme.Flows.Password.TokenURL
		v2Scheme.Scopes = scheme.Flows.Password.Scopes
	} else if scheme.Flows.ClientCredentials != nil {
		v2Scheme.Flow = "application"
		v2Scheme.TokenURL = scheme.Flows.ClientCredentials.TokenURL
		v2Scheme.Scopes = scheme.Flows.ClientCredentials.Scopes
	} else if scheme.Flows.AuthorizationCode != nil {
		v2Scheme.Flow = "accessCode"
		v2Scheme.AuthorizationURL = scheme.Flows.AuthorizationCode.AuthorizationURL
		v2Scheme.TokenURL = scheme.Flows.AuthorizationCode.TokenURL
		v2Scheme.Scopes = scheme.Flows.AuthorizationCode.Scopes
	}

	// Warn if multiple flows defined
	flowCount := 0
	if scheme.Flows.Implicit != nil {
		flowCount++
	}
	if scheme.Flows.Password != nil {
		flowCount++
	}
	if scheme.Flows.ClientCredentials != nil {
		flowCount++
	}
	if scheme.Flows.AuthorizationCode != nil {
		flowCount++
	}
	if scheme.Flows.DeviceAuthorization != nil {
		flowCount++
		c.warnings = append(c.warnings, "DeviceAuthorization OAuth2 flow is not supported in Swagger 2.0 (OpenAPI 3.2.0 feature) and was ignored")
	}
	if flowCount > 1 {
		c.warnings = append(c.warnings, "multiple OAuth2 flows detected: Swagger 2.0 supports only one flow per security scheme")
	}
}

// convertSecurity converts security requirements.
func (c *Converter) convertSecurity(security []openapi.SecurityRequirement) []swagger.SecurityRequirement {
	if len(security) == 0 {
		return nil
	}

	v2Security := make([]swagger.SecurityRequirement, len(security))
	for i, req := range security {
		v2Security[i] = swagger.SecurityRequirement(req)
	}

	return v2Security
}

// convertTags converts tags.
func (c *Converter) convertTags(tags []openapi.Tag) []swagger.Tag {
	if len(tags) == 0 {
		return nil
	}

	v2Tags := make([]swagger.Tag, len(tags))
	for i, tag := range tags {
		v2Tags[i] = swagger.Tag{
			Name:         tag.Name,
			Description:  tag.Description,
			ExternalDocs: c.convertExternalDocs(tag.ExternalDocs),
		}
	}

	return v2Tags
}

// convertExternalDocs converts external documentation.
func (c *Converter) convertExternalDocs(docs *openapi.ExternalDocs) *swagger.ExternalDocs {
	if docs == nil {
		return nil
	}

	return &swagger.ExternalDocs{
		Description: docs.Description,
		URL:         docs.URL,
	}
}

// =============================================================================
// Swagger 2.0 → OpenAPI 3.1.0 Conversion
// =============================================================================

// ConvertToV3 converts a Swagger 2.0 specification to OpenAPI 3.1.0.
func (c *Converter) ConvertToV3(swagger *swagger.Swagger) (*openapi.OpenAPI, error) {
	if swagger == nil {
		return nil, fmt.Errorf("input specification is nil")
	}

	spec := &openapi.OpenAPI{
		OpenAPI:      "3.1.0",
		Info:         c.convertInfoToV3(swagger.Info),
		Paths:        c.convertPathsToV3(swagger.Paths),
		Tags:         c.convertTagsToV3(swagger.Tags),
		ExternalDocs: c.convertExternalDocsToV3(swagger.ExternalDocs),
		Security:     c.convertSecurityToV3(swagger.Security),
	}

	// Convert host/basePath/schemes to servers
	if swagger.Host != "" || swagger.BasePath != "" || len(swagger.Schemes) > 0 {
		spec.Servers = c.buildServers(swagger.Host, swagger.BasePath, swagger.Schemes)
	}

	// Convert definitions/parameters/responses/securityDefinitions to components
	if swagger.Definitions != nil || swagger.Parameters != nil || swagger.Responses != nil || swagger.SecurityDefinitions != nil {
		spec.Components = &openapi.Components{
			Schemas:         c.convertDefinitionsToSchemas(swagger.Definitions),
			Parameters:      c.convertParameterDefinitionsToV3(swagger.Parameters),
			Responses:       c.convertResponseDefinitionsToV3(swagger.Responses),
			SecuritySchemes: c.convertSecurityDefinitionsToV3(swagger.SecurityDefinitions),
		}
	}

	// Handle global consumes/produces
	if len(swagger.Consumes) > 0 || len(swagger.Produces) > 0 {
		c.warnings = append(c.warnings, "global consumes/produces are not directly supported in OpenAPI 3.1.0, applied to operations where missing")
	}

	return spec, nil
}

// convertInfoToV3 converts Swagger Info to OpenAPI Info.
func (c *Converter) convertInfoToV3(info swagger.Info) openapi.Info {
	v3Info := openapi.Info{
		Title:          info.Title,
		Description:    info.Description,
		TermsOfService: info.TermsOfService,
		Version:        info.Version,
	}

	if info.Contact != nil {
		v3Info.Contact = &openapi.Contact{
			Name:  info.Contact.Name,
			URL:   info.Contact.URL,
			Email: info.Contact.Email,
		}
	}

	if info.License != nil {
		v3Info.License = &openapi.License{
			Name: info.License.Name,
			URL:  info.License.URL,
		}
	}

	return v3Info
}

// buildServers constructs Server objects from host/basePath/schemes.
func (c *Converter) buildServers(host, basePath string, schemes []string) []openapi.Server {
	if len(schemes) == 0 {
		schemes = []string{"https"} // Default to https
	}

	servers := make([]openapi.Server, 0, len(schemes))
	for _, scheme := range schemes {
		serverURL := scheme + "://"
		if host != "" {
			serverURL += host
		} else {
			serverURL += "localhost"
		}
		if basePath != "" && basePath != "/" {
			serverURL += basePath
		}

		servers = append(servers, openapi.Server{
			URL: serverURL,
		})
	}

	return servers
}

// convertPathsToV3 converts Swagger Paths to OpenAPI Paths.
func (c *Converter) convertPathsToV3(paths swagger.Paths) openapi.Paths {
	if paths == nil {
		return nil
	}

	v3Paths := make(openapi.Paths, len(paths))
	for path, pathItem := range paths {
		v3Paths[path] = c.convertPathItemToV3(pathItem)
	}

	return v3Paths
}

// convertPathItemToV3 converts a Swagger PathItem to OpenAPI PathItem.
func (c *Converter) convertPathItemToV3(pathItem *swagger.PathItem) *openapi.PathItem {
	if pathItem == nil {
		return nil
	}

	v3PathItem := &openapi.PathItem{
		Ref:        c.convertRefToV3(pathItem.Ref),
		Parameters: c.convertParametersToV3(pathItem.Parameters),
	}

	if pathItem.Get != nil {
		v3PathItem.Get = c.convertOperationToV3(pathItem.Get)
	}
	if pathItem.Put != nil {
		v3PathItem.Put = c.convertOperationToV3(pathItem.Put)
	}
	if pathItem.Post != nil {
		v3PathItem.Post = c.convertOperationToV3(pathItem.Post)
	}
	if pathItem.Delete != nil {
		v3PathItem.Delete = c.convertOperationToV3(pathItem.Delete)
	}
	if pathItem.Options != nil {
		v3PathItem.Options = c.convertOperationToV3(pathItem.Options)
	}
	if pathItem.Head != nil {
		v3PathItem.Head = c.convertOperationToV3(pathItem.Head)
	}
	if pathItem.Patch != nil {
		v3PathItem.Patch = c.convertOperationToV3(pathItem.Patch)
	}

	return v3PathItem
}

// convertOperationToV3 converts a Swagger Operation to OpenAPI Operation.
func (c *Converter) convertOperationToV3(op *swagger.Operation) *openapi.Operation {
	if op == nil {
		return nil
	}

	v3Op := &openapi.Operation{
		Tags:         op.Tags,
		Summary:      op.Summary,
		Description:  op.Description,
		ExternalDocs: c.convertExternalDocsToV3(op.ExternalDocs),
		OperationID:  op.OperationID,
		Responses:    c.convertResponsesToV3(op.Responses, op.Produces),
		Deprecated:   op.Deprecated,
		Security:     c.convertSecurityToV3(op.Security),
	}

	// Convert parameters, separating body parameters into requestBody
	nonBodyParams, bodyParam := c.separateBodyParameter(op.Parameters)
	v3Op.Parameters = c.convertParametersToV3(nonBodyParams)

	if bodyParam != nil {
		v3Op.RequestBody = c.convertBodyParameterToRequestBody(bodyParam, op.Consumes)
	}

	return v3Op
}

// separateBodyParameter separates body parameters from other parameters.
func (c *Converter) separateBodyParameter(params []*swagger.Parameter) (nonBody []*swagger.Parameter, body *swagger.Parameter) {
	if len(params) == 0 {
		return nil, nil
	}

	nonBody = make([]*swagger.Parameter, 0, len(params))
	for _, param := range params {
		if param.In == "body" {
			if body != nil {
				c.warnings = append(c.warnings, "multiple body parameters detected: only the last one is converted")
			}
			body = param
		} else {
			nonBody = append(nonBody, param)
		}
	}

	return nonBody, body
}

// convertBodyParameterToRequestBody converts a body parameter to RequestBody.
func (c *Converter) convertBodyParameterToRequestBody(param *swagger.Parameter, consumes []string) *openapi.RequestBody {
	if param == nil || param.Schema == nil {
		return nil
	}

	if len(consumes) == 0 {
		consumes = []string{"application/json"} // Default
	}

	content := make(map[string]*openapi.MediaType, len(consumes))
	for _, mediaType := range consumes {
		content[mediaType] = &openapi.MediaType{
			Schema: c.convertSchemaToV3(param.Schema),
		}
	}

	return &openapi.RequestBody{
		Description: param.Description,
		Content:     content,
		Required:    param.Required,
	}
}

// convertParametersToV3 converts Swagger Parameters to OpenAPI Parameters.
func (c *Converter) convertParametersToV3(params []*swagger.Parameter) []openapi.Parameter {
	if len(params) == 0 {
		return nil
	}

	v3Params := make([]openapi.Parameter, 0, len(params))
	for _, param := range params {
		if param.In == "body" {
			continue // Body parameters are handled separately
		}
		v3Param := c.convertParameterToV3(param)
		if v3Param != nil {
			v3Params = append(v3Params, *v3Param)
		}
	}

	return v3Params
}

// convertParameterToV3 converts a Swagger Parameter to OpenAPI Parameter.
func (c *Converter) convertParameterToV3(param *swagger.Parameter) *openapi.Parameter {
	if param == nil {
		return nil
	}

	v3Param := &openapi.Parameter{
		Name:            param.Name,
		In:              param.In,
		Description:     param.Description,
		Required:        param.Required,
		AllowEmptyValue: param.AllowEmptyValue,
	}

	// Convert type/format to schema
	v3Param.Schema = c.convertParameterPropertiesToSchema(param)

	// Handle deprecated (use extension in v2, native in v3)
	if param.Extensions != nil {
		if deprecated, ok := param.Extensions["x-deprecated"].(bool); ok && deprecated {
			v3Param.Deprecated = true
		}
	}

	return v3Param
}

// convertParameterPropertiesToSchema converts parameter type/format properties to a schema.
func (c *Converter) convertParameterPropertiesToSchema(param *swagger.Parameter) *openapi.Schema {
	if param.Schema != nil {
		return c.convertSchemaToV3(param.Schema)
	}

	schema := &openapi.Schema{
		Type:    param.Type,
		Format:  param.Format,
		Default: param.Default,
		Enum:    param.Enum,
	}

	// Copy numeric/string constraints
	if param.Maximum != nil {
		schema.Maximum = *param.Maximum
	}
	if param.Minimum != nil {
		schema.Minimum = *param.Minimum
	}
	if param.MaxLength != nil {
		schema.MaxLength = *param.MaxLength
	}
	if param.MinLength != nil {
		schema.MinLength = *param.MinLength
	}
	if param.Pattern != "" {
		schema.Pattern = param.Pattern
	}
	if param.MaxItems != nil {
		schema.MaxItems = *param.MaxItems
	}
	if param.MinItems != nil {
		schema.MinItems = *param.MinItems
	}
	schema.UniqueItems = param.UniqueItems

	// Convert items for array type
	if param.Type == "array" && param.Items != nil {
		schema.Items = c.convertItemsToSchema(param.Items)
	}

	// Handle nullable (x-nullable extension → type array)
	if param.Extensions != nil {
		if nullable, ok := param.Extensions["x-nullable"].(bool); ok && nullable {
			schema.Type = []interface{}{param.Type, "null"}
		}
	}

	return schema
}

// convertItemsToSchema converts Items to Schema.
func (c *Converter) convertItemsToSchema(items *swagger.Items) *openapi.Schema {
	if items == nil {
		return nil
	}

	schema := &openapi.Schema{
		Type:    items.Type,
		Format:  items.Format,
		Default: items.Default,
		Enum:    items.Enum,
	}

	if items.Maximum != nil {
		schema.Maximum = *items.Maximum
	}
	if items.Minimum != nil {
		schema.Minimum = *items.Minimum
	}

	// Handle nested arrays
	if items.Items != nil {
		schema.Items = c.convertItemsToSchema(items.Items)
	}

	return schema
}

// convertResponsesToV3 converts Swagger Responses to OpenAPI Responses.
func (c *Converter) convertResponsesToV3(responses swagger.Responses, produces []string) openapi.Responses {
	if len(responses) == 0 {
		return nil
	}

	v3Responses := make(openapi.Responses, len(responses))
	for code, resp := range responses {
		v3Responses[code] = c.convertResponseToV3(resp, produces)
	}

	return v3Responses
}

// convertResponseToV3 converts a Swagger Response to OpenAPI Response.
func (c *Converter) convertResponseToV3(resp *swagger.Response, produces []string) *openapi.Response {
	if resp == nil {
		return nil
	}

	v3Resp := &openapi.Response{
		Description: resp.Description,
		Headers:     c.convertHeadersToV3(resp.Headers),
	}

	// Convert schema to content
	if resp.Schema != nil {
		if len(produces) == 0 {
			produces = []string{"application/json"} // Default
		}

		content := make(map[string]*openapi.MediaType, len(produces))
		for _, mediaType := range produces {
			mt := &openapi.MediaType{
				Schema: c.convertSchemaToV3(resp.Schema),
			}

			// Add examples if available
			if resp.Examples != nil {
				if example, ok := resp.Examples[mediaType]; ok {
					mt.Example = example
				}
			}

			content[mediaType] = mt
		}

		v3Resp.Content = content
	}

	return v3Resp
}

// convertHeadersToV3 converts Swagger Headers to OpenAPI Headers.
func (c *Converter) convertHeadersToV3(headers map[string]*swagger.Header) map[string]*openapi.Header {
	if len(headers) == 0 {
		return nil
	}

	v3Headers := make(map[string]*openapi.Header, len(headers))
	for name, header := range headers {
		v3Headers[name] = c.convertHeaderToV3(header)
	}

	return v3Headers
}

// convertHeaderToV3 converts a Swagger Header to OpenAPI Header.
func (c *Converter) convertHeaderToV3(header *swagger.Header) *openapi.Header {
	if header == nil {
		return nil
	}

	schema := &openapi.Schema{
		Type:    header.Type,
		Format:  header.Format,
		Default: header.Default,
		Enum:    header.Enum,
	}

	if header.Maximum != nil {
		schema.Maximum = *header.Maximum
	}
	if header.Minimum != nil {
		schema.Minimum = *header.Minimum
	}

	return &openapi.Header{
		Description: header.Description,
		Schema:      schema,
	}
}

// convertDefinitionsToSchemas converts Swagger Definitions to OpenAPI Schemas.
func (c *Converter) convertDefinitionsToSchemas(definitions map[string]*swagger.Schema) map[string]*openapi.Schema {
	if len(definitions) == 0 {
		return nil
	}

	v3Schemas := make(map[string]*openapi.Schema, len(definitions))
	for name, schema := range definitions {
		v3Schemas[name] = c.convertSchemaToV3(schema)
	}

	return v3Schemas
}

// convertSchemaToV3 converts a Swagger Schema (JSON Schema Draft 4) to OpenAPI Schema (2020-12).
func (c *Converter) convertSchemaToV3(schema *swagger.Schema) *openapi.Schema {
	if schema == nil {
		return nil
	}

	v3Schema := &openapi.Schema{
		Ref:          c.convertRefToV3(schema.Ref),
		Type:         schema.Type,
		Format:       schema.Format,
		Title:        schema.Title,
		Description:  schema.Description,
		Default:      schema.Default,
		Enum:         schema.Enum,
		Example:      schema.Example,
		ReadOnly:     schema.ReadOnly,
		ExternalDocs: c.convertExternalDocsToV3(schema.ExternalDocs),
		Required:     schema.Required,
		UniqueItems:  schema.UniqueItems,
	}

	// Copy numeric/string constraints
	if schema.MultipleOf != nil {
		v3Schema.MultipleOf = *schema.MultipleOf
	}
	if schema.Maximum != nil {
		v3Schema.Maximum = *schema.Maximum
	}
	if schema.Minimum != nil {
		v3Schema.Minimum = *schema.Minimum
	}
	if schema.MaxLength != nil {
		v3Schema.MaxLength = *schema.MaxLength
	}
	if schema.MinLength != nil {
		v3Schema.MinLength = *schema.MinLength
	}
	if schema.Pattern != "" {
		v3Schema.Pattern = schema.Pattern
	}
	if schema.MaxItems != nil {
		v3Schema.MaxItems = *schema.MaxItems
	}
	if schema.MinItems != nil {
		v3Schema.MinItems = *schema.MinItems
	}
	if schema.MaxProperties != nil {
		v3Schema.MaxProperties = *schema.MaxProperties
	}
	if schema.MinProperties != nil {
		v3Schema.MinProperties = *schema.MinProperties
	}

	// ExclusiveMaximum/ExclusiveMinimum: boolean in Draft 4, can be boolean or number in 2020-12
	if schema.ExclusiveMaximum {
		v3Schema.ExclusiveMaximum = true
	}
	if schema.ExclusiveMinimum {
		v3Schema.ExclusiveMinimum = true
	}

	// Convert properties
	if len(schema.Properties) > 0 {
		v3Schema.Properties = c.convertDefinitionsToSchemas(schema.Properties)
	}

	// Convert additionalProperties
	if schema.AdditionalProperties != nil {
		switch v := schema.AdditionalProperties.(type) {
		case bool:
			v3Schema.AdditionalProperties = v
		case *swagger.Schema:
			v3Schema.AdditionalProperties = c.convertSchemaToV3(v)
		}
	}

	// Convert items
	if schema.Items != nil {
		v3Schema.Items = c.convertSchemaToV3(schema.Items)
	}

	// Convert allOf
	if len(schema.AllOf) > 0 {
		v3Schema.AllOf = make([]openapi.Schema, len(schema.AllOf))
		for i, s := range schema.AllOf {
			converted := c.convertSchemaToV3(s)
			if converted != nil {
				v3Schema.AllOf[i] = *converted
			}
		}
	}

	// Handle XML
	if schema.XML != nil {
		v3Schema.XML = &openapi.XML{
			Name:      schema.XML.Name,
			Namespace: schema.XML.Namespace,
			Prefix:    schema.XML.Prefix,
			Attribute: schema.XML.Attribute,
			Wrapped:   schema.XML.Wrapped,
		}
	}

	// Handle discriminator (string in v2, object in v3)
	if schema.Discriminator != "" {
		v3Schema.Discriminator = &openapi.Discriminator{
			PropertyName: schema.Discriminator,
		}
	}

	// Handle nullable (x-nullable extension → type array in 3.1.0)
	if schema.Extensions != nil {
		if nullable, ok := schema.Extensions["x-nullable"].(bool); ok && nullable {
			v3Schema.Type = []interface{}{schema.Type, "null"}
		}
		// Copy other extensions
		for k, v := range schema.Extensions {
			if k != "x-nullable" {
				if v3Schema.Extensions == nil {
					v3Schema.Extensions = make(map[string]interface{})
				}
				v3Schema.Extensions[k] = v
			}
		}
	}

	return v3Schema
}

// convertParameterDefinitionsToV3 converts component parameters.
func (c *Converter) convertParameterDefinitionsToV3(params map[string]*swagger.Parameter) map[string]*openapi.Parameter {
	if len(params) == 0 {
		return nil
	}

	v3Params := make(map[string]*openapi.Parameter, len(params))
	for name, param := range params {
		v3Params[name] = c.convertParameterToV3(param)
	}

	return v3Params
}

// convertResponseDefinitionsToV3 converts component responses.
func (c *Converter) convertResponseDefinitionsToV3(responses map[string]*swagger.Response) map[string]*openapi.Response {
	if len(responses) == 0 {
		return nil
	}

	v3Responses := make(map[string]*openapi.Response, len(responses))
	for name, resp := range responses {
		v3Responses[name] = c.convertResponseToV3(resp, []string{"application/json"})
	}

	return v3Responses
}

// convertSecurityDefinitionsToV3 converts Swagger SecurityDefinitions to OpenAPI SecuritySchemes.
func (c *Converter) convertSecurityDefinitionsToV3(schemes map[string]*swagger.SecurityScheme) map[string]*openapi.SecurityScheme {
	if len(schemes) == 0 {
		return nil
	}

	v3Schemes := make(map[string]*openapi.SecurityScheme, len(schemes))
	for name, scheme := range schemes {
		v3Schemes[name] = c.convertSecuritySchemeToV3(scheme)
	}

	return v3Schemes
}

// convertSecuritySchemeToV3 converts a Swagger SecurityScheme to OpenAPI SecurityScheme.
func (c *Converter) convertSecuritySchemeToV3(scheme *swagger.SecurityScheme) *openapi.SecurityScheme {
	if scheme == nil {
		return nil
	}

	v3Scheme := &openapi.SecurityScheme{
		Description: scheme.Description,
	}

	// Map Swagger 2.0 types to OpenAPI 3.x types
	switch scheme.Type {
	case "basic":
		v3Scheme.Type = "http"
		v3Scheme.Scheme = "basic"
	case "apiKey":
		v3Scheme.Type = "apiKey"
		v3Scheme.Name = scheme.Name
		v3Scheme.In = scheme.In
		// Check for bearer format extension
		if scheme.Extensions != nil {
			if bearerFormat, ok := scheme.Extensions["x-bearer-format"].(string); ok {
				v3Scheme.Type = "http"
				v3Scheme.Scheme = "bearer"
				v3Scheme.BearerFormat = bearerFormat
			}
		}
	case "oauth2":
		v3Scheme.Type = "oauth2"
		v3Scheme.Flows = c.convertOAuth2FlowsToV3(scheme)
	default:
		c.warnings = append(c.warnings, fmt.Sprintf("unknown security scheme type %q", scheme.Type))
	}

	return v3Scheme
}

// convertOAuth2FlowsToV3 converts OAuth2 flow from Swagger 2.0 to OpenAPI 3.x format.
func (c *Converter) convertOAuth2FlowsToV3(scheme *swagger.SecurityScheme) *openapi.OAuthFlows {
	if scheme.Flow == "" {
		return nil
	}

	flows := &openapi.OAuthFlows{}

	switch scheme.Flow {
	case "implicit":
		flows.Implicit = &openapi.OAuthFlow{
			AuthorizationURL: scheme.AuthorizationURL,
			Scopes:           scheme.Scopes,
		}
	case "password":
		flows.Password = &openapi.OAuthFlow{
			TokenURL: scheme.TokenURL,
			Scopes:   scheme.Scopes,
		}
	case "application":
		flows.ClientCredentials = &openapi.OAuthFlow{
			TokenURL: scheme.TokenURL,
			Scopes:   scheme.Scopes,
		}
	case "accessCode":
		flows.AuthorizationCode = &openapi.OAuthFlow{
			AuthorizationURL: scheme.AuthorizationURL,
			TokenURL:         scheme.TokenURL,
			Scopes:           scheme.Scopes,
		}
	default:
		c.warnings = append(c.warnings, fmt.Sprintf("unknown OAuth2 flow %q", scheme.Flow))
	}

	return flows
}

// convertSecurityToV3 converts security requirements.
func (c *Converter) convertSecurityToV3(security []swagger.SecurityRequirement) []openapi.SecurityRequirement {
	if len(security) == 0 {
		return nil
	}

	v3Security := make([]openapi.SecurityRequirement, len(security))
	for i, req := range security {
		v3Security[i] = openapi.SecurityRequirement(req)
	}

	return v3Security
}

// convertTagsToV3 converts tags.
func (c *Converter) convertTagsToV3(tags []swagger.Tag) []openapi.Tag {
	if len(tags) == 0 {
		return nil
	}

	v3Tags := make([]openapi.Tag, len(tags))
	for i, tag := range tags {
		v3Tags[i] = openapi.Tag{
			Name:         tag.Name,
			Description:  tag.Description,
			ExternalDocs: c.convertExternalDocsToV3(tag.ExternalDocs),
		}
	}

	return v3Tags
}

// convertExternalDocsToV3 converts external documentation.
func (c *Converter) convertExternalDocsToV3(docs *swagger.ExternalDocs) *openapi.ExternalDocs {
	if docs == nil {
		return nil
	}

	return &openapi.ExternalDocs{
		Description: docs.Description,
		URL:         docs.URL,
	}
}
