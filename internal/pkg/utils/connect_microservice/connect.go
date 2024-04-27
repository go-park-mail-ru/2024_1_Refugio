package connect_microservice

import (
	"fmt"
	"google.golang.org/grpc"
)

// OpenGRPCConnection opens a connection to the rpc server on the specified port.
func OpenGRPCConnection(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("0.0.0.0:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	return conn, nil
}

// 89.208.223.140
