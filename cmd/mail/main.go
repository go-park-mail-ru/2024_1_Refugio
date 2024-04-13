// 89.208.223.140
package main

import (
	"database/sql"
	"fmt"
	"log"
	"mail/internal/pkg/logger"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"context"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/requestid"
	"github.com/rs/cors"

	"mail/internal/pkg/middleware"
	"mail/internal/pkg/session"

	migrate "github.com/rubenv/sql-migrate"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	userHand "mail/internal/pkg/auth/delivery/http"
	sessionRepo "mail/internal/pkg/auth/repository"
	userRepo "mail/internal/pkg/auth/repository"
	sessionUc "mail/internal/pkg/auth/usecase"
	userUc "mail/internal/pkg/auth/usecase"
	emailHand "mail/internal/pkg/email/delivery/http"
	emailRepo "mail/internal/pkg/email/repository"
	emailUc "mail/internal/pkg/email/usecase"

	_ "mail/docs"
)

// @title API Mail
// @version 1.0
// @description API server for mail

// @host localhost:8080
// @BasePath /
func main() {
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	migrateDatabase(db)

	LoggerAcces := initializeLogger()

	sessionsManager := initializeSessionsManager(db)
	emailHandler := initializeEmailHandler(db, sessionsManager)
	userHandler := initializeUserHandler(db, sessionsManager)

	router := setupRouter(userHandler, emailHandler, LoggerAcces)

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

func migrateDatabase(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	_, errMigration := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if errMigration != nil {
		log.Fatalf("Failed to apply migrations: %v", errMigration)
	}
}

func initializeSessionsManager(db *sql.DB) *session.SessionsManager {
	sessionRepository := sessionRepo.NewSessionRepository(sqlx.NewDb(db, "pgx"))
	sessionUsaCase := sessionUc.NewSessionUseCase(sessionRepository)
	sessionsManager := session.NewSessionsManager(sessionUsaCase)
	session.InitializationGlobalSeaaionManager(sessionsManager)

	StartSessionCleaner(sessionUsaCase, 24*time.Hour)

	return sessionsManager
}

func initializeEmailHandler(db *sql.DB, sessionsManager *session.SessionsManager) *emailHand.EmailHandler {
	emailRepository := emailRepo.NewEmailRepository(sqlx.NewDb(db, "pgx"))
	emailUseCase := emailUc.NewEmailUseCase(emailRepository)

	return &emailHand.EmailHandler{
		EmailUseCase: emailUseCase,
		Sessions:     sessionsManager,
	}
}

func initializeUserHandler(db *sql.DB, sessionsManager *session.SessionsManager) *userHand.UserHandler {
	userRepository := userRepo.NewUserRepository(sqlx.NewDb(db, "pgx"))
	userUseCase := userUc.NewUserUseCase(userRepository)

	return &userHand.UserHandler{
		UserUseCase: userUseCase,
		Sessions:    sessionsManager,
	}
}

func initializeLogger() *middleware.Logger {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}
	defer f.Close()
	LogrusAcces := logger.InitializationAccesLog(f)
	LoggerAcces := new(middleware.Logger)
	LoggerAcces.Logger = LogrusAcces

	return LoggerAcces
}

func setupRouter(userHandler *userHand.UserHandler, emailHandler *emailHand.EmailHandler, logger *middleware.Logger) http.Handler {
	router := mux.NewRouter()

	auth := setupAuthRouter(userHandler, emailHandler, logger)
	router.PathPrefix("/api/v1/auth").Handler(auth)

	logRouter := setupLogRouter(emailHandler, userHandler, logger)
	router.PathPrefix("/api/v1").Handler(logRouter)

	staticDir := "/media/"
	staticFileServer := http.StripPrefix(staticDir, http.FileServer(http.Dir("./avatars")))
	router.PathPrefix(staticDir).Handler(staticFileServer)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return logger.AccessLogMiddleware(router)
}

func setupAuthRouter(userHandler *userHand.UserHandler, emailHandler *emailHand.EmailHandler, logger *middleware.Logger) http.Handler {
	auth := mux.NewRouter().PathPrefix("/api/v1/auth").Subrouter()
	auth.Use(logger.AccessLogMiddleware, middleware.PanicMiddleware)

	auth.HandleFunc("/login", userHandler.Login).Methods("POST", "OPTIONS")
	auth.HandleFunc("/signup", userHandler.Signup).Methods("POST", "OPTIONS")
	auth.HandleFunc("/logout", userHandler.Logout).Methods("POST", "OPTIONS")
	auth.HandleFunc("/sendOther", emailHandler.SendFromAnotherDomain).Methods("POST", "OPTIONS")

	return auth
}

func setupLogRouter(emailHandler *emailHand.EmailHandler, userHandler *userHand.UserHandler, logger *middleware.Logger) http.Handler {
	logRouter := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	logRouter.Use(logger.AccessLogMiddleware, middleware.PanicMiddleware, middleware.AuthMiddleware)

	logRouter.HandleFunc("/verify-auth", userHandler.VerifyAuth).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/user/get", userHandler.GetUserBySession).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/user/update", userHandler.UpdateUserData).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/user/delete/{id}", userHandler.DeleteUserData).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/user/avatar/upload", userHandler.UploadUserAvatar).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/emails/incoming", emailHandler.Incoming).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/emails/sent", emailHandler.Sent).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/email/{id}", emailHandler.GetByID).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/email/update/{id}", emailHandler.Update).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/email/delete/{id}", emailHandler.Delete).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/email/send", emailHandler.Send).Methods("POST", "OPTIONS")

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

func StartSessionCleaner(sessionCleaner *sessionUc.SessionUseCase, interval time.Duration) {
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

				err = sessionCleaner.CleanupExpiredSessions(ctx)
				if err != nil {
					fmt.Printf("Error cleaning expired sessions: %v\n", err)
				}
			}
		}
	}()
}
