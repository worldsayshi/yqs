package main

import (
	"testing"
)

func TestSuggestContinuationsWithContent(t *testing.T) {
	yamlContent := `
name: "test"
users:
  - name: "alice"
    age: 25
  - name: "bob"
    age: 30
config:
  enabled: true
  timeout: 300
`

	tests := []struct {
		name           string
		baseExpression string
		yamlContent    string
		expectedCount  int // We'll check that we get a reasonable number of suggestions
	}{
		{
			name:           "root expression",
			baseExpression: ".",
			yamlContent:    yamlContent,
			expectedCount:  10, // Should include base operators plus key-based suggestions
		},
		{
			name:           "users array expression",
			baseExpression: ".users",
			yamlContent:    yamlContent,
			expectedCount:  5, // Should include array operators
		},
		{
			name:           "config object expression",
			baseExpression: ".config",
			yamlContent:    yamlContent,
			expectedCount:  10, // Should include object keys (enabled, timeout)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := suggestContinuationsWithContent(tt.baseExpression, tt.yamlContent)

			if len(suggestions) < tt.expectedCount {
				t.Errorf("Expected at least %d suggestions, got %d", tt.expectedCount, len(suggestions))
			}

			// Check that we have some basic continuations
			hasBasicContinuation := false
			for _, suggestion := range suggestions {
				if suggestion == tt.baseExpression+" | keys" ||
					suggestion == tt.baseExpression+" | length" {
					hasBasicContinuation = true
					break
				}
			}

			if !hasBasicContinuation {
				t.Error("Expected to find basic continuations like '| keys' or '| length'")
			}
		})
	}
}

func TestGetKeysFromYQOutputWithContent(t *testing.T) {
	yamlContent := `
name: "test"
users:
  - name: "alice"
    age: 25
  - name: "bob"
    age: 30
config:
  enabled: true
  timeout: 300
`

	tests := []struct {
		name           string
		baseExpression string
		yamlContent    string
		expectedKeys   []string
	}{
		{
			name:           "root keys",
			baseExpression: ".",
			yamlContent:    yamlContent,
			expectedKeys:   []string{"name", "users", "config"},
		},
		{
			name:           "config keys",
			baseExpression: ".config",
			yamlContent:    yamlContent,
			expectedKeys:   []string{"enabled", "timeout"},
		},
		{
			name:           "first user keys",
			baseExpression: ".users[0]",
			yamlContent:    yamlContent,
			expectedKeys:   []string{"name", "age"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := getKeysFromYQOutputWithContent(tt.yamlContent, tt.baseExpression)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(keys) != len(tt.expectedKeys) {
				t.Errorf("Expected %d keys, got %d: %v", len(tt.expectedKeys), len(keys), keys)
				return
			}

			// Check that all expected keys are present (order doesn't matter)
			keyMap := make(map[string]bool)
			for _, key := range keys {
				keyMap[key] = true
			}

			for _, expectedKey := range tt.expectedKeys {
				if !keyMap[expectedKey] {
					t.Errorf("Expected key '%s' not found in results: %v", expectedKey, keys)
				}
			}
		})
	}
}

func TestTestYQExpressionWithContent(t *testing.T) {
	yamlContent := `
name: "test"
users:
  - name: "alice"
    age: 25
  - name: "bob"
    age: 30
`

	tests := []struct {
		name        string
		expression  string
		yamlContent string
		expected    string
		shouldError bool
	}{
		{
			name:        "get name",
			expression:  ".name",
			yamlContent: yamlContent,
			expected:    "test",
			shouldError: false,
		},
		{
			name:        "get first user name",
			expression:  ".users[0].name",
			yamlContent: yamlContent,
			expected:    "alice",
			shouldError: false,
		},
		{
			name:        "count users",
			expression:  ".users | length",
			yamlContent: yamlContent,
			expected:    "2",
			shouldError: false,
		},
		{
			name:        "invalid expression",
			expression:  ".nonexistent.deeply.nested",
			yamlContent: yamlContent,
			expected:    "null",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := testYQExpressionWithContent(tt.yamlContent, tt.expression)

			if tt.shouldError && err == nil {
				t.Error("Expected an error but got none")
				return
			}

			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestJoinYQExpressions(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		subPath  string
		expected string
	}{
		{
			name:     "root with key",
			base:     ".",
			subPath:  "name",
			expected: ".name",
		},
		{
			name:     "root with dotted key",
			base:     ".",
			subPath:  ".name",
			expected: ".name",
		},
		{
			name:     "base with key",
			base:     ".users",
			subPath:  "name",
			expected: ".users.name",
		},
		{
			name:     "base with pipe operation",
			base:     ".users",
			subPath:  "| length",
			expected: ".users | length",
		},
		{
			name:     "base with array selector",
			base:     ".users",
			subPath:  "[0]",
			expected: ".users[0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinYQExpressions(tt.base, tt.subPath)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
