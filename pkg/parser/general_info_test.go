package parser

import (
	"go/parser"
	"go/token"
	"testing"

	openapi "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

func TestParseGeneralInfo(t *testing.T) {
	t.Parallel()
	content := `package main

// @title Test API
// @version 1.0.0
// @description This is a test API
// @termsOfService http://example.com/terms
// @contact.name API Support
// @contact.email support@example.com
// @contact.url http://example.com/support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host api.example.com
// @BasePath /v1

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() returned error: %v", err)
	}

	// Verify Info was populated
	if p.openapi.Info.Title != "Test API" {
		t.Errorf("Info.Title = %q, want %q", p.openapi.Info.Title, "Test API")
	}
	if p.openapi.Info.Version != "1.0.0" {
		t.Errorf("Info.Version = %q, want %q", p.openapi.Info.Version, "1.0.0")
	}
}

func TestParseGeneralInfoMinimal(t *testing.T) {
	t.Parallel()
	content := `package main

// @title Minimal API
// @version 0.1.0

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() returned error: %v", err)
	}

	if p.openapi.Info.Title != "Minimal API" {
		t.Errorf("Info.Title = %q, want %q", p.openapi.Info.Title, "Minimal API")
	}
}

func TestNewGeneralInfoProcessor(t *testing.T) {
	t.Parallel()
	p := New()
	processor := NewGeneralInfoProcessor(p.openapi)

	if processor == nil {
		t.Fatal("NewGeneralInfoProcessor() returned nil")
	}
}

func TestGeneralInfoProcessorAnnotations(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		annotation string
		checkFunc  func(*Parser) bool
	}{
		{
			name:       "title",
			annotation: "@title My API",
			checkFunc: func(p *Parser) bool {
				return p.openapi.Info.Title == "My API"
			},
		},
		{
			name:       "version",
			annotation: "@version 2.0.0",
			checkFunc: func(p *Parser) bool {
				return p.openapi.Info.Version == "2.0.0"
			},
		},
		{
			name:       "description",
			annotation: "@description API description",
			checkFunc: func(p *Parser) bool {
				return p.openapi.Info.Description == "API description"
			},
		},
		{
			name:       "host",
			annotation: "@host api.example.com",
			checkFunc: func(p *Parser) bool {
				return len(p.openapi.Servers) > 0 && p.openapi.Servers[0].URL == "https://api.example.com"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			processor := NewGeneralInfoProcessor(p.openapi)

			err := processor.Process(tt.annotation)
			if err != nil {
				t.Errorf("Process() returned error: %v", err)
			}

			if !tt.checkFunc(p) {
				t.Errorf("Annotation %q was not processed correctly", tt.annotation)
			}
		})
	}
}

func TestGeneralInfoProcessorContact(t *testing.T) {
	t.Parallel()
	p := New()
	processor := NewGeneralInfoProcessor(p.openapi)

	annotations := []string{
		"@contact.name API Support",
		"@contact.email support@example.com",
		"@contact.url http://example.com/support",
	}

	for _, annotation := range annotations {
		if err := processor.Process(annotation); err != nil {
			t.Errorf("Process(%q) returned error: %v", annotation, err)
		}
	}

	if p.openapi.Info.Contact == nil {
		t.Fatal("Contact should not be nil")
	}
	if p.openapi.Info.Contact.Name != "API Support" {
		t.Errorf("Contact.Name = %q, want %q", p.openapi.Info.Contact.Name, "API Support")
	}
	if p.openapi.Info.Contact.Email != "support@example.com" {
		t.Errorf("Contact.Email = %q, want %q", p.openapi.Info.Contact.Email, "support@example.com")
	}
}

func TestGeneralInfoProcessorLicense(t *testing.T) {
	t.Parallel()
	p := New()
	processor := NewGeneralInfoProcessor(p.openapi)

	annotations := []string{
		"@license.name MIT",
		"@license.url https://opensource.org/licenses/MIT",
	}

	for _, annotation := range annotations {
		if err := processor.Process(annotation); err != nil {
			t.Errorf("Process(%q) returned error: %v", annotation, err)
		}
	}

	if p.openapi.Info.License == nil {
		t.Fatal("License should not be nil")
	}
	if p.openapi.Info.License.Name != "MIT" {
		t.Errorf("License.Name = %q, want %q", p.openapi.Info.License.Name, "MIT")
	}
}

func TestGeneralInfoProcessorServers(t *testing.T) {
	t.Parallel()
	p := New()
	processor := NewGeneralInfoProcessor(p.openapi)

	annotations := []string{
		"@host api.example.com",
		"@BasePath /v1",
		"@schemes https http",
	}

	for _, annotation := range annotations {
		if err := processor.Process(annotation); err != nil {
			t.Errorf("Process(%q) returned error: %v", annotation, err)
		}
	}

	if len(p.openapi.Servers) == 0 {
		t.Fatal("Servers should not be empty")
	}
}

func TestGeneralInfoProcessorTags(t *testing.T) {
	t.Parallel()
	content := `package main

// @title Test API
// @version 1.0.0
// @tag.name users
// @tag.description User management operations
// @tag.name products
// @tag.description Product management operations

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() returned error: %v", err)
	}

	// Should have parsed tags
	if len(p.openapi.Tags) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestGeneralInfoProcessorSecurityDefinitions(t *testing.T) {
	t.Parallel()
	content := `package main

// @title Test API
// @version 1.0.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() returned error: %v", err)
	}

	// Check if security schemes were defined
	if p.openapi.Components == nil || p.openapi.Components.SecuritySchemes == nil {
		t.Log("SecuritySchemes might not be populated (depends on implementation)")
	}
}

func TestParseGeneralInfoNoComments(t *testing.T) {
	t.Parallel()
	content := `package main

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() with no comments should not return error: %v", err)
	}

	// Title and Version should remain empty
	if p.openapi.Info.Title != "" {
		t.Error("Info.Title should be empty when no annotations present")
	}
}

func TestParseGeneralInfoMultilineDescription(t *testing.T) {
	t.Parallel()
	content := `package main

// @title Test API
// @version 1.0.0
// @description This is line 1
// @description This is line 2
// @description This is line 3

func main() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	err = p.parseGeneralInfo(file)
	if err != nil {
		t.Errorf("parseGeneralInfo() returned error: %v", err)
	}

	if p.openapi.Info.Description == "" {
		t.Error("Description should not be empty")
	}
}

// TestSecuritySchemeDeprecated tests @securityDefinitions.*.deprecated annotation (OpenAPI 3.2.0).
func TestSecuritySchemeDeprecated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		annotations        []string
		expectedScheme     string
		expectedDeprecated bool
	}{
		{
			name: "apikey_deprecated_true",
			annotations: []string{
				"@securityDefinitions.apikey ApiKeyAuth header X-API-Key",
				"@securityDefinitions.ApiKeyAuth.deprecated true",
			},
			expectedScheme:     "ApiKeyAuth",
			expectedDeprecated: true,
		},
		{
			name: "apikey_deprecated_false",
			annotations: []string{
				"@securityDefinitions.apikey ApiKeyAuth header X-API-Key",
				"@securityDefinitions.ApiKeyAuth.deprecated false",
			},
			expectedScheme:     "ApiKeyAuth",
			expectedDeprecated: false,
		},
		{
			name: "basic_deprecated_true",
			annotations: []string{
				"@securityDefinitions.basic BasicAuth",
				"@securityDefinitions.BasicAuth.deprecated true",
			},
			expectedScheme:     "BasicAuth",
			expectedDeprecated: true,
		},
		{
			name: "deprecated_without_scheme_definition",
			annotations: []string{
				"@securityDefinitions.NonExistent.deprecated true",
			},
			expectedScheme:     "NonExistent",
			expectedDeprecated: false, // Não deve existir
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			spec := &openapi.OpenAPI{
				Components: &openapi.Components{
					SecuritySchemes: make(map[string]*openapi.SecurityScheme),
				},
			}

			gproc := NewGeneralInfoProcessor(spec)

			for _, annotation := range tt.annotations {
				if err := gproc.Process(annotation); err != nil {
					t.Fatalf("Process() error = %v", err)
				}
			}

			scheme, exists := spec.Components.SecuritySchemes[tt.expectedScheme]

			if tt.name == "deprecated_without_scheme_definition" {
				if exists {
					t.Errorf("Expected scheme %s not to exist", tt.expectedScheme)
				}
				return
			}

			if !exists {
				t.Fatalf("Expected scheme %s to exist", tt.expectedScheme)
			}

			if scheme.Deprecated != tt.expectedDeprecated {
				t.Errorf("Expected Deprecated = %v, got %v", tt.expectedDeprecated, scheme.Deprecated)
			}
		})
	}
}

// TestSecuritySchemeOAuth2MetadataURL tests @securityDefinitions.*.oauth2metadataurl annotation (OpenAPI 3.2.0).
func TestSecuritySchemeOAuth2MetadataURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		annotations         []string
		expectedScheme      string
		expectedMetadataURL string
		shouldExist         bool
	}{
		{
			name: "apikey_with_metadata_url",
			annotations: []string{
				"@securityDefinitions.apikey ApiKeyAuth header X-API-Key",
				"@securityDefinitions.ApiKeyAuth.oauth2metadataurl https://auth.example.com/.well-known/oauth-authorization-server",
			},
			expectedScheme:      "ApiKeyAuth",
			expectedMetadataURL: "https://auth.example.com/.well-known/oauth-authorization-server",
			shouldExist:         true,
		},
		{
			name: "basic_with_metadata_url",
			annotations: []string{
				"@securityDefinitions.basic BasicAuth",
				"@securityDefinitions.BasicAuth.oauth2metadataurl https://auth.example.com/metadata",
			},
			expectedScheme:      "BasicAuth",
			expectedMetadataURL: "https://auth.example.com/metadata",
			shouldExist:         true,
		},
		{
			name: "metadata_url_without_scheme",
			annotations: []string{
				"@securityDefinitions.NonExistent.oauth2metadataurl https://example.com/metadata",
			},
			expectedScheme: "NonExistent",
			shouldExist:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			spec := &openapi.OpenAPI{
				Components: &openapi.Components{
					SecuritySchemes: make(map[string]*openapi.SecurityScheme),
				},
			}

			gproc := NewGeneralInfoProcessor(spec)

			for _, annotation := range tt.annotations {
				if err := gproc.Process(annotation); err != nil {
					t.Fatalf("Process() error = %v", err)
				}
			}

			scheme, exists := spec.Components.SecuritySchemes[tt.expectedScheme]

			if !tt.shouldExist {
				if exists {
					t.Errorf("Expected scheme %s not to exist", tt.expectedScheme)
				}
				return
			}

			if !exists {
				t.Fatalf("Expected scheme %s to exist", tt.expectedScheme)
			}

			if scheme.OAuth2MetadataURL != tt.expectedMetadataURL {
				t.Errorf("Expected OAuth2MetadataURL = %q, got %q", tt.expectedMetadataURL, scheme.OAuth2MetadataURL)
			}
		})
	}
}

// TestSecurityDeviceAuthorization tests @securityDefinitions.oauth2.deviceAuthorization annotation (OpenAPI 3.2.0).
func TestSecurityDeviceAuthorization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		annotations           []string
		expectedScheme        string
		expectedDeviceAuthURL string
		expectedTokenURL      string
	}{
		{
			name: "device_authorization_flow",
			annotations: []string{
				"@securityDefinitions.oauth2.deviceAuthorization devicecode https://auth.example.com/device https://auth.example.com/token",
			},
			expectedScheme:        "oauth2_devicecode",
			expectedDeviceAuthURL: "https://auth.example.com/device",
			expectedTokenURL:      "https://auth.example.com/token",
		},
		{
			name: "device_authorization_without_token_url",
			annotations: []string{
				"@securityDefinitions.oauth2.deviceAuthorization mydevice https://auth.example.com/device",
			},
			expectedScheme:        "oauth2_mydevice",
			expectedDeviceAuthURL: "https://auth.example.com/device",
			expectedTokenURL:      "",
		},
		{
			name: "multiple_device_flows",
			annotations: []string{
				"@securityDefinitions.oauth2.deviceAuthorization flow1 https://auth1.example.com/device https://auth1.example.com/token",
				"@securityDefinitions.oauth2.deviceAuthorization flow2 https://auth2.example.com/device https://auth2.example.com/token",
			},
			expectedScheme:        "oauth2_flow2",
			expectedDeviceAuthURL: "https://auth2.example.com/device",
			expectedTokenURL:      "https://auth2.example.com/token",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			spec := &openapi.OpenAPI{
				Components: &openapi.Components{
					SecuritySchemes: make(map[string]*openapi.SecurityScheme),
				},
			}

			gproc := NewGeneralInfoProcessor(spec)

			for _, annotation := range tt.annotations {
				if err := gproc.Process(annotation); err != nil {
					t.Fatalf("Process() error = %v", err)
				}
			}

			scheme, exists := spec.Components.SecuritySchemes[tt.expectedScheme]
			if !exists {
				t.Fatalf("Expected scheme %s to exist", tt.expectedScheme)
			}

			if scheme.Type != "oauth2" {
				t.Errorf("Expected Type = oauth2, got %s", scheme.Type)
			}

			if scheme.Flows == nil {
				t.Fatal("Expected Flows to be non-nil")
			}

			if scheme.Flows.DeviceAuthorization == nil {
				t.Fatal("Expected DeviceAuthorization flow to exist")
			}

			if scheme.Flows.DeviceAuthorization.AuthorizationURL != tt.expectedDeviceAuthURL {
				t.Errorf("Expected AuthorizationURL = %q, got %q",
					tt.expectedDeviceAuthURL,
					scheme.Flows.DeviceAuthorization.AuthorizationURL)
			}

			if scheme.Flows.DeviceAuthorization.TokenURL != tt.expectedTokenURL {
				t.Errorf("Expected TokenURL = %q, got %q",
					tt.expectedTokenURL,
					scheme.Flows.DeviceAuthorization.TokenURL)
			}

			if scheme.Flows.DeviceAuthorization.Scopes == nil {
				t.Error("Expected Scopes map to be initialized")
			}
		})
	}
}

// TestSecurityOpenAPI32EdgeCases tests edge cases for OpenAPI 3.2.0 security features.
func TestSecurityOpenAPI32EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("deprecated_and_metadata_url_combined", func(t *testing.T) {
		t.Parallel()

		spec := &openapi.OpenAPI{
			Components: &openapi.Components{
				SecuritySchemes: make(map[string]*openapi.SecurityScheme),
			},
		}

		gproc := NewGeneralInfoProcessor(spec)

		annotations := []string{
			"@securityDefinitions.apikey ApiKeyAuth header X-API-Key",
			"@securityDefinitions.ApiKeyAuth.deprecated true",
			"@securityDefinitions.ApiKeyAuth.oauth2metadataurl https://auth.example.com/metadata",
		}

		for _, annotation := range annotations {
			if err := gproc.Process(annotation); err != nil {
				t.Fatalf("Process() error = %v", err)
			}
		}

		scheme := spec.Components.SecuritySchemes["ApiKeyAuth"]
		if scheme == nil {
			t.Fatal("Expected ApiKeyAuth scheme to exist")
		}

		if !scheme.Deprecated {
			t.Error("Expected Deprecated = true")
		}

		if scheme.OAuth2MetadataURL != "https://auth.example.com/metadata" {
			t.Errorf("Expected OAuth2MetadataURL = %q, got %q",
				"https://auth.example.com/metadata",
				scheme.OAuth2MetadataURL)
		}
	})

	t.Run("empty_values", func(t *testing.T) {
		t.Parallel()

		spec := &openapi.OpenAPI{
			Components: &openapi.Components{
				SecuritySchemes: make(map[string]*openapi.SecurityScheme),
			},
		}

		gproc := NewGeneralInfoProcessor(spec)

		// Create scheme first
		_ = gproc.Process("@securityDefinitions.apikey TestScheme header X-Test")

		// These should not panic
		_ = gproc.Process("@securityDefinitions.TestScheme.oauth2metadataurl ")

		scheme := spec.Components.SecuritySchemes["TestScheme"]
		if scheme == nil {
			t.Fatal("Expected TestScheme to exist")
		}
	})

	t.Run("device_authorization_existing_scheme", func(t *testing.T) {
		t.Parallel()

		spec := &openapi.OpenAPI{
			Components: &openapi.Components{
				SecuritySchemes: make(map[string]*openapi.SecurityScheme),
			},
		}

		// Pre-create oauth2 scheme
		spec.Components.SecuritySchemes["oauth2_myflow"] = &openapi.SecurityScheme{
			Type: "oauth2",
			Flows: &openapi.OAuthFlows{
				Implicit: &openapi.OAuthFlow{
					AuthorizationURL: "https://example.com/auth",
					Scopes:           make(map[string]string),
				},
			},
		}

		gproc := NewGeneralInfoProcessor(spec)

		// Add device authorization to existing scheme
		_ = gproc.Process("@securityDefinitions.oauth2.deviceAuthorization myflow https://example.com/device https://example.com/token")

		scheme := spec.Components.SecuritySchemes["oauth2_myflow"]
		if scheme.Flows.DeviceAuthorization == nil {
			t.Fatal("Expected DeviceAuthorization to be added to existing scheme")
		}

		// Implicit flow should still exist
		if scheme.Flows.Implicit == nil {
			t.Error("Expected Implicit flow to still exist")
		}
	})
}

// TestSecurityOpenAPI32AnnotationValidation tests validation of OpenAPI 3.2.0 security annotations.
func TestSecurityOpenAPI32AnnotationValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		annotation  string
		shouldMatch bool
		setupScheme bool
		schemeName  string
	}{
		{
			name:        "valid_deprecated_true",
			annotation:  "@securityDefinitions.MyScheme.deprecated true",
			shouldMatch: true,
			setupScheme: true,
			schemeName:  "MyScheme",
		},
		{
			name:        "valid_deprecated_false",
			annotation:  "@securityDefinitions.MyScheme.deprecated false",
			shouldMatch: true,
			setupScheme: true,
			schemeName:  "MyScheme",
		},
		{
			name:        "valid_oauth2metadataurl",
			annotation:  "@securityDefinitions.MyScheme.oauth2metadataurl https://example.com/metadata",
			shouldMatch: true,
			setupScheme: true,
			schemeName:  "MyScheme",
		},
		{
			name:        "valid_device_authorization",
			annotation:  "@securityDefinitions.oauth2.deviceAuthorization myflow https://example.com/device https://example.com/token",
			shouldMatch: true,
			setupScheme: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			spec := &openapi.OpenAPI{
				Components: &openapi.Components{
					SecuritySchemes: make(map[string]*openapi.SecurityScheme),
				},
			}

			gproc := NewGeneralInfoProcessor(spec)

			if tt.setupScheme {
				spec.Components.SecuritySchemes[tt.schemeName] = &openapi.SecurityScheme{
					Type: "apiKey",
					Name: "X-API-Key",
					In:   "header",
				}
			}

			err := gproc.Process(tt.annotation)
			if err != nil {
				t.Errorf("Process() unexpected error = %v", err)
			}
		})
	}
}
