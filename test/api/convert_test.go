package api_test

import (
	"testing"

	"html2go-converter/api"
)

func TestRemoveBodyWrapper(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic case with var n = Body() wrapper",
			input:    "var n = Body(h.Div(h.P(\"Hello world\")))",
			expected: "h.Div(h.P(\"Hello world\"))",
		},
		{
			name:     "With whitespace",
			input:    "  var n = Body(  h.Div(h.P(\"Hello\"))  )  ",
			expected: "h.Div(h.P(\"Hello\"))",
		},
		{
			name:     "No body wrapper",
			input:    "h.Div(h.P(\"No wrapper\"))",
			expected: "h.Div(h.P(\"No wrapper\"))",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only whitespace",
			input:    "   ",
			expected: "",
		},
		{
			name:     "Multiline content",
			input:    "var n = Body(\n  h.Div(\n    h.P(\"Multiline\")\n  )\n)",
			expected: "h.Div(\n    h.P(\"Multiline\")\n  )",
		},
		{
			name:     "With package prefix wrapper",
			input:    "var n = mypackage.Body(h.Div(h.P(\"With package\")))",
			expected: "h.Div(h.P(\"With package\"))",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := api.RemoveBodyWrapper(tc.input)
			if result != tc.expected {
				t.Errorf("Expected: %q, got: %q", tc.expected, result)
			}
		})
	}
}
