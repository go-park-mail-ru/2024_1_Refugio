package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/requestid"
	"github.com/rs/cors"
	migrate "github.com/rubenv/sql-migrate"
	"log"
	"mail/pkg/delivery/middleware"
	"mail/pkg/delivery/session"
	"mail/pkg/domain/logger"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	emailHand "mail/pkg/delivery/email"
	userHand "mail/pkg/delivery/user"
	emailRepo "mail/pkg/repository/postgres/email"
	sessionRepo "mail/pkg/repository/postgres/session"
	userRepo "mail/pkg/repository/postgres/user"
	emailUc "mail/pkg/usecase/email"
	sessionUc "mail/pkg/usecase/session"
	userUc "mail/pkg/usecase/user"

	"github.com/gorilla/mux"
	_ "mail/docs"
)

// @title API Mail
// @version 1.0
// @description API server for mail

// @host 89.208.223.140:8080
// @BasePath /
func main() {
	// dsn := "user=postgres dbname=Mail password=postgres host=localhost port=5432 sslmode=disable"
	dsn := "user=postgres dbname=Mail password=postgres host=89.208.223.140 port=5432 sslmode=disable"
	db, errDb := sql.Open("pgx", dsn)
	if errDb != nil {
		log.Fatalln("Can't parse config", errDb)
	}
	defer db.Close()

	errDb = db.Ping()
	if errDb != nil {
		log.Fatalln(errDb)
	}
	db.SetMaxOpenConns(10)

	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}
	_, errMigration := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if errMigration != nil {
		log.Fatalf("Failed to apply migrations: %v", errMigration)
	}
	dbx := sqlx.NewDb(db, "pgx")

	sessionRepository := sessionRepo.NewSessionRepository(dbx)
	sessionUsaCase := sessionUc.NewSessionUseCase(sessionRepository)
	sessionsManager := session.NewSessionsManager(sessionUsaCase)
	session.InitializationGlobalSeaaionManager(sessionsManager)

	StartSessionCleaner(sessionUsaCase, 24*time.Hour)

	emailRepository := emailRepo.NewEmailRepository(dbx)
	emailUseCase := emailUc.NewEmailUseCase(emailRepository)
	emailHandler := &emailHand.EmailHandler{
		EmailUseCase: emailUseCase,
		Sessions:     sessionsManager,
	}

	userRepository := userRepo.NewUserRepository(dbx)
	userUseCase := userUc.NewUserUseCase(userRepository)
	userHandler := &userHand.UserHandler{
		UserUseCase: userUseCase,
		Sessions:    sessionsManager,
	}

	Logrus := logger.InitializationAccesLog()
	Logger := new(middleware.Logger)
	Logger.Logger = Logrus

	router := mux.NewRouter()

	auth := mux.NewRouter().PathPrefix("/api/v1/auth").Subrouter()
	auth.Use(Logger.AccessLogMiddleware, middleware.PanicMiddleware)
	router.PathPrefix("/api/v1/auth").Handler(auth)

	auth.HandleFunc("/login", userHandler.Login).Methods("POST", "OPTIONS")
	auth.HandleFunc("/signup", userHandler.Signup).Methods("POST", "OPTIONS")
	auth.HandleFunc("/logout", userHandler.Logout).Methods("POST", "OPTIONS")
	auth.HandleFunc("/sendOther", emailHandler.SendFromAnotherDomain).Methods("POST", "OPTIONS")

	logRouter := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	logRouter.Use(Logger.AccessLogMiddleware, middleware.PanicMiddleware, middleware.AuthMiddleware)
	router.PathPrefix("/api/v1").Handler(logRouter)

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

	staticDir := "/media/"
	staticFileServer := http.StripPrefix(staticDir, http.FileServer(http.Dir("./avatars")))
	router.PathPrefix(staticDir).Handler(staticFileServer)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8081", "http://89.208.223.140:8081", "http://mailhub.su:8081", "http://mailhub.su:8080", "http://localhost:8080", "http://localhost:8081", "http://89.208.223.140:8080"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
	})

	corsHandler := c.Handler(router)

	port := 8080
	fmt.Printf("The server is running on http://localhost:%d\n", port)
	fmt.Printf("Swagger is running on http://localhost:%d/swagger/index.html\n", port)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), requestid.Handler(corsHandler))
	if err != nil {
		fmt.Println("Error when starting the server:", err)
	}
	// 89.208.223.140
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
