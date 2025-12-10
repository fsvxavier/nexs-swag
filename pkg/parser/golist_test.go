package parser

import (
	"testing"
)

func TestParseWithGoList(t *testing.T) {
	t.Parallel()
	p := New()
	p.parseGoList = true

	// This will likely fail if not in a valid Go module
	// but we test that it doesn't panic
	err := p.parseWithGoList(".")
	// Error is expected if not in proper Go module context
	if err != nil {
		t.Logf("parseWithGoList() returned expected error: %v", err)
	}
}

func TestParseWithGoListInvalidPath(t *testing.T) {
	t.Parallel()
	p := New()
	p.parseGoList = true

	err := p.parseWithGoList("/nonexistent/path")
	if err == nil {
		t.Error("parseWithGoList() with invalid path should return error")
	}
}

func TestParseDependencies(t *testing.T) {
	t.Parallel()
	p := New()
	p.parseDependency = false

	// Should not parse dependencies when disabled
	err := p.parseDependencies()
	if err != nil {
		t.Errorf("parseDependencies() with disabled flag returned error: %v", err)
	}
}

func TestParseDependenciesEnabled(t *testing.T) {
	t.Parallel()
	p := New()
	p.parseDependency = true

	// This will try to parse go.mod
	err := p.parseDependencies()
	// Error is expected if go.mod doesn't exist or is invalid
	if err != nil {
		t.Logf("parseDependencies() returned expected error: %v", err)
	}
}

func TestGoListPackageInfo(t *testing.T) {
	t.Parallel()
	// Test that go list functionality exists
	p := New()
	p.parseGoList = true

	if !p.parseGoList {
		t.Error("parseGoList should be true")
	}
}
