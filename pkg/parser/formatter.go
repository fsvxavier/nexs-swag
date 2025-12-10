package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"strings"
	"text/tabwriter"
)

// Formatter implements a formatter for Go source files with swagger comments.
type Formatter struct {
}

// NewFormatter creates a new formatter instance.
func NewFormatter() *Formatter {
	return &Formatter{}
}

// Format formats swag comments in contents. It uses fileName to report errors.
func (f *Formatter) Format(fileName string, contents []byte) ([]byte, error) {
	fset := token.NewFileSet()

	astFile, err := goparser.ParseFile(fset, fileName, contents, goparser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var formatted bytes.Buffer

	if err := f.formatFile(&formatted, fset, astFile); err != nil {
		return nil, err
	}

	return formatted.Bytes(), nil
}

func (f *Formatter) formatFile(output *bytes.Buffer, fset *token.FileSet, file *ast.File) error {
	// Write package declaration
	fmt.Fprintf(output, "package %s\n\n", file.Name.Name)

	// Process imports
	if len(file.Imports) > 0 {
		if err := f.formatImports(output, file.Imports); err != nil {
			return err
		}
	}

	// Process declarations (functions, types, etc.)
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if err := f.formatGenDecl(output, fset, genDecl); err != nil {
				return err
			}
		} else if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if err := f.formatFuncDecl(output, fset, funcDecl); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *Formatter) formatImports(output *bytes.Buffer, imports []*ast.ImportSpec) error {
	output.WriteString("import (\n")
	for _, imp := range imports {
		if imp.Name != nil {
			fmt.Fprintf(output, "\t%s %s\n", imp.Name.Name, imp.Path.Value)
		} else {
			fmt.Fprintf(output, "\t%s\n", imp.Path.Value)
		}
	}
	output.WriteString(")\n\n")
	return nil
}

func (f *Formatter) formatGenDecl(output *bytes.Buffer, fset *token.FileSet, genDecl *ast.GenDecl) error {
	// Format comments if present
	if genDecl.Doc != nil {
		f.formatCommentGroup(output, fset, genDecl.Doc)
	}

	// Write the declaration (type, const, var)
	// For now, we'll skip the actual code and focus on comments
	output.WriteString("\n")

	return nil
}

func (f *Formatter) formatFuncDecl(output *bytes.Buffer, fset *token.FileSet, funcDecl *ast.FuncDecl) error {
	// Format function comments
	if funcDecl.Doc != nil {
		if err := f.formatSwaggerComments(output, fset, funcDecl.Doc); err != nil {
			return err
		}
	}

	output.WriteString("\n")

	return nil
}

func (f *Formatter) formatCommentGroup(output *bytes.Buffer, fset *token.FileSet, cg *ast.CommentGroup) {
	for _, comment := range cg.List {
		output.WriteString(comment.Text)
		output.WriteString("\n")
	}
}

var swaggerAnnotations = map[string]bool{
	"@summary":       true,
	"@description":   true,
	"@tags":          true,
	"@accept":        true,
	"@produce":       true,
	"@param":         true,
	"@success":       true,
	"@failure":       true,
	"@response":      true,
	"@header":        true,
	"@router":        true,
	"@security":      true,
	"@deprecated":    true,
	"@id":            true,
	"@state":         true,
	"@x-codesamples": true,
}

func (f *Formatter) formatSwaggerComments(output *bytes.Buffer, fset *token.FileSet, cg *ast.CommentGroup) error {
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 4, 1, ' ', 0)

	for _, comment := range cg.List {
		text := comment.Text

		// Remove comment markers
		text = strings.TrimPrefix(text, "//")
		text = strings.TrimPrefix(text, "/*")
		text = strings.TrimSuffix(text, "*/")
		text = strings.TrimSpace(text)

		// Check if it's a swagger annotation
		if strings.HasPrefix(text, "@") {
			parts := strings.Fields(text)
			if len(parts) > 0 {
				annotation := strings.ToLower(parts[0])
				if swaggerAnnotations[annotation] {
					// Format swagger annotations with consistent spacing
					fmt.Fprintf(tw, "// %s\t%s\n", parts[0], strings.Join(parts[1:], " "))
					continue
				}
			}
		}

		// Regular comment
		fmt.Fprintf(tw, "// %s\n", text)
	}

	tw.Flush()
	output.Write(buf.Bytes())

	return nil
}
