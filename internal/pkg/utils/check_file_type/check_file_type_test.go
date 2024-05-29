package check_file_type

import (
	"testing"
)

func TestGetFileType(t *testing.T) {
	tests := []struct {
		contentType string
		expected    string
	}{
		{"image/jpeg", "Image"},
		{"video/mp4", "Video"},
		{"audio/mpeg", "Audio"},
		{"application/pdf", "PDF Document"},
		{"text/plain", "Text Document"},
	}

	for _, test := range tests {
		t.Run(test.contentType, func(t *testing.T) {
			result := GetFileType(test.contentType)
			if result != test.expected {
				t.Errorf("Expected type %s for content type %s, but got %s", test.expected, test.contentType, result)
			}
		})
	}
}
