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
