package v2

import (
	"encoding/json"
	"testing"
)

func TestNewSwagger(t *testing.T) {
	swagger := &Swagger{
		Swagger: "2.0",
		Info: Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	if swagger.Swagger != "2.0" {
		t.Errorf("Swagger version = %v, want 2.0", swagger.Swagger)
	}

	if swagger.Info.Title != "Test API" {
		t.Errorf("Info.Title = %v, want Test API", swagger.Info.Title)
	}
}

func TestSwaggerGetVersion(t *testing.T) {
	swagger := &Swagger{
		Swagger: "2.0",
		Info: Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	version := swagger.GetVersion()
	if version != "2.0" {
		t.Errorf("GetVersion() = %v, want 2.0", version)
	}
}

func TestSwaggerGetTitle(t *testing.T) {
	swagger := &Swagger{
		Swagger: "2.0",
		Info: Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	title := swagger.GetTitle()
	if title != "Test API" {
		t.Errorf("GetTitle() = %v, want Test API", title)
	}
}

func TestSwaggerGetInfo(t *testing.T) {
	swagger := &Swagger{
		Swagger: "2.0",
		Info: Info{
			Title:       "Test API",
			Description: "Test Description",
			Version:     "1.0.0",
		},
	}

	info := swagger.GetInfo()
	if info == nil {
		t.Fatal("GetInfo() returned nil")
	}

	swaggerInfo, ok := info.(Info)
	if !ok {
		t.Fatal("GetInfo() did not return Info type")
	}

	if swaggerInfo.Title != "Test API" {
		t.Errorf("Info.Title = %v, want Test API", swaggerInfo.Title)
	}

	if swaggerInfo.Description != "Test Description" {
		t.Errorf("Info.Description = %v, want Test Description", swaggerInfo.Description)
	}
}

func TestInfoFields(t *testing.T) {
	info := Info{
		Title:          "My API",
		Description:    "My API Description",
		TermsOfService: "https://example.com/terms",
		Version:        "2.0.0",
		Contact: &Contact{
			Name:  "API Support",
			URL:   "https://example.com/support",
			Email: "support@example.com",
		},
		License: &License{
			Name: "MIT",
			URL:  "https://opensource.org/licenses/MIT",
		},
	}

	if info.Title != "My API" {
		t.Errorf("Title = %v, want My API", info.Title)
	}

	if info.Contact.Name != "API Support" {
		t.Errorf("Contact.Name = %v, want API Support", info.Contact.Name)
	}

	if info.License.Name != "MIT" {
		t.Errorf("License.Name = %v, want MIT", info.License.Name)
	}
}

func TestSwaggerValidate(t *testing.T) {
	tests := []struct {
		name    string
		swagger *Swagger
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid swagger",
			swagger: &Swagger{
				Swagger: "2.0",
				Info: Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
				Paths: make(Paths),
			},
			wantErr: false,
		},
		{
			name: "invalid swagger version",
			swagger: &Swagger{
				Swagger: "3.0",
				Info: Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
				Paths: make(Paths),
			},
			wantErr: true,
			errMsg:  "invalid swagger version",
		},
		{
			name: "missing title",
			swagger: &Swagger{
				Swagger: "2.0",
				Info: Info{
					Version: "1.0.0",
				},
				Paths: make(Paths),
			},
			wantErr: true,
			errMsg:  "info.title is required",
		},
		{
			name: "missing version",
			swagger: &Swagger{
				Swagger: "2.0",
				Info: Info{
					Title: "Test API",
				},
				Paths: make(Paths),
			},
			wantErr: true,
			errMsg:  "info.version is required",
		},
		{
			name: "missing paths",
			swagger: &Swagger{
				Swagger: "2.0",
				Info: Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
			},
			wantErr: true,
			errMsg:  "paths is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.swagger.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
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

func TestSwaggerMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		swagger *Swagger
		wantKey string
	}{
		{
			name: "swagger without extensions",
			swagger: &Swagger{
				Swagger: "2.0",
				Info: Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
			},
			wantKey: "swagger",
		},
		{
			name: "swagger with extensions",
			swagger: &Swagger{
				Swagger: "2.0",
				Info: Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
				Extensions: map[string]interface{}{
					"x-custom": "value",
					"x-logo": map[string]string{
						"url": "https://example.com/logo.png",
					},
				},
			},
			wantKey: "x-custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.swagger.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Unmarshal error = %v", err)
			}

			if _, ok := result[tt.wantKey]; !ok {
				t.Errorf("MarshalJSON() missing key %q", tt.wantKey)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestPathItem(t *testing.T) {
	pathItem := &PathItem{
		Get: &Operation{
			Summary:     "Get item",
			Description: "Retrieves an item",
			OperationID: "getItem",
			Tags:        []string{"items"},
			Responses: map[string]*Response{
				"200": {
					Description: "Success",
				},
			},
		},
		Post: &Operation{
			Summary:     "Create item",
			Description: "Creates a new item",
			OperationID: "createItem",
			Parameters: []*Parameter{
				{
					Name:        "body",
					In:          "body",
					Description: "Item to create",
					Required:    true,
					Schema: &Schema{
						Type: "object",
					},
				},
			},
			Responses: map[string]*Response{
				"201": {
					Description: "Created",
				},
			},
		},
	}

	if pathItem.Get.Summary != "Get item" {
		t.Errorf("Get.Summary = %v, want Get item", pathItem.Get.Summary)
	}

	if pathItem.Post.OperationID != "createItem" {
		t.Errorf("Post.OperationID = %v, want createItem", pathItem.Post.OperationID)
	}

	if len(pathItem.Post.Parameters) != 1 {
		t.Errorf("len(Post.Parameters) = %v, want 1", len(pathItem.Post.Parameters))
	}
}

func TestSchema(t *testing.T) {
	schema := &Schema{
		Type:        "object",
		Description: "User object",
		Required:    []string{"id", "name"},
		Properties: map[string]*Schema{
			"id": {
				Type:        "integer",
				Format:      "int64",
				Description: "User ID",
			},
			"name": {
				Type:        "string",
				Description: "User name",
			},
			"email": {
				Type:   "string",
				Format: "email",
			},
		},
	}

	if schema.Type != "object" {
		t.Errorf("Type = %v, want object", schema.Type)
	}

	if len(schema.Required) != 2 {
		t.Errorf("len(Required) = %v, want 2", len(schema.Required))
	}

	if len(schema.Properties) != 3 {
		t.Errorf("len(Properties) = %v, want 3", len(schema.Properties))
	}

	if schema.Properties["id"].Format != "int64" {
		t.Errorf("Properties[id].Format = %v, want int64", schema.Properties["id"].Format)
	}
}

func TestParameter(t *testing.T) {
	tests := []struct {
		name  string
		param *Parameter
	}{
		{
			name: "query parameter",
			param: &Parameter{
				Name:        "limit",
				In:          "query",
				Description: "Maximum number of results",
				Required:    false,
				Type:        "integer",
				Format:      "int32",
				Default:     10,
			},
		},
		{
			name: "path parameter",
			param: &Parameter{
				Name:        "id",
				In:          "path",
				Description: "Item ID",
				Required:    true,
				Type:        "string",
			},
		},
		{
			name: "body parameter",
			param: &Parameter{
				Name:        "body",
				In:          "body",
				Description: "Request body",
				Required:    true,
				Schema: &Schema{
					Type: "object",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.param.Name == "" {
				t.Error("Parameter name should not be empty")
			}

			if tt.param.In == "" {
				t.Error("Parameter 'in' should not be empty")
			}

			if tt.param.In == "body" && tt.param.Schema == nil {
				t.Error("Body parameter should have schema")
			}
		})
	}
}

func TestResponse(t *testing.T) {
	response := &Response{
		Description: "Successful response",
		Schema: &Schema{
			Type: "array",
			Items: &Schema{
				Type: "object",
				Properties: map[string]*Schema{
					"id":   {Type: "integer"},
					"name": {Type: "string"},
				},
			},
		},
		Headers: map[string]*Header{
			"X-Rate-Limit": {
				Description: "Rate limit",
				Type:        "integer",
			},
		},
	}

	if response.Description != "Successful response" {
		t.Errorf("Description = %v, want Successful response", response.Description)
	}

	if response.Schema.Type != "array" {
		t.Errorf("Schema.Type = %v, want array", response.Schema.Type)
	}

	if len(response.Headers) != 1 {
		t.Errorf("len(Headers) = %v, want 1", len(response.Headers))
	}
}

func TestSecurityScheme(t *testing.T) {
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
			name: "OAuth2",
			scheme: &SecurityScheme{
				Type:             "oauth2",
				Flow:             "implicit",
				AuthorizationURL: "https://example.com/oauth/authorize",
				Scopes: map[string]string{
					"read":  "Read access",
					"write": "Write access",
				},
			},
		},
		{
			name: "Basic Auth",
			scheme: &SecurityScheme{
				Type: "basic",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.scheme.Type == "" {
				t.Error("SecurityScheme type should not be empty")
			}

			if tt.scheme.Type == "apiKey" && tt.scheme.Name == "" {
				t.Error("API Key scheme should have name")
			}

			if tt.scheme.Type == "apiKey" && tt.scheme.In == "" {
				t.Error("API Key scheme should have 'in' field")
			}

			if tt.scheme.Type == "oauth2" && tt.scheme.Flow == "" {
				t.Error("OAuth2 scheme should have flow")
			}
		})
	}
}

func TestSwaggerJSONSerialization(t *testing.T) {
	swagger := &Swagger{
		Swagger: "2.0",
		Info: Info{
			Title:       "Test API",
			Description: "API Description",
			Version:     "1.0.0",
		},
		Host:     "api.example.com",
		BasePath: "/v1",
		Schemes:  []string{"https"},
		Paths: map[string]*PathItem{
			"/users": {
				Get: &Operation{
					Summary: "List users",
					Responses: map[string]*Response{
						"200": {
							Description: "Success",
						},
					},
				},
			},
		},
	}

	// Serialize to JSON
	data, err := json.Marshal(swagger)
	if err != nil {
		t.Fatalf("Failed to marshal Swagger: %v", err)
	}

	// Deserialize from JSON
	var unmarshaled Swagger
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Swagger: %v", err)
	}

	if unmarshaled.Swagger != "2.0" {
		t.Errorf("Swagger version = %v, want 2.0", unmarshaled.Swagger)
	}

	if unmarshaled.Info.Title != "Test API" {
		t.Errorf("Info.Title = %v, want Test API", unmarshaled.Info.Title)
	}

	if unmarshaled.Host != "api.example.com" {
		t.Errorf("Host = %v, want api.example.com", unmarshaled.Host)
	}
}

func TestTag(t *testing.T) {
	tag := Tag{
		Name:        "users",
		Description: "User operations",
		ExternalDocs: &ExternalDocs{
			Description: "Find more info here",
			URL:         "https://example.com/docs",
		},
	}

	if tag.Name != "users" {
		t.Errorf("Name = %v, want users", tag.Name)
	}

	if tag.ExternalDocs.URL != "https://example.com/docs" {
		t.Errorf("ExternalDocs.URL = %v, want https://example.com/docs", tag.ExternalDocs.URL)
	}
}
