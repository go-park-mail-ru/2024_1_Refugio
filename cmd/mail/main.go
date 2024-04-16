// 89.208.223.140
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/requestid"
	"github.com/rs/cors"

	"mail/pkg/delivery/middleware"
	"mail/pkg/delivery/session"
	"mail/pkg/domain/logger"

	migrate "github.com/rubenv/sql-migrate"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	authHand "mail/pkg/delivery/auth"
	emailHand "mail/pkg/delivery/email"
	userHand "mail/pkg/delivery/user"
	emailRepo "mail/pkg/repository/postgres/email"
	sessionRepo "mail/pkg/repository/postgres/session"
	userRepo "mail/pkg/repository/postgres/user"
	emailUc "mail/pkg/usecase/email"
	sessionUc "mail/pkg/usecase/session"
	userUc "mail/pkg/usecase/user"

	_ "mail/docs"
)

// @title API Mail
// @version 1.0
// @description API server for mail

// @host mailhub.su:8080
// @BasePath /
func main() {
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	migrateDatabase(db)

	sessionsManager := initializeSessionsManager(db)
	authHandler := initializeAuthHandler(db, sessionsManager)
	emailHandler := initializeEmailHandler(db, sessionsManager)
	userHandler := initializeUserHandler(db, sessionsManager)

	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}
	defer f.Close()

	Logger := initializeLogger(f)

	router := setupRouter(authHandler, emailHandler, userHandler, Logger)

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

func initializeSessionsManager(db *sql.DB) *session.SessionsManager {
	sessionRepository := sessionRepo.NewSessionRepository(sqlx.NewDb(db, "pgx"))
	sessionUsaCase := sessionUc.NewSessionUseCase(sessionRepository)
	sessionsManager := session.NewSessionsManager(sessionUsaCase)
	session.InitializationGlobalSeaaionManager(sessionsManager)

	StartSessionCleaner(sessionUsaCase, 24*time.Hour)

	return sessionsManager
}

func initializeAuthHandler(db *sql.DB, sessionsManager *session.SessionsManager) *authHand.AuthHandler {
	userRepository := userRepo.NewUserRepository(sqlx.NewDb(db, "pgx"))
	userUseCase := userUc.NewUserUseCase(userRepository)

	return &authHand.AuthHandler{
		UserUseCase: userUseCase,
		Sessions:    sessionsManager,
	}
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

func initializeLogger(f *os.File) *middleware.Logger {
	Logrus := logger.InitializationAccesLog(f)
	Logger := new(middleware.Logger)
	Logger.Logger = Logrus

	return Logger
}

func setupRouter(authHandler *authHand.AuthHandler, emailHandler *emailHand.EmailHandler, userHandler *userHand.UserHandler, logger *middleware.Logger) http.Handler {
	router := mux.NewRouter()

	auth := setupAuthRouter(authHandler, emailHandler, logger)
	router.PathPrefix("/api/v1/auth").Handler(auth)

	logRouter := setupLogRouter(emailHandler, userHandler, logger)
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
		AllowedOrigins:   []string{
			"http://127.0.0.1:8081", "http://89.208.223.140:8081", "http://mailhub.su:8081", "http://mailhub.su:8080", "http://localhost:8080", "http://localhost:8081", "http://89.208.223.140:8080",
			"https://127.0.0.1:8081", "https://89.208.223.140:8081", "https://mailhub.su:8081", "https://mailhub.su:8080", "https://localhost:8080", "https://localhost:8081", "https://89.208.223.140:8080",
			"http://127.0.0.1", "http://89.208.223.140", "http://mailhub.su", "http://mailhub.su", "http://localhost", "http://localhost", "http://89.208.223.140"},
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
				err := sessionCleaner.CleanupExpiredSessions()
				if err != nil {
					fmt.Printf("Error cleaning expired sessions: %v\n", err)
				}
			}
		}
	}()
}
