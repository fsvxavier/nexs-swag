package parser

import (
	"testing"
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
