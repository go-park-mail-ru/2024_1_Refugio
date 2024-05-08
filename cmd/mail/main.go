// 89.208.223.140
package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/kataras/requestid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"

	"mail/cmd/configs"
	"mail/internal/models/microservice_ports"
	"mail/internal/monitoring"
	"mail/internal/pkg/logger"
	"mail/internal/pkg/middleware"
	"mail/internal/pkg/session"
	"mail/internal/pkg/utils/connect_microservice"

	migrate "github.com/rubenv/sql-migrate"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	auth_proto "mail/internal/microservice/auth/proto"
	email_proto "mail/internal/microservice/email/proto"
	folder_proto "mail/internal/microservice/folder/proto"
	question_proto "mail/internal/microservice/questionnaire/proto"
	session_proto "mail/internal/microservice/session/proto"
	user_proto "mail/internal/microservice/user/proto"
	authHand "mail/internal/pkg/auth/delivery/http"
	emailHand "mail/internal/pkg/email/delivery/http"
	folderHand "mail/internal/pkg/folder/delivery/http"
	questionHand "mail/internal/pkg/questionnairy/delivery/http"
	userHand "mail/internal/pkg/user/delivery/http"

	_ "mail/docs"
)

// @title API Mailhub
// @version 1.0
// @description API server for Mailhub

// @host mailhub.su
// @BasePath /
func main() {
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	migrateDatabase(db)

	loggerMiddlewareAccess := initializeMiddlewareLogger()

	sessionManagerServiceConn, err := connect_microservice.
		OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		log.Fatalf("connection with microservice auth fail")
	}
	defer sessionManagerServiceConn.Close()
	sessionsManager := initializeSessionsManager(session_proto.NewSessionServiceClient(sessionManagerServiceConn))

	authServiceConn, err := connect_microservice.
		OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.AuthService))
	if err != nil {
		log.Fatalf("connection with microservice auth fail")
	}
	defer authServiceConn.Close()
	authHandler := initializeAuthHandler(sessionsManager, auth_proto.NewAuthServiceClient(authServiceConn))

	userServiceConn, err := connect_microservice.
		OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		log.Fatalf("connection with microservice user fail")
	}
	defer userServiceConn.Close()
	userHandler := initializeUserHandler(sessionsManager, user_proto.NewUserServiceClient(userServiceConn))

	emailServiceConn, err := connect_microservice.
		OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
	if err != nil {
		log.Fatalf("connection with microservice user fail")
	}
	defer emailServiceConn.Close()
	emailHandler := initializeEmailHandler(sessionsManager, email_proto.NewEmailServiceClient(emailServiceConn))

	folderServiceConn, err := connect_microservice.
		OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.FolderService))
	if err != nil {
		log.Fatalf("connection with microservice user fail")
	}
	defer folderServiceConn.Close()
	folderHandler := initializeFolderHandler(sessionsManager, folder_proto.NewFolderServiceClient(folderServiceConn))

	questionServiceConn, err := connect_microservice.
		OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.QuestionService))
	if err != nil {
		log.Fatalf("connection with microservice question fail")
	}
	defer questionServiceConn.Close()
	questionHandler := initializeQuestionHandler(sessionsManager, question_proto.NewQuestionServiceClient(questionServiceConn))

	router := setupRouter(authHandler, userHandler, emailHandler, folderHandler, questionHandler, loggerMiddlewareAccess)

	startServer(router)
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

// migrateDatabase applying database migration
func migrateDatabase(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	_, errMigration := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if errMigration != nil {
		log.Fatalf("Failed to apply migrations: %v", errMigration)
	}
}

// initializeSessionsManager initializing session manager
func initializeSessionsManager(sessionServiceClient session_proto.SessionServiceClient) *session.SessionsManager {
	sessionsManager := session.NewSessionsManager(sessionServiceClient)
	session.InitializationGlobalSessionManager(sessionsManager)

	startSessionCleaner(24*time.Hour, sessionServiceClient)

	return sessionsManager
}

// initializeAuthHandler initializing authorization handler
func initializeAuthHandler(sessionsManager *session.SessionsManager, authServiceClient auth_proto.AuthServiceClient) *authHand.AuthHandler {
	return &authHand.AuthHandler{
		Sessions:          sessionsManager,
		AuthServiceClient: authServiceClient,
	}
}

// initializeEmailHandler initializing email handler
func initializeEmailHandler(sessionsManager *session.SessionsManager, emailServiceClient email_proto.EmailServiceClient) *emailHand.EmailHandler {
	return &emailHand.EmailHandler{
		Sessions:           sessionsManager,
		EmailServiceClient: emailServiceClient,
	}
}

// initializeUserHandler initializing user handler
func initializeUserHandler(sessionsManager *session.SessionsManager, userServiceClient user_proto.UserServiceClient) *userHand.UserHandler {
	return &userHand.UserHandler{
		Sessions:          sessionsManager,
		UserServiceClient: userServiceClient,
	}
}

// initializeFolderHandler initializing folder handler
func initializeFolderHandler(sessionsManager *session.SessionsManager, folderServiceClient folder_proto.FolderServiceClient) *folderHand.FolderHandler {
	return &folderHand.FolderHandler{
		Sessions:            sessionsManager,
		FolderServiceClient: folderServiceClient,
	}
}

// initializeQuestionHandler initializing question handler
func initializeQuestionHandler(sessionsManager *session.SessionsManager, questionServiceClient question_proto.QuestionServiceClient) *questionHand.QuestionHandler {
	return &questionHand.QuestionHandler{
		Sessions:              sessionsManager,
		QuestionServiceClient: questionServiceClient,
	}
}

// initializeMiddlewareLogger initializing logger
func initializeMiddlewareLogger() *middleware.Logger {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}

	logrusAccess := logger.InitializationAccesLog(f)
	loggerAccess := new(middleware.Logger)
	loggerAccess.Logger = logrusAccess

	return loggerAccess
}

// setupRouter configuring routers
func setupRouter(authHandler *authHand.AuthHandler, userHandler *userHand.UserHandler, emailHandler *emailHand.EmailHandler, folderHandler *folderHand.FolderHandler, questionHandler *questionHand.QuestionHandler, logger *middleware.Logger) http.Handler {
	router := mux.NewRouter()

	auth := setupAuthRouter(authHandler, emailHandler, logger)
	router.PathPrefix("/api/v1/auth").Handler(auth)

	logRouter := setupLogRouter(emailHandler, userHandler, folderHandler, questionHandler, logger)
	router.PathPrefix("/api/v1").Handler(logRouter)

	staticDir := "/media/"
	staticFileServer := http.StripPrefix(staticDir, http.FileServer(http.Dir("./avatars")))
	router.PathPrefix(staticDir).Handler(staticFileServer)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	logger.Metrics = monitoring.RegisterMonitoring(router)

	return logger.AccessLogMiddleware(router)
}

// setupAuthRouter configuring authorization router
func setupAuthRouter(authHandler *authHand.AuthHandler, emailHandler *emailHand.EmailHandler, logger *middleware.Logger) http.Handler {
	auth := mux.NewRouter().PathPrefix("/api/v1/auth").Subrouter()
	auth.Use(logger.AccessLogMiddleware, middleware.PanicMiddleware)

	auth.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	auth.HandleFunc("/signup", authHandler.Signup).Methods("POST", "OPTIONS")
	auth.HandleFunc("/logout", authHandler.Logout).Methods("POST", "OPTIONS")
	auth.HandleFunc("/sendOther", emailHandler.SendFromAnotherDomain).Methods("POST", "OPTIONS")

	return auth
}

// setupLogRouter configuring router with logger
func setupLogRouter(emailHandler *emailHand.EmailHandler, userHandler *userHand.UserHandler, folderHandler *folderHand.FolderHandler, questionHandler *questionHand.QuestionHandler, logger *middleware.Logger) http.Handler {
	logRouter := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	logRouter.Use(logger.AccessLogMiddleware, middleware.PanicMiddleware, middleware.AuthMiddleware)

	logRouter.HandleFunc("/verify-auth", userHandler.VerifyAuth).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/user/get", userHandler.GetUserBySession).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/user/update", userHandler.UpdateUserData).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/user/delete/{id}", userHandler.DeleteUserData).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/user/avatar/upload", userHandler.UploadUserAvatar).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/user/avatar/delete", userHandler.DeleteUserAvatar).Methods("DELETE", "OPTIONS")

	logRouter.HandleFunc("/emails/incoming", emailHandler.Incoming).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/emails/sent", emailHandler.Sent).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/emails/draft", emailHandler.Draft).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/emails/spam", emailHandler.Spam).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/email/{id}", emailHandler.GetByID).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/email/update/{id}", emailHandler.Update).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/email/delete/{id}", emailHandler.Delete).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/email/send", emailHandler.Send).Methods("POST", "OPTIONS")

	logRouter.HandleFunc("/questions", questionHandler.GetAllQuestions).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/questions", questionHandler.AddQuestion).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/answers", questionHandler.AddAnswer).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/statistics", questionHandler.GetStatistics).Methods("GET", "OPTIONS")

	logRouter.HandleFunc("/folder/add", folderHandler.Add).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/folder/all", folderHandler.GetAll).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/folder/delete/{id}", folderHandler.Delete).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/folder/update/{id}", folderHandler.Update).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/folder/add_email", folderHandler.AddEmailInFolder).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/folder/delete_email", folderHandler.DeleteEmailInFolder).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/folder/all_emails/{id}", folderHandler.GetAllEmailsInFolder).Methods("GET", "OPTIONS")

	return logRouter
}

// startServer starting server
func startServer(router http.Handler) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1:8081", "http://89.208.223.140:8081", "http://mailhub.su:8081", "http://mailhub.su:8080", "http://localhost:8080", "http://localhost:8081", "http://89.208.223.140:8080",
			"https://127.0.0.1:8081", "https://89.208.223.140:8081", "https://mailhub.su:8081", "https://mailhub.su:8080", "https://localhost:8080", "https://localhost:8081", "https://89.208.223.140:8080",
			"https://127.0.0.1", "https://89.208.223.140", "https://mailhub.su", "https://mailhub.su", "https://localhost", "https://localhost", "https://89.208.223.140"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Csrf-Token", "Content-Type"},
		ExposedHeaders:   []string{"X-Csrf-Token"},
	})

	corsHandler := c.Handler(router)

	port := 8080
	fmt.Printf("The server is running on http://localhost:%d\n", port)
	fmt.Printf("Swagger is running on http://localhost:%d/swagger/index.html\n", port)

	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), requestid.Handler(corsHandler))
	if err != nil {
		fmt.Println("Error when starting the server:", err)
	}
}

// StartSessionCleaner starting session cleanup
func startSessionCleaner(interval time.Duration, sessionServiceClient session_proto.SessionServiceClient) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					fmt.Println("Failed to create logfile" + "log.txt")
				}
				defer f.Close()

				c := context.WithValue(context.Background(), "logger", logger.InitializationBdLog(f))
				ctx := context.WithValue(c, "requestID", "DeleteExpiredSessionsNULL")

				req, err := sessionServiceClient.CleanupExpiredSessions(
					metadata.NewOutgoingContext(ctx,
						metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
					&session_proto.CleanupExpiredSessionsRequest{},
				)
				if err != nil {
					fmt.Printf("Error cleaning expired sessions: %v\n", err)
					return
				}
				fmt.Println(req)
			}
		}
	}()
}
