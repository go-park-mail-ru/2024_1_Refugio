// 89.208.223.140
package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"mail/internal/models/microservice_ports"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/kataras/requestid"
	"github.com/rs/cors"

	"mail/internal/pkg/logger"
	"mail/internal/pkg/middleware"
	"mail/internal/pkg/session"
	"mail/internal/pkg/utils/connect_microservice"

	migrate "github.com/rubenv/sql-migrate"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	session_proto "mail/internal/microservice/session/proto"
	authHand "mail/internal/pkg/auth/delivery/http"
	emailHand "mail/internal/pkg/email/delivery/http"
	folderHand "mail/internal/pkg/folder/delivery/http"
	userHand "mail/internal/pkg/user/delivery/http"

	_ "mail/docs"
)

// @title API Mail
// @version 1.0
// @description API server for mail

// @host mailhub.su
// @BasePath /
func main() {
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	migrateDatabase(db)

	loggerMiddlewareAccess := initializeMiddlewareLogger()

	sessionsManager := initializeSessionsManager()
	authHandler := initializeAuthHandler(sessionsManager)
	emailHandler := initializeEmailHandler(sessionsManager)
	userHandler := initializeUserHandler(sessionsManager)
	folderHandler := initializeFolderHandler(sessionsManager)

	router := setupRouter(authHandler, userHandler, emailHandler, folderHandler, loggerMiddlewareAccess)

	startServer(router)
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

func migrateDatabase(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	_, errMigration := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if errMigration != nil {
		log.Fatalf("Failed to apply migrations: %v", errMigration)
	}
}

func initializeSessionsManager() *session.SessionsManager {
	sessionsManager := session.NewSessionsManager()
	session.InitializationGlobalSeaaionManager(sessionsManager)

	StartSessionCleaner(24 * time.Hour)

	return sessionsManager
}

func initializeAuthHandler(sessionsManager *session.SessionsManager) *authHand.AuthHandler {
	return &authHand.AuthHandler{
		Sessions: sessionsManager,
	}
}

func initializeEmailHandler(sessionsManager *session.SessionsManager) *emailHand.EmailHandler {

	return &emailHand.EmailHandler{
		Sessions: sessionsManager,
	}
}

func initializeUserHandler(sessionsManager *session.SessionsManager) *userHand.UserHandler {
	return &userHand.UserHandler{
		Sessions: sessionsManager,
	}
}

func initializeFolderHandler(sessionsManager *session.SessionsManager) *folderHand.FolderHandler {
	return &folderHand.FolderHandler{
		Sessions: sessionsManager,
	}
}

func initializeMiddlewareLogger() *middleware.Logger {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}

	LogrusAcces := logger.InitializationAccesLog(f)
	LoggerAcces := new(middleware.Logger)
	LoggerAcces.Logger = LogrusAcces

	return LoggerAcces
}

func setupRouter(authHandler *authHand.AuthHandler, userHandler *userHand.UserHandler, emailHandler *emailHand.EmailHandler, folderHandler *folderHand.FolderHandler, logger *middleware.Logger) http.Handler {
	router := mux.NewRouter()

	auth := setupAuthRouter(authHandler, emailHandler, logger)
	router.PathPrefix("/api/v1/auth").Handler(auth)

	logRouter := setupLogRouter(emailHandler, userHandler, folderHandler, logger)
	router.PathPrefix("/api/v1").Handler(logRouter)

	staticDir := "/media/"
	staticFileServer := http.StripPrefix(staticDir, http.FileServer(http.Dir("./avatars")))
	router.PathPrefix(staticDir).Handler(staticFileServer)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return logger.AccessLogMiddleware(router)
}

func setupAuthRouter(authHandler *authHand.AuthHandler, emailHandler *emailHand.EmailHandler, logger *middleware.Logger) http.Handler {
	auth := mux.NewRouter().PathPrefix("/api/v1/auth").Subrouter()
	auth.Use(logger.AccessLogMiddleware, middleware.PanicMiddleware)

	auth.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	auth.HandleFunc("/signup", authHandler.Signup).Methods("POST", "OPTIONS")
	auth.HandleFunc("/logout", authHandler.Logout).Methods("POST", "OPTIONS")
	auth.HandleFunc("/sendOther", emailHandler.SendFromAnotherDomain).Methods("POST", "OPTIONS")

	return auth
}

func setupLogRouter(emailHandler *emailHand.EmailHandler, userHandler *userHand.UserHandler, folderHandler *folderHand.FolderHandler, logger *middleware.Logger) http.Handler {
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

	logRouter.HandleFunc("/folder/add", folderHandler.Add).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/folder/all", folderHandler.GetAll).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/folder/delete/{id}", folderHandler.Delete).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/folder/update/{id}", folderHandler.Update).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/folder/add_email", folderHandler.AddEmailInFolder).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/folder/delete_email", folderHandler.DeleteEmailInFolder).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/folder/all_emails/{id}", folderHandler.GetAllEmailsInFolder).Methods("GET", "OPTIONS")

	return logRouter
}

func startServer(router http.Handler) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8081", "http://89.208.223.140:8081", "http://mailhub.su:8081", "http://mailhub.su:8080", "http://localhost:8080", "http://localhost:8081", "http://89.208.223.140:8080"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Csrf-Token", "Content-Type"},
		ExposedHeaders:   []string{"X-Csrf-Token"},
	})

	corsHandler := c.Handler(router)

	port := 8080
	fmt.Printf("The server is running on http://localhost:%d\n", port)
	fmt.Printf("Swagger is running on http://localhost:%d/swagger/index.html\n", port)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), requestid.Handler(corsHandler))
	if err != nil {
		fmt.Println("Error when starting the server:", err)
	}
}

func StartSessionCleaner(interval time.Duration) {
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

				conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
				if err != nil {
					fmt.Printf("connection fail")
					conn.Close()
					return
				}

				sessionServiceClient := session_proto.NewSessionServiceClient(conn)
				req, err := sessionServiceClient.CleanupExpiredSessions(
					metadata.NewOutgoingContext(ctx,
						metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
					&session_proto.CleanupExpiredSessionsRequest{},
				)
				if err != nil {
					fmt.Printf("Error cleaning expired sessions: %v\n", err)
					conn.Close()
					return
				}
				fmt.Println(req)

				conn.Close()
			}
		}
	}()
}
