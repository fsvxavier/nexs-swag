package parser

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "camelCase",
			input:    "userName",
			expected: "user_name",
		},
		{
			name:     "PascalCase",
			input:    "UserName",
			expected: "user_name",
		},
		{
			name:     "with numbers",
			input:    "user1Name2",
			expected: "user1_name2",
		},
		{
			name:     "already snake_case",
			input:    "user_name",
			expected: "user_name",
		},
		{
			name:     "single word lowercase",
			input:    "user",
			expected: "user",
		},
		{
			name:     "single word uppercase",
			input:    "USER",
			expected: "user",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "ID field",
			input:    "ID",
			expected: "id",
		},
		{
			name:     "HTTPResponse",
			input:    "HTTPResponse",
			expected: "httpresponse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "snake_case - only lowercases first letter",
			input:    "User_name",
			expected: "user_name",
		},
		{
			name:     "single word",
			input:    "user",
			expected: "user",
		},
		{
			name:     "already camelCase",
			input:    "userName",
			expected: "userName",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "PascalCase input",
			input:    "UserName",
			expected: "userName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
