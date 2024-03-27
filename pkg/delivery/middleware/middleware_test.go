package middleware

import (
	"testing"
)

func TestInitializationAcceslog(t *testing.T) {
	port := 8080

	logger := InitializationAcceslog(port)

	if logger == nil {
		t.Error("Expected non-nil logger, got nil")
	}

	expectedFields := map[string]interface{}{
		"logger": "Logrus",
		"host":   "localhost",
		"port":   port,
	}

	for key, value := range expectedFields {
		if fieldValue := logger.LogrusLogger.Data[key]; fieldValue != value {
			t.Errorf("Expected field '%s' to be '%v', got '%v'", key, value, fieldValue)
		}
	}
}
