package connect_microservice

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mail/cmd/configs"
)

// OpenGRPCConnection opens a connection to the RPC server on the specified port.
func OpenGRPCConnection(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(configs.IP_ADDRESS+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	return conn, nil
}
