package connect_microservice

import (
	"testing"
)

func TestOpenGRPCConnection_Success(t *testing.T) {
	port := "50051"
	expectedAddress := "0.0.0.0:50051"

	conn, err := OpenGRPCConnection(port)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if conn == nil {
		t.Error("Connection is nil")
	}

	if conn.Target() != expectedAddress {
		t.Errorf("Unexpected connection target address: got %s, want %s", conn.Target(), expectedAddress)
	}

	conn.Close()
}
