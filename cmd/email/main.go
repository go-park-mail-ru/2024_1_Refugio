package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	grpcEmail "mail/internal/microservice/email/server"
	"net"
	"time"

	"mail/internal/microservice/email/proto"
	"mail/internal/microservice/interceptors"

	emailRepo "mail/internal/microservice/email/repository"
	emailUc "mail/internal/microservice/email/usecase"
)

func main() {
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	emailGrpc := initializeEmail(db)

	startServer(emailGrpc)
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

func initializeEmail(db *sql.DB) *grpcEmail.EmailServer {
	emailRepository := emailRepo.NewEmailRepository(sqlx.NewDb(db, "pgx"))
	emailUseCase := emailUc.NewEmailUseCase(emailRepository)

	return grpcEmail.NewEmailServer(emailUseCase)
}

func startServer(emailGrpc *grpcEmail.EmailServer) {
	listen, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8002", err.Error())
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.PanicRecoveryInterceptor,
		),
	}
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterEmailServiceServer(grpcServer, emailGrpc)

	fmt.Printf("The server is running\n")

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8002", err.Error())
	}
}