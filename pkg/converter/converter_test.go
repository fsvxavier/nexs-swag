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
