package main

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"

	"mail/internal/microservice/auth/proto"
	"mail/internal/microservice/interceptors"

	grpcAuth "mail/internal/microservice/auth/server"
)

func main() {
	settingTime()

	authGrpc := initializeSession()

	startServer(authGrpc)
}

func settingTime() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error loc time")
	}

	time.Local = loc
}

func initializeSession() *grpcAuth.AuthServer {
	return grpcAuth.NewAuthServer()
}

func startServer(authGrpc *grpcAuth.AuthServer) {
	listen, err := net.Listen("tcp", ":8004")
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8004", err.Error())
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.PanicRecoveryWithoutLoggerInterceptor,
		),
	}
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterAuthServiceServer(grpcServer, authGrpc)

	fmt.Printf("The server is running in port 8004\n")

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8004", err.Error())
	}
}
