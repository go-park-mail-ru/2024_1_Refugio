package generate_filename

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateUniqueFileName(t *testing.T) {
	expectedFormats := []string{".jpg", ".png", ".txt"}

	for _, format := range expectedFormats {
		fileName := GenerateUniqueFileName(format)

		if !strings.HasSuffix(fileName, format) {
			t.Errorf("Generated file name %q does not have the expected format %q", fileName, format)
		}

		currentTime := time.Now().Format("20060102_150405")
		if !strings.Contains(fileName, currentTime) {
			t.Errorf("Generated file name %q does not contain the current time", fileName)
		}

		if !strings.Contains(fileName, "_") {
			t.Errorf("Generated file name %q does not contain a random number", fileName)
		}
	}
}
