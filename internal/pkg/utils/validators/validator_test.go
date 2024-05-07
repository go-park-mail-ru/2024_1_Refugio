package validators

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"  ", true},
		{"   hello   ", false},
		{"hello", false},
		{"   hello", false},
		{"hello   ", false},
		{"   hello   world   ", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := IsEmpty(test.input)
			if result != test.expected {
				t.Errorf("Expected IsEmpty(%q) to return %t, but got %t", test.input, test.expected, result)
			}
		})
	}
}

func TestIsValidEmailFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"user@mailhub.su", true},
		{"user@mailhub.ru", false},
		{"user@mailhub", false},
		{"user@mailhub.su.com", false},
		{"user@domain.com", false},
		{"", false},
		{"user@", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := IsValidEmailFormat(test.input)
			if result != test.expected {
				t.Errorf("Expected IsValidEmailFormat(%q) to return %t, but got %t", test.input, test.expected, result)
			}
		})
	}
}
