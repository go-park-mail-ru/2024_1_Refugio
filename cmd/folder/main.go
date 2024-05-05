package main

import (
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"mail/cmd/configs"
	"mail/internal/microservice/folder/proto"
	"mail/internal/microservice/interceptors"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	folderRepo "mail/internal/microservice/folder/repository"
	grpcFolder "mail/internal/microservice/folder/server"
	folderUc "mail/internal/microservice/folder/usecase"
)

func main() {
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	folderGrpc := initializeFolder(db)

	loggerInterceptorAccess := initializationInterceptorLogger()

	startServer(folderGrpc, loggerInterceptorAccess)
}

// settingTime setting local time on server
func settingTime() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error in location detection")
	}

	time.Local = loc
}

// initializeDatabase database initialization
func initializeDatabase() *sql.DB {
	db, err := sql.Open("pgx", configs.DSN)
	if err != nil {
		log.Fatalln("Can't parse config", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("Database is not available", err)
	}

	db.SetMaxOpenConns(10)

	return db
}

// initializeFolder initializing folder server
func initializeFolder(db *sql.DB) *grpcFolder.FolderServer {
	folderRepository := folderRepo.NewFolderRepository(sqlx.NewDb(db, "pgx"))
	folderUseCase := folderUc.NewFolderUseCase(folderRepository)

	return grpcFolder.NewFolderServer(folderUseCase)
}

// initializationInterceptorLogger initializing logger
func initializationInterceptorLogger() *interceptors.Logger {
	f, err := os.OpenFile("logInterFolder.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}

	logrusAccess := interceptors.InitializationAccessLogInterceptor(f)
	loggerAccess := new(interceptors.Logger)
	loggerAccess.Logger = logrusAccess

	return loggerAccess
}

// startServer starting server
func startServer(folderGrpc *grpcFolder.FolderServer, interceptorsLogger *interceptors.Logger) {
	listen, err := net.Listen("tcp", ":8005")
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8005", err.Error())
	}

	grpc_prometheus.EnableHandlingTimeHistogram()

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptorsLogger.AccessLogInterceptor,
			interceptors.PanicRecoveryInterceptor,
			grpc_prometheus.UnaryServerInterceptor,
		),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	}
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterFolderServiceServer(grpcServer, folderGrpc)

	fmt.Printf("The server is running  in port 8005\n")

	grpc_prometheus.Register(grpcServer)
	http.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Addr:    ":9095",
		Handler: nil,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Printf("Failed to start Prometheus metrics server: %s\n", err)
		}
	}()

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8005", err.Error())
	}
}
