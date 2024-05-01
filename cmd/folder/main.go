package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	grpcFolder "mail/internal/microservice/folder/server"
	"mail/internal/models/configs"
	"net"
	"net/http"
	"os"
	"time"

	"mail/internal/microservice/folder/proto"
	"mail/internal/microservice/interceptors"

	folderRepo "mail/internal/microservice/folder/repository"
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

func settingTime() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error loc time")
	}

	time.Local = loc
}

func initializeDatabase() *sql.DB {
	db, err := sql.Open("pgx", configs.DSN)
	if err != nil {
		log.Fatalln("Can't parse config", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxOpenConns(10)

	return db
}

func initializationInterceptorLogger() *interceptors.Logger {
	f, err := os.OpenFile("logInterFolder.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}

	LogrusAcces := interceptors.InitializationAccessLogInterceptor(f)
	LoggerAcces := new(interceptors.Logger)
	LoggerAcces.Logger = LogrusAcces

	return LoggerAcces
}

func initializeFolder(db *sql.DB) *grpcFolder.FolderServer {
	folderRepository := folderRepo.NewFolderRepository(sqlx.NewDb(db, "pgx"))
	folderUseCase := folderUc.NewFolderUseCase(folderRepository)

	return grpcFolder.NewFolderServer(folderUseCase)
}

func startServer(folderGrpc *grpcFolder.FolderServer, interceptorsLogger *interceptors.Logger) {
	listen, err := net.Listen("tcp", ":8005")
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8005", err.Error())
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptorsLogger.AccessLogInterceptor,
			interceptors.PanicRecoveryInterceptor,
		),
	}
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterFolderServiceServer(grpcServer, folderGrpc)

	fmt.Printf("The server is running  in port 8005\n")

	http.Handle("/metrics", promhttp.Handler())

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8005", err.Error())
	}
}
