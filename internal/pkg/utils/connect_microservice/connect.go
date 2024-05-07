package connect_microservice

import (
	"fmt"
	"google.golang.org/grpc"
	"mail/cmd/configs"
)

// OpenGRPCConnection opens a connection to the rpc server on the specified port.
func OpenGRPCConnection(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(configs.IP_ADDRESS+":"+port, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	return conn, nil
}
