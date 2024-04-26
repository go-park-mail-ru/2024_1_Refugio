package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"

	"mail/internal/microservice/interceptors"
	"mail/internal/microservice/user/proto"

	userRepo "mail/internal/microservice/user/repository"
	grpcUser "mail/internal/microservice/user/server"
	userUc "mail/internal/microservice/user/usecase"
)

func main() {
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	userGrpc := initializeUser(db)

	loggerInterceptorAccess := initializationInterceptorLogger()

	startServer(userGrpc, loggerInterceptorAccess)
}

func settingTime() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error loc time")
	}

	time.Local = loc
}

func initializeDatabase() *sql.DB {
	// dsn := "user=postgres dbname=Mail password=postgres host=localhost port=5432 sslmode=disable"
	dsn := "user=postgres dbname=Mail password=postgres host=89.208.223.140 port=5432 sslmode=disable"
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

func initializeUser(db *sql.DB) *grpcUser.UserServer {
	userRepository := userRepo.NewUserRepository(sqlx.NewDb(db, "pgx"))
	userUseCase := userUc.NewUserUseCase(userRepository)

	return grpcUser.NewUserServer(userUseCase)
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

func startServer(userGrpc *grpcUser.UserServer, interceptorsLogger *interceptors.Logger) {
	listen, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8001", err.Error())
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptorsLogger.AccessLogInterceptor,
			interceptors.PanicRecoveryInterceptor,
		),
	}
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterUserServiceServer(grpcServer, userGrpc)

	fmt.Printf("The server is running in port 8001\n")

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8001", err.Error())
	}
}
