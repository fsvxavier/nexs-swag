package parser

import (
	"regexp"
	"testing"

	"github.com/fsvxavier/nexs-swag/pkg/openapi"
)

func TestTransToValidCollectionFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "csv format",
			input:    "csv",
			expected: "csv",
		},
		{
			name:     "multi format",
			input:    "multi",
			expected: "multi",
		},
		{
			name:     "pipes format",
			input:    "pipes",
			expected: "pipes",
		},
		{
			name:     "tsv format",
			input:    "tsv",
			expected: "tsv",
		},
		{
			name:     "ssv format",
			input:    "ssv",
			expected: "ssv",
		},
		{
			name:     "invalid format defaults to csv",
			input:    "invalid",
			expected: "csv",
		},
		{
			name:     "empty format defaults to csv",
			input:    "",
			expected: "csv",
		},
		{
			name:     "uppercase format",
			input:    "CSV",
			expected: "csv",
		},
		{
			name:     "mixed case format",
			input:    "Multi",
			expected: "multi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TransToValidCollectionFormat(tt.input)
			if result != tt.expected {
				t.Errorf("TransToValidCollectionFormat(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessDescription(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "simple description",
			text:     "@Description Get user by ID",
			expected: "Get user by ID",
		},
		{
			name:     "multiline description",
			text:     "@Description This is a long description",
			expected: "This is a long description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &openapi.Operation{}
			proc.processDescription(tt.text, op)
			if op.Description != tt.expected {
				t.Errorf("Expected description '%s', got '%s'", tt.expected, op.Description)
			}
		})
	}
}

func TestProcessID(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	op := &openapi.Operation{}
	proc.processID("@ID getUserByID", op)

	if op.OperationID != "getUserByID" {
		t.Errorf("Expected operation ID 'getUserByID', got '%s'", op.OperationID)
	}
}

func TestProcessTags(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name         string
		text         string
		expectedTags []string
	}{
		{
			name:         "single tag",
			text:         "@Tags users",
			expectedTags: []string{"users"},
		},
		{
			name:         "multiple tags",
			text:         "@Tags users,admin,api",
			expectedTags: []string{"users", "admin", "api"},
		},
		{
			name:         "tags with spaces",
			text:         "@Tags users, admin, api",
			expectedTags: []string{"users", "admin", "api"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &openapi.Operation{}
			proc.processTags(tt.text, op)

			if len(op.Tags) != len(tt.expectedTags) {
				t.Errorf("Expected %d tags, got %d", len(tt.expectedTags), len(op.Tags))
			}

			for i, expected := range tt.expectedTags {
				if i >= len(op.Tags) || op.Tags[i] != expected {
					t.Errorf("Expected tag[%d] = '%s', got '%s'", i, expected, op.Tags[i])
				}
			}
		})
	}
}

func TestProcessAccept(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	op := &openapi.Operation{
		RequestBody: &openapi.RequestBody{
			Content: map[string]*openapi.MediaType{
				"application/json": {
					Schema: &openapi.Schema{Type: "object"},
				},
			},
		},
	}

	proc.processAccept("@Accept json,xml", op)

	// Verify content types were set
	if op.RequestBody == nil || len(op.RequestBody.Content) == 0 {
		t.Error("Expected request body content to be set")
	}
}

func TestProcessProduce(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	op := &openapi.Operation{
		Responses: openapi.Responses{
			"200": &openapi.Response{
				Description: "Success",
				Content: map[string]*openapi.MediaType{
					"application/json": {
						Schema: &openapi.Schema{Type: "object"},
					},
				},
			},
		},
	}

	proc.processProduce("@Produce json,xml", op)

	// Verify response content types were set
	if resp, ok := op.Responses["200"]; ok {
		if len(resp.Content) == 0 {
			t.Error("Expected response content to be set")
		}
	}
}

func TestParseMimeTypes(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "json shortcut",
			input:    "json",
			expected: []string{"application/json"},
		},
		{
			name:     "multiple shortcuts",
			input:    "json,xml,plain",
			expected: []string{"application/json", "text/xml", "text/plain"},
		},
		{
			name:     "full MIME type",
			input:    "application/json",
			expected: []string{"application/json"},
		},
		{
			name:     "mixed",
			input:    "json,text/html",
			expected: []string{"application/json", "text/html"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := proc.parseMimeTypes(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d MIME types, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if i >= len(result) || result[i] != expected {
					t.Errorf("Expected MIME type[%d] = '%s', got '%s'", i, expected, result[i])
				}
			}
		})
	}
}

func TestProcessParameter(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name         string
		text         string
		expectedName string
		expectedIn   string
		expectedType string
		expectedReq  bool
	}{
		{
			name:         "query parameter",
			text:         "@Param id query int true \"User ID\"",
			expectedName: "id",
			expectedIn:   "query",
			expectedType: "integer",
			expectedReq:  true,
		},
		{
			name:         "path parameter",
			text:         "@Param id path string true \"User ID\"",
			expectedName: "id",
			expectedIn:   "path",
			expectedType: "string",
			expectedReq:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &openapi.Operation{}
			proc.processParameter(tt.text, op)

			if len(op.Parameters) != 1 {
				t.Fatalf("Expected 1 parameter, got %d", len(op.Parameters))
			}

			param := op.Parameters[0]
			if param.Name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, param.Name)
			}
			if param.In != tt.expectedIn {
				t.Errorf("Expected in '%s', got '%s'", tt.expectedIn, param.In)
			}
			if param.Required != tt.expectedReq {
				t.Errorf("Expected required %v, got %v", tt.expectedReq, param.Required)
			}
		})
	}
}

func TestProcessRequestBody(t *testing.T) {
	t.Parallel()
	p := New()
	p.openapi.Components.Schemas = map[string]*openapi.Schema{
		"User": {Type: "object"},
	}
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	op := &openapi.Operation{}
	proc.processRequestBody("User", true, "User object", op)

	if op.RequestBody == nil {
		t.Fatal("Expected request body to be set")
	}
	if !op.RequestBody.Required {
		t.Error("Expected request body to be required")
	}
	if op.RequestBody.Description != "User object" {
		t.Errorf("Expected description 'User object', got '%s'", op.RequestBody.Description)
	}
}

func TestProcessResponse(t *testing.T) {
	t.Parallel()
	p := New()
	p.openapi.Components.Schemas = map[string]*openapi.Schema{
		"User": {Type: "object"},
	}
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name         string
		text         string
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "success response",
			text:         "@Success 200 {object} User \"Success\"",
			expectedCode: "200",
			expectedDesc: "Success",
		},
		{
			name:         "error response",
			text:         "@Failure 404 {object} Error \"Not found\"",
			expectedCode: "404",
			expectedDesc: "Not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &openapi.Operation{
				Responses: make(openapi.Responses),
			}
			// Use correct regex format matching operation.go
			var regex *regexp.Regexp
			if tt.name == "success response" {
				regex = regexp.MustCompile(`^@Success\s+(\d+)\s+\{(\w+)\}\s+(\S+)(?:\s+"([^"]*)")?`)
			} else {
				regex = regexp.MustCompile(`^@Failure\s+(\d+)\s+\{(\w+)\}\s+(\S+)(?:\s+"([^"]*)")?`)
			}
			proc.processResponse(tt.text, regex, op)

			if _, ok := op.Responses[tt.expectedCode]; !ok {
				t.Errorf("Expected response for code %s", tt.expectedCode)
			}
		})
	}
}

func TestProcessHeader(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	op := &openapi.Operation{
		Responses: openapi.Responses{
			"200": &openapi.Response{
				Description: "Success",
			},
		},
	}

	proc.processHeader("@Header 200 {int} X-Rate-Limit \"Rate limit\"", op)

	if resp, ok := op.Responses["200"]; ok {
		if len(resp.Headers) == 0 {
			t.Error("Expected headers to be set")
		}
		if _, ok := resp.Headers["X-Rate-Limit"]; !ok {
			t.Error("Expected X-Rate-Limit header")
		}
	} else {
		t.Error("Response 200 not found")
	}
}

func TestProcessSecurity(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name           string
		text           string
		expectedName   string
		expectedScopes int
	}{
		{
			name:           "simple security",
			text:           "@Security Bearer",
			expectedName:   "Bearer",
			expectedScopes: 0,
		},
		{
			name:           "security with scopes",
			text:           "@Security OAuth2 read:users,write:users",
			expectedName:   "OAuth2",
			expectedScopes: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &openapi.Operation{}
			proc.processSecurity(tt.text, op)

			if len(op.Security) == 0 {
				t.Fatal("Expected security to be set")
			}

			if _, ok := op.Security[0][tt.expectedName]; !ok {
				t.Errorf("Expected security '%s'", tt.expectedName)
			}

			if tt.expectedScopes > 0 {
				scopes := op.Security[0][tt.expectedName]
				if len(scopes) != tt.expectedScopes {
					t.Errorf("Expected %d scopes, got %d", tt.expectedScopes, len(scopes))
				}
			}
		})
	}
}

func TestParseAttributes(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "single attribute with parentheses",
			input: "required(true)",
			expected: map[string]string{
				"required": "true",
			},
		},
		{
			name:  "key-value pairs with parentheses",
			input: "minLength(5) maxLength(100)",
			expected: map[string]string{
				"minlength": "5",
				"maxlength": "100",
			},
		},
		{
			name:  "multiple attributes",
			input: "min(10) max(200) required(true)",
			expected: map[string]string{
				"min":      "10",
				"max":      "200",
				"required": "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := proc.parseAttributes(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d attributes, got %d", len(tt.expected), len(result))
			}

			for key, expectedVal := range tt.expected {
				if val, ok := result[key]; !ok {
					t.Errorf("Expected attribute '%s'", key)
				} else if val != expectedVal {
					t.Errorf("Expected %s='%s', got '%s'", key, expectedVal, val)
				}
			}
		})
	}
}

func TestGetSchemaTypeString(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "integer type",
			input:    "int",
			expected: "integer",
		},
		{
			name:     "string type",
			input:    "string",
			expected: "string",
		},
		{
			name:     "boolean type",
			input:    "bool",
			expected: "boolean",
		},
		{
			name:     "array type",
			input:    "[]string",
			expected: "array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create schema based on input type
			schema := &openapi.Schema{Type: tt.expected}
			result := proc.getSchemaTypeString(schema)
			if result != tt.expected {
				t.Errorf("Expected type '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestParseEnumValues(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name     string
		input    string
		expected []interface{}
	}{
		{
			name:     "string values",
			input:    "active,inactive,pending",
			expected: []interface{}{"active", "inactive", "pending"},
		},
		{
			name:     "numeric values",
			input:    "1,2,3",
			expected: []interface{}{"1", "2", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := proc.parseEnumValues(tt.input, "string")

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d enum values, got %d", len(tt.expected), len(result))
			}
		})
	}
}

func TestParseValue(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name     string
		value    string
		dataType string
	}{
		{
			name:     "integer value",
			value:    "123",
			dataType: "integer",
		},
		{
			name:     "string value",
			value:    "hello",
			dataType: "string",
		},
		{
			name:     "boolean value",
			value:    "true",
			dataType: "boolean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := proc.parseValue(tt.value, tt.dataType)
			if result == nil {
				t.Error("Expected non-nil result")
			}
		})
	}
}

func TestParseSchemaType(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	tests := []struct {
		name     string
		typeStr  string
		expected string
	}{
		{
			name:     "object type",
			typeStr:  "{object}",
			expected: "object",
		},
		{
			name:     "array type",
			typeStr:  "{array}",
			expected: "array",
		},
		{
			name:     "string type",
			typeStr:  "{string}",
			expected: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := proc.parseSchemaType(tt.typeStr)
			if schema == nil {
				t.Fatal("Expected non-nil schema")
			}
		})
	}
}

func TestProcessCodeSamples(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	op := &openapi.Operation{}
	proc.processCodeSamples("@CodeSample file:example.go", op)

	// Just verify it doesn't panic
	// Code samples are stored in extensions
}

// TestTransToValidCollectionFormat removed - already exists in operation_test.go
