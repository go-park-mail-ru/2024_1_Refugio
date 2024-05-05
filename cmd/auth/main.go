package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"mail/internal/microservice/auth/proto"
	"mail/internal/microservice/interceptors"
	"mail/internal/models/microservice_ports"
	"mail/internal/pkg/utils/connect_microservice"

	grpcAuth "mail/internal/microservice/auth/server"
	session_proto "mail/internal/microservice/session/proto"
	user_proto "mail/internal/microservice/user/proto"
)

func main() {
	settingTime()

	sessionServiceConn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		log.Fatalf("connection with microservice session fail")
	}
	defer sessionServiceConn.Close()
	sessionServiceClient := session_proto.NewSessionServiceClient(sessionServiceConn)

	userServiceConn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		log.Fatalf("connection with microservice user fail")
	}
	defer userServiceConn.Close()
	userServiceClient := user_proto.NewUserServiceClient(userServiceConn)

	authGrpc := initializeAuth(sessionServiceClient, userServiceClient)

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
func initializeAuth(sessionServiceClient session_proto.SessionServiceClient, userServiceClient user_proto.UserServiceClient) *grpcAuth.AuthServer {
	return grpcAuth.NewAuthServer(sessionServiceClient, userServiceClient)
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

	http.Handle("/metrics", promhttp.Handler())

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8004", err.Error())
	}
}
