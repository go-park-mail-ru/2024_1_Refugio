package main

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"

	"mail/internal/microservice/auth/proto"
	"mail/internal/microservice/interceptors"

	grpcAuth "mail/internal/microservice/auth/server"
)

func main() {
	settingTime()

	authGrpc := initializeAuth()

	loggerInterceptorAccess := initializationInterceptorLogger()

	startServer(authGrpc, loggerInterceptorAccess)
}

// settingTime setting local time on server
func settingTime() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error in location detection")
	}

	time.Local = loc
}

// initializeAuth initializing authorization server
func initializeAuth() *grpcAuth.AuthServer {
	return grpcAuth.NewAuthServer()
}

// initializationInterceptorLogger initializing logger
func initializationInterceptorLogger() *interceptors.Logger {
	f, err := os.OpenFile("logInterEmail.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}

	logrusAccess := interceptors.InitializationAccessLogInterceptor(f)
	loggerAccess := new(interceptors.Logger)
	loggerAccess.Logger = logrusAccess

	return loggerAccess
}

// startServer starting server
func startServer(authGrpc *grpcAuth.AuthServer, interceptorsLogger *interceptors.Logger) {
	listen, err := net.Listen("tcp", ":8004")
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8004", err.Error())
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptorsLogger.AccessLogInterceptor,
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
