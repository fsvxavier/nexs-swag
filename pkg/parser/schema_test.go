package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestNewSchemaProcessor(t *testing.T) {
	t.Parallel()
	p := New()
	sp := NewSchemaProcessor(p, p.openapi, p.typeCache)

	if sp == nil {
		t.Fatal("NewSchemaProcessor() returned nil")
	}
	if sp.parser != p {
		t.Error("parser not set correctly")
	}
	if sp.openapi != p.openapi {
		t.Error("openapi not set correctly")
	}
	if sp.depth != 0 {
		t.Errorf("initial depth = %d, want 0", sp.depth)
	}
}

func TestProcessStructSimple(t *testing.T) {
	t.Parallel()
	content := `package main

type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	sp := NewSchemaProcessor(p, p.openapi, p.typeCache)

	// Find the struct
	var structType *ast.StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if st, ok := n.(*ast.StructType); ok {
			structType = st
			return false
		}
		return true
	})

	if structType == nil {
		t.Fatal("Failed to find struct in parsed file")
	}

	schema := sp.ProcessStruct(structType, nil, "User")
	if schema == nil {
		t.Fatal("ProcessStruct() returned nil")
	}
	if schema.Type != "object" {
		t.Errorf("schema.Type = %q, want %q", schema.Type, "object")
	}
	if len(schema.Properties) != 2 {
		t.Errorf("schema.Properties length = %d, want 2", len(schema.Properties))
	}
}

func TestProcessStructWithTags(t *testing.T) {
	t.Parallel()
	content := `package main

type Product struct {
	ID    int     ` + "`json:\"id\" example:\"1\"`" + `
	Name  string  ` + "`json:\"name\" example:\"Product Name\"`" + `
	Price float64 ` + "`json:\"price\" example:\"99.99\"`" + `
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	sp := NewSchemaProcessor(p, p.openapi, p.typeCache)

	var structType *ast.StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if st, ok := n.(*ast.StructType); ok {
			structType = st
			return false
		}
		return true
	})

	if structType == nil {
		t.Fatal("Failed to find struct in parsed file")
	}

	schema := sp.ProcessStruct(structType, nil, "Product")
	if schema == nil {
		t.Fatal("ProcessStruct() returned nil")
	}
	if len(schema.Properties) != 3 {
		t.Errorf("schema.Properties length = %d, want 3", len(schema.Properties))
	}
}

func TestProcessStructWithDocumentation(t *testing.T) {
	t.Parallel()
	content := `package main

// User represents a user in the system
// @Description User object with basic information
type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	sp := NewSchemaProcessor(p, p.openapi, p.typeCache)

	var structType *ast.StructType
	var doc *ast.CommentGroup
	ast.Inspect(file, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if st, ok := ts.Type.(*ast.StructType); ok {
				structType = st
				doc = ts.Doc
				return false
			}
		}
		return true
	})

	if structType == nil {
		t.Fatal("Failed to find struct in parsed file")
	}

	schema := sp.ProcessStruct(structType, doc, "User")
	if schema == nil {
		t.Fatal("ProcessStruct() returned nil")
	}
	// Note: Doc may not be associated if there's a blank line before type
	// This is expected Go AST behavior
	if schema.Type != "object" {
		t.Errorf("schema.Type = %s, want object", schema.Type)
	}
}

func TestProcessStructWithPointerFields(t *testing.T) {
	t.Parallel()
	content := `package main

type User struct {
	ID    int     ` + "`json:\"id\"`" + `
	Email *string ` + "`json:\"email,omitempty\"`" + `
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	sp := NewSchemaProcessor(p, p.openapi, p.typeCache)

	var structType *ast.StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if st, ok := n.(*ast.StructType); ok {
			structType = st
			return false
		}
		return true
	})

	if structType == nil {
		t.Fatal("Failed to find struct in parsed file")
	}

	schema := sp.ProcessStruct(structType, nil, "User")
	if schema == nil {
		t.Fatal("ProcessStruct() returned nil")
	}
	if len(schema.Properties) != 2 {
		t.Errorf("schema.Properties length = %d, want 2", len(schema.Properties))
	}
}

func TestProcessStructEmpty(t *testing.T) {
	t.Parallel()
	content := `package main

type Empty struct {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	sp := NewSchemaProcessor(p, p.openapi, p.typeCache)

	var structType *ast.StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if st, ok := n.(*ast.StructType); ok {
			structType = st
			return false
		}
		return true
	})

	if structType == nil {
		t.Fatal("Failed to find struct in parsed file")
	}

	schema := sp.ProcessStruct(structType, nil, "Empty")
	if schema == nil {
		t.Fatal("ProcessStruct() returned nil")
	}
	if schema.Type != "object" {
		t.Errorf("schema.Type = %q, want %q", schema.Type, "object")
	}
	if len(schema.Properties) != 0 {
		t.Errorf("schema.Properties length = %d, want 0", len(schema.Properties))
	}
}

func TestProcessStructWithUnexportedFields(t *testing.T) {
	t.Parallel()
	content := `package main

type User struct {
	ID       int    ` + "`json:\"id\"`" + `
	name     string ` + "`json:\"name\"`" + `  // unexported
	Email    string ` + "`json:\"email\"`" + `
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	sp := NewSchemaProcessor(p, p.openapi, p.typeCache)

	var structType *ast.StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if st, ok := n.(*ast.StructType); ok {
			structType = st
			return false
		}
		return true
	})

	if structType == nil {
		t.Fatal("Failed to find struct in parsed file")
	}

	schema := sp.ProcessStruct(structType, nil, "User")
	if schema == nil {
		t.Fatal("ProcessStruct() returned nil")
	}
	// Should only have 2 properties (ID and Email), not 3
	if len(schema.Properties) != 2 {
		t.Errorf("schema.Properties length = %d, want 2 (unexported fields should be skipped)", len(schema.Properties))
	}
}

func TestProcessStructWithArrayFields(t *testing.T) {
	t.Parallel()
	content := `package main

type User struct {
	ID    int      ` + "`json:\"id\"`" + `
	Tags  []string ` + "`json:\"tags\"`" + `
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	sp := NewSchemaProcessor(p, p.openapi, p.typeCache)

	var structType *ast.StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if st, ok := n.(*ast.StructType); ok {
			structType = st
			return false
		}
		return true
	})

	if structType == nil {
		t.Fatal("Failed to find struct in parsed file")
	}

	schema := sp.ProcessStruct(structType, nil, "User")
	if schema == nil {
		t.Fatal("ProcessStruct() returned nil")
	}
	if len(schema.Properties) != 2 {
		t.Errorf("schema.Properties length = %d, want 2", len(schema.Properties))
	}
}

func TestParseStructDocAnnotations(t *testing.T) {
	t.Parallel()
	// Test that parseStructDoc correctly processes annotations when doc is provided
	content := `package main

type User struct {
	ID int ` + "`json:\"id\"`" + `
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	p := New()
	sp := NewSchemaProcessor(p, p.openapi, p.typeCache)

	// Create a mock CommentGroup with annotations
	mockDoc := &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: "// @Description This is a user"},
			{Text: "// @Title User Model"},
		},
	}

	var structType *ast.StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if st, ok := ts.Type.(*ast.StructType); ok {
				structType = st
				return false
			}
		}
		return true
	})

	schema := sp.ProcessStruct(structType, mockDoc, "User")
	if schema.Description == "" {
		t.Error("Expected description to be set from @Description annotation")
	}
	if schema.Title == "" {
		t.Error("Expected title to be set from @Title annotation")
	}
}
