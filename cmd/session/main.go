package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"

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

	startServer(sessionGrpc)
}

func settingTime() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error loc time")
	}

	time.Local = loc
}

func initializeDatabase() *sql.DB {
	dsn := "user=postgres dbname=Mail password=postgres host=localhost port=5432 sslmode=disable"
	// dsn := "user=postgres dbname=Mail password=postgres host=89.208.223.140 port=5432 sslmode=disable"
	db, err := sql.Open("pgx", dsn)
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

func initializeSession(db *sql.DB) *grpcSession.SessionServer {
	sessionRepository := sessionRepo.NewSessionRepository(sqlx.NewDb(db, "pgx"))
	sessionUseCase := sessionUc.NewSessionUseCase(sessionRepository)

	return grpcSession.NewSessionServer(sessionUseCase)
}

func startServer(sessionGrpc *grpcSession.SessionServer) {
	listen, err := net.Listen("tcp", ":8003")
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8003", err.Error())
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
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
