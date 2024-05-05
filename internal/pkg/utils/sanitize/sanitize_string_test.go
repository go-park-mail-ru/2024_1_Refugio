package sanitize

import (
	"testing"
)

func TestSanitizeString(t *testing.T) {
	input := "<script>alert('Hello World')</script>"
	expected := ""

	result := SanitizeString(input)

	if result != expected {
		t.Errorf("Expected SanitizeString(%q) to return %q, but got %q", input, expected, result)
	}
}
