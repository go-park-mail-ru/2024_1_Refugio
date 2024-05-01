package main

import (
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"mail/cmd/configs"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"mail/internal/microservice/email/proto"
	"mail/internal/microservice/interceptors"

	emailRepo "mail/internal/microservice/email/repository"
	grpcEmail "mail/internal/microservice/email/server"
	emailUc "mail/internal/microservice/email/usecase"
)

func main() {
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	emailGrpc := initializeEmail(db)

	loggerInterceptorAccess := initializationInterceptorLogger()

	startServer(emailGrpc, loggerInterceptorAccess)
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

// initializeEmail initializing email server
func initializeEmail(db *sql.DB) *grpcEmail.EmailServer {
	emailRepository := emailRepo.NewEmailRepository(sqlx.NewDb(db, "pgx"))
	emailUseCase := emailUc.NewEmailUseCase(emailRepository)

	return grpcEmail.NewEmailServer(emailUseCase)
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
func startServer(emailGrpc *grpcEmail.EmailServer, interceptorsLogger *interceptors.Logger) {
	listen, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8002", err.Error())
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptorsLogger.AccessLogInterceptor,
			interceptors.PanicRecoveryInterceptor,
		),
	}
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterEmailServiceServer(grpcServer, emailGrpc)

	fmt.Printf("The server is running  in port 8002\n")

	http.Handle("/metrics", promhttp.Handler())

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8002", err.Error())
	}
}
