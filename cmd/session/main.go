package main

import (
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"mail/cmd/configs"
	"net"
	"os"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"mail/internal/microservice/interceptors"
	"mail/internal/microservice/session/proto"
	sessionRepo "mail/internal/microservice/session/repository"
	grpcSession "mail/internal/microservice/session/server"
	sessionUc "mail/internal/microservice/session/usecase"
)

func main() {
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	sessionGrpc := initializeSession(db)

	loggerInterceptorAccess := initializationInterceptorLogger()

	startServer(sessionGrpc, loggerInterceptorAccess)
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

// initializeSession initializing session server
func initializeSession(db *sql.DB) *grpcSession.SessionServer {
	sessionRepository := sessionRepo.NewSessionRepository(sqlx.NewDb(db, "pgx"))
	sessionUseCase := sessionUc.NewSessionUseCase(sessionRepository)

	return grpcSession.NewSessionServer(sessionUseCase)
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
func startServer(sessionGrpc *grpcSession.SessionServer, interceptorsLogger *interceptors.Logger) {
	listen, err := net.Listen("tcp", ":8003")
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8003", err.Error())
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptorsLogger.AccessLogInterceptor,
			interceptors.PanicRecoveryInterceptor,
		),
	}
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterSessionServiceServer(grpcServer, sessionGrpc)

	fmt.Printf("The server is running in port 8003\n")

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8003", err.Error())
	}
}
