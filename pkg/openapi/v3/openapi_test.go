package v3

import (
	"encoding/json"
	"testing"
)

// Test helper to verify JSON marshaling.
func verifyJSONMarshal(t *testing.T, v interface{}) {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Errorf("json.Marshal() error = %v", err)
		return
	}
	if len(data) == 0 {
		t.Error("json.Marshal() returned empty data")
	}
}

func TestNewOpenAPI(t *testing.T) {
	t.Parallel()
	api := &OpenAPI{
		OpenAPI: "3.1.0",
	}
	if api.OpenAPI != "3.1.0" {
		t.Errorf("OpenAPI version = %q, want %q", api.OpenAPI, "3.1.0")
	}
}

func TestOpenAPIGetVersion(t *testing.T) {
	t.Parallel()
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	version := api.GetVersion()
	if version != "3.1.0" {
		t.Errorf("GetVersion() = %v, want 3.1.0", version)
	}
}

func TestOpenAPIGetTitle(t *testing.T) {
	t.Parallel()
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:   "My API",
			Version: "1.0.0",
		},
	}

	title := api.GetTitle()
	if title != "My API" {
		t.Errorf("GetTitle() = %v, want 'My API'", title)
	}
}

func TestOpenAPIGetInfo(t *testing.T) {
	t.Parallel()
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:       "Test API",
			Description: "Test Description",
			Version:     "2.0.0",
			Summary:     "API Summary",
		},
	}

	info := api.GetInfo()
	if info == nil {
		t.Fatal("GetInfo() returned nil")
	}

	apiInfo, ok := info.(Info)
	if !ok {
		t.Fatal("GetInfo() did not return Info type")
	}

	if apiInfo.Title != "Test API" {
		t.Errorf("Info.Title = %v, want 'Test API'", apiInfo.Title)
	}

	if apiInfo.Summary != "API Summary" {
		t.Errorf("Info.Summary = %v, want 'API Summary'", apiInfo.Summary)
	}
}

func TestOpenAPIValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		api     *OpenAPI
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid openapi",
			api: &OpenAPI{
				OpenAPI: "3.1.0",
				Info: Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
			},
			wantErr: false,
		},
		{
			name: "missing openapi version",
			api: &OpenAPI{
				Info: Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
			},
			wantErr: true,
			errMsg:  "openapi version is required",
		},
		{
			name: "missing title",
			api: &OpenAPI{
				OpenAPI: "3.1.0",
				Info: Info{
					Version: "1.0.0",
				},
			},
			wantErr: true,
			errMsg:  "info.title is required",
		},
		{
			name: "missing info version",
			api: &OpenAPI{
				OpenAPI: "3.1.0",
				Info: Info{
					Title: "Test API",
				},
			},
			wantErr: true,
			errMsg:  "info.version is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.api.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !containsString(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			}
		})
	}
}

func TestOpenAPIMarshalJSON(t *testing.T) {
	t.Parallel()
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Servers: []Server{
			{
				URL:         "https://api.example.com",
				Description: "Production server",
			},
		},
	}

	data, err := api.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if result["openapi"] != "3.1.0" {
		t.Errorf("openapi field = %v, want '3.1.0'", result["openapi"])
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstr(s, substr)))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestInfoFields(t *testing.T) {
	t.Parallel()
	info := &Info{
		Title:       "Test API",
		Description: "Test Description",
		Version:     "1.0.0",
		Summary:     "Short summary",
	}
	if info.Title != "Test API" {
		t.Errorf("Info.Title = %q, want %q", info.Title, "Test API")
	}
	if info.Summary != "Short summary" {
		t.Errorf("Info.Summary = %q, want %q", info.Summary, "Short summary")
	}
	if info.Version != "1.0.0" {
		t.Errorf("Info.Version = %q, want %q", info.Version, "1.0.0")
	}
}

func TestContactFields(t *testing.T) {
	t.Parallel()
	contact := &Contact{
		Name:  "Test Contact",
		Email: "test@example.com",
		URL:   "https://example.com",
	}
	if contact.Name != "Test Contact" {
		t.Errorf("Contact.Name = %q, want %q", contact.Name, "Test Contact")
	}
	if contact.Email != "test@example.com" {
		t.Errorf("Contact.Email = %q, want %q", contact.Email, "test@example.com")
	}
	if contact.URL != "https://example.com" {
		t.Errorf("Contact.URL = %q, want %q", contact.URL, "https://example.com")
	}
}

func TestLicenseFields(t *testing.T) {
	t.Parallel()
	license := &License{
		Name:       "MIT",
		URL:        "https://opensource.org/licenses/MIT",
		Identifier: "MIT",
	}
	if license.Name != "MIT" {
		t.Errorf("License.Name = %q, want %q", license.Name, "MIT")
	}
	if license.Identifier != "MIT" {
		t.Errorf("License.Identifier = %q, want %q", license.Identifier, "MIT")
	}
}

func TestServerFields(t *testing.T) {
	t.Parallel()
	server := &Server{
		URL:         "https://api.example.com",
		Description: "Production Server",
	}
	if server.URL != "https://api.example.com" {
		t.Errorf("Server.URL = %q, want %q", server.URL, "https://api.example.com")
	}
}

func TestPathItemOperations(t *testing.T) {
	pathItem := &PathItem{
		Get:  &Operation{OperationID: "getTest"},
		Post: &Operation{OperationID: "postTest"},
	}
	if pathItem.Get == nil {
		t.Error("PathItem.Get should not be nil")
	}
	if pathItem.Get.OperationID != "getTest" {
		t.Errorf("PathItem.Get.OperationID = %q, want %q", pathItem.Get.OperationID, "getTest")
	}
}

func TestOperationFields(t *testing.T) {
	op := &Operation{
		OperationID: "testOp",
		Summary:     "Test Operation",
		Tags:        []string{"test"},
	}
	if op.OperationID != "testOp" {
		t.Errorf("Operation.OperationID = %q, want %q", op.OperationID, "testOp")
	}
}

func TestParameterFields(t *testing.T) {
	param := &Parameter{
		Name:     "id",
		In:       "path",
		Required: true,
	}
	if param.Name != "id" {
		t.Errorf("Parameter.Name = %q, want %q", param.Name, "id")
	}
	if !param.Required {
		t.Error("Parameter.Required should be true")
	}
}

func TestSchemaTypes(t *testing.T) {
	schema := &Schema{
		Type:   "string",
		Format: "email",
	}
	if schema.Type != "string" {
		t.Errorf("Schema.Type = %q, want %q", schema.Type, "string")
	}
}

func TestComponentsStructure(t *testing.T) {
	t.Parallel()
	components := &Components{
		Schemas: make(map[string]*Schema),
	}
	components.Schemas["User"] = &Schema{Type: "object"}
	if len(components.Schemas) != 1 {
		t.Errorf("Components.Schemas length = %d, want 1", len(components.Schemas))
	}
}

func TestSchemaMarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		schema  *Schema
		wantErr bool
	}{
		{
			name: "simple schema without extensions",
			schema: &Schema{
				Type:        "string",
				Description: "A simple string",
			},
			wantErr: false,
		},
		{
			name: "schema with extensions",
			schema: &Schema{
				Type: "object",
				Extensions: map[string]interface{}{
					"x-custom":     "value",
					"x-另一个feature": true,
				},
			},
			wantErr: false,
		},
		{
			name: "complex schema",
			schema: &Schema{
				Type:     "object",
				Required: []string{"id", "name"},
				Properties: map[string]*Schema{
					"id":   {Type: "integer"},
					"name": {Type: "string"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(data) == 0 {
				t.Error("MarshalJSON() returned empty data")
			}

			// Verify it can be unmarshalled back
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Errorf("Failed to unmarshal result: %v", err)
			}
		})
	}
}

func TestOperationMarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		operation *Operation
		wantErr   bool
	}{
		{
			name: "operation without extensions",
			operation: &Operation{
				Summary:     "Test operation",
				OperationID: "testOp",
			},
			wantErr: false,
		},
		{
			name: "operation with extensions",
			operation: &Operation{
				Summary: "Test with extensions",
				Extensions: map[string]interface{}{
					"x-code-samples": []string{"sample1", "sample2"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.operation)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(data) == 0 {
				t.Error("MarshalJSON() returned empty data")
			}
		})
	}
}

func TestServerVariable(t *testing.T) {
	t.Parallel()
	sv := &ServerVariable{
		Default:     "v1",
		Enum:        []string{"v1", "v2"},
		Description: "API version",
	}

	if sv.Default != "v1" {
		t.Errorf("ServerVariable.Default = %q, want %q", sv.Default, "v1")
	}
	if len(sv.Enum) != 2 {
		t.Errorf("ServerVariable.Enum length = %d, want 2", len(sv.Enum))
	}
}

func TestRequestBody(t *testing.T) {
	t.Parallel()
	rb := &RequestBody{
		Description: "User object",
		Required:    true,
		Content: map[string]*MediaType{
			"application/json": {
				Schema: &Schema{Type: "object"},
			},
		},
	}

	if !rb.Required {
		t.Error("RequestBody.Required should be true")
	}
	if rb.Content == nil {
		t.Fatal("RequestBody.Content should not be nil")
	}
	if _, ok := rb.Content["application/json"]; !ok {
		t.Error("RequestBody should have application/json content")
	}
}

func TestMediaType(t *testing.T) {
	t.Parallel()
	mt := &MediaType{
		Schema: &Schema{Type: "object"},
		Example: map[string]interface{}{
			"id":   1,
			"name": "test",
		},
	}

	if mt.Schema == nil {
		t.Error("MediaType.Schema should not be nil")
	}
	if mt.Example == nil {
		t.Error("MediaType.Example should not be nil")
	}
}

func TestResponse(t *testing.T) {
	t.Parallel()
	resp := &Response{
		Description: "Success response",
		Content: map[string]*MediaType{
			"application/json": {
				Schema: &Schema{Type: "object"},
			},
		},
	}

	if resp.Description != "Success response" {
		t.Errorf("Response.Description = %q, want %q", resp.Description, "Success response")
	}
	if len(resp.Content) != 1 {
		t.Errorf("Response.Content length = %d, want 1", len(resp.Content))
	}
}

func TestHeader(t *testing.T) {
	t.Parallel()
	header := &Header{
		Description: "Request ID",
		Schema:      &Schema{Type: "string"},
		Required:    true,
	}

	if !header.Required {
		t.Error("Header.Required should be true")
	}
	if header.Schema == nil {
		t.Error("Header.Schema should not be nil")
	}
}

func TestSecurityScheme(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		scheme *SecurityScheme
	}{
		{
			name: "API Key",
			scheme: &SecurityScheme{
				Type: "apiKey",
				Name: "api_key",
				In:   "header",
			},
		},
		{
			name: "HTTP Bearer",
			scheme: &SecurityScheme{
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		},
		{
			name: "OAuth2",
			scheme: &SecurityScheme{
				Type: "oauth2",
				Flows: &OAuthFlows{
					Implicit: &OAuthFlow{
						AuthorizationURL: "https://example.com/oauth/authorize",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.scheme.Type == "" {
				t.Error("SecurityScheme.Type should not be empty")
			}
		})
	}
}

func TestTag(t *testing.T) {
	t.Parallel()
	tag := &Tag{
		Name:        "users",
		Description: "User operations",
		ExternalDocs: &ExternalDocs{
			Description: "Find more info",
			URL:         "https://example.com/docs",
		},
	}

	if tag.Name != "users" {
		t.Errorf("Tag.Name = %q, want %q", tag.Name, "users")
	}
	if tag.ExternalDocs == nil {
		t.Error("Tag.ExternalDocs should not be nil")
	}
}

func TestExternalDocs(t *testing.T) {
	t.Parallel()
	docs := &ExternalDocs{
		Description: "More information",
		URL:         "https://docs.example.com",
	}

	if docs.URL != "https://docs.example.com" {
		t.Errorf("ExternalDocs.URL = %q, want %q", docs.URL, "https://docs.example.com")
	}
}

func TestDiscriminator(t *testing.T) {
	t.Parallel()
	disc := &Discriminator{
		PropertyName: "type",
		Mapping: map[string]string{
			"dog": "#/components/schemas/Dog",
			"cat": "#/components/schemas/Cat",
		},
	}

	if disc.PropertyName != "type" {
		t.Errorf("Discriminator.PropertyName = %q, want %q", disc.PropertyName, "type")
	}
	if len(disc.Mapping) != 2 {
		t.Errorf("Discriminator.Mapping length = %d, want 2", len(disc.Mapping))
	}
}

func TestXML(t *testing.T) {
	t.Parallel()
	xml := &XML{
		Name:      "user",
		Namespace: "http://example.com/schema",
		Prefix:    "ex",
		Attribute: false,
		Wrapped:   true,
	}

	if xml.Name != "user" {
		t.Errorf("XML.Name = %q, want %q", xml.Name, "user")
	}
	if !xml.Wrapped {
		t.Error("XML.Wrapped should be true")
	}
}

func TestExample(t *testing.T) {
	t.Parallel()
	example := &Example{
		Summary:     "A user example",
		Description: "Example of a user object",
		Value: map[string]interface{}{
			"id":   1,
			"name": "John Doe",
		},
	}

	if example.Summary != "A user example" {
		t.Errorf("Example.Summary = %q, want %q", example.Summary, "A user example")
	}
	if example.Value == nil {
		t.Error("Example.Value should not be nil")
	}
}

func TestLink(t *testing.T) {
	t.Parallel()
	link := &Link{
		OperationID: "getUserById",
		Parameters: map[string]interface{}{
			"userId": "$response.body#/id",
		},
		Description: "Link to user details",
	}

	if link.OperationID != "getUserById" {
		t.Errorf("Link.OperationID = %q, want %q", link.OperationID, "getUserById")
	}
	if len(link.Parameters) == 0 {
		t.Error("Link.Parameters should not be empty")
	}
}

func TestSchemaValidation(t *testing.T) {
	t.Parallel()
	schema := &Schema{
		Type:      "string",
		MinLength: 3,
		MaxLength: 50,
		Pattern:   "^[a-zA-Z]+$",
	}

	if schema.MinLength != 3 {
		t.Errorf("Schema.MinLength = %d, want 3", schema.MinLength)
	}
	if schema.MaxLength != 50 {
		t.Errorf("Schema.MaxLength = %d, want 50", schema.MaxLength)
	}
}

func TestSchemaNumericValidation(t *testing.T) {
	t.Parallel()
	schema := &Schema{
		Type:       "number",
		Minimum:    0,
		Maximum:    100,
		MultipleOf: 5,
	}

	if schema.Minimum != 0 {
		t.Errorf("Schema.Minimum = %f, want 0", schema.Minimum)
	}
	if schema.Maximum != 100 {
		t.Errorf("Schema.Maximum = %f, want 100", schema.Maximum)
	}
	if schema.MultipleOf != 5 {
		t.Errorf("Schema.MultipleOf = %f, want 5", schema.MultipleOf)
	}
}

func TestSchemaArrayValidation(t *testing.T) {
	t.Parallel()
	schema := &Schema{
		Type:        "array",
		Items:       &Schema{Type: "string"},
		MinItems:    1,
		MaxItems:    10,
		UniqueItems: true,
	}

	if schema.MinItems != 1 {
		t.Errorf("Schema.MinItems = %d, want 1", schema.MinItems)
	}
	if schema.MaxItems != 10 {
		t.Errorf("Schema.MaxItems = %d, want 10", schema.MaxItems)
	}
	if !schema.UniqueItems {
		t.Error("Schema.UniqueItems should be true")
	}
}

func TestSchemaObjectValidation(t *testing.T) {
	t.Parallel()
	schema := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"id":   {Type: "integer"},
			"name": {Type: "string"},
		},
		Required:      []string{"id"},
		MinProperties: 1,
		MaxProperties: 10,
	}

	if len(schema.Required) != 1 {
		t.Errorf("Schema.Required length = %d, want 1", len(schema.Required))
	}
	if schema.MinProperties != 1 {
		t.Errorf("Schema.MinProperties = %d, want 1", schema.MinProperties)
	}
}

func TestSchemaComposition(t *testing.T) {
	t.Parallel()
	schema := &Schema{
		AllOf: []Schema{
			{Type: "object", Properties: map[string]*Schema{"id": {Type: "integer"}}},
			{Type: "object", Properties: map[string]*Schema{"name": {Type: "string"}}},
		},
	}

	if len(schema.AllOf) != 2 {
		t.Errorf("Schema.AllOf length = %d, want 2", len(schema.AllOf))
	}
}

func TestEncoding(t *testing.T) {
	t.Parallel()
	encoding := &Encoding{
		ContentType:   "application/json",
		Style:         "form",
		Explode:       true,
		AllowReserved: false,
	}

	if encoding.ContentType != "application/json" {
		t.Errorf("Encoding.ContentType = %q, want %q", encoding.ContentType, "application/json")
	}
	if !encoding.Explode {
		t.Error("Encoding.Explode should be true")
	}
}

func TestOAuthFlows(t *testing.T) {
	t.Parallel()
	flows := &OAuthFlows{
		AuthorizationCode: &OAuthFlow{
			AuthorizationURL: "https://example.com/oauth/authorize",
			TokenURL:         "https://example.com/oauth/token",
			Scopes: map[string]string{
				"read":  "Read access",
				"write": "Write access",
			},
		},
	}

	if flows.AuthorizationCode == nil {
		t.Fatal("OAuthFlows.AuthorizationCode should not be nil")
	}
	if len(flows.AuthorizationCode.Scopes) != 2 {
		t.Errorf("Scopes length = %d, want 2", len(flows.AuthorizationCode.Scopes))
	}
}

func TestCompleteOpenAPIStructure(t *testing.T) {
	t.Parallel()
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:   "Complete API",
			Version: "1.0.0",
			Contact: &Contact{
				Name:  "API Team",
				Email: "api@example.com",
			},
			License: &License{
				Name: "Apache 2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		Servers: []Server{
			{
				URL:         "https://api.example.com",
				Description: "Production",
			},
		},
		Paths: Paths{
			"/users": &PathItem{
				Get: &Operation{
					Summary:     "List users",
					OperationID: "listUsers",
					Responses: Responses{
						"200": &Response{
							Description: "Success",
						},
					},
				},
			},
		},
		Components: &Components{
			Schemas: map[string]*Schema{
				"User": {
					Type: "object",
					Properties: map[string]*Schema{
						"id":   {Type: "integer"},
						"name": {Type: "string"},
					},
				},
			},
			SecuritySchemes: map[string]*SecurityScheme{
				"bearerAuth": {
					Type:   "http",
					Scheme: "bearer",
				},
			},
		},
		Tags: []Tag{
			{
				Name:        "users",
				Description: "User operations",
			},
		},
	}

	if api.OpenAPI != "3.1.0" {
		t.Errorf("OpenAPI version = %q, want %q", api.OpenAPI, "3.1.0")
	}
	if len(api.Servers) != 1 {
		t.Errorf("Servers length = %d, want 1", len(api.Servers))
	}
	if len(api.Paths) != 1 {
		t.Errorf("Paths length = %d, want 1", len(api.Paths))
	}
	if len(api.Components.Schemas) != 1 {
		t.Errorf("Schemas length = %d, want 1", len(api.Components.Schemas))
	}
	if len(api.Tags) != 1 {
		t.Errorf("Tags length = %d, want 1", len(api.Tags))
	}
}
