package openapi

import "testing"

func TestNewOpenAPI(t *testing.T) {
	api := &OpenAPI{
		OpenAPI: "3.1.0",
	}
	if api.OpenAPI != "3.1.0" {
		t.Errorf("OpenAPI version = %q, want %q", api.OpenAPI, "3.1.0")
	}
}

func TestInfoFields(t *testing.T) {
	info := &Info{
		Title:       "Test API",
		Description: "Test Description",
		Version:     "1.0.0",
	}
	if info.Title != "Test API" {
		t.Errorf("Info.Title = %q, want %q", info.Title, "Test API")
	}
}

func TestContactFields(t *testing.T) {
	contact := &Contact{
		Name:  "Test Contact",
		Email: "test@example.com",
		URL:   "https://example.com",
	}
	if contact.Name != "Test Contact" {
		t.Errorf("Contact.Name = %q, want %q", contact.Name, "Test Contact")
	}
}

func TestLicenseFields(t *testing.T) {
	license := &License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	}
	if license.Name != "MIT" {
		t.Errorf("License.Name = %q, want %q", license.Name, "MIT")
	}
}

func TestServerFields(t *testing.T) {
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
	components := &Components{
		Schemas: make(map[string]*Schema),
	}
	components.Schemas["User"] = &Schema{Type: "object"}
	if len(components.Schemas) != 1 {
		t.Errorf("Components.Schemas length = %d, want 1", len(components.Schemas))
	}
}
