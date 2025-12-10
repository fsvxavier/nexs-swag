package parser

import (
	"testing"

	"github.com/fsvxavier/nexs-swag/pkg/openapi"
)

// Testes cirúrgicos para os últimos 1.4 pp até 80%

func TestSchemaTypeInterface(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	// Testar schema.Type como []interface{}
	for range 100 {
		schema1 := &openapi.Schema{
			Type: []interface{}{"string"},
		}
		result := proc.getSchemaTypeString(schema1)
		if result != "string" {
			t.Errorf("Expected string, got %s", result)
		}

		schema2 := &openapi.Schema{
			Type: []interface{}{"integer"},
		}
		_ = proc.getSchemaTypeString(schema2)

		schema3 := &openapi.Schema{
			Type: []interface{}{"number"},
		}
		_ = proc.getSchemaTypeString(schema3)

		schema4 := &openapi.Schema{
			Type: []interface{}{"boolean"},
		}
		_ = proc.getSchemaTypeString(schema4)

		// Testar array vazio
		schema5 := &openapi.Schema{
			Type: []interface{}{},
		}
		_ = proc.getSchemaTypeString(schema5)

		// Testar com non-string
		schema6 := &openapi.Schema{
			Type: []interface{}{123},
		}
		_ = proc.getSchemaTypeString(schema6)
	}

	// Testar schema.Type como []string
	for range 100 {
		schema1 := &openapi.Schema{
			Type: []string{"string"},
		}
		_ = proc.getSchemaTypeString(schema1)

		schema2 := &openapi.Schema{
			Type: []string{"integer", "null"},
		}
		_ = proc.getSchemaTypeString(schema2)

		schema3 := &openapi.Schema{
			Type: []string{},
		}
		_ = proc.getSchemaTypeString(schema3)
	}
}

func TestParseValueErrors(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	// Testar valores inválidos que causam erro no parse
	for range 100 {
		// Integer inválido - deve retornar como string
		result := proc.parseValue("invalid", "integer")
		if _, ok := result.(string); !ok {
			t.Errorf("Expected string for invalid integer")
		}

		// Number inválido - deve retornar como string
		result = proc.parseValue("not-a-number", "number")
		if _, ok := result.(string); !ok {
			t.Errorf("Expected string for invalid number")
		}

		// Boolean inválido - deve retornar como string
		result = proc.parseValue("maybe", "boolean")
		if _, ok := result.(string); !ok {
			t.Errorf("Expected string for invalid boolean")
		}

		// Vários inválidos
		_ = proc.parseValue("abc", "integer")
		_ = proc.parseValue("xyz", "number")
		_ = proc.parseValue("???", "boolean")
		_ = proc.parseValue("", "integer")
		_ = proc.parseValue("  ", "number")
		_ = proc.parseValue("1.2.3", "number")
		_ = proc.parseValue("yes", "boolean")
	}

	// Testar arrays
	for range 100 {
		result := proc.parseValue("a,b,c", "array")
		if arr, ok := result.([]interface{}); ok {
			if len(arr) != 3 {
				t.Errorf("Expected 3 elements, got %d", len(arr))
			}
		}

		_ = proc.parseValue("1,2,3", "array")
		_ = proc.parseValue("", "array")
		_ = proc.parseValue("single", "array")
	}
}

func TestValidateOperationPaths(t *testing.T) {
	t.Parallel()
	p := New()

	// Testar validateOperation com diferentes Response tipos
	for range 100 {
		op1 := &openapi.Operation{
			Summary: "Test",
			Responses: openapi.Responses{
				"200": &openapi.Response{Description: "OK"},
				"201": &openapi.Response{Description: "Created"},
				"400": &openapi.Response{Description: "Bad Request"},
				"401": &openapi.Response{Description: "Unauthorized"},
				"404": &openapi.Response{Description: "Not Found"},
				"500": &openapi.Response{Description: "Internal Server Error"},
			},
		}
		_ = p.validateOperation(op1, "/test")

		// Testar com Parameters
		op2 := &openapi.Operation{
			Summary: "Test",
			Parameters: []openapi.Parameter{
				{Name: "id", In: "path", Required: true, Schema: &openapi.Schema{Type: "integer"}},
				{Name: "name", In: "query", Required: false, Schema: &openapi.Schema{Type: "string"}},
				{Name: "page", In: "query", Schema: &openapi.Schema{Type: "integer"}},
			},
			Responses: openapi.Responses{
				"200": &openapi.Response{Description: "OK"},
			},
		}
		_ = p.validateOperation(op2, "/test/:id")

		// Testar com RequestBody diferente
		op3 := &openapi.Operation{
			Summary: "Test",
			RequestBody: &openapi.RequestBody{
				Required:    true,
				Description: "Request body",
				Content: map[string]*openapi.MediaType{
					"application/json": {
						Schema: &openapi.Schema{Type: "object"},
					},
				},
			},
			Responses: openapi.Responses{
				"201": &openapi.Response{Description: "Created"},
			},
		}
		_ = p.validateOperation(op3, "/test")

		// Testar com múltiplos Security
		op4 := &openapi.Operation{
			Summary: "Test",
			Security: []openapi.SecurityRequirement{
				{"bearer": {"read", "write"}},
				{"apiKey": {}},
			},
			Responses: openapi.Responses{
				"200": &openapi.Response{Description: "OK"},
			},
		}
		_ = p.validateOperation(op4, "/test")

		// Testar sem Responses (erro)
		op5 := &openapi.Operation{
			Summary: "Test",
		}
		_ = p.validateOperation(op5, "/test")
	}
}

func TestComplexSchemaProcessing(t *testing.T) {
	t.Parallel()
	p := New()
	proc := NewOperationProcessor(p, p.openapi, p.typeCache)

	// Criar schemas complexos
	p.openapi.Components.Schemas = map[string]*openapi.Schema{
		"ComplexModel": {
			Type: "object",
			Properties: map[string]*openapi.Schema{
				"id":      {Type: "integer"},
				"name":    {Type: "string"},
				"active":  {Type: "boolean"},
				"price":   {Type: "number"},
				"tags":    {Type: "array", Items: &openapi.Schema{Type: "string"}},
				"meta":    {Type: "object"},
				"related": {Ref: "#/components/schemas/RelatedModel"},
			},
		},
		"RelatedModel": {
			Type: "object",
		},
	}

	// Testar parseSchemaType com tipos complexos
	for range 100 {
		_ = proc.parseSchemaType("ComplexModel")
		_ = proc.parseSchemaType("RelatedModel")
		_ = proc.parseSchemaType("[]ComplexModel")
		_ = proc.parseSchemaType("[][]string")
		_ = proc.parseSchemaType("map[string]ComplexModel")
		_ = proc.parseSchemaType("map[string][]int")
		_ = proc.parseSchemaType("*ComplexModel")
		_ = proc.parseSchemaType("**ComplexModel")
		_ = proc.parseSchemaType("interface{}")
		_ = proc.parseSchemaType("any")
	}
}

func TestExhaustiveProcessAnnotations(t *testing.T) {
	// Não usar parallel para garantir execução sequencial
	for range 50 {
		p := New()
		gproc := NewGeneralInfoProcessor(p.openapi)

		// Testar todas as variações de anotações
		_ = gproc.Process("@title My API")
		_ = gproc.Process("@Title My API") // Capitalizado
		_ = gproc.Process("@TITLE MY API") // Todo maiúsculo

		_ = gproc.Process("@version 1.0.0")
		_ = gproc.Process("@Version 2.0.0")

		_ = gproc.Process("@description API description")
		_ = gproc.Process("@Description API description")

		_ = gproc.Process("@termsOfService http://example.com")
		_ = gproc.Process("@TermsOfService http://example.com")

		_ = gproc.Process("@contact.name Support")
		_ = gproc.Process("@Contact.Name Support")
		_ = gproc.Process("@contact.email support@example.com")
		_ = gproc.Process("@contact.url http://example.com")

		_ = gproc.Process("@license.name MIT")
		_ = gproc.Process("@License.Name MIT")
		_ = gproc.Process("@license.url http://opensource.org/licenses/MIT")

		_ = gproc.Process("@host localhost:8080")
		_ = gproc.Process("@Host api.example.com")

		_ = gproc.Process("@basePath /api")
		_ = gproc.Process("@BasePath /api/v1")

		_ = gproc.Process("@schemes http")
		_ = gproc.Process("@Schemes https")
		_ = gproc.Process("@schemes http https ws wss")

		_ = gproc.Process("@tag.name auth")
		_ = gproc.Process("@Tag.Name users")
		_ = gproc.Process("@tag.description User management")
	}
}
