package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	migrate "github.com/rubenv/sql-migrate"
	"log"
	"mail/pkg/delivery/middleware"
	"mail/pkg/delivery/session"
	"net/http"

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

// @host localhost:8080
// @BasePath /
func main() {
	dsn := "user=postgres dbname=Mail password=postgres host=localhost port=5432 sslmode=disable"
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
		Dir: "migrations",
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

	port := 8080
	Logrus := middleware.InitializationAcceslog(port)

	router := mux.NewRouter()

	auth := mux.NewRouter().PathPrefix("/api/v1/auth").Subrouter()
	auth.Use(Logrus.AccessLogMiddleware, middleware.PanicMiddleware, middleware.AuthMiddleware)
	router.PathPrefix("/api/v1/auth").Handler(auth)

	auth.HandleFunc("/verify-auth", userHandler.VerifyAuth).Methods("GET", "OPTIONS")
	auth.HandleFunc("/get-user", userHandler.GetUserBySession).Methods("GET", "OPTIONS") //??
	auth.HandleFunc("/emails", emailHandler.List).Methods("GET", "OPTIONS")
	auth.HandleFunc("/email/{id}", emailHandler.GetByID).Methods("GET", "OPTIONS")
	auth.HandleFunc("/email/add", emailHandler.Add).Methods("POST", "OPTIONS")
	auth.HandleFunc("/email/update/{id}", emailHandler.Update).Methods("PUT", "OPTIONS")
	auth.HandleFunc("/email/delete/{id}", emailHandler.Delete).Methods("DELETE", "OPTIONS")

	logRouter := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	logRouter.Use(Logrus.AccessLogMiddleware, middleware.PanicMiddleware)
	router.PathPrefix("/api/v1").Handler(logRouter)

	logRouter.HandleFunc("/login", userHandler.Login).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/signup", userHandler.Signup).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/logout", userHandler.Logout).Methods("POST", "OPTIONS")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8081", "http://89.208.223.140:8081", "http://localhost:8080", "http://localhost:8081", "http://89.208.223.140:8080"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
	})

	corsHandler := c.Handler(router)

	fmt.Printf("The server is running on http://localhost:%d\n", port)
	fmt.Printf("Swagger is running on http://localhost:%d/swagger/index.html\n", port)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), corsHandler)
	if err != nil {
		fmt.Println("Error when starting the server:", err)
	}
	// 89.208.223.140
}
