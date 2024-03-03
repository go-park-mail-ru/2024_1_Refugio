package main

import (
	"fmt"
	"github.com/rs/cors"
	"net/http"

	"mail/pkg/email"
	"mail/pkg/handlers"
	"mail/pkg/session"
	"mail/pkg/user"

	_ "mail/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title API Mail
// @version 1.0
// @description API server for mail

// @host localhost:8080
// @BasePath /
func main() {
	sessionsManager := session.NewSessionsManager()

	emailRepository := email.NewEmailMemoryRepository()
	userRepository := user.NewInMemoryUserRepository()

	emailHandler := &handlers.EmailHandler{
		EmailRepository: emailRepository,
		Sessions:        sessionsManager,
	}

	userHandler := &handlers.UserHandler{
		UserRepository: userRepository,
		Sessions:       sessionsManager,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/emails", emailHandler.List).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/email/{id}", emailHandler.GetByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/email/add", emailHandler.Add).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/email/update/{id}", emailHandler.Update).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/email/delete/{id}", emailHandler.Delete).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/v1/verify-auth", userHandler.VerifyAuth).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/login", userHandler.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/signup", userHandler.Signup).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/logout", userHandler.Logout).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/get-user", userHandler.GetUserBySession).Methods("GET", "OPTIONS")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8081", "http://localhost:8080", "http://89.208.223.140:8080", "http://localhost:8081"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
	})

	corsHandler := c.Handler(router)

	port := 8080
	fmt.Printf("The server is running on http://0.0.0.0:%d\n", port)
	fmt.Printf("Swagger is running on http://0.0.0.0:%d/swagger/index.html\n", port)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), corsHandler)
	if err != nil {
		fmt.Println("Error when starting the server:", err)
	}
	// 89.208.223.140
}
