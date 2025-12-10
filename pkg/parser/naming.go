// Package parser - Property naming utilities
package parser

import (
	"strings"
	"unicode"
)

// applyPropertyNaming applies the property naming strategy to a field name.
// Only applies if no JSON tag is specified.
func (s *SchemaProcessor) applyPropertyNaming(fieldName string, jsonTag string) string {
	// If JSON tag is explicitly set, use it
	if jsonTag != "" && jsonTag != "-" {
		parts := strings.Split(jsonTag, ",")
		if parts[0] != "" {
			return parts[0]
		}
	}

	// Apply strategy based on parser configuration
	strategy := s.parser.propertyStrategy

	switch strategy {
	case "snakecase":
		return toSnakeCase(fieldName)
	case "pascalcase":
		return fieldName // PascalCase is the original Go field name
	case "camelcase":
		fallthrough
	default:
		return toCamelCase(fieldName)
	}
}

// toSnakeCase converts a string to snake_case.
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			// Add underscore before uppercase letters (except at start)
			if i > 0 && !unicode.IsUpper(rune(s[i-1])) {
				result = append(result, '_')
			}
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// toCamelCase converts a string to camelCase.
func toCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	// First letter lowercase
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// shouldBeRequired determines if a field should be required.
func (s *SchemaProcessor) shouldBeRequired(tags StructTags, isPointer bool) bool {
	// If explicitly required via tag
	if tags.Required {
		return true
	}

	// If omitempty is set, it's optional
	if tags.OmitEmpty {
		return false
	}

	// If it's a pointer, it's optional by default
	if isPointer {
		return false
	}

	// Apply requiredByDefault configuration
	return s.parser.requiredByDefault
}
