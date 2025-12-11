// Package v2 implements Swagger 2.0 / OpenAPI 2.0 specification structures.
// Spec: https://swagger.io/specification/v2/
package v2

import (
	"encoding/json"
	"fmt"
)

// Swagger represents the root document object of the Swagger 2.0 specification.
type Swagger struct {
	Swagger             string                     `json:"swagger"                       yaml:"swagger"`                       // REQUIRED. Swagger version (always "2.0")
	Info                Info                       `json:"info"                          yaml:"info"`                          // REQUIRED. Metadata about the API
	Host                string                     `json:"host,omitempty"                yaml:"host,omitempty"`                // Host (hostname or ip:port)
	BasePath            string                     `json:"basePath,omitempty"            yaml:"basePath,omitempty"`            // Base path (must start with /)
	Schemes             []string                   `json:"schemes,omitempty"             yaml:"schemes,omitempty"`             // Transfer protocols: http, https, ws, wss
	Consumes            []string                   `json:"consumes,omitempty"            yaml:"consumes,omitempty"`            // MIME types the APIs can consume
	Produces            []string                   `json:"produces,omitempty"            yaml:"produces,omitempty"`            // MIME types the APIs can produce
	Paths               Paths                      `json:"paths"                         yaml:"paths"`                         // REQUIRED. Available paths and operations
	Definitions         map[string]*Schema         `json:"definitions,omitempty"         yaml:"definitions,omitempty"`         // Data types definitions
	Parameters          map[string]*Parameter      `json:"parameters,omitempty"          yaml:"parameters,omitempty"`          // Reusable parameters
	Responses           map[string]*Response       `json:"responses,omitempty"           yaml:"responses,omitempty"`           // Reusable responses
	SecurityDefinitions map[string]*SecurityScheme `json:"securityDefinitions,omitempty" yaml:"securityDefinitions,omitempty"` // Security scheme definitions
	Security            []SecurityRequirement      `json:"security,omitempty"            yaml:"security,omitempty"`            // Security requirements
	Tags                []Tag                      `json:"tags,omitempty"                yaml:"tags,omitempty"`                // List of tags
	ExternalDocs        *ExternalDocs              `json:"externalDocs,omitempty"        yaml:"externalDocs,omitempty"`        // Additional external documentation
	Extensions          map[string]interface{}     `json:"-"                             yaml:"-"`                             // Custom extensions (x-*)
}

// GetVersion returns the Swagger specification version.
func (s *Swagger) GetVersion() string {
	return "2.0"
}

// Validate performs basic validation of the Swagger specification.
func (s *Swagger) Validate() error {
	if s.Swagger != "2.0" {
		return fmt.Errorf("invalid swagger version: %s (must be 2.0)", s.Swagger)
	}
	if s.Info.Title == "" {
		return fmt.Errorf("info.title is required")
	}
	if s.Info.Version == "" {
		return fmt.Errorf("info.version is required")
	}
	if s.Paths == nil {
		return fmt.Errorf("paths is required")
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling with extensions support.
func (s *Swagger) MarshalJSON() ([]byte, error) {
	type Alias Swagger
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	data, err := json.Marshal(aux)
	if err != nil {
		return nil, err
	}

	// Merge extensions if any
	if len(s.Extensions) > 0 {
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		for k, v := range s.Extensions {
			m[k] = v
		}
		return json.Marshal(m)
	}

	return data, nil
}

// Info provides metadata about the API.
type Info struct {
	Title          string   `json:"title"                    yaml:"title"`                    // REQUIRED. Application title
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
	Name string `json:"name"           yaml:"name"`          // REQUIRED. License name
	URL  string `json:"url,omitempty"  yaml:"url,omitempty"` // URL to license
}

// Paths holds the relative paths to the individual endpoints.
type Paths map[string]*PathItem

// PathItem describes operations available on a single path.
type PathItem struct {
	Ref        string                 `json:"$ref,omitempty"        yaml:"$ref,omitempty"`       // Reference to another PathItem
	Get        *Operation             `json:"get,omitempty"         yaml:"get,omitempty"`        // GET operation
	Put        *Operation             `json:"put,omitempty"         yaml:"put,omitempty"`        // PUT operation
	Post       *Operation             `json:"post,omitempty"        yaml:"post,omitempty"`       // POST operation
	Delete     *Operation             `json:"delete,omitempty"      yaml:"delete,omitempty"`     // DELETE operation
	Options    *Operation             `json:"options,omitempty"     yaml:"options,omitempty"`    // OPTIONS operation
	Head       *Operation             `json:"head,omitempty"        yaml:"head,omitempty"`       // HEAD operation
	Patch      *Operation             `json:"patch,omitempty"       yaml:"patch,omitempty"`      // PATCH operation
	Parameters []*Parameter           `json:"parameters,omitempty"  yaml:"parameters,omitempty"` // Common parameters
	Extensions map[string]interface{} `json:"-"           yaml:"-"`                              // Custom extensions (x-*)
}

// Operation describes a single API operation on a path.
type Operation struct {
	Tags         []string               `json:"tags,omitempty"         yaml:"tags,omitempty"`         // Tags for API documentation control
	Summary      string                 `json:"summary,omitempty"      yaml:"summary,omitempty"`      // Short summary
	Description  string                 `json:"description,omitempty"  yaml:"description,omitempty"`  // Verbose explanation
	ExternalDocs *ExternalDocs          `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"` // Additional external documentation
	OperationID  string                 `json:"operationId,omitempty"  yaml:"operationId,omitempty"`  // Unique operation identifier
	Consumes     []string               `json:"consumes,omitempty"     yaml:"consumes,omitempty"`     // MIME types operation can consume
	Produces     []string               `json:"produces,omitempty"     yaml:"produces,omitempty"`     // MIME types operation can produce
	Parameters   []*Parameter           `json:"parameters,omitempty"   yaml:"parameters,omitempty"`   // List of parameters
	Responses    Responses              `json:"responses"              yaml:"responses"`              // REQUIRED. Possible responses
	Schemes      []string               `json:"schemes,omitempty"      yaml:"schemes,omitempty"`      // Transfer protocols
	Deprecated   bool                   `json:"deprecated,omitempty"   yaml:"deprecated,omitempty"`   // Operation is deprecated
	Security     []SecurityRequirement  `json:"security,omitempty"     yaml:"security,omitempty"`     // Security requirements
	Extensions   map[string]interface{} `json:"-"                     yaml:"-"`                       // Custom extensions (x-*)
}

// Parameter describes a single operation parameter.
type Parameter struct {
	Name        string `json:"name"                       yaml:"name"`                  // REQUIRED. Parameter name
	In          string `json:"in"                         yaml:"in"`                    // REQUIRED. Location: query, header, path, formData, body
	Description string `json:"description,omitempty"      yaml:"description,omitempty"` // Parameter description
	Required    bool   `json:"required,omitempty"         yaml:"required,omitempty"`    // Required parameter

	// For in != "body"
	Type             string        `json:"type,omitempty"             yaml:"type,omitempty"`             // Type: string, number, integer, boolean, array, file
	Format           string        `json:"format,omitempty"           yaml:"format,omitempty"`           // Format modifier
	AllowEmptyValue  bool          `json:"allowEmptyValue,omitempty"  yaml:"allowEmptyValue,omitempty"`  // Allow empty value
	Items            *Items        `json:"items,omitempty"            yaml:"items,omitempty"`            // Items definition (for type=array)
	CollectionFormat string        `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"` // Collection format: csv, ssv, tsv, pipes, multi
	Default          interface{}   `json:"default,omitempty"          yaml:"default,omitempty"`          // Default value
	Maximum          *float64      `json:"maximum,omitempty"          yaml:"maximum,omitempty"`          // Maximum value
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"` // Exclusive maximum
	Minimum          *float64      `json:"minimum,omitempty"          yaml:"minimum,omitempty"`          // Minimum value
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"` // Exclusive minimum
	MaxLength        *int          `json:"maxLength,omitempty"        yaml:"maxLength,omitempty"`        // Maximum length
	MinLength        *int          `json:"minLength,omitempty"        yaml:"minLength,omitempty"`        // Minimum length
	Pattern          string        `json:"pattern,omitempty"          yaml:"pattern,omitempty"`          // Regex pattern
	MaxItems         *int          `json:"maxItems,omitempty"         yaml:"maxItems,omitempty"`         // Maximum array items
	MinItems         *int          `json:"minItems,omitempty"         yaml:"minItems,omitempty"`         // Minimum array items
	UniqueItems      bool          `json:"uniqueItems,omitempty"      yaml:"uniqueItems,omitempty"`      // Unique array items
	Enum             []interface{} `json:"enum,omitempty"             yaml:"enum,omitempty"`             // Enumeration of values
	MultipleOf       *float64      `json:"multipleOf,omitempty"       yaml:"multipleOf,omitempty"`       // Multiple of

	// For in = "body"
	Schema *Schema `json:"schema,omitempty"           yaml:"schema,omitempty"` // Schema definition

	// Reference
	Ref string `json:"$ref,omitempty"             yaml:"$ref,omitempty"` // Reference to parameter definition

	Extensions map[string]interface{} `json:"-"                 yaml:"-"` // Custom extensions (x-*)
}

// Items describes the type of items in an array.
type Items struct {
	Type             string        `json:"type,omitempty"             yaml:"type,omitempty"`             // Type: string, number, integer, boolean, array
	Format           string        `json:"format,omitempty"           yaml:"format,omitempty"`           // Format modifier
	Items            *Items        `json:"items,omitempty"            yaml:"items,omitempty"`            // Nested items (for nested arrays)
	CollectionFormat string        `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"` // Collection format
	Default          interface{}   `json:"default,omitempty"          yaml:"default,omitempty"`          // Default value
	Maximum          *float64      `json:"maximum,omitempty"          yaml:"maximum,omitempty"`          // Maximum value
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"` // Exclusive maximum
	Minimum          *float64      `json:"minimum,omitempty"          yaml:"minimum,omitempty"`          // Minimum value
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"` // Exclusive minimum
	MaxLength        *int          `json:"maxLength,omitempty"        yaml:"maxLength,omitempty"`        // Maximum length
	MinLength        *int          `json:"minLength,omitempty"        yaml:"minLength,omitempty"`        // Minimum length
	Pattern          string        `json:"pattern,omitempty"          yaml:"pattern,omitempty"`          // Regex pattern
	MaxItems         *int          `json:"maxItems,omitempty"         yaml:"maxItems,omitempty"`         // Maximum array items
	MinItems         *int          `json:"minItems,omitempty"         yaml:"minItems,omitempty"`         // Minimum array items
	UniqueItems      bool          `json:"uniqueItems,omitempty"      yaml:"uniqueItems,omitempty"`      // Unique array items
	Enum             []interface{} `json:"enum,omitempty"             yaml:"enum,omitempty"`             // Enumeration of values
	MultipleOf       *float64      `json:"multipleOf,omitempty"       yaml:"multipleOf,omitempty"`       // Multiple of
	Ref              string        `json:"$ref,omitempty"             yaml:"$ref,omitempty"`             // Reference
}

// Schema represents a data type definition (JSON Schema Draft 4 subset).
type Schema struct {
	Ref              string        `json:"$ref,omitempty"                  yaml:"$ref,omitempty"`             // Reference to another schema
	Type             string        `json:"type,omitempty"                  yaml:"type,omitempty"`             // Type: string, number, integer, boolean, array, object
	Format           string        `json:"format,omitempty"                yaml:"format,omitempty"`           // Format modifier
	Title            string        `json:"title,omitempty"                 yaml:"title,omitempty"`            // Schema title
	Description      string        `json:"description,omitempty"           yaml:"description,omitempty"`      // Schema description
	Default          interface{}   `json:"default,omitempty"               yaml:"default,omitempty"`          // Default value
	MultipleOf       *float64      `json:"multipleOf,omitempty"            yaml:"multipleOf,omitempty"`       // Multiple of
	Maximum          *float64      `json:"maximum,omitempty"               yaml:"maximum,omitempty"`          // Maximum value
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"      yaml:"exclusiveMaximum,omitempty"` // Exclusive maximum
	Minimum          *float64      `json:"minimum,omitempty"               yaml:"minimum,omitempty"`          // Minimum value
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"      yaml:"exclusiveMinimum,omitempty"` // Exclusive minimum
	MaxLength        *int          `json:"maxLength,omitempty"             yaml:"maxLength,omitempty"`        // Maximum string length
	MinLength        *int          `json:"minLength,omitempty"             yaml:"minLength,omitempty"`        // Minimum string length
	Pattern          string        `json:"pattern,omitempty"               yaml:"pattern,omitempty"`          // Regex pattern
	MaxItems         *int          `json:"maxItems,omitempty"              yaml:"maxItems,omitempty"`         // Maximum array items
	MinItems         *int          `json:"minItems,omitempty"              yaml:"minItems,omitempty"`         // Minimum array items
	UniqueItems      bool          `json:"uniqueItems,omitempty"           yaml:"uniqueItems,omitempty"`      // Unique array items
	MaxProperties    *int          `json:"maxProperties,omitempty"         yaml:"maxProperties,omitempty"`    // Maximum object properties
	MinProperties    *int          `json:"minProperties,omitempty"         yaml:"minProperties,omitempty"`    // Minimum object properties
	Required         []string      `json:"required,omitempty"              yaml:"required,omitempty"`         // Required properties
	Enum             []interface{} `json:"enum,omitempty"                  yaml:"enum,omitempty"`             // Enumeration of values

	// Object properties
	Properties           map[string]*Schema `json:"properties,omitempty"            yaml:"properties,omitempty"`           // Object properties
	AdditionalProperties interface{}        `json:"additionalProperties,omitempty"  yaml:"additionalProperties,omitempty"` // Additional properties (bool or *Schema)

	// Array items
	Items *Schema `json:"items,omitempty"                 yaml:"items,omitempty"` // Array items schema

	// Composition
	AllOf []*Schema `json:"allOf,omitempty"                 yaml:"allOf,omitempty"` // AllOf composition

	// Discriminator
	Discriminator string        `json:"discriminator,omitempty"         yaml:"discriminator,omitempty"` // Discriminator property
	ReadOnly      bool          `json:"readOnly,omitempty"              yaml:"readOnly,omitempty"`      // Read-only property
	XML           *XML          `json:"xml,omitempty"                   yaml:"xml,omitempty"`           // XML representation
	ExternalDocs  *ExternalDocs `json:"externalDocs,omitempty"          yaml:"externalDocs,omitempty"`  // External documentation
	Example       interface{}   `json:"example,omitempty"               yaml:"example,omitempty"`       // Example value

	Extensions map[string]interface{} `json:"-"                               yaml:"-"` // Custom extensions (x-*)
}

// XML describes XML representation of a schema.
type XML struct {
	Name      string `json:"name,omitempty"      yaml:"name,omitempty"`      // XML element name
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"` // XML namespace URI
	Prefix    string `json:"prefix,omitempty"    yaml:"prefix,omitempty"`    // XML namespace prefix
	Attribute bool   `json:"attribute,omitempty" yaml:"attribute,omitempty"` // Translate to XML attribute
	Wrapped   bool   `json:"wrapped,omitempty"   yaml:"wrapped,omitempty"`   // Wrap array elements
}

// Responses is a container for the expected responses of an operation.
type Responses map[string]*Response

// Response describes a single response from an API operation.
type Response struct {
	Description string                 `json:"description"          yaml:"description"`        // REQUIRED. Response description
	Schema      *Schema                `json:"schema,omitempty"     yaml:"schema,omitempty"`   // Response schema
	Headers     map[string]*Header     `json:"headers,omitempty"    yaml:"headers,omitempty"`  // Response headers
	Examples    map[string]interface{} `json:"examples,omitempty"   yaml:"examples,omitempty"` // Response examples (MIME type → example)
	Ref         string                 `json:"$ref,omitempty"       yaml:"$ref,omitempty"`     // Reference to response definition
	Extensions  map[string]interface{} `json:"-"                    yaml:"-"`                  // Custom extensions (x-*)
}

// Header represents a single HTTP header.
type Header struct {
	Description      string        `json:"description,omitempty"      yaml:"description,omitempty"`      // Header description
	Type             string        `json:"type"                       yaml:"type"`                       // REQUIRED. Type
	Format           string        `json:"format,omitempty"           yaml:"format,omitempty"`           // Format modifier
	Items            *Items        `json:"items,omitempty"            yaml:"items,omitempty"`            // Items definition (for type=array)
	CollectionFormat string        `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"` // Collection format
	Default          interface{}   `json:"default,omitempty"          yaml:"default,omitempty"`          // Default value
	Maximum          *float64      `json:"maximum,omitempty"          yaml:"maximum,omitempty"`          // Maximum value
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"` // Exclusive maximum
	Minimum          *float64      `json:"minimum,omitempty"          yaml:"minimum,omitempty"`          // Minimum value
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"` // Exclusive minimum
	MaxLength        *int          `json:"maxLength,omitempty"        yaml:"maxLength,omitempty"`        // Maximum length
	MinLength        *int          `json:"minLength,omitempty"        yaml:"minLength,omitempty"`        // Minimum length
	Pattern          string        `json:"pattern,omitempty"          yaml:"pattern,omitempty"`          // Regex pattern
	MaxItems         *int          `json:"maxItems,omitempty"         yaml:"maxItems,omitempty"`         // Maximum array items
	MinItems         *int          `json:"minItems,omitempty"         yaml:"minItems,omitempty"`         // Minimum array items
	UniqueItems      bool          `json:"uniqueItems,omitempty"      yaml:"uniqueItems,omitempty"`      // Unique array items
	Enum             []interface{} `json:"enum,omitempty"             yaml:"enum,omitempty"`             // Enumeration of values
	MultipleOf       *float64      `json:"multipleOf,omitempty"       yaml:"multipleOf,omitempty"`       // Multiple of
}

// SecurityScheme defines a security scheme that can be used by operations.
type SecurityScheme struct {
	Type             string                 `json:"type"                       yaml:"type"`                       // REQUIRED. Type: basic, apiKey, oauth2
	Description      string                 `json:"description,omitempty"      yaml:"description,omitempty"`      // Security scheme description
	Name             string                 `json:"name,omitempty"             yaml:"name,omitempty"`             // Name of header/query parameter (for apiKey)
	In               string                 `json:"in,omitempty"               yaml:"in,omitempty"`               // Location: query, header (for apiKey)
	Flow             string                 `json:"flow,omitempty"             yaml:"flow,omitempty"`             // OAuth2 flow: implicit, password, application, accessCode
	AuthorizationURL string                 `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"` // Authorization URL (for oauth2)
	TokenURL         string                 `json:"tokenUrl,omitempty"         yaml:"tokenUrl,omitempty"`         // Token URL (for oauth2)
	Scopes           map[string]string      `json:"scopes,omitempty"           yaml:"scopes,omitempty"`           // Available scopes (for oauth2)
	Extensions       map[string]interface{} `json:"-"                          yaml:"-"`                          // Custom extensions (x-*)
}

// SecurityRequirement lists the required security schemes for an operation.
type SecurityRequirement map[string][]string

// Tag allows adding metadata to a single tag.
type Tag struct {
	Name         string        `json:"name"                   yaml:"name"`                   // REQUIRED. Tag name
	Description  string        `json:"description,omitempty"  yaml:"description,omitempty"`  // Tag description
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"` // Additional external documentation
}

// ExternalDocs allows referencing an external resource for extended documentation.
type ExternalDocs struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"` // Documentation description
	URL         string `json:"url"                   yaml:"url"`                   // REQUIRED. Documentation URL
}
