package converter

import (
	"strings"
	"testing"

	swagger "github.com/fsvxavier/nexs-swag/pkg/openapi/v2"
	openapi "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

func TestNewConverter(t *testing.T) {
	conv := New()

	if conv == nil {
		t.Fatal("New() returned nil")
	}

	warnings := conv.GetWarnings()
	if warnings == nil {
		t.Error("warnings slice should be initialized")
	}

	if len(warnings) != 0 {
		t.Errorf("warnings should be empty initially, got %d", len(warnings))
	}
}

func TestConvertToV2Basic(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:       "Test API",
			Description: "Test Description",
			Version:     "1.0.0",
		},
		Servers: []openapi.Server{
			{
				URL:         "https://api.example.com/v1",
				Description: "Production server",
			},
		},
		Paths: make(map[string]*openapi.PathItem),
	}

	conv := New()
	result, err := conv.ConvertToV2(spec)

	if err != nil {
		t.Fatalf("ConvertToV2() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("ConvertToV2() returned nil")
	}

	if result.Swagger != "2.0" {
		t.Errorf("Swagger version = %v, want 2.0", result.Swagger)
	}

	if result.Info.Title != "Test API" {
		t.Errorf("Info.Title = %v, want Test API", result.Info.Title)
	}

	if result.Host != "api.example.com" {
		t.Errorf("Host = %v, want api.example.com", result.Host)
	}

	if result.BasePath != "/v1" {
		t.Errorf("BasePath = %v, want /v1", result.BasePath)
	}

	if len(result.Schemes) == 0 || result.Schemes[0] != "https" {
		t.Errorf("Schemes = %v, want [https]", result.Schemes)
	}
}

func TestConvertToV3Basic(t *testing.T) {
	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:       "Test API",
			Description: "Test Description",
			Version:     "1.0.0",
		},
		Host:     "api.example.com",
		BasePath: "/v1",
		Schemes:  []string{"https"},
		Paths:    make(map[string]*swagger.PathItem),
	}

	conv := New()
	result, err := conv.ConvertToV3(spec)

	if err != nil {
		t.Fatalf("ConvertToV3() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("ConvertToV3() returned nil")
	}

	if result.OpenAPI != "3.1.0" {
		t.Errorf("OpenAPI version = %v, want 3.1.0", result.OpenAPI)
	}

	if result.Info.Title != "Test API" {
		t.Errorf("Info.Title = %v, want Test API", result.Info.Title)
	}

	if len(result.Servers) == 0 {
		t.Fatal("Servers should not be empty")
	}

	if result.Servers[0].URL != "https://api.example.com/v1" {
		t.Errorf("Server URL = %v, want https://api.example.com/v1", result.Servers[0].URL)
	}
}

func TestConvertServerToHostBasePath(t *testing.T) {
	tests := []struct {
		name         string
		serverURL    string
		wantHost     string
		wantBasePath string
		wantSchemes  []string
	}{
		{
			name:         "https with path",
			serverURL:    "https://api.example.com/v1",
			wantHost:     "api.example.com",
			wantBasePath: "/v1",
			wantSchemes:  []string{"https"},
		},
		{
			name:         "http no path",
			serverURL:    "http://api.example.com",
			wantHost:     "api.example.com",
			wantBasePath: "/",
			wantSchemes:  []string{"http"},
		},
		{
			name:         "with port",
			serverURL:    "https://api.example.com:8080/api",
			wantHost:     "api.example.com:8080",
			wantBasePath: "/api",
			wantSchemes:  []string{"https"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &openapi.OpenAPI{
				OpenAPI: "3.1.0",
				Info:    openapi.Info{Title: "Test", Version: "1.0.0"},
				Servers: []openapi.Server{{URL: tt.serverURL}},
				Paths:   make(map[string]*openapi.PathItem),
			}

			conv := New()
			result, err := conv.ConvertToV2(spec)
			if err != nil {
				t.Fatalf("ConvertToV2() error = %v", err)
			}

			if result.Host != tt.wantHost {
				t.Errorf("Host = %v, want %v", result.Host, tt.wantHost)
			}

			if result.BasePath != tt.wantBasePath {
				t.Errorf("BasePath = %v, want %v", result.BasePath, tt.wantBasePath)
			}

			if len(result.Schemes) == 0 || result.Schemes[0] != tt.wantSchemes[0] {
				t.Errorf("Schemes = %v, want %v", result.Schemes, tt.wantSchemes)
			}
		})
	}
}

func TestConvertHostBasePathToServer(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		basePath string
		schemes  []string
		wantURL  string
	}{
		{
			name:     "https with path",
			host:     "api.example.com",
			basePath: "/v1",
			schemes:  []string{"https"},
			wantURL:  "https://api.example.com/v1",
		},
		{
			name:     "http no path",
			host:     "api.example.com",
			basePath: "",
			schemes:  []string{"http"},
			wantURL:  "http://api.example.com",
		},
		{
			name:     "multiple schemes",
			host:     "api.example.com",
			basePath: "/api",
			schemes:  []string{"http", "https"},
			wantURL:  "http://api.example.com/api", // First scheme is used
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &swagger.Swagger{
				Swagger:  "2.0",
				Info:     swagger.Info{Title: "Test", Version: "1.0.0"},
				Host:     tt.host,
				BasePath: tt.basePath,
				Schemes:  tt.schemes,
				Paths:    make(map[string]*swagger.PathItem),
			}

			conv := New()
			result, err := conv.ConvertToV3(spec)
			if err != nil {
				t.Fatalf("ConvertToV3() error = %v", err)
			}

			if len(result.Servers) == 0 {
				t.Fatal("Servers should not be empty")
			}

			if result.Servers[0].URL != tt.wantURL {
				t.Errorf("Server URL = %v, want %v", result.Servers[0].URL, tt.wantURL)
			}
		})
	}
}

func TestConvertWarnings(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI:           "3.1.0",
		Info:              openapi.Info{Title: "Test", Version: "1.0.0"},
		JSONSchemaDialect: "https://json-schema.org/draft/2020-12/schema",
		Paths:             make(map[string]*openapi.PathItem),
	}

	conv := New()
	_, err := conv.ConvertToV2(spec)
	if err != nil {
		t.Fatalf("ConvertToV2() error = %v", err)
	}

	warnings := conv.GetWarnings()
	if len(warnings) == 0 {
		t.Error("Expected warnings for jsonSchemaDialect")
	}

	found := false
	for _, w := range warnings {
		if strings.Contains(w, "jsonSchemaDialect") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected warning about jsonSchemaDialect not found")
	}
}

func TestConvertInfo(t *testing.T) {
	v3Info := openapi.Info{
		Title:          "My API",
		Description:    "API Description",
		TermsOfService: "https://example.com/terms",
		Contact: &openapi.Contact{
			Name:  "Support",
			URL:   "https://example.com",
			Email: "support@example.com",
		},
		License: &openapi.License{
			Name: "MIT",
			URL:  "https://opensource.org/licenses/MIT",
		},
		Version: "2.0.0",
	}

	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info:    v3Info,
		Paths:   make(map[string]*openapi.PathItem),
	}

	conv := New()
	result, err := conv.ConvertToV2(spec)
	if err != nil {
		t.Fatalf("ConvertToV2() error = %v", err)
	}

	if result.Info.Title != v3Info.Title {
		t.Errorf("Title = %v, want %v", result.Info.Title, v3Info.Title)
	}

	if result.Info.Description != v3Info.Description {
		t.Errorf("Description = %v, want %v", result.Info.Description, v3Info.Description)
	}

	if result.Info.Contact.Email != v3Info.Contact.Email {
		t.Errorf("Contact.Email = %v, want %v", result.Info.Contact.Email, v3Info.Contact.Email)
	}

	if result.Info.License.Name != v3Info.License.Name {
		t.Errorf("License.Name = %v, want %v", result.Info.License.Name, v3Info.License.Name)
	}
}

func TestConvertPaths(t *testing.T) {
	v3Paths := map[string]*openapi.PathItem{
		"/users": {
			Get: &openapi.Operation{
				Summary:     "List users",
				Description: "Returns list of users",
				OperationID: "listUsers",
				Tags:        []string{"users"},
				Responses: map[string]*openapi.Response{
					"200": {
						Description: "Success",
					},
				},
			},
		},
	}

	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info:    openapi.Info{Title: "Test", Version: "1.0.0"},
		Paths:   v3Paths,
	}

	conv := New()
	result, err := conv.ConvertToV2(spec)
	if err != nil {
		t.Fatalf("ConvertToV2() error = %v", err)
	}

	if result.Paths == nil {
		t.Fatal("Paths should not be nil")
	}

	userPath, exists := result.Paths["/users"]
	if !exists {
		t.Fatal("/users path not found")
	}

	if userPath.Get == nil {
		t.Fatal("GET operation not found")
	}

	if userPath.Get.Summary != "List users" {
		t.Errorf("GET Summary = %v, want List users", userPath.Get.Summary)
	}
}

func TestConvertComponents(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info:    openapi.Info{Title: "Test", Version: "1.0.0"},
		Paths:   make(map[string]*openapi.PathItem),
		Components: &openapi.Components{
			Schemas: map[string]*openapi.Schema{
				"User": {
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"id":   {Type: "integer"},
						"name": {Type: "string"},
					},
					Required: []string{"id", "name"},
				},
			},
		},
	}

	conv := New()
	result, err := conv.ConvertToV2(spec)
	if err != nil {
		t.Fatalf("ConvertToV2() error = %v", err)
	}

	if result.Definitions == nil {
		t.Fatal("Definitions should not be nil")
	}

	userSchema, exists := result.Definitions["User"]
	if !exists {
		t.Fatal("User schema not found in definitions")
	}

	if userSchema.Type != "object" {
		t.Errorf("User schema Type = %v, want object", userSchema.Type)
	}

	if len(userSchema.Properties) != 2 {
		t.Errorf("len(Properties) = %v, want 2", len(userSchema.Properties))
	}
}

func TestGetWarnings(t *testing.T) {
	conv := New()

	if len(conv.GetWarnings()) != 0 {
		t.Error("Initial warnings should be empty")
	}

	conv.warnings = append(conv.warnings, "test warning")

	warnings := conv.GetWarnings()
	if len(warnings) != 1 {
		t.Errorf("len(warnings) = %v, want 1", len(warnings))
	}

	if warnings[0] != "test warning" {
		t.Errorf("warning = %v, want test warning", warnings[0])
	}
}

func TestClearWarnings(t *testing.T) {
	conv := New()
	conv.warnings = append(conv.warnings, "warning 1", "warning 2")

	if len(conv.GetWarnings()) != 2 {
		t.Error("Should have 2 warnings before clear")
	}

	conv.ClearWarnings()

	if len(conv.GetWarnings()) != 0 {
		t.Error("Warnings should be empty after clear")
	}
}

func TestConvertRoundTrip(t *testing.T) {
	// Create V2 spec
	originalSwagger := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Round Trip API",
			Version: "1.0.0",
		},
		Host:     "api.example.com",
		BasePath: "/v1",
		Schemes:  []string{"https"},
		Paths: map[string]*swagger.PathItem{
			"/test": {
				Get: &swagger.Operation{
					Summary: "Test",
					Responses: map[string]*swagger.Response{
						"200": {Description: "OK"},
					},
				},
			},
		},
	}

	conv := New()

	// Convert V2 -> V3
	spec3, err := conv.ConvertToV3(originalSwagger)
	if err != nil {
		t.Fatalf("V2->V3 conversion error: %v", err)
	}

	// Convert V3 -> V2
	conv2 := New()
	spec2, err := conv2.ConvertToV2(spec3)
	if err != nil {
		t.Fatalf("V3->V2 conversion error: %v", err)
	}

	// Verify basic properties are preserved
	if spec2.Info.Title != originalSwagger.Info.Title {
		t.Errorf("Title = %v, want %v", spec2.Info.Title, originalSwagger.Info.Title)
	}

	if spec2.Host != originalSwagger.Host {
		t.Errorf("Host = %v, want %v", spec2.Host, originalSwagger.Host)
	}

	if spec2.BasePath != originalSwagger.BasePath {
		t.Errorf("BasePath = %v, want %v", spec2.BasePath, originalSwagger.BasePath)
	}
}

func TestConvertNilInputs(t *testing.T) {
	conv := New()

	_, err := conv.ConvertToV2(nil)
	if err == nil {
		t.Error("ConvertToV2(nil) should return error")
	}

	_, err = conv.ConvertToV3(nil)
	if err == nil {
		t.Error("ConvertToV3(nil) should return error")
	}
}

// TestConvertQueryMethodToV2 tests conversion of QUERY method to Swagger 2.0
func TestConvertQueryMethodToV2(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: map[string]*openapi.PathItem{
			"/users": {
				Query: &openapi.Operation{
					Summary:     "Query users",
					Description: "Search users by criteria",
					OperationID: "queryUsers",
					Responses: map[string]*openapi.Response{
						"200": {Description: "Success"},
					},
				},
			},
		},
	}

	conv := New()
	result, err := conv.ConvertToV2(spec)

	if err != nil {
		t.Fatalf("ConvertToV2() error: %v", err)
	}

	// Verify QUERY method triggers warning
	warnings := conv.GetWarnings()
	found := false
	for _, w := range warnings {
		if strings.Contains(w, "QUERY") && strings.Contains(w, "not supported in Swagger 2.0") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected warning about QUERY method not supported in Swagger 2.0")
	}

	// Verify QUERY method is NOT in converted spec
	pathItem := result.Paths["/users"]
	if pathItem != nil {
		// PathItem in Swagger 2.0 doesn't have Query field, so this is just to ensure
		// the conversion completes without panic
		if pathItem.Get != nil {
			t.Error("QUERY should not be converted to GET")
		}
	}
}

// TestConvertSecuritySchemeDeprecatedToV2 tests deprecated field conversion
func TestConvertSecuritySchemeDeprecatedToV2(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Components: &openapi.Components{
			SecuritySchemes: map[string]*openapi.SecurityScheme{
				"OldApiKey": {
					Type:       "apiKey",
					Name:       "X-API-Key",
					In:         "header",
					Deprecated: true,
				},
			},
		},
	}

	conv := New()
	result, err := conv.ConvertToV2(spec)

	if err != nil {
		t.Fatalf("ConvertToV2() error: %v", err)
	}

	// Verify deprecated field triggers warning
	warnings := conv.GetWarnings()
	found := false
	for _, w := range warnings {
		if strings.Contains(w, "deprecated") && strings.Contains(w, "x-deprecated") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected warning about deprecated field conversion to x-deprecated extension")
	}

	// Verify x-deprecated extension is added
	scheme := result.SecurityDefinitions["OldApiKey"]
	if scheme == nil {
		t.Fatal("SecurityScheme OldApiKey not found")
	}

	if scheme.Extensions == nil {
		t.Fatal("Extensions should be set")
	}

	deprecated, ok := scheme.Extensions["x-deprecated"].(bool)
	if !ok || !deprecated {
		t.Error("x-deprecated extension should be true")
	}
}

// TestConvertOAuth2MetadataURLToV2 tests OAuth2MetadataURL conversion
func TestConvertOAuth2MetadataURLToV2(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Components: &openapi.Components{
			SecuritySchemes: map[string]*openapi.SecurityScheme{
				"OAuth2": {
					Type:              "oauth2",
					OAuth2MetadataURL: "https://example.com/.well-known/oauth-authorization-server",
					Flows: &openapi.OAuthFlows{
						AuthorizationCode: &openapi.OAuthFlow{
							AuthorizationURL: "https://example.com/oauth/authorize",
							TokenURL:         "https://example.com/oauth/token",
							Scopes:           map[string]string{"read": "Read access"},
						},
					},
				},
			},
		},
	}

	conv := New()
	_, err := conv.ConvertToV2(spec)

	if err != nil {
		t.Fatalf("ConvertToV2() error: %v", err)
	}

	// Verify warning about OAuth2MetadataURL
	warnings := conv.GetWarnings()
	found := false
	for _, w := range warnings {
		if strings.Contains(w, "OAuth2MetadataURL") && strings.Contains(w, "not supported") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected warning about OAuth2MetadataURL not supported in Swagger 2.0")
	}
}

// TestConvertDeviceAuthorizationFlowToV2 tests device authorization flow conversion
func TestConvertDeviceAuthorizationFlowToV2(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Components: &openapi.Components{
			SecuritySchemes: map[string]*openapi.SecurityScheme{
				"OAuth2": {
					Type: "oauth2",
					Flows: &openapi.OAuthFlows{
						DeviceAuthorization: &openapi.OAuthFlow{
							TokenURL: "https://example.com/oauth/device/token",
							Scopes:   map[string]string{"read": "Read access"},
						},
					},
				},
			},
		},
	}

	conv := New()
	_, err := conv.ConvertToV2(spec)

	if err != nil {
		t.Fatalf("ConvertToV2() error: %v", err)
	}

	// Verify warning about DeviceAuthorization flow
	warnings := conv.GetWarnings()
	found := false
	for _, w := range warnings {
		if strings.Contains(w, "DeviceAuthorization") && strings.Contains(w, "not supported") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected warning about DeviceAuthorization flow not supported in Swagger 2.0")
	}
}

// TestConvertItemSchemaToV2 tests streaming ItemSchema conversion
func TestConvertItemSchemaToV2(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: map[string]*openapi.PathItem{
			"/stream": {
				Get: &openapi.Operation{
					Summary: "Stream data",
					Responses: map[string]*openapi.Response{
						"200": {
							Description: "Streaming response",
							Content: map[string]*openapi.MediaType{
								"text/event-stream": {
									Schema: &openapi.Schema{
										Type: "string",
									},
									ItemSchema: &openapi.Schema{
										Type: "object",
										Properties: map[string]*openapi.Schema{
											"id":   {Type: "integer"},
											"data": {Type: "string"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	conv := New()
	_, err := conv.ConvertToV2(spec)

	if err != nil {
		t.Fatalf("ConvertToV2() error: %v", err)
	}

	// Verify warning about ItemSchema
	warnings := conv.GetWarnings()
	found := false
	for _, w := range warnings {
		if strings.Contains(w, "itemSchema") && strings.Contains(w, "streaming") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected warning about itemSchema for streaming not supported in Swagger 2.0")
	}
}

// TestConvertItemEncodingToV2 tests streaming ItemEncoding conversion
func TestConvertItemEncodingToV2(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: map[string]*openapi.PathItem{
			"/upload": {
				Post: &openapi.Operation{
					Summary: "Upload streaming data",
					RequestBody: &openapi.RequestBody{
						Content: map[string]*openapi.MediaType{
							"multipart/form-data": {
								Schema: &openapi.Schema{
									Type: "object",
								},
								ItemEncoding: map[string]*openapi.Encoding{
									"file": {
										ContentType: "application/octet-stream",
									},
								},
							},
						},
					},
					Responses: map[string]*openapi.Response{
						"200": {Description: "Success"},
					},
				},
			},
		},
	}

	conv := New()
	_, err := conv.ConvertToV2(spec)

	if err != nil {
		t.Fatalf("ConvertToV2() error: %v", err)
	}

	// Verify warning about ItemEncoding
	warnings := conv.GetWarnings()
	found := false
	for _, w := range warnings {
		if strings.Contains(w, "itemEncoding") && strings.Contains(w, "streaming") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected warning about itemEncoding for streaming not supported in Swagger 2.0")
	}
}

// TestMultipleOpenAPI32Features tests multiple 3.2.0 features together
func TestMultipleOpenAPI32Features(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: map[string]*openapi.PathItem{
			"/users": {
				Query: &openapi.Operation{
					Summary: "Query users",
					Responses: map[string]*openapi.Response{
						"200": {
							Description: "Streaming results",
							Content: map[string]*openapi.MediaType{
								"text/event-stream": {
									ItemSchema: &openapi.Schema{Type: "object"},
								},
							},
						},
					},
				},
			},
		},
		Components: &openapi.Components{
			SecuritySchemes: map[string]*openapi.SecurityScheme{
				"OAuth2": {
					Type:              "oauth2",
					Deprecated:        true,
					OAuth2MetadataURL: "https://example.com/.well-known/oauth",
					Flows: &openapi.OAuthFlows{
						DeviceAuthorization: &openapi.OAuthFlow{
							TokenURL: "https://example.com/token",
							Scopes:   map[string]string{"read": "Read"},
						},
					},
				},
			},
		},
	}

	conv := New()
	result, err := conv.ConvertToV2(spec)

	if err != nil {
		t.Fatalf("ConvertToV2() error: %v", err)
	}

	if result == nil {
		t.Fatal("ConvertToV2() returned nil")
	}

	// Verify all warnings are generated
	warnings := conv.GetWarnings()

	expectedWarnings := []string{
		"QUERY",
		"itemSchema",
		"deprecated",
		"OAuth2MetadataURL",
		"DeviceAuthorization",
	}

	for _, expected := range expectedWarnings {
		found := false
		for _, w := range warnings {
			if strings.Contains(w, expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected warning containing %q, got warnings: %v", expected, warnings)
		}
	}

	// Should have at least 5 warnings (one for each feature)
	if len(warnings) < 5 {
		t.Errorf("Expected at least 5 warnings, got %d: %v", len(warnings), warnings)
	}
}

// TestConvertParameter tests the convertParameter function
func TestConvertParameter(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		param    *openapi.Parameter
		expected *swagger.Parameter
	}{
		{
			name:     "nil parameter",
			param:    nil,
			expected: nil,
		},
		{
			name: "simple string parameter",
			param: &openapi.Parameter{
				Name:        "username",
				In:          "query",
				Description: "Username parameter",
				Required:    true,
				Schema: &openapi.Schema{
					Type:   "string",
					Format: "email",
				},
			},
			expected: &swagger.Parameter{
				Name:        "username",
				In:          "query",
				Description: "Username parameter",
				Required:    true,
				Type:        "string",
				Format:      "email",
			},
		},
		{
			name: "deprecated parameter",
			param: &openapi.Parameter{
				Name:       "oldParam",
				In:         "header",
				Deprecated: true,
				Schema: &openapi.Schema{
					Type: "string",
				},
			},
			expected: &swagger.Parameter{
				Name: "oldParam",
				In:   "header",
				Type: "string",
				Extensions: map[string]interface{}{
					"x-deprecated": true,
				},
			},
		},
		{
			name: "integer parameter with constraints",
			param: &openapi.Parameter{
				Name:     "limit",
				In:       "query",
				Required: false,
				Schema: &openapi.Schema{
					Type:    "integer",
					Format:  "int32",
					Minimum: 1,
					Maximum: 100,
					Default: 10,
				},
			},
			expected: &swagger.Parameter{
				Name:    "limit",
				In:      "query",
				Type:    "integer",
				Format:  "int32",
				Default: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertParameter(tt.param)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected parameter, got nil")
			}

			if result.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", result.Name, tt.expected.Name)
			}

			if result.In != tt.expected.In {
				t.Errorf("In = %v, want %v", result.In, tt.expected.In)
			}

			if result.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", result.Type, tt.expected.Type)
			}

			if tt.param != nil && tt.param.Deprecated {
				if result.Extensions == nil {
					t.Error("Expected Extensions for deprecated parameter")
				} else if result.Extensions["x-deprecated"] != true {
					t.Error("Expected x-deprecated extension to be true")
				}
			}
		})
	}
}

// TestConvertSchemaToParameter tests schema to parameter conversion
func TestConvertSchemaToParameter(t *testing.T) {
	conv := New()

	tests := []struct {
		name   string
		schema *openapi.Schema
		verify func(*testing.T, *swagger.Parameter)
	}{
		{
			name:   "nil schema",
			schema: nil,
			verify: func(t *testing.T, p *swagger.Parameter) {
				if p.Type != "" {
					t.Errorf("Expected empty type for nil schema, got %v", p.Type)
				}
			},
		},
		{
			name: "string with pattern",
			schema: &openapi.Schema{
				Type:    "string",
				Format:  "email",
				Pattern: "^[a-z]+@[a-z]+\\.[a-z]+$",
			},
			verify: func(t *testing.T, p *swagger.Parameter) {
				if p.Type != "string" {
					t.Errorf("Type = %v, want string", p.Type)
				}
				if p.Format != "email" {
					t.Errorf("Format = %v, want email", p.Format)
				}
				if p.Pattern != "^[a-z]+@[a-z]+\\.[a-z]+$" {
					t.Errorf("Pattern = %v, want pattern", p.Pattern)
				}
			},
		},
		{
			name: "number with min/max",
			schema: &openapi.Schema{
				Type:    "number",
				Format:  "double",
				Minimum: 1.0,
				Maximum: 100.0,
			},
			verify: func(t *testing.T, p *swagger.Parameter) {
				if p.Type != "number" {
					t.Errorf("Type = %v, want number", p.Type)
				}
				if p.Minimum == nil || *p.Minimum != 1.0 {
					t.Errorf("Minimum = %v, want 1.0", p.Minimum)
				}
				if p.Maximum == nil || *p.Maximum != 100.0 {
					t.Errorf("Maximum = %v, want 100.0", p.Maximum)
				}
			},
		},
		{
			name: "string with length constraints",
			schema: &openapi.Schema{
				Type:      "string",
				MinLength: 5,
				MaxLength: 50,
			},
			verify: func(t *testing.T, p *swagger.Parameter) {
				if p.MinLength == nil || *p.MinLength != 5 {
					t.Errorf("MinLength = %v, want 5", p.MinLength)
				}
				if p.MaxLength == nil || *p.MaxLength != 50 {
					t.Errorf("MaxLength = %v, want 50", p.MaxLength)
				}
			},
		},
		{
			name: "array with items",
			schema: &openapi.Schema{
				Type:        "array",
				MinItems:    1,
				MaxItems:    10,
				UniqueItems: true,
				Items: &openapi.Schema{
					Type:   "string",
					Format: "uuid",
				},
			},
			verify: func(t *testing.T, p *swagger.Parameter) {
				if p.Type != "array" {
					t.Errorf("Type = %v, want array", p.Type)
				}
				if p.MinItems == nil || *p.MinItems != 1 {
					t.Errorf("MinItems = %v, want 1", p.MinItems)
				}
				if p.MaxItems == nil || *p.MaxItems != 10 {
					t.Errorf("MaxItems = %v, want 10", p.MaxItems)
				}
				if !p.UniqueItems {
					t.Error("UniqueItems should be true")
				}
				if p.Items == nil {
					t.Fatal("Items should not be nil")
				}
				if p.Items.Type != "string" {
					t.Errorf("Items.Type = %v, want string", p.Items.Type)
				}
				if p.Items.Format != "uuid" {
					t.Errorf("Items.Format = %v, want uuid", p.Items.Format)
				}
			},
		},
		{
			name: "enum values",
			schema: &openapi.Schema{
				Type: "string",
				Enum: []interface{}{"active", "inactive", "pending"},
			},
			verify: func(t *testing.T, p *swagger.Parameter) {
				if p.Enum == nil || len(p.Enum) != 3 {
					t.Errorf("Enum = %v, want 3 values", p.Enum)
				}
			},
		},
		{
			name: "default value",
			schema: &openapi.Schema{
				Type:    "integer",
				Default: 42,
			},
			verify: func(t *testing.T, p *swagger.Parameter) {
				if p.Default != 42 {
					t.Errorf("Default = %v, want 42", p.Default)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := &swagger.Parameter{}
			conv.convertSchemaToParameter(tt.schema, param)
			tt.verify(t, param)
		})
	}
}

// TestConvertSchemaToItems tests schema to items conversion
func TestConvertSchemaToItems(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		schema   *openapi.Schema
		expected *swagger.Items
	}{
		{
			name:     "nil schema",
			schema:   nil,
			expected: nil,
		},
		{
			name: "simple string items",
			schema: &openapi.Schema{
				Type:   "string",
				Format: "date-time",
			},
			expected: &swagger.Items{
				Type:   "string",
				Format: "date-time",
			},
		},
		{
			name: "integer items with default",
			schema: &openapi.Schema{
				Type:    "integer",
				Format:  "int64",
				Default: 100,
			},
			expected: &swagger.Items{
				Type:    "integer",
				Format:  "int64",
				Default: 100,
			},
		},
		{
			name: "boolean items",
			schema: &openapi.Schema{
				Type:    "boolean",
				Default: true,
			},
			expected: &swagger.Items{
				Type:    "boolean",
				Default: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertSchemaToItems(tt.schema)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected items, got nil")
			}

			if result.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", result.Type, tt.expected.Type)
			}

			if result.Format != tt.expected.Format {
				t.Errorf("Format = %v, want %v", result.Format, tt.expected.Format)
			}

			if result.Default != tt.expected.Default {
				t.Errorf("Default = %v, want %v", result.Default, tt.expected.Default)
			}
		})
	}
}

// TestExtractType tests type extraction from various formats
func TestExtractType(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		typeVal  interface{}
		expected string
	}{
		{
			name:     "nil type",
			typeVal:  nil,
			expected: "",
		},
		{
			name:     "simple string",
			typeVal:  "string",
			expected: "string",
		},
		{
			name:     "integer type",
			typeVal:  "integer",
			expected: "integer",
		},
		{
			name:     "array type",
			typeVal:  "array",
			expected: "array",
		},
		{
			name:     "nullable string (interface slice)",
			typeVal:  []interface{}{"string", "null"},
			expected: "string",
		},
		{
			name:     "nullable integer (interface slice)",
			typeVal:  []interface{}{"integer", "null"},
			expected: "integer",
		},
		{
			name:     "only null (interface slice)",
			typeVal:  []interface{}{"null"},
			expected: "",
		},
		{
			name:     "nullable string (string slice)",
			typeVal:  []string{"string", "null"},
			expected: "string",
		},
		{
			name:     "nullable number (string slice)",
			typeVal:  []string{"number", "null"},
			expected: "number",
		},
		{
			name:     "multiple types (first non-null)",
			typeVal:  []interface{}{"null", "string", "integer"},
			expected: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.extractType(tt.typeVal)

			if result != tt.expected {
				t.Errorf("extractType(%v) = %v, want %v", tt.typeVal, result, tt.expected)
			}
		})
	}
}

// TestConvertHeader tests header conversion
func TestConvertHeader(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		header   *openapi.Header
		expected *swagger.Header
	}{
		{
			name:     "nil header",
			header:   nil,
			expected: nil,
		},
		{
			name: "simple string header",
			header: &openapi.Header{
				Description: "API version header",
				Schema: &openapi.Schema{
					Type:   "string",
					Format: "version",
				},
			},
			expected: &swagger.Header{
				Description: "API version header",
				Type:        "string",
				Format:      "version",
			},
		},
		{
			name: "integer header with default",
			header: &openapi.Header{
				Description: "Rate limit",
				Schema: &openapi.Schema{
					Type:    "integer",
					Format:  "int64",
					Default: 1000,
				},
			},
			expected: &swagger.Header{
				Description: "Rate limit",
				Type:        "integer",
				Format:      "int64",
				Default:     1000,
			},
		},
		{
			name: "header without schema",
			header: &openapi.Header{
				Description: "Custom header",
			},
			expected: &swagger.Header{
				Description: "Custom header",
			},
		},
		{
			name: "boolean header",
			header: &openapi.Header{
				Description: "Feature flag",
				Schema: &openapi.Schema{
					Type:    "boolean",
					Default: false,
				},
			},
			expected: &swagger.Header{
				Description: "Feature flag",
				Type:        "boolean",
				Default:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertHeader(tt.header)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected header, got nil")
			}

			if result.Description != tt.expected.Description {
				t.Errorf("Description = %v, want %v", result.Description, tt.expected.Description)
			}

			if result.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", result.Type, tt.expected.Type)
			}

			if result.Format != tt.expected.Format {
				t.Errorf("Format = %v, want %v", result.Format, tt.expected.Format)
			}

			if result.Default != tt.expected.Default {
				t.Errorf("Default = %v, want %v", result.Default, tt.expected.Default)
			}
		})
	}
}

// TestConvertHeaders tests headers conversion
func TestConvertHeaders(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		headers  map[string]*openapi.Header
		expected int
	}{
		{
			name:     "nil headers",
			headers:  nil,
			expected: 0,
		},
		{
			name:     "empty headers",
			headers:  map[string]*openapi.Header{},
			expected: 0,
		},
		{
			name: "single header",
			headers: map[string]*openapi.Header{
				"X-RateLimit": {
					Description: "Rate limit",
					Schema: &openapi.Schema{
						Type: "integer",
					},
				},
			},
			expected: 1,
		},
		{
			name: "multiple headers",
			headers: map[string]*openapi.Header{
				"X-RateLimit": {
					Schema: &openapi.Schema{Type: "integer"},
				},
				"X-API-Version": {
					Schema: &openapi.Schema{Type: "string"},
				},
				"X-Request-ID": {
					Schema: &openapi.Schema{Type: "string", Format: "uuid"},
				},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertHeaders(tt.headers)

			if tt.headers == nil {
				if result != nil {
					t.Errorf("Expected nil for nil input, got %v", result)
				}
				return
			}

			if len(result) != tt.expected {
				t.Errorf("Got %d headers, want %d", len(result), tt.expected)
			}

			for name := range tt.headers {
				if _, ok := result[name]; !ok {
					t.Errorf("Header %q not found in result", name)
				}
			}
		})
	}
}

// TestConvertExamples tests examples conversion
func TestConvertExamples(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		content  map[string]*openapi.MediaType
		expected map[string]interface{}
	}{
		{
			name:     "nil content",
			content:  nil,
			expected: nil,
		},
		{
			name:     "empty content",
			content:  map[string]*openapi.MediaType{},
			expected: nil,
		},
		{
			name: "single example",
			content: map[string]*openapi.MediaType{
				"application/json": {
					Example: map[string]interface{}{
						"id":   1,
						"name": "Test",
					},
				},
			},
			expected: map[string]interface{}{
				"application/json": map[string]interface{}{
					"id":   1,
					"name": "Test",
				},
			},
		},
		{
			name: "multiple examples",
			content: map[string]*openapi.MediaType{
				"application/json": {
					Example: map[string]interface{}{"id": 1},
				},
				"application/xml": {
					Example: "<root><id>1</id></root>",
				},
			},
			expected: map[string]interface{}{
				"application/json": map[string]interface{}{"id": 1},
				"application/xml":  "<root><id>1</id></root>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertExamples(tt.content)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Got %d examples, want %d", len(result), len(tt.expected))
			}

			for key := range tt.expected {
				if _, ok := result[key]; !ok {
					t.Errorf("Example for %q not found", key)
				}
			}
		})
	}
}

// TestConvertRefFunctions tests ref conversion functions
func TestConvertRefToV2(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{
			name:     "empty ref",
			ref:      "",
			expected: "",
		},
		{
			name:     "v3 schema ref",
			ref:      "#/components/schemas/User",
			expected: "#/definitions/User",
		},
		{
			name:     "v3 parameter ref",
			ref:      "#/components/parameters/limit",
			expected: "#/parameters/limit",
		},
		{
			name:     "v3 response ref",
			ref:      "#/components/responses/NotFound",
			expected: "#/responses/NotFound",
		},
		{
			name:     "already v2 format",
			ref:      "#/definitions/Product",
			expected: "#/definitions/Product",
		},
		{
			name:     "external ref",
			ref:      "common.yaml#/components/schemas/Error",
			expected: "common.yaml#/definitions/Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertRefToV2(tt.ref)

			if result != tt.expected {
				t.Errorf("convertRefToV2(%q) = %q, want %q", tt.ref, result, tt.expected)
			}
		})
	}
}

func TestConvertRefToV3(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{
			name:     "empty ref",
			ref:      "",
			expected: "",
		},
		{
			name:     "v2 definition ref",
			ref:      "#/definitions/User",
			expected: "#/components/schemas/User",
		},
		{
			name:     "v2 parameter ref",
			ref:      "#/parameters/limit",
			expected: "#/components/parameters/limit",
		},
		{
			name:     "v2 response ref",
			ref:      "#/responses/NotFound",
			expected: "#/components/responses/NotFound",
		},
		{
			name:     "already v3 format",
			ref:      "#/components/schemas/Product",
			expected: "#/components/schemas/Product",
		},
		{
			name:     "external ref",
			ref:      "common.yaml#/definitions/Error",
			expected: "common.yaml#/components/schemas/Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertRefToV3(tt.ref)

			if result != tt.expected {
				t.Errorf("convertRefToV3(%q) = %q, want %q", tt.ref, result, tt.expected)
			}
		})
	}
}

// TestConvertSchemaComplex tests complex schema conversion scenarios
func TestConvertSchemaComplex(t *testing.T) {
	conv := New()

	tests := []struct {
		name   string
		schema *openapi.Schema
		verify func(*testing.T, *swagger.Schema)
	}{
		{
			name: "schema with all numeric constraints",
			schema: &openapi.Schema{
				Type:          "number",
				MultipleOf:    0.5,
				Minimum:       float64(1.0),
				Maximum:       float64(100.0),
				MinLength:     int(5),
				MaxLength:     int(50),
				MinItems:      int(1),
				MaxItems:      int(10),
				MinProperties: int(2),
				MaxProperties: int(20),
			},
			verify: func(t *testing.T, s *swagger.Schema) {
				if s.MultipleOf == nil || *s.MultipleOf != 0.5 {
					t.Errorf("MultipleOf = %v, want 0.5", s.MultipleOf)
				}
				if s.Minimum == nil || *s.Minimum != 1.0 {
					t.Errorf("Minimum = %v, want 1.0", s.Minimum)
				}
				if s.Maximum == nil || *s.Maximum != 100.0 {
					t.Errorf("Maximum = %v, want 100.0", s.Maximum)
				}
			},
		},
		{
			name: "schema with properties",
			schema: &openapi.Schema{
				Type: "object",
				Properties: map[string]*openapi.Schema{
					"name": {Type: "string"},
					"age":  {Type: "integer"},
				},
				Required: []string{"name"},
			},
			verify: func(t *testing.T, s *swagger.Schema) {
				if len(s.Properties) != 2 {
					t.Errorf("Properties count = %d, want 2", len(s.Properties))
				}
				if len(s.Required) != 1 || s.Required[0] != "name" {
					t.Errorf("Required = %v, want [name]", s.Required)
				}
			},
		},
		{
			name: "schema with allOf",
			schema: &openapi.Schema{
				AllOf: []openapi.Schema{
					{Ref: "#/components/schemas/Base"},
					{Type: "object", Properties: map[string]*openapi.Schema{
						"extra": {Type: "string"},
					}},
				},
			},
			verify: func(t *testing.T, s *swagger.Schema) {
				if len(s.AllOf) != 2 {
					t.Errorf("AllOf count = %d, want 2", len(s.AllOf))
				}
				if s.AllOf[0].Ref != "#/definitions/Base" {
					t.Errorf("AllOf[0].Ref = %v, want #/definitions/Base", s.AllOf[0].Ref)
				}
			},
		},
		{
			name: "schema with items",
			schema: &openapi.Schema{
				Type: "array",
				Items: &openapi.Schema{
					Type:   "string",
					Format: "uuid",
				},
			},
			verify: func(t *testing.T, s *swagger.Schema) {
				if s.Items == nil {
					t.Fatal("Items should not be nil")
				}
				if s.Items.Type != "string" {
					t.Errorf("Items.Type = %v, want string", s.Items.Type)
				}
			},
		},
		{
			name: "schema with additionalProperties (bool)",
			schema: &openapi.Schema{
				Type:                 "object",
				AdditionalProperties: true,
			},
			verify: func(t *testing.T, s *swagger.Schema) {
				if b, ok := s.AdditionalProperties.(bool); !ok || !b {
					t.Errorf("AdditionalProperties = %v, want true", s.AdditionalProperties)
				}
			},
		},
		{
			name: "schema with additionalProperties (schema)",
			schema: &openapi.Schema{
				Type: "object",
				AdditionalProperties: &openapi.Schema{
					Type: "integer",
				},
			},
			verify: func(t *testing.T, s *swagger.Schema) {
				if schema, ok := s.AdditionalProperties.(*swagger.Schema); !ok {
					t.Error("AdditionalProperties should be a schema")
				} else if schema.Type != "integer" {
					t.Errorf("AdditionalProperties.Type = %v, want integer", schema.Type)
				}
			},
		},
		{
			name: "nullable schema",
			schema: &openapi.Schema{
				Type: []interface{}{"string", "null"},
			},
			verify: func(t *testing.T, s *swagger.Schema) {
				if s.Type != "string" {
					t.Errorf("Type = %v, want string", s.Type)
				}
				if s.Extensions == nil || s.Extensions["x-nullable"] != true {
					t.Error("Expected x-nullable extension")
				}
			},
		},
		{
			name: "schema with XML",
			schema: &openapi.Schema{
				Type: "object",
				XML: &openapi.XML{
					Name:      "user",
					Namespace: "http://example.com",
					Prefix:    "ex",
					Attribute: false,
					Wrapped:   true,
				},
			},
			verify: func(t *testing.T, s *swagger.Schema) {
				if s.XML == nil {
					t.Fatal("XML should not be nil")
				}
				if s.XML.Name != "user" {
					t.Errorf("XML.Name = %v, want user", s.XML.Name)
				}
			},
		},
		{
			name: "schema with discriminator",
			schema: &openapi.Schema{
				Discriminator: &openapi.Discriminator{
					PropertyName: "type",
				},
			},
			verify: func(t *testing.T, s *swagger.Schema) {
				if s.Discriminator != "type" {
					t.Errorf("Discriminator = %v, want type", s.Discriminator)
				}
			},
		},
		{
			name: "schema with exclusive maximum/minimum (bool)",
			schema: &openapi.Schema{
				Type:             "number",
				Maximum:          100.0,
				ExclusiveMaximum: true,
				Minimum:          0.0,
				ExclusiveMinimum: true,
			},
			verify: func(t *testing.T, s *swagger.Schema) {
				if !s.ExclusiveMaximum {
					t.Error("ExclusiveMaximum should be true")
				}
				if !s.ExclusiveMinimum {
					t.Error("ExclusiveMinimum should be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertSchema(tt.schema)
			if result == nil {
				t.Fatal("Expected schema, got nil")
			}
			tt.verify(t, result)
		})
	}
}

// TestConvertSecuritySchemeComplex tests complex security scheme conversions
func TestConvertSecuritySchemeComplex(t *testing.T) {
	conv := New()

	tests := []struct {
		name   string
		scheme *openapi.SecurityScheme
		verify func(*testing.T, *swagger.SecurityScheme)
	}{
		{
			name: "http basic auth",
			scheme: &openapi.SecurityScheme{
				Type:   "http",
				Scheme: "basic",
			},
			verify: func(t *testing.T, s *swagger.SecurityScheme) {
				if s.Type != "basic" {
					t.Errorf("Type = %v, want basic", s.Type)
				}
			},
		},
		{
			name: "http bearer auth",
			scheme: &openapi.SecurityScheme{
				Type:   "http",
				Scheme: "bearer",
			},
			verify: func(t *testing.T, s *swagger.SecurityScheme) {
				if s.Type != "apiKey" {
					t.Errorf("Type = %v, want apiKey", s.Type)
				}
				if s.In != "header" {
					t.Errorf("In = %v, want header", s.In)
				}
				if s.Name != "Authorization" {
					t.Errorf("Name = %v, want Authorization", s.Name)
				}
			},
		},
		{
			name: "apiKey scheme",
			scheme: &openapi.SecurityScheme{
				Type: "apiKey",
				Name: "X-API-Key",
				In:   "header",
			},
			verify: func(t *testing.T, s *swagger.SecurityScheme) {
				if s.Type != "apiKey" {
					t.Errorf("Type = %v, want apiKey", s.Type)
				}
				if s.Name != "X-API-Key" {
					t.Errorf("Name = %v, want X-API-Key", s.Name)
				}
			},
		},
		{
			name: "oauth2 with flows",
			scheme: &openapi.SecurityScheme{
				Type: "oauth2",
				Flows: &openapi.OAuthFlows{
					Implicit: &openapi.OAuthFlow{
						AuthorizationURL: "https://example.com/auth",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
				},
			},
			verify: func(t *testing.T, s *swagger.SecurityScheme) {
				if s.Type != "oauth2" {
					t.Errorf("Type = %v, want oauth2", s.Type)
				}
				if s.AuthorizationURL != "https://example.com/auth" {
					t.Errorf("AuthorizationURL = %v", s.AuthorizationURL)
				}
				if len(s.Scopes) != 2 {
					t.Errorf("Scopes count = %d, want 2", len(s.Scopes))
				}
			},
		},
		{
			name: "openIdConnect",
			scheme: &openapi.SecurityScheme{
				Type:             "openIdConnect",
				OpenIDConnectURL: "https://example.com/.well-known/openid-configuration",
			},
			verify: func(t *testing.T, s *swagger.SecurityScheme) {
				// openIdConnect is converted to oauth2 in Swagger 2.0
				if s.Type != "oauth2" {
					t.Errorf("Type = %v, want oauth2", s.Type)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertSecurityScheme(tt.scheme)
			if result == nil {
				t.Fatal("Expected security scheme, got nil")
			}
			tt.verify(t, result)
		})
	}
}

// TestConvertOAuth2FlowsComplex tests OAuth2 flows conversion
func TestConvertOAuth2FlowsComplex(t *testing.T) {
	conv := New()

	tests := []struct {
		name   string
		flows  *openapi.OAuthFlows
		verify func(*testing.T, *swagger.SecurityScheme)
	}{
		{
			name: "implicit flow",
			flows: &openapi.OAuthFlows{
				Implicit: &openapi.OAuthFlow{
					AuthorizationURL: "https://example.com/oauth/authorize",
					Scopes: map[string]string{
						"read":  "Read scope",
						"write": "Write scope",
					},
				},
			},
			verify: func(t *testing.T, s *swagger.SecurityScheme) {
				if s.Flow != "implicit" {
					t.Errorf("Flow = %v, want implicit", s.Flow)
				}
				if s.AuthorizationURL != "https://example.com/oauth/authorize" {
					t.Errorf("AuthorizationURL = %v", s.AuthorizationURL)
				}
			},
		},
		{
			name: "password flow",
			flows: &openapi.OAuthFlows{
				Password: &openapi.OAuthFlow{
					TokenURL: "https://example.com/oauth/token",
					Scopes: map[string]string{
						"admin": "Admin scope",
					},
				},
			},
			verify: func(t *testing.T, s *swagger.SecurityScheme) {
				if s.Flow != "password" {
					t.Errorf("Flow = %v, want password", s.Flow)
				}
				if s.TokenURL != "https://example.com/oauth/token" {
					t.Errorf("TokenURL = %v", s.TokenURL)
				}
			},
		},
		{
			name: "application (client credentials) flow",
			flows: &openapi.OAuthFlows{
				ClientCredentials: &openapi.OAuthFlow{
					TokenURL: "https://example.com/oauth/token",
					Scopes: map[string]string{
						"api": "API access",
					},
				},
			},
			verify: func(t *testing.T, s *swagger.SecurityScheme) {
				if s.Flow != "application" {
					t.Errorf("Flow = %v, want application", s.Flow)
				}
				if s.TokenURL != "https://example.com/oauth/token" {
					t.Errorf("TokenURL = %v", s.TokenURL)
				}
			},
		},
		{
			name: "authorization code flow",
			flows: &openapi.OAuthFlows{
				AuthorizationCode: &openapi.OAuthFlow{
					AuthorizationURL: "https://example.com/oauth/authorize",
					TokenURL:         "https://example.com/oauth/token",
					Scopes: map[string]string{
						"openid": "OpenID",
					},
				},
			},
			verify: func(t *testing.T, s *swagger.SecurityScheme) {
				if s.Flow != "accessCode" {
					t.Errorf("Flow = %v, want accessCode", s.Flow)
				}
				if s.AuthorizationURL != "https://example.com/oauth/authorize" {
					t.Errorf("AuthorizationURL = %v", s.AuthorizationURL)
				}
				if s.TokenURL != "https://example.com/oauth/token" {
					t.Errorf("TokenURL = %v", s.TokenURL)
				}
			},
		},
		{
			name: "multiple flows (uses first)",
			flows: &openapi.OAuthFlows{
				Implicit: &openapi.OAuthFlow{
					AuthorizationURL: "https://example.com/auth",
					Scopes:           map[string]string{"read": "Read"},
				},
				Password: &openapi.OAuthFlow{
					TokenURL: "https://example.com/token",
					Scopes:   map[string]string{"write": "Write"},
				},
			},
			verify: func(t *testing.T, s *swagger.SecurityScheme) {
				if s.Flow != "implicit" {
					t.Errorf("Flow = %v, want implicit (first flow)", s.Flow)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v2Scheme := &swagger.SecurityScheme{Type: "oauth2"}
			v3Scheme := &openapi.SecurityScheme{Type: "oauth2", Flows: tt.flows}
			conv.convertOAuth2Flows(v3Scheme, v2Scheme)
			tt.verify(t, v2Scheme)
		})
	}
}

// TestConvertParametersV2ToV3 tests parameter array conversion
func TestConvertParametersV2ToV3(t *testing.T) {
	conv := New()

	params := []*swagger.Parameter{
		{
			Name:        "id",
			In:          "path",
			Description: "User ID",
			Required:    true,
			Type:        "integer",
		},
		{
			Name:        "name",
			In:          "query",
			Description: "User name",
			Type:        "string",
		},
	}

	result := conv.convertParametersToV3(params)

	if len(result) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(result))
	}

	if result[0].Name != "id" {
		t.Errorf("Parameter[0].Name = %v, want id", result[0].Name)
	}

	if result[0].Schema == nil {
		t.Fatal("Parameter[0].Schema should not be nil")
	}

	if result[0].Schema.Type != "integer" {
		t.Errorf("Parameter[0].Schema.Type = %v, want integer", result[0].Schema.Type)
	}
}

// TestConvertSecurity tests security requirements conversion
func TestConvertSecurity(t *testing.T) {
	conv := New()

	security := []openapi.SecurityRequirement{
		{"api_key": []string{}},
		{"oauth2": []string{"read", "write"}},
	}

	result := conv.convertSecurity(security)

	if len(result) != 2 {
		t.Errorf("Expected 2 security requirements, got %d", len(result))
	}

	if _, ok := result[0]["api_key"]; !ok {
		t.Error("Expected api_key in first requirement")
	}

	if scopes, ok := result[1]["oauth2"]; !ok {
		t.Error("Expected oauth2 in second requirement")
	} else if len(scopes) != 2 {
		t.Errorf("Expected 2 scopes for oauth2, got %d", len(scopes))
	}
}

// TestConvertTags tests tags conversion
func TestConvertTags(t *testing.T) {
	conv := New()

	tags := []openapi.Tag{
		{
			Name:        "users",
			Description: "User operations",
		},
		{
			Name:        "products",
			Description: "Product operations",
			ExternalDocs: &openapi.ExternalDocs{
				Description: "External docs",
				URL:         "https://example.com/docs",
			},
		},
	}

	result := conv.convertTags(tags)

	if len(result) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(result))
	}

	if result[0].Name != "users" {
		t.Errorf("Tag[0].Name = %v, want users", result[0].Name)
	}

	if result[1].ExternalDocs == nil {
		t.Error("Tag[1].ExternalDocs should not be nil")
	}
}

// TestConvertBodyParameterToRequestBody tests body parameter to request body conversion
func TestConvertBodyParameterToRequestBody(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		param    *swagger.Parameter
		consumes []string
		verify   func(*testing.T, *openapi.RequestBody)
	}{
		{
			name: "simple body parameter",
			param: &swagger.Parameter{
				Name:        "body",
				In:          "body",
				Description: "User object",
				Required:    true,
				Schema: &swagger.Schema{
					Type: "object",
					Properties: map[string]*swagger.Schema{
						"name": {Type: "string"},
						"age":  {Type: "integer"},
					},
				},
			},
			consumes: []string{"application/json"},
			verify: func(t *testing.T, rb *openapi.RequestBody) {
				if !rb.Required {
					t.Error("RequestBody should be required")
				}
				if rb.Description != "User object" {
					t.Errorf("Description = %v, want 'User object'", rb.Description)
				}
				if rb.Content == nil {
					t.Fatal("Content should not be nil")
				}
				if _, ok := rb.Content["application/json"]; !ok {
					t.Error("Expected application/json media type")
				}
			},
		},
		{
			name: "body parameter with example",
			param: &swagger.Parameter{
				Name:     "body",
				In:       "body",
				Required: false,
				Schema: &swagger.Schema{
					Type:    "string",
					Example: "test data",
				},
			},
			consumes: []string{"application/xml"},
			verify: func(t *testing.T, rb *openapi.RequestBody) {
				if rb.Required {
					t.Error("RequestBody should not be required")
				}
				if rb.Content == nil {
					t.Fatal("Content should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertBodyParameterToRequestBody(tt.param, tt.consumes)
			if result == nil {
				t.Fatal("Expected request body, got nil")
			}
			tt.verify(t, result)
		})
	}
}

// TestConvertParameterToV3 tests single parameter conversion to V3
func TestConvertParameterToV3(t *testing.T) {
	conv := New()

	tests := []struct {
		name   string
		param  *swagger.Parameter
		verify func(*testing.T, *openapi.Parameter)
	}{
		{
			name: "simple query parameter",
			param: &swagger.Parameter{
				Name:        "limit",
				In:          "query",
				Description: "Limit results",
				Required:    false,
				Type:        "integer",
				Format:      "int32",
				Default:     10,
			},
			verify: func(t *testing.T, p *openapi.Parameter) {
				if p.Name != "limit" {
					t.Errorf("Name = %v, want limit", p.Name)
				}
				if p.In != "query" {
					t.Errorf("In = %v, want query", p.In)
				}
				if p.Schema == nil {
					t.Fatal("Schema should not be nil")
				}
				if p.Schema.Type != "integer" {
					t.Errorf("Schema.Type = %v, want integer", p.Schema.Type)
				}
			},
		},
		{
			name: "path parameter",
			param: &swagger.Parameter{
				Name:     "id",
				In:       "path",
				Required: true,
				Type:     "string",
				Format:   "uuid",
			},
			verify: func(t *testing.T, p *openapi.Parameter) {
				if !p.Required {
					t.Error("Path parameter should be required")
				}
				if p.Schema.Format != "uuid" {
					t.Errorf("Schema.Format = %v, want uuid", p.Schema.Format)
				}
			},
		},
		{
			name: "header parameter",
			param: &swagger.Parameter{
				Name:        "X-Request-ID",
				In:          "header",
				Description: "Request ID",
				Type:        "string",
			},
			verify: func(t *testing.T, p *openapi.Parameter) {
				if p.In != "header" {
					t.Errorf("In = %v, want header", p.In)
				}
			},
		},
		{
			name: "parameter with enum",
			param: &swagger.Parameter{
				Name: "status",
				In:   "query",
				Type: "string",
				Enum: []interface{}{"active", "inactive"},
			},
			verify: func(t *testing.T, p *openapi.Parameter) {
				if len(p.Schema.Enum) != 2 {
					t.Errorf("Schema.Enum count = %d, want 2", len(p.Schema.Enum))
				}
			},
		},
		{
			name: "array parameter",
			param: &swagger.Parameter{
				Name:             "tags",
				In:               "query",
				Type:             "array",
				CollectionFormat: "csv",
				Items: &swagger.Items{
					Type: "string",
				},
			},
			verify: func(t *testing.T, p *openapi.Parameter) {
				if p.Schema.Type != "array" {
					t.Errorf("Schema.Type = %v, want array", p.Schema.Type)
				}
				if p.Schema.Items == nil {
					t.Fatal("Schema.Items should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertParameterToV3(tt.param)
			if result == nil {
				t.Fatal("Expected parameter, got nil")
			}
			tt.verify(t, result)
		})
	}
}

// TestIsNullable tests nullable type detection
func TestIsNullable(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		typ      interface{}
		expected bool
	}{
		{
			name:     "null type only",
			typ:      []interface{}{"null"},
			expected: true,
		},
		{
			name:     "string and null",
			typ:      []interface{}{"string", "null"},
			expected: true,
		},
		{
			name:     "integer and null",
			typ:      []interface{}{"integer", "null"},
			expected: true,
		},
		{
			name:     "multiple types with null",
			typ:      []interface{}{"string", "integer", "null"},
			expected: true,
		},
		{
			name:     "no null type",
			typ:      []interface{}{"string", "integer"},
			expected: false,
		},
		{
			name:     "string only",
			typ:      "string",
			expected: false,
		},
		{
			name:     "empty array",
			typ:      []interface{}{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.isNullable(tt.typ)
			if result != tt.expected {
				t.Errorf("isNullable(%v) = %v, want %v", tt.typ, result, tt.expected)
			}
		})
	}
}

// TestConvertParameterDefinitions tests parameter definitions conversion
func TestConvertParameterDefinitions(t *testing.T) {
	conv := New()

	defs := map[string]*openapi.Parameter{
		"limitParam": {
			Name:        "limit",
			In:          "query",
			Description: "Limit results",
			Schema: &openapi.Schema{
				Type: "integer",
			},
		},
		"offsetParam": {
			Name:        "offset",
			In:          "query",
			Description: "Offset results",
			Schema: &openapi.Schema{
				Type: "integer",
			},
		},
	}

	result := conv.convertParameterDefinitions(defs)

	if len(result) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(result))
	}

	if _, ok := result["limitParam"]; !ok {
		t.Error("Expected limitParam in result")
	}

	if _, ok := result["offsetParam"]; !ok {
		t.Error("Expected offsetParam in result")
	}
}

// TestConvertResponseDefinitions tests response definitions conversion
func TestConvertResponseDefinitions(t *testing.T) {
	conv := New()

	defs := map[string]*openapi.Response{
		"NotFound": {
			Description: "Resource not found",
			Content: map[string]*openapi.MediaType{
				"application/json": {
					Schema: &openapi.Schema{
						Type: "object",
						Properties: map[string]*openapi.Schema{
							"error": {Type: "string"},
						},
					},
				},
			},
		},
		"Success": {
			Description: "Success response",
			Content: map[string]*openapi.MediaType{
				"application/json": {
					Schema: &openapi.Schema{
						Type: "string",
					},
				},
			},
		},
	}

	result := conv.convertResponseDefinitions(defs)

	if len(result) != 2 {
		t.Errorf("Expected 2 responses, got %d", len(result))
	}

	if _, ok := result["NotFound"]; !ok {
		t.Error("Expected NotFound in result")
	}

	notFound := result["NotFound"]
	if notFound.Description != "Resource not found" {
		t.Errorf("NotFound.Description = %v, want 'Resource not found'", notFound.Description)
	}
}

// TestSeparateBodyParameter tests body parameter separation
func TestSeparateBodyParameter(t *testing.T) {
	conv := New()

	params := []*swagger.Parameter{
		{Name: "id", In: "path", Type: "string"},
		{Name: "body", In: "body", Schema: &swagger.Schema{Type: "object"}},
		{Name: "limit", In: "query", Type: "integer"},
	}

	nonBody, bodyParam := conv.separateBodyParameter(params)

	if bodyParam == nil {
		t.Fatal("Expected body parameter, got nil")
	}

	if bodyParam.Name != "body" {
		t.Errorf("bodyParam.Name = %v, want body", bodyParam.Name)
	}

	if len(nonBody) != 2 {
		t.Errorf("Expected 2 other parameters, got %d", len(nonBody))
	}

	// Test with no body parameter
	noBodyParams := []*swagger.Parameter{
		{Name: "id", In: "path", Type: "string"},
		{Name: "limit", In: "query", Type: "integer"},
	}

	nonBody2, bodyParam2 := conv.separateBodyParameter(noBodyParams)

	if bodyParam2 != nil {
		t.Error("Expected nil body parameter")
	}

	if len(nonBody2) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(nonBody2))
	}
}

// TestConvertResponseToV3 tests response conversion to V3
func TestConvertResponseToV3(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		resp     *swagger.Response
		produces []string
		verify   func(*testing.T, *openapi.Response)
	}{
		{
			name: "response with schema",
			resp: &swagger.Response{
				Description: "User response",
				Schema: &swagger.Schema{
					Type: "object",
					Properties: map[string]*swagger.Schema{
						"id":   {Type: "integer"},
						"name": {Type: "string"},
					},
				},
			},
			produces: []string{"application/json"},
			verify: func(t *testing.T, r *openapi.Response) {
				if r.Description != "User response" {
					t.Errorf("Description = %v, want 'User response'", r.Description)
				}
				if r.Content == nil {
					t.Fatal("Content should not be nil")
				}
				if _, ok := r.Content["application/json"]; !ok {
					t.Error("Expected application/json media type")
				}
			},
		},
		{
			name: "response with headers",
			resp: &swagger.Response{
				Description: "Success",
				Headers: map[string]*swagger.Header{
					"X-Rate-Limit": {
						Type:        "integer",
						Description: "Rate limit",
					},
				},
			},
			produces: []string{"application/json"},
			verify: func(t *testing.T, r *openapi.Response) {
				if r.Headers == nil {
					t.Fatal("Headers should not be nil")
				}
				if _, ok := r.Headers["X-Rate-Limit"]; !ok {
					t.Error("Expected X-Rate-Limit header")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertResponseToV3(tt.resp, tt.produces)
			if result == nil {
				t.Fatal("Expected response, got nil")
			}
			tt.verify(t, result)
		})
	}
}

// TestConvertHeadersToV3 tests headers conversion to V3
func TestConvertHeadersToV3(t *testing.T) {
	conv := New()

	headers := map[string]*swagger.Header{
		"X-Rate-Limit": {
			Type:        "integer",
			Description: "Rate limit remaining",
			Format:      "int32",
		},
		"X-Request-ID": {
			Type:        "string",
			Description: "Request identifier",
			Format:      "uuid",
		},
	}

	result := conv.convertHeadersToV3(headers)

	if len(result) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(result))
	}

	if rateLimit, ok := result["X-Rate-Limit"]; !ok {
		t.Error("Expected X-Rate-Limit header")
	} else {
		if rateLimit.Schema == nil {
			t.Fatal("X-Rate-Limit.Schema should not be nil")
		}
		if rateLimit.Schema.Type != "integer" {
			t.Errorf("X-Rate-Limit.Schema.Type = %v, want integer", rateLimit.Schema.Type)
		}
	}
}

// TestConvertItemsToSchema tests items conversion to schema
func TestConvertItemsToSchema(t *testing.T) {
	conv := New()

	tests := []struct {
		name   string
		items  *swagger.Items
		verify func(*testing.T, *openapi.Schema)
	}{
		{
			name: "simple string items",
			items: &swagger.Items{
				Type:   "string",
				Format: "uuid",
			},
			verify: func(t *testing.T, s *openapi.Schema) {
				if s.Type != "string" {
					t.Errorf("Type = %v, want string", s.Type)
				}
				if s.Format != "uuid" {
					t.Errorf("Format = %v, want uuid", s.Format)
				}
			},
		},
		{
			name: "items with enum",
			items: &swagger.Items{
				Type: "string",
				Enum: []interface{}{"red", "green", "blue"},
			},
			verify: func(t *testing.T, s *openapi.Schema) {
				if len(s.Enum) != 3 {
					t.Errorf("Enum count = %d, want 3", len(s.Enum))
				}
			},
		},
		{
			name: "items with min/max",
			items: &swagger.Items{
				Type:      "integer",
				Minimum:   &[]float64{0.0}[0],
				Maximum:   &[]float64{100.0}[0],
				MinLength: &[]int{5}[0],
				MaxLength: &[]int{50}[0],
			},
			verify: func(t *testing.T, s *openapi.Schema) {
				if s.Minimum != 0 {
					t.Errorf("Minimum = %v, want 0", s.Minimum)
				}
				if s.Maximum != 100 {
					t.Errorf("Maximum = %v, want 100", s.Maximum)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertItemsToSchema(tt.items)
			if result == nil {
				t.Fatal("Expected schema, got nil")
			}
			tt.verify(t, result)
		})
	}
}

// TestConvertParameterPropertiesToSchema tests parameter properties to schema conversion
func TestConvertParameterPropertiesToSchema(t *testing.T) {
	conv := New()

	tests := []struct {
		name   string
		param  *swagger.Parameter
		verify func(*testing.T, *openapi.Schema)
	}{
		{
			name: "simple type parameter",
			param: &swagger.Parameter{
				Type:    "string",
				Format:  "email",
				Default: "test@example.com",
			},
			verify: func(t *testing.T, s *openapi.Schema) {
				if s.Type != "string" {
					t.Errorf("Type = %v, want string", s.Type)
				}
				if s.Format != "email" {
					t.Errorf("Format = %v, want email", s.Format)
				}
			},
		},
		{
			name: "array parameter",
			param: &swagger.Parameter{
				Type: "array",
				Items: &swagger.Items{
					Type: "integer",
				},
			},
			verify: func(t *testing.T, s *openapi.Schema) {
				if s.Type != "array" {
					t.Errorf("Type = %v, want array", s.Type)
				}
				if s.Items == nil {
					t.Fatal("Items should not be nil")
				}
			},
		},
		{
			name: "parameter with constraints",
			param: &swagger.Parameter{
				Type:      "integer",
				Minimum:   &[]float64{1.0}[0],
				Maximum:   &[]float64{100.0}[0],
				MinLength: &[]int{5}[0],
				MaxLength: &[]int{50}[0],
			},
			verify: func(t *testing.T, s *openapi.Schema) {
				if s.Minimum != 1 {
					t.Errorf("Minimum = %v, want 1", s.Minimum)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertParameterPropertiesToSchema(tt.param)
			if result == nil {
				t.Fatal("Expected schema, got nil")
			}
			tt.verify(t, result)
		})
	}
}

// TestConvertPathItem tests path item conversion
func TestConvertPathItem(t *testing.T) {
	conv := New()

	pathItem := &openapi.PathItem{
		Summary:     "User operations",
		Description: "Operations for users",
		Get: &openapi.Operation{
			Summary:     "Get user",
			Description: "Get user by ID",
			OperationID: "getUser",
			Responses: openapi.Responses{
				"200": &openapi.Response{
					Description: "Success",
				},
			},
		},
		Post: &openapi.Operation{
			Summary:     "Create user",
			Description: "Create new user",
			OperationID: "createUser",
		},
	}

	result := conv.convertPathItem(pathItem)

	if result.Get == nil {
		t.Error("Expected Get operation")
	}
	if result.Post == nil {
		t.Error("Expected Post operation")
	}
}

// TestConvertOperation tests operation conversion
func TestConvertOperation(t *testing.T) {
	conv := New()

	op := &openapi.Operation{
		Summary:     "Get items",
		Description: "Get all items",
		OperationID: "getItems",
		Tags:        []string{"items"},
		Deprecated:  true,
		Responses: openapi.Responses{
			"200": &openapi.Response{
				Description: "Success",
				Content: map[string]*openapi.MediaType{
					"application/json": {
						Schema: &openapi.Schema{
							Type: "array",
						},
					},
				},
			},
		},
		Security: []openapi.SecurityRequirement{
			{"api_key": []string{}},
		},
	}

	result := conv.convertOperation(op)

	if result.Summary != "Get items" {
		t.Errorf("Summary = %v, want 'Get items'", result.Summary)
	}
	if !result.Deprecated {
		t.Error("Operation should be deprecated")
	}
}

// TestConvertParameters tests parameters array conversion V3 to V2
func TestConvertParameters(t *testing.T) {
	conv := New()

	params := []openapi.Parameter{
		{
			Name:        "id",
			In:          "path",
			Required:    true,
			Description: "Item ID",
			Schema: &openapi.Schema{
				Type: "string",
			},
		},
		{
			Name: "filter",
			In:   "query",
			Schema: &openapi.Schema{
				Type: "string",
			},
		},
	}

	result := conv.convertParameters(params)

	if len(result) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(result))
	}
	if result[0].Name != "id" {
		t.Errorf("Parameter[0].Name = %v, want id", result[0].Name)
	}
	if result[0].Type != "string" {
		t.Errorf("Parameter[0].Type = %v, want string", result[0].Type)
	}
}

// TestConvertRequestBodyToParameter tests request body to body parameter conversion
func TestConvertRequestBodyToParameter(t *testing.T) {
	conv := New()

	tests := []struct {
		name string
		rb   *openapi.RequestBody
		want string
	}{
		{
			name: "JSON request body",
			rb: &openapi.RequestBody{
				Description: "User data",
				Required:    true,
				Content: map[string]*openapi.MediaType{
					"application/json": {
						Schema: &openapi.Schema{
							Type: "object",
							Properties: map[string]*openapi.Schema{
								"name": {Type: "string"},
							},
						},
					},
				},
			},
			want: "body",
		},
		{
			name: "XML request body",
			rb: &openapi.RequestBody{
				Description: "XML data",
				Content: map[string]*openapi.MediaType{
					"application/xml": {
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
				},
			},
			want: "body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertRequestBodyToParameter(tt.rb)
			if result == nil {
				t.Fatal("Expected parameter, got nil")
			}
			if result.Name != tt.want {
				t.Errorf("Name = %v, want %v", result.Name, tt.want)
			}
			if result.In != "body" {
				t.Errorf("In = %v, want body", result.In)
			}
		})
	}
}

// TestConvertExamplesMediaType tests media type examples conversion
func TestConvertExamplesMediaType(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		examples map[string]*openapi.MediaType
		expected interface{}
	}{
		{
			name: "single media type with example",
			examples: map[string]*openapi.MediaType{
				"application/json": {
					Example: "test value",
				},
			},
			expected: "test value",
		},
		{
			name: "multiple media types",
			examples: map[string]*openapi.MediaType{
				"application/json": {Example: "value1"},
				"application/xml":  {Example: "value2"},
			},
			expected: nil, // First value found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertExamples(tt.examples)
			if tt.expected != nil && result == nil {
				t.Error("Expected example value, got nil")
			}
		})
	}
}

// TestIsNullableEdgeCases tests edge cases for nullable detection
func TestIsNullableEdgeCases(t *testing.T) {
	conv := New()

	tests := []struct {
		name     string
		typ      interface{}
		expected bool
	}{
		{
			name:     "nil type",
			typ:      nil,
			expected: false,
		},
		{
			name:     "non-array type",
			typ:      "string",
			expected: false,
		},
		{
			name:     "array with numbers",
			typ:      []interface{}{1, 2, 3},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.isNullable(tt.typ)
			if result != tt.expected {
				t.Errorf("isNullable(%v) = %v, want %v", tt.typ, result, tt.expected)
			}
		})
	}
}

// TestConvertInfoToV3 tests info conversion to V3
func TestConvertInfoToV3(t *testing.T) {
	conv := New()

	info := &swagger.Info{
		Title:          "Test API",
		Description:    "API description",
		Version:        "1.0.0",
		TermsOfService: "https://example.com/terms",
		Contact: &swagger.Contact{
			Name:  "Support",
			Email: "support@example.com",
			URL:   "https://example.com",
		},
		License: &swagger.License{
			Name: "MIT",
			URL:  "https://opensource.org/licenses/MIT",
		},
	}

	result := conv.convertInfoToV3(*info)

	if result.Title != "Test API" {
		t.Errorf("Title = %v, want 'Test API'", result.Title)
	}
	if result.Version != "1.0.0" {
		t.Errorf("Version = %v, want '1.0.0'", result.Version)
	}
	if result.Contact == nil {
		t.Error("Contact should not be nil")
	}
	if result.License == nil {
		t.Error("License should not be nil")
	}
}

// TestConvertPathItemToV3 tests path item conversion to V3
func TestConvertPathItemToV3(t *testing.T) {
	conv := New()

	pathItem := &swagger.PathItem{
		Get: &swagger.Operation{
			Summary:     "Get user",
			Description: "Get user by ID",
			OperationID: "getUser",
			Produces:    []string{"application/json"},
			Responses: swagger.Responses{
				"200": &swagger.Response{
					Description: "Success",
				},
			},
		},
		Post: &swagger.Operation{
			Summary:     "Create user",
			Description: "Create new user",
			Consumes:    []string{"application/json"},
		},
		Parameters: []*swagger.Parameter{
			{
				Name:     "id",
				In:       "path",
				Required: true,
				Type:     "string",
			},
		},
	}

	result := conv.convertPathItemToV3(pathItem)

	if result.Get == nil {
		t.Error("Expected Get operation")
	}
	if result.Post == nil {
		t.Error("Expected Post operation")
	}
	if len(result.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(result.Parameters))
	}
}

// TestConvertParameterDefinitionsToV3 tests parameter definitions conversion to V3
func TestConvertParameterDefinitionsToV3(t *testing.T) {
	conv := New()

	defs := map[string]*swagger.Parameter{
		"limitParam": {
			Name:     "limit",
			In:       "query",
			Type:     "integer",
			Required: false,
		},
		"pageParam": {
			Name:     "page",
			In:       "query",
			Type:     "integer",
			Required: false,
		},
	}

	result := conv.convertParameterDefinitionsToV3(defs)

	if len(result) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(result))
	}
	if _, ok := result["limitParam"]; !ok {
		t.Error("Expected limitParam in result")
	}
	if _, ok := result["pageParam"]; !ok {
		t.Error("Expected pageParam in result")
	}
}

// TestConvertResponseDefinitionsToV3 tests response definitions conversion to V3
func TestConvertResponseDefinitionsToV3(t *testing.T) {
	conv := New()

	defs := map[string]*swagger.Response{
		"ErrorResponse": {
			Description: "Error occurred",
			Schema: &swagger.Schema{
				Type: "object",
				Properties: map[string]*swagger.Schema{
					"message": {Type: "string"},
				},
			},
		},
		"SuccessResponse": {
			Description: "Operation successful",
		},
	}

	result := conv.convertResponseDefinitionsToV3(defs)

	if len(result) != 2 {
		t.Errorf("Expected 2 responses, got %d", len(result))
	}
	if _, ok := result["ErrorResponse"]; !ok {
		t.Error("Expected ErrorResponse in result")
	}
}

// TestConvertSecurityDefinitionsToV3 tests security definitions conversion to V3
func TestConvertSecurityDefinitionsToV3(t *testing.T) {
	conv := New()

	defs := map[string]*swagger.SecurityScheme{
		"api_key": {
			Type: "apiKey",
			Name: "X-API-Key",
			In:   "header",
		},
		"basic_auth": {
			Type: "basic",
		},
		"oauth": {
			Type:             "oauth2",
			Flow:             "implicit",
			AuthorizationURL: "https://example.com/oauth/authorize",
			Scopes: map[string]string{
				"read":  "Read access",
				"write": "Write access",
			},
		},
	}

	result := conv.convertSecurityDefinitionsToV3(defs)

	if len(result) != 3 {
		t.Errorf("Expected 3 security schemes, got %d", len(result))
	}
	if _, ok := result["api_key"]; !ok {
		t.Error("Expected api_key in result")
	}
	if _, ok := result["basic_auth"]; !ok {
		t.Error("Expected basic_auth in result")
	}
	if _, ok := result["oauth"]; !ok {
		t.Error("Expected oauth in result")
	}

	// Verify api_key conversion
	if result["api_key"].Type != "apiKey" {
		t.Errorf("api_key type = %v, want apiKey", result["api_key"].Type)
	}

	// Verify basic_auth conversion
	if result["basic_auth"].Type != "http" {
		t.Errorf("basic_auth type = %v, want http", result["basic_auth"].Type)
	}
	if result["basic_auth"].Scheme != "basic" {
		t.Errorf("basic_auth scheme = %v, want basic", result["basic_auth"].Scheme)
	}

	// Verify oauth conversion
	if result["oauth"].Type != "oauth2" {
		t.Errorf("oauth type = %v, want oauth2", result["oauth"].Type)
	}
}

// TestConvertSecuritySchemeToV3 tests individual security scheme conversion to V3
func TestConvertSecuritySchemeToV3(t *testing.T) {
	conv := New()

	tests := []struct {
		name   string
		scheme *swagger.SecurityScheme
		verify func(*testing.T, *openapi.SecurityScheme)
	}{
		{
			name: "basic auth",
			scheme: &swagger.SecurityScheme{
				Type:        "basic",
				Description: "Basic authentication",
			},
			verify: func(t *testing.T, s *openapi.SecurityScheme) {
				if s.Type != "http" {
					t.Errorf("Type = %v, want http", s.Type)
				}
				if s.Scheme != "basic" {
					t.Errorf("Scheme = %v, want basic", s.Scheme)
				}
			},
		},
		{
			name: "apiKey in header",
			scheme: &swagger.SecurityScheme{
				Type: "apiKey",
				Name: "Authorization",
				In:   "header",
			},
			verify: func(t *testing.T, s *openapi.SecurityScheme) {
				if s.Type != "apiKey" {
					t.Errorf("Type = %v, want apiKey", s.Type)
				}
				if s.In != "header" {
					t.Errorf("In = %v, want header", s.In)
				}
			},
		},
		{
			name: "oauth2 implicit flow",
			scheme: &swagger.SecurityScheme{
				Type:             "oauth2",
				Flow:             "implicit",
				AuthorizationURL: "https://example.com/oauth",
				Scopes: map[string]string{
					"read": "Read access",
				},
			},
			verify: func(t *testing.T, s *openapi.SecurityScheme) {
				if s.Type != "oauth2" {
					t.Errorf("Type = %v, want oauth2", s.Type)
				}
				if s.Flows == nil {
					t.Fatal("Flows should not be nil")
				}
				if s.Flows.Implicit == nil {
					t.Fatal("Implicit flow should not be nil")
				}
			},
		},
		{
			name: "oauth2 password flow",
			scheme: &swagger.SecurityScheme{
				Type:     "oauth2",
				Flow:     "password",
				TokenURL: "https://example.com/token",
				Scopes: map[string]string{
					"admin": "Admin access",
				},
			},
			verify: func(t *testing.T, s *openapi.SecurityScheme) {
				if s.Flows == nil || s.Flows.Password == nil {
					t.Fatal("Password flow should not be nil")
				}
			},
		},
		{
			name: "oauth2 application flow",
			scheme: &swagger.SecurityScheme{
				Type:     "oauth2",
				Flow:     "application",
				TokenURL: "https://example.com/token",
				Scopes:   map[string]string{},
			},
			verify: func(t *testing.T, s *openapi.SecurityScheme) {
				if s.Flows == nil || s.Flows.ClientCredentials == nil {
					t.Fatal("ClientCredentials flow should not be nil")
				}
			},
		},
		{
			name: "oauth2 accessCode flow",
			scheme: &swagger.SecurityScheme{
				Type:             "oauth2",
				Flow:             "accessCode",
				AuthorizationURL: "https://example.com/oauth",
				TokenURL:         "https://example.com/token",
				Scopes:           map[string]string{"openid": "OpenID"},
			},
			verify: func(t *testing.T, s *openapi.SecurityScheme) {
				if s.Flows == nil || s.Flows.AuthorizationCode == nil {
					t.Fatal("AuthorizationCode flow should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.convertSecuritySchemeToV3(tt.scheme)
			if result == nil {
				t.Fatal("Expected security scheme, got nil")
			}
			tt.verify(t, result)
		})
	}
}

// TestConvertOperationToV3 tests operation conversion to V3
func TestConvertOperationToV3(t *testing.T) {
	conv := New()

	op := &swagger.Operation{
		Summary:     "List items",
		Description: "Get all items",
		OperationID: "listItems",
		Tags:        []string{"items"},
		Deprecated:  false,
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json", "application/xml"},
		Parameters: []*swagger.Parameter{
			{
				Name: "limit",
				In:   "query",
				Type: "integer",
			},
		},
		Responses: swagger.Responses{
			"200": &swagger.Response{
				Description: "Success",
				Schema: &swagger.Schema{
					Type: "array",
				},
			},
		},
	}

	result := conv.convertOperationToV3(op)

	if result.Summary != "List items" {
		t.Errorf("Summary = %v, want 'List items'", result.Summary)
	}
	if result.OperationID != "listItems" {
		t.Errorf("OperationID = %v, want 'listItems'", result.OperationID)
	}
	if len(result.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(result.Parameters))
	}
	if result.Responses == nil {
		t.Fatal("Responses should not be nil")
	}
}

// TestRoundtripV2ToV3ToV2 tests complete roundtrip conversion V2->V3->V2
func TestRoundtripV2ToV3ToV2(t *testing.T) {
	tests := []struct {
		name     string
		original *swagger.Swagger
		verify   func(t *testing.T, original, final *swagger.Swagger)
	}{
		{
			name: "basic API with paths",
			original: &swagger.Swagger{
				Swagger: "2.0",
				Info: swagger.Info{
					Title:          "Roundtrip Test API",
					Description:    "Testing roundtrip conversion",
					Version:        "2.0.0",
					TermsOfService: "https://example.com/terms",
					Contact: &swagger.Contact{
						Name:  "API Support",
						Email: "support@example.com",
						URL:   "https://example.com/support",
					},
					License: &swagger.License{
						Name: "MIT",
						URL:  "https://opensource.org/licenses/MIT",
					},
				},
				Host:     "api.example.com",
				BasePath: "/v2",
				Schemes:  []string{"https", "http"},
				Paths: map[string]*swagger.PathItem{
					"/users": {
						Get: &swagger.Operation{
							Summary:     "List users",
							Description: "Get all users",
							OperationID: "listUsers",
							Tags:        []string{"users"},
							Produces:    []string{"application/json"},
							Responses: swagger.Responses{
								"200": {
									Description: "Success",
									Schema: &swagger.Schema{
										Type: "array",
										Items: &swagger.Schema{
											Ref: "#/definitions/User",
										},
									},
								},
							},
						},
						Post: &swagger.Operation{
							Summary:     "Create user",
							Description: "Create a new user",
							OperationID: "createUser",
							Tags:        []string{"users"},
							Consumes:    []string{"application/json"},
							Produces:    []string{"application/json"},
							Parameters: []*swagger.Parameter{
								{
									Name:     "body",
									In:       "body",
									Required: true,
									Schema: &swagger.Schema{
										Ref: "#/definitions/User",
									},
								},
							},
							Responses: swagger.Responses{
								"201": {
									Description: "Created",
									Schema: &swagger.Schema{
										Ref: "#/definitions/User",
									},
								},
							},
						},
					},
				},
				Definitions: map[string]*swagger.Schema{
					"User": {
						Type: "object",
						Properties: map[string]*swagger.Schema{
							"id": {
								Type:   "integer",
								Format: "int64",
							},
							"name": {
								Type: "string",
							},
							"email": {
								Type:   "string",
								Format: "email",
							},
						},
						Required: []string{"name", "email"},
					},
				},
			},
			verify: func(t *testing.T, original, final *swagger.Swagger) {
				if final.Info.Title != original.Info.Title {
					t.Errorf("Info.Title = %v, want %v", final.Info.Title, original.Info.Title)
				}
				if final.Info.Version != original.Info.Version {
					t.Errorf("Info.Version = %v, want %v", final.Info.Version, original.Info.Version)
				}
				if final.Host != original.Host {
					t.Errorf("Host = %v, want %v", final.Host, original.Host)
				}
				if final.BasePath != original.BasePath {
					t.Errorf("BasePath = %v, want %v", final.BasePath, original.BasePath)
				}
				if len(final.Paths) != len(original.Paths) {
					t.Errorf("Paths count = %d, want %d", len(final.Paths), len(original.Paths))
				}
				if len(final.Definitions) != len(original.Definitions) {
					t.Errorf("Definitions count = %d, want %d", len(final.Definitions), len(original.Definitions))
				}
			},
		},
		{
			name: "API with security schemes",
			original: &swagger.Swagger{
				Swagger: "2.0",
				Info: swagger.Info{
					Title:   "Secure API",
					Version: "1.0.0",
				},
				Host:     "secure.example.com",
				BasePath: "/api",
				Schemes:  []string{"https"},
				Paths: map[string]*swagger.PathItem{
					"/protected": {
						Get: &swagger.Operation{
							Summary:  "Protected endpoint",
							Security: []swagger.SecurityRequirement{{"apiKey": {}}},
							Responses: swagger.Responses{
								"200": {Description: "OK"},
							},
						},
					},
				},
				SecurityDefinitions: map[string]*swagger.SecurityScheme{
					"apiKey": {
						Type: "apiKey",
						Name: "X-API-Key",
						In:   "header",
					},
					"oauth2": {
						Type:             "oauth2",
						Flow:             "accessCode",
						AuthorizationURL: "https://auth.example.com/authorize",
						TokenURL:         "https://auth.example.com/token",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
				},
			},
			verify: func(t *testing.T, original, final *swagger.Swagger) {
				if len(final.SecurityDefinitions) != len(original.SecurityDefinitions) {
					t.Errorf("SecurityDefinitions count = %d, want %d",
						len(final.SecurityDefinitions), len(original.SecurityDefinitions))
				}
				if final.SecurityDefinitions["apiKey"].Type != "apiKey" {
					t.Error("apiKey security scheme not preserved")
				}
			},
		},
		{
			name: "API with parameters and responses",
			original: &swagger.Swagger{
				Swagger: "2.0",
				Info: swagger.Info{
					Title:   "Parameters API",
					Version: "1.0.0",
				},
				Host:     "params.example.com",
				BasePath: "/",
				Schemes:  []string{"https"},
				Paths: map[string]*swagger.PathItem{
					"/items/{id}": {
						Parameters: []*swagger.Parameter{
							{
								Name:     "id",
								In:       "path",
								Required: true,
								Type:     "string",
							},
						},
						Get: &swagger.Operation{
							Summary: "Get item",
							Parameters: []*swagger.Parameter{
								{
									Name: "include",
									In:   "query",
									Type: "array",
									Items: &swagger.Items{
										Type: "string",
									},
									CollectionFormat: "csv",
								},
							},
							Responses: swagger.Responses{
								"200": {
									Description: "Success",
									Schema: &swagger.Schema{
										Ref: "#/definitions/Item",
									},
								},
								"404": {
									Description: "Not found",
								},
							},
						},
					},
				},
				Definitions: map[string]*swagger.Schema{
					"Item": {
						Type: "object",
						Properties: map[string]*swagger.Schema{
							"id":   {Type: "string"},
							"name": {Type: "string"},
						},
					},
				},
			},
			verify: func(t *testing.T, original, final *swagger.Swagger) {
				origPath := original.Paths["/items/{id}"]
				finalPath := final.Paths["/items/{id}"]
				if finalPath == nil {
					t.Fatal("Path /items/{id} not preserved")
				}
				if len(finalPath.Parameters) != len(origPath.Parameters) {
					t.Errorf("Path parameters count = %d, want %d",
						len(finalPath.Parameters), len(origPath.Parameters))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: V2 -> V3
			conv1 := New()
			v3Spec, err := conv1.ConvertToV3(tt.original)
			if err != nil {
				t.Fatalf("V2->V3 conversion failed: %v", err)
			}
			if v3Spec == nil {
				t.Fatal("V3 spec is nil")
			}

			// Step 2: V3 -> V2
			conv2 := New()
			v2Spec, err := conv2.ConvertToV2(v3Spec)
			if err != nil {
				t.Fatalf("V3->V2 conversion failed: %v", err)
			}
			if v2Spec == nil {
				t.Fatal("Final V2 spec is nil")
			}

			// Step 3: Verify
			tt.verify(t, tt.original, v2Spec)
		})
	}
}

// TestRoundtripV3ToV2ToV3 tests complete roundtrip conversion V3->V2->V3
func TestRoundtripV3ToV2ToV3(t *testing.T) {
	tests := []struct {
		name     string
		original *openapi.OpenAPI
		verify   func(t *testing.T, original, final *openapi.OpenAPI)
	}{
		{
			name: "OpenAPI 3.0 with components",
			original: &openapi.OpenAPI{
				OpenAPI: "3.0.0",
				Info: openapi.Info{
					Title:       "Component API",
					Version:     "1.0.0",
					Description: "API with reusable components",
				},
				Servers: []openapi.Server{
					{
						URL:         "https://api.example.com/v3",
						Description: "Production server",
					},
				},
				Paths: map[string]*openapi.PathItem{
					"/products": {
						Get: &openapi.Operation{
							Summary:     "List products",
							OperationID: "listProducts",
							Responses: map[string]*openapi.Response{
								"200": {
									Description: "Success",
									Content: map[string]*openapi.MediaType{
										"application/json": {
											Schema: &openapi.Schema{
												Type: "array",
												Items: &openapi.Schema{
													Ref: "#/components/schemas/Product",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Components: &openapi.Components{
					Schemas: map[string]*openapi.Schema{
						"Product": {
							Type: "object",
							Properties: map[string]*openapi.Schema{
								"id":    {Type: "integer"},
								"name":  {Type: "string"},
								"price": {Type: "number", Format: "float"},
							},
							Required: []string{"name", "price"},
						},
					},
				},
			},
			verify: func(t *testing.T, original, final *openapi.OpenAPI) {
				if final.Info.Title != original.Info.Title {
					t.Errorf("Info.Title = %v, want %v", final.Info.Title, original.Info.Title)
				}
				if len(final.Paths) != len(original.Paths) {
					t.Errorf("Paths count = %d, want %d", len(final.Paths), len(original.Paths))
				}
				if final.Components == nil {
					t.Fatal("Components should not be nil")
				}
				if len(final.Components.Schemas) != len(original.Components.Schemas) {
					t.Errorf("Schemas count = %d, want %d",
						len(final.Components.Schemas), len(original.Components.Schemas))
				}
			},
		},
		{
			name: "OpenAPI 3.1 with webhooks",
			original: &openapi.OpenAPI{
				OpenAPI: "3.1.0",
				Info: openapi.Info{
					Title:   "Webhook API",
					Version: "1.0.0",
				},
				Servers: []openapi.Server{
					{URL: "https://webhook.example.com"},
				},
				Paths: map[string]*openapi.PathItem{
					"/subscribe": {
						Post: &openapi.Operation{
							Summary: "Subscribe to events",
							Responses: map[string]*openapi.Response{
								"200": {Description: "Subscribed"},
							},
						},
					},
				},
				Webhooks: map[string]*openapi.PathItem{
					"newOrder": {
						Post: &openapi.Operation{
							Summary:     "New order notification",
							Description: "Called when a new order is created",
							RequestBody: &openapi.RequestBody{
								Content: map[string]*openapi.MediaType{
									"application/json": {
										Schema: &openapi.Schema{
											Type: "object",
											Properties: map[string]*openapi.Schema{
												"orderId": {Type: "string"},
												"amount":  {Type: "number"},
											},
										},
									},
								},
							},
							Responses: map[string]*openapi.Response{
								"200": {Description: "Acknowledged"},
							},
						},
					},
				},
			},
			verify: func(t *testing.T, original, final *openapi.OpenAPI) {
				if final.Info.Title != original.Info.Title {
					t.Errorf("Info.Title = %v, want %v", final.Info.Title, original.Info.Title)
				}
				// Webhooks may be lost in V2 conversion, but should be preserved if possible
				if len(final.Paths) == 0 {
					t.Error("All paths were lost in roundtrip")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: V3 -> V2
			conv1 := New()
			v2Spec, err := conv1.ConvertToV2(tt.original)
			if err != nil {
				t.Fatalf("V3->V2 conversion failed: %v", err)
			}
			if v2Spec == nil {
				t.Fatal("V2 spec is nil")
			}

			// Step 2: V2 -> V3
			conv2 := New()
			v3Spec, err := conv2.ConvertToV3(v2Spec)
			if err != nil {
				t.Fatalf("V2->V3 conversion failed: %v", err)
			}
			if v3Spec == nil {
				t.Fatal("Final V3 spec is nil")
			}

			// Step 3: Verify
			tt.verify(t, tt.original, v3Spec)
		})
	}
}

// TestRoundtripDataIntegrity verifies data integrity across multiple roundtrips
func TestRoundtripDataIntegrity(t *testing.T) {
	original := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:       "Data Integrity Test",
			Description: "Testing data preservation",
			Version:     "1.0.0",
		},
		Host:     "integrity.example.com",
		BasePath: "/api",
		Schemes:  []string{"https"},
		Tags: []swagger.Tag{
			{Name: "users", Description: "User operations"},
			{Name: "admin", Description: "Admin operations"},
		},
		Paths: map[string]*swagger.PathItem{
			"/health": {
				Get: &swagger.Operation{
					Summary:  "Health check",
					Tags:     []string{"admin"},
					Produces: []string{"application/json"},
					Responses: swagger.Responses{
						"200": {Description: "Healthy"},
					},
				},
			},
		},
	}

	// Perform 3 complete roundtrips
	current := original
	for i := 0; i < 3; i++ {
		t.Run(strings.Join([]string{"roundtrip", string(rune('1' + i))}, "_"), func(t *testing.T) {
			// V2 -> V3
			conv1 := New()
			v3, err := conv1.ConvertToV3(current)
			if err != nil {
				t.Fatalf("Roundtrip %d: V2->V3 failed: %v", i+1, err)
			}

			// V3 -> V2
			conv2 := New()
			v2, err := conv2.ConvertToV2(v3)
			if err != nil {
				t.Fatalf("Roundtrip %d: V3->V2 failed: %v", i+1, err)
			}

			// Verify core data preserved
			if v2.Info.Title != original.Info.Title {
				t.Errorf("Roundtrip %d: Title changed to %v", i+1, v2.Info.Title)
			}
			if v2.Host != original.Host {
				t.Errorf("Roundtrip %d: Host changed to %v", i+1, v2.Host)
			}
			if len(v2.Paths) != len(original.Paths) {
				t.Errorf("Roundtrip %d: Paths count changed to %d", i+1, len(v2.Paths))
			}

			current = v2
		})
	}
}

// Benchmark tests
func BenchmarkConvertToV2Simple(b *testing.B) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.0.0",
		Info: openapi.Info{
			Title:   "Benchmark API",
			Version: "1.0.0",
		},
		Servers: []openapi.Server{{URL: "https://api.example.com"}},
		Paths: map[string]*openapi.PathItem{
			"/test": {
				Get: &openapi.Operation{
					Summary: "Test endpoint",
					Responses: map[string]*openapi.Response{
						"200": {Description: "OK"},
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv := New()
		_, err := conv.ConvertToV2(spec)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertToV3Simple(b *testing.B) {
	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Benchmark API",
			Version: "1.0.0",
		},
		Host:     "api.example.com",
		BasePath: "/v1",
		Schemes:  []string{"https"},
		Paths: map[string]*swagger.PathItem{
			"/test": {
				Get: &swagger.Operation{
					Summary: "Test endpoint",
					Responses: swagger.Responses{
						"200": {Description: "OK"},
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv := New()
		_, err := conv.ConvertToV3(spec)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertToV2Complex(b *testing.B) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.0.0",
		Info: openapi.Info{
			Title:       "Complex API",
			Description: "API with many endpoints and schemas",
			Version:     "2.0.0",
		},
		Servers: []openapi.Server{{URL: "https://api.example.com/v2"}},
		Paths:   make(map[string]*openapi.PathItem),
		Components: &openapi.Components{
			Schemas: make(map[string]*openapi.Schema),
		},
	}

	// Create 50 paths with operations
	for i := 0; i < 50; i++ {
		path := strings.Join([]string{"/endpoint", string(rune('0' + i))}, "")
		spec.Paths[path] = &openapi.PathItem{
			Get: &openapi.Operation{
				Summary:     strings.Join([]string{"Get endpoint", string(rune('0' + i))}, " "),
				OperationID: strings.Join([]string{"get", string(rune('0' + i))}, ""),
				Responses: map[string]*openapi.Response{
					"200": {
						Description: "Success",
						Content: map[string]*openapi.MediaType{
							"application/json": {
								Schema: &openapi.Schema{Type: "object"},
							},
						},
					},
				},
			},
			Post: &openapi.Operation{
				Summary:     strings.Join([]string{"Post endpoint", string(rune('0' + i))}, " "),
				OperationID: strings.Join([]string{"post", string(rune('0' + i))}, ""),
				RequestBody: &openapi.RequestBody{
					Content: map[string]*openapi.MediaType{
						"application/json": {
							Schema: &openapi.Schema{Type: "object"},
						},
					},
				},
				Responses: map[string]*openapi.Response{
					"201": {Description: "Created"},
				},
			},
		}
	}

	// Create 20 schemas
	for i := 0; i < 20; i++ {
		schemaName := strings.Join([]string{"Schema", string(rune('0' + i))}, "")
		spec.Components.Schemas[schemaName] = &openapi.Schema{
			Type: "object",
			Properties: map[string]*openapi.Schema{
				"id":   {Type: "integer"},
				"name": {Type: "string"},
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv := New()
		_, err := conv.ConvertToV2(spec)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertToV3Complex(b *testing.B) {
	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:       "Complex API",
			Description: "API with many endpoints and definitions",
			Version:     "2.0.0",
		},
		Host:        "api.example.com",
		BasePath:    "/v2",
		Schemes:     []string{"https"},
		Paths:       make(map[string]*swagger.PathItem),
		Definitions: make(map[string]*swagger.Schema),
	}

	// Create 50 paths with operations
	for i := 0; i < 50; i++ {
		path := strings.Join([]string{"/endpoint", string(rune('0' + i))}, "")
		spec.Paths[path] = &swagger.PathItem{
			Get: &swagger.Operation{
				Summary:     strings.Join([]string{"Get endpoint", string(rune('0' + i))}, " "),
				OperationID: strings.Join([]string{"get", string(rune('0' + i))}, ""),
				Produces:    []string{"application/json"},
				Responses: swagger.Responses{
					"200": {
						Description: "Success",
						Schema:      &swagger.Schema{Type: "object"},
					},
				},
			},
			Post: &swagger.Operation{
				Summary:     strings.Join([]string{"Post endpoint", string(rune('0' + i))}, " "),
				OperationID: strings.Join([]string{"post", string(rune('0' + i))}, ""),
				Consumes:    []string{"application/json"},
				Produces:    []string{"application/json"},
				Parameters: []*swagger.Parameter{
					{
						Name:     "body",
						In:       "body",
						Required: true,
						Schema:   &swagger.Schema{Type: "object"},
					},
				},
				Responses: swagger.Responses{
					"201": {Description: "Created"},
				},
			},
		}
	}

	// Create 20 definitions
	for i := 0; i < 20; i++ {
		defName := strings.Join([]string{"Definition", string(rune('0' + i))}, "")
		spec.Definitions[defName] = &swagger.Schema{
			Type: "object",
			Properties: map[string]*swagger.Schema{
				"id":   {Type: "integer"},
				"name": {Type: "string"},
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv := New()
		_, err := conv.ConvertToV3(spec)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRoundtripV2ToV3ToV2(b *testing.B) {
	spec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Roundtrip Benchmark",
			Version: "1.0.0",
		},
		Host:     "api.example.com",
		BasePath: "/v1",
		Schemes:  []string{"https"},
		Paths: map[string]*swagger.PathItem{
			"/users": {
				Get: &swagger.Operation{
					Summary: "List users",
					Responses: swagger.Responses{
						"200": {Description: "OK"},
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv1 := New()
		v3, err := conv1.ConvertToV3(spec)
		if err != nil {
			b.Fatal(err)
		}

		conv2 := New()
		_, err = conv2.ConvertToV2(v3)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRoundtripV3ToV2ToV3(b *testing.B) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.0.0",
		Info: openapi.Info{
			Title:   "Roundtrip Benchmark",
			Version: "1.0.0",
		},
		Servers: []openapi.Server{{URL: "https://api.example.com/v1"}},
		Paths: map[string]*openapi.PathItem{
			"/users": {
				Get: &openapi.Operation{
					Summary: "List users",
					Responses: map[string]*openapi.Response{
						"200": {Description: "OK"},
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv1 := New()
		v2, err := conv1.ConvertToV2(spec)
		if err != nil {
			b.Fatal(err)
		}

		conv2 := New()
		_, err = conv2.ConvertToV3(v2)
		if err != nil {
			b.Fatal(err)
		}
	}
}
