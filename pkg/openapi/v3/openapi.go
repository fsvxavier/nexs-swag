// Package v3 implements OpenAPI 3.1.x specification structures.
package v3

import (
	"encoding/json"
	"fmt"
)

// OpenAPI represents the root document object of the OpenAPI 3.1.x document.
// Spec: https://spec.openapis.org/oas/v3.1.0#openapi-object
type OpenAPI struct {
	OpenAPI           string                `json:"openapi"                     yaml:"openapi"`                     // REQUIRED. OpenAPI Specification version (3.1.0)
	Info              Info                  `json:"info"                        yaml:"info"`                        // REQUIRED. Metadata about the API
	JSONSchemaDialect string                `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"` // Default JSON Schema dialect
	Servers           []Server              `json:"servers,omitempty"           yaml:"servers,omitempty"`           // Array of Server Objects
	Paths             Paths                 `json:"paths,omitempty"             yaml:"paths,omitempty"`             // Available paths and operations
	Webhooks          map[string]*PathItem  `json:"webhooks,omitempty"          yaml:"webhooks,omitempty"`          // Webhooks (new in 3.1)
	Components        *Components           `json:"components,omitempty"        yaml:"components,omitempty"`        // Reusable components
	Security          []SecurityRequirement `json:"security,omitempty"          yaml:"security,omitempty"`          // Security requirements
	Tags              []Tag                 `json:"tags,omitempty"              yaml:"tags,omitempty"`              // List of tags
	ExternalDocs      *ExternalDocs         `json:"externalDocs,omitempty"      yaml:"externalDocs,omitempty"`      // Additional external documentation
}

// GetVersion returns the OpenAPI specification version.
func (o *OpenAPI) GetVersion() string {
	return "3.1.0"
}

// GetTitle returns the API title.
func (o *OpenAPI) GetTitle() string {
	return o.Info.Title
}

// GetInfo returns the Info object.
func (o *OpenAPI) GetInfo() interface{} {
	return o.Info
}

// Validate performs basic validation of the OpenAPI specification.
func (o *OpenAPI) Validate() error {
	if o.OpenAPI == "" {
		return fmt.Errorf("openapi version is required")
	}
	if o.Info.Title == "" {
		return fmt.Errorf("info.title is required")
	}
	if o.Info.Version == "" {
		return fmt.Errorf("info.version is required")
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling.
func (o *OpenAPI) MarshalJSON() ([]byte, error) {
	type Alias OpenAPI
	return json.Marshal((*Alias)(o))
}

// Info provides metadata about the API.
type Info struct {
	Title          string   `json:"title"                    yaml:"title"`                    // REQUIRED. Title of the API
	Summary        string   `json:"summary,omitempty"        yaml:"summary,omitempty"`        // Short summary (new in 3.1)
	Description    string   `json:"description,omitempty"    yaml:"description,omitempty"`    // Description (CommonMark syntax)
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"` // URL to terms of service
	Contact        *Contact `json:"contact,omitempty"        yaml:"contact,omitempty"`        // Contact information
	License        *License `json:"license,omitempty"        yaml:"license,omitempty"`        // License information
	Version        string   `json:"version"                  yaml:"version"`                  // REQUIRED. API version
}

// Contact information for the API.
type Contact struct {
	Name  string `json:"name,omitempty"  yaml:"name,omitempty"`  // Contact name
	URL   string `json:"url,omitempty"   yaml:"url,omitempty"`   // Contact URL
	Email string `json:"email,omitempty" yaml:"email,omitempty"` // Contact email
}

// License information for the API.
type License struct {
	Name       string `json:"name"                 yaml:"name"`                 // REQUIRED. License name
	Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"` // SPDX license identifier (new in 3.1)
	URL        string `json:"url,omitempty"        yaml:"url,omitempty"`        // URL to license
}

// Server represents a server.
type Server struct {
	URL         string                     `json:"url"                   yaml:"url"`                   // REQUIRED. Server URL
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"` // Server description
	Variables   map[string]*ServerVariable `json:"variables,omitempty"   yaml:"variables,omitempty"`   // Server variables
}

// ServerVariable for server URL template substitution.
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"        yaml:"enum,omitempty"`        // Enumeration of values
	Default     string   `json:"default"               yaml:"default"`               // REQUIRED. Default value
	Description string   `json:"description,omitempty" yaml:"description,omitempty"` // Description
}

// Paths holds the relative paths to the individual endpoints.
type Paths map[string]*PathItem

// PathItem describes operations available on a single path.
type PathItem struct {
	Ref         string      `json:"$ref,omitempty"        yaml:"$ref,omitempty"`        // Reference to another PathItem
	Summary     string      `json:"summary,omitempty"     yaml:"summary,omitempty"`     // Summary for all operations
	Description string      `json:"description,omitempty" yaml:"description,omitempty"` // Description for all operations
	Get         *Operation  `json:"get,omitempty"         yaml:"get,omitempty"`         // GET operation
	Put         *Operation  `json:"put,omitempty"         yaml:"put,omitempty"`         // PUT operation
	Post        *Operation  `json:"post,omitempty"        yaml:"post,omitempty"`        // POST operation
	Delete      *Operation  `json:"delete,omitempty"      yaml:"delete,omitempty"`      // DELETE operation
	Options     *Operation  `json:"options,omitempty"     yaml:"options,omitempty"`     // OPTIONS operation
	Head        *Operation  `json:"head,omitempty"        yaml:"head,omitempty"`        // HEAD operation
	Patch       *Operation  `json:"patch,omitempty"       yaml:"patch,omitempty"`       // PATCH operation
	Trace       *Operation  `json:"trace,omitempty"       yaml:"trace,omitempty"`       // TRACE operation
	Query       *Operation  `json:"query,omitempty"       yaml:"query,omitempty"`       // QUERY operation (new in 3.2.0)
	Servers     []Server    `json:"servers,omitempty"     yaml:"servers,omitempty"`     // Alternative servers
	Parameters  []Parameter `json:"parameters,omitempty"  yaml:"parameters,omitempty"`  // Common parameters
}

// Operation describes a single API operation on a path.
type Operation struct {
	Tags         []string               `json:"tags,omitempty"         yaml:"tags,omitempty"`         // Tags for API documentation control
	Summary      string                 `json:"summary,omitempty"      yaml:"summary,omitempty"`      // Short summary
	Description  string                 `json:"description,omitempty"  yaml:"description,omitempty"`  // Verbose explanation
	ExternalDocs *ExternalDocs          `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"` // Additional external documentation
	OperationID  string                 `json:"operationId,omitempty"  yaml:"operationId,omitempty"`  // Unique operation identifier
	Parameters   []Parameter            `json:"parameters,omitempty"   yaml:"parameters,omitempty"`   // List of parameters
	RequestBody  *RequestBody           `json:"requestBody,omitempty"  yaml:"requestBody,omitempty"`  // Request body
	Responses    Responses              `json:"responses"              yaml:"responses"`              // REQUIRED. Possible responses
	Callbacks    map[string]*Callback   `json:"callbacks,omitempty"    yaml:"callbacks,omitempty"`    // Callbacks (webhooks)
	Deprecated   bool                   `json:"deprecated,omitempty"   yaml:"deprecated,omitempty"`   // Operation is deprecated
	Security     []SecurityRequirement  `json:"security,omitempty"     yaml:"security,omitempty"`     // Security requirements
	Servers      []Server               `json:"servers,omitempty"      yaml:"servers,omitempty"`      // Alternative servers
	Extensions   map[string]interface{} `json:"-"                      yaml:"-"`                      // Custom extensions (x-*)
}

// Parameter describes a single operation parameter.
type Parameter struct {
	Name            string              `json:"name"                      yaml:"name"`                      // REQUIRED. Parameter name
	In              string              `json:"in"                        yaml:"in"`                        // REQUIRED. Location: query, header, path, cookie
	Description     string              `json:"description,omitempty"     yaml:"description,omitempty"`     // Parameter description
	Required        bool                `json:"required,omitempty"        yaml:"required,omitempty"`        // Required (true for path parameters)
	Deprecated      bool                `json:"deprecated,omitempty"      yaml:"deprecated,omitempty"`      // Parameter is deprecated
	AllowEmptyValue bool                `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"` // Allow empty value
	Schema          *Schema             `json:"schema,omitempty"          yaml:"schema,omitempty"`          // Parameter schema
	Example         interface{}         `json:"example,omitempty"         yaml:"example,omitempty"`         // Example value
	Examples        map[string]*Example `json:"examples,omitempty"        yaml:"examples,omitempty"`        // Multiple examples
}

// RequestBody describes a single request body.
type RequestBody struct {
	Description string                `json:"description,omitempty" yaml:"description,omitempty"` // Description
	Content     map[string]*MediaType `json:"content"               yaml:"content"`               // REQUIRED. Content (MIME types)
	Required    bool                  `json:"required,omitempty"    yaml:"required,omitempty"`    // Request body is required
}

// MediaType provides schema and examples for the media type.
type MediaType struct {
	Schema       *Schema              `json:"schema,omitempty"       yaml:"schema,omitempty"`       // Schema
	Example      interface{}          `json:"example,omitempty"      yaml:"example,omitempty"`      // Example value
	Examples     map[string]*Example  `json:"examples,omitempty"     yaml:"examples,omitempty"`     // Multiple examples
	Encoding     map[string]*Encoding `json:"encoding,omitempty"     yaml:"encoding,omitempty"`     // Encoding for multipart
	ItemSchema   *Schema              `json:"itemSchema,omitempty"   yaml:"itemSchema,omitempty"`   // Schema for streaming items (new in 3.2.0)
	ItemEncoding map[string]*Encoding `json:"itemEncoding,omitempty" yaml:"itemEncoding,omitempty"` // Encoding for streaming items (new in 3.2.0)
}

// Encoding for request body properties.
type Encoding struct {
	ContentType   string             `json:"contentType,omitempty"   yaml:"contentType,omitempty"`   // Content-Type
	Headers       map[string]*Header `json:"headers,omitempty"       yaml:"headers,omitempty"`       // Headers
	Style         string             `json:"style,omitempty"         yaml:"style,omitempty"`         // Serialization style
	Explode       bool               `json:"explode,omitempty"       yaml:"explode,omitempty"`       // Explode parameter
	AllowReserved bool               `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"` // Allow reserved characters
}

// Responses container for the expected responses of an operation.
type Responses map[string]*Response

// Response describes a single response from an API operation.
type Response struct {
	Description string                `json:"description"       yaml:"description"`       // REQUIRED. Response description
	Headers     map[string]*Header    `json:"headers,omitempty" yaml:"headers,omitempty"` // Response headers
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"` // Response content
	Links       map[string]*Link      `json:"links,omitempty"   yaml:"links,omitempty"`   // Links to other operations
}

// Header describes a single header.
type Header struct {
	Description string              `json:"description,omitempty" yaml:"description,omitempty"` // Header description
	Required    bool                `json:"required,omitempty"    yaml:"required,omitempty"`    // Header is required
	Deprecated  bool                `json:"deprecated,omitempty"  yaml:"deprecated,omitempty"`  // Header is deprecated
	Schema      *Schema             `json:"schema,omitempty"      yaml:"schema,omitempty"`      // Header schema
	Example     interface{}         `json:"example,omitempty"     yaml:"example,omitempty"`     // Example value
	Examples    map[string]*Example `json:"examples,omitempty"    yaml:"examples,omitempty"`    // Multiple examples
}

// Example for parameter, request body, or response.
type Example struct {
	Summary       string      `json:"summary,omitempty"       yaml:"summary,omitempty"`       // Short description
	Description   string      `json:"description,omitempty"   yaml:"description,omitempty"`   // Long description
	Value         interface{} `json:"value,omitempty"         yaml:"value,omitempty"`         // Embedded example value
	ExternalValue string      `json:"externalValue,omitempty" yaml:"externalValue,omitempty"` // URL to external example
}

// Link represents a possible design-time link for a response.
type Link struct {
	OperationRef string                 `json:"operationRef,omitempty" yaml:"operationRef,omitempty"` // Reference to operation
	OperationID  string                 `json:"operationId,omitempty"  yaml:"operationId,omitempty"`  // Operation ID
	Parameters   map[string]interface{} `json:"parameters,omitempty"   yaml:"parameters,omitempty"`   // Parameters
	RequestBody  interface{}            `json:"requestBody,omitempty"  yaml:"requestBody,omitempty"`  // Request body value
	Description  string                 `json:"description,omitempty"  yaml:"description,omitempty"`  // Description
	Server       *Server                `json:"server,omitempty"       yaml:"server,omitempty"`       // Server object
}

// Callback represents a callback request.
type Callback map[string]*PathItem

// Schema represents a JSON Schema (draft 2020-12 compatible with OpenAPI 3.1).
type Schema struct {
	// JSON Schema fields
	Type        interface{}   `json:"type,omitempty"        yaml:"type,omitempty"`        // Type (string or array)
	Format      string        `json:"format,omitempty"      yaml:"format,omitempty"`      // Format
	Title       string        `json:"title,omitempty"       yaml:"title,omitempty"`       // Title
	Description string        `json:"description,omitempty" yaml:"description,omitempty"` // Description
	Default     interface{}   `json:"default,omitempty"     yaml:"default,omitempty"`     // Default value
	Enum        []interface{} `json:"enum,omitempty"        yaml:"enum,omitempty"`        // Enum values
	Const       interface{}   `json:"const,omitempty"       yaml:"const,omitempty"`       // Const value (JSON Schema 2020-12)

	// Number validation
	MultipleOf       float64     `json:"multipleOf,omitempty"       yaml:"multipleOf,omitempty"`       // Multiple of
	Maximum          float64     `json:"maximum,omitempty"          yaml:"maximum,omitempty"`          // Maximum value
	ExclusiveMaximum interface{} `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"` // Exclusive maximum
	Minimum          float64     `json:"minimum,omitempty"          yaml:"minimum,omitempty"`          // Minimum value
	ExclusiveMinimum interface{} `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"` // Exclusive minimum

	// String validation
	MaxLength int    `json:"maxLength,omitempty" yaml:"maxLength,omitempty"` // Maximum length
	MinLength int    `json:"minLength,omitempty" yaml:"minLength,omitempty"` // Minimum length
	Pattern   string `json:"pattern,omitempty"   yaml:"pattern,omitempty"`   // Regex pattern

	// Array validation
	Items       *Schema   `json:"items,omitempty"       yaml:"items,omitempty"`       // Array items
	MaxItems    int       `json:"maxItems,omitempty"    yaml:"maxItems,omitempty"`    // Maximum items
	MinItems    int       `json:"minItems,omitempty"    yaml:"minItems,omitempty"`    // Minimum items
	UniqueItems bool      `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"` // Unique items
	PrefixItems []*Schema `json:"prefixItems,omitempty" yaml:"prefixItems,omitempty"` // Prefix items (JSON Schema 2020-12)

	// Object validation
	Properties           map[string]*Schema `json:"properties,omitempty"           yaml:"properties,omitempty"`           // Object properties
	Required             []string           `json:"required,omitempty"             yaml:"required,omitempty"`             // Required properties
	AdditionalProperties interface{}        `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"` // Additional properties (bool or Schema)
	MaxProperties        int                `json:"maxProperties,omitempty"        yaml:"maxProperties,omitempty"`        // Maximum properties
	MinProperties        int                `json:"minProperties,omitempty"        yaml:"minProperties,omitempty"`        // Minimum properties

	// Composition
	AllOf []Schema `json:"allOf,omitempty" yaml:"allOf,omitempty"` // All of (intersection)
	OneOf []Schema `json:"oneOf,omitempty" yaml:"oneOf,omitempty"` // One of (union)
	AnyOf []Schema `json:"anyOf,omitempty" yaml:"anyOf,omitempty"` // Any of
	Not   *Schema  `json:"not,omitempty"   yaml:"not,omitempty"`   // Not

	// OpenAPI specific
	Nullable      bool                   `json:"nullable,omitempty"      yaml:"nullable,omitempty"`      // Deprecated in 3.1 (use type: [type, "null"])
	Discriminator *Discriminator         `json:"discriminator,omitempty" yaml:"discriminator,omitempty"` // Discriminator
	ReadOnly      bool                   `json:"readOnly,omitempty"      yaml:"readOnly,omitempty"`      // Read only
	WriteOnly     bool                   `json:"writeOnly,omitempty"     yaml:"writeOnly,omitempty"`     // Write only
	XML           *XML                   `json:"xml,omitempty"           yaml:"xml,omitempty"`           // XML representation
	ExternalDocs  *ExternalDocs          `json:"externalDocs,omitempty"  yaml:"externalDocs,omitempty"`  // External documentation
	Example       interface{}            `json:"example,omitempty"       yaml:"example,omitempty"`       // Example value (deprecated, use examples)
	Deprecated    bool                   `json:"deprecated,omitempty"    yaml:"deprecated,omitempty"`    // Deprecated
	Extensions    map[string]interface{} `json:"-"                       yaml:"-"`                       // Custom extensions (x-*)

	// Reference
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"` // JSON Reference
}

// Discriminator for polymorphism.
type Discriminator struct {
	PropertyName string            `json:"propertyName"      yaml:"propertyName"`      // REQUIRED. Property name
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"` // Mapping of values to schemas
}

// XML representation hints.
type XML struct {
	Name      string `json:"name,omitempty"      yaml:"name,omitempty"`      // XML element name
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"` // XML namespace URI
	Prefix    string `json:"prefix,omitempty"    yaml:"prefix,omitempty"`    // XML namespace prefix
	Attribute bool   `json:"attribute,omitempty" yaml:"attribute,omitempty"` // XML attribute
	Wrapped   bool   `json:"wrapped,omitempty"   yaml:"wrapped,omitempty"`   // XML wrapped
}

// Components holds reusable objects.
type Components struct {
	Schemas         map[string]*Schema         `json:"schemas,omitempty"         yaml:"schemas,omitempty"`         // Reusable schemas
	Responses       map[string]*Response       `json:"responses,omitempty"       yaml:"responses,omitempty"`       // Reusable responses
	Parameters      map[string]*Parameter      `json:"parameters,omitempty"      yaml:"parameters,omitempty"`      // Reusable parameters
	Examples        map[string]*Example        `json:"examples,omitempty"        yaml:"examples,omitempty"`        // Reusable examples
	RequestBodies   map[string]*RequestBody    `json:"requestBodies,omitempty"   yaml:"requestBodies,omitempty"`   // Reusable request bodies
	Headers         map[string]*Header         `json:"headers,omitempty"         yaml:"headers,omitempty"`         // Reusable headers
	SecuritySchemes map[string]*SecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"` // Security schemes
	Links           map[string]*Link           `json:"links,omitempty"           yaml:"links,omitempty"`           // Reusable links
	Callbacks       map[string]*Callback       `json:"callbacks,omitempty"       yaml:"callbacks,omitempty"`       // Reusable callbacks
	PathItems       map[string]*PathItem       `json:"pathItems,omitempty"       yaml:"pathItems,omitempty"`       // Reusable path items (new in 3.1)
}

// SecurityScheme defines a security scheme.
type SecurityScheme struct {
	Type              string      `json:"type"                       yaml:"type"`                         // REQUIRED. Type: apiKey, http, oauth2, openIdConnect, mutualTLS
	Description       string      `json:"description,omitempty"      yaml:"description,omitempty"`        // Description
	Name              string      `json:"name,omitempty"             yaml:"name,omitempty"`               // Name (for apiKey)
	In                string      `json:"in,omitempty"               yaml:"in,omitempty"`                 // Location (for apiKey): query, header, cookie
	Scheme            string      `json:"scheme,omitempty"           yaml:"scheme,omitempty"`             // HTTP scheme (for http): basic, bearer, etc.
	BearerFormat      string      `json:"bearerFormat,omitempty"     yaml:"bearerFormat,omitempty"`       // Bearer token format (for http bearer)
	Flows             *OAuthFlows `json:"flows,omitempty"            yaml:"flows,omitempty"`              // OAuth flows (for oauth2)
	OpenIDConnectURL  string      `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`   // OpenID Connect URL (for openIdConnect)
	Deprecated        bool        `json:"deprecated,omitempty"       yaml:"deprecated,omitempty"`         // Deprecated (new in 3.2.0)
	OAuth2MetadataURL string      `json:"oauth2MetadataUrl,omitempty" yaml:"oauth2MetadataUrl,omitempty"` // OAuth2 metadata URL (new in 3.2.0)
}

// OAuthFlows configuration for OAuth 2.0.
type OAuthFlows struct {
	Implicit            *OAuthFlow `json:"implicit,omitempty"          yaml:"implicit,omitempty"`              // Implicit flow
	Password            *OAuthFlow `json:"password,omitempty"          yaml:"password,omitempty"`              // Password flow
	ClientCredentials   *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`     // Client credentials flow
	AuthorizationCode   *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`     // Authorization code flow
	DeviceAuthorization *OAuthFlow `json:"deviceAuthorization,omitempty" yaml:"deviceAuthorization,omitempty"` // Device authorization flow (new in 3.2.0)
}

// OAuthFlow configuration for a single OAuth flow.
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"` // Authorization URL
	TokenURL         string            `json:"tokenUrl,omitempty"         yaml:"tokenUrl,omitempty"`         // Token URL
	RefreshURL       string            `json:"refreshUrl,omitempty"       yaml:"refreshUrl,omitempty"`       // Refresh URL
	Scopes           map[string]string `json:"scopes"                     yaml:"scopes"`                     // REQUIRED. Available scopes
}

// SecurityRequirement lists required security schemes.
type SecurityRequirement map[string][]string

// Tag for API documentation organization.
type Tag struct {
	Name         string        `json:"name"                   yaml:"name"`                   // REQUIRED. Tag name
	Description  string        `json:"description,omitempty"  yaml:"description,omitempty"`  // Tag description
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"` // Additional documentation
}

// ExternalDocs references external documentation.
type ExternalDocs struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"` // Description
	URL         string `json:"url"                   yaml:"url"`                   // REQUIRED. URL
}

// MarshalJSON customizes JSON encoding to include extensions as top-level fields.
func (s *Schema) MarshalJSON() ([]byte, error) {
	// Create a type alias to avoid infinite recursion
	type Alias Schema

	// Marshal the schema normally
	base, err := json.Marshal((*Alias)(s))
	if err != nil {
		return nil, err
	}

	// If no extensions, return as is
	if len(s.Extensions) == 0 {
		return base, nil
	}

	// Unmarshal to map to add extensions
	var result map[string]interface{}
	if err := json.Unmarshal(base, &result); err != nil {
		return nil, err
	}

	// Add extensions to top level
	for k, v := range s.Extensions {
		result[k] = v
	}

	return json.Marshal(result)
}

// MarshalJSON customizes JSON encoding to include extensions as top-level fields.
func (o *Operation) MarshalJSON() ([]byte, error) {
	// Create a type alias to avoid infinite recursion
	type Alias Operation

	// Marshal the operation normally
	base, err := json.Marshal((*Alias)(o))
	if err != nil {
		return nil, err
	}

	// If no extensions, return as is
	if len(o.Extensions) == 0 {
		return base, nil
	}

	// Unmarshal to map to add extensions
	var result map[string]interface{}
	if err := json.Unmarshal(base, &result); err != nil {
		return nil, err
	}

	// Add extensions to top level
	for k, v := range o.Extensions {
		result[k] = v
	}

	return json.Marshal(result)
}
