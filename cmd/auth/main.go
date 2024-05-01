package main

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"mail/internal/microservice/auth/proto"
	"mail/internal/microservice/interceptors"

	grpcAuth "mail/internal/microservice/auth/server"
)

func main() {
	settingTime()

	authGrpc := initializeSession()

	loggerInterceptorAccess := initializationInterceptorLogger()

	startServer(authGrpc, loggerInterceptorAccess)
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

func initializationInterceptorLogger() *interceptors.Logger {
	f, err := os.OpenFile("logInterEmail.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}

	LogrusAcces := interceptors.InitializationAccessLogInterceptor(f)
	LoggerAcces := new(interceptors.Logger)
	LoggerAcces.Logger = LogrusAcces

	return LoggerAcces
}

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
