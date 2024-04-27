package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	questionnaireRepo "mail/internal/microservice/questionnaire/repository"
	grpcQuestionnaire "mail/internal/microservice/questionnaire/server"
	questionnaireUc "mail/internal/microservice/questionnaire/usecase"
	"net"
	"os"
	"time"

	migrate "github.com/rubenv/sql-migrate"
	"mail/internal/microservice/interceptors"
	"mail/internal/microservice/questionnaire/proto"
)

func main() {
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	migrateDatabase(db)

	questionGrpc := initializeQuestion(db)

	loggerInterceptorAccess := initializationInterceptorLogger()

	startServer(questionGrpc, loggerInterceptorAccess)
}

func settingTime() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error loc time")
	}

	time.Local = loc
}

func initializeDatabase() *sql.DB {
	dsn := "user=postgres dbname=Question password=postgres host=localhost port=5432 sslmode=disable"
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

func initializeQuestion(db *sql.DB) *grpcQuestionnaire.QuestionAnswerServer {
	questionnaireRepository := questionnaireRepo.NewQuestionRepository(sqlx.NewDb(db, "pgx"))
	questionnaireUseCase := questionnaireUc.NewQuestionAnswerUseCase(questionnaireRepository)

	return grpcQuestionnaire.NewQestionAnswerServer(questionnaireUseCase)
}

func migrateDatabase(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: ".",
	}

	_, errMigration := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if errMigration != nil {
		log.Fatalf("Failed to apply migrations: %v", errMigration)
	}
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

func startServer(questionnaireGrpc *grpcQuestionnaire.QuestionAnswerServer, interceptorsLogger *interceptors.Logger) {
	listen, err := net.Listen("tcp", ":8006")
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8006", err.Error())
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptorsLogger.AccessLogInterceptor,
			interceptors.PanicRecoveryInterceptor,
		),
	}
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterQuestionServiceServer(grpcServer, questionnaireGrpc)

	fmt.Printf("The server is running  in port 8006\n")

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Cannot listen port: %s. Err: %s", "8006", err.Error())
	}
}
