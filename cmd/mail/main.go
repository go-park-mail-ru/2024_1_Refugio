package main

import (
	"fmt"
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
	router.HandleFunc("/emails", emailHandler.List).Methods("GET")
	router.HandleFunc("/email/{id}", emailHandler.GetByID).Methods("GET")
	router.HandleFunc("/email/add", emailHandler.Add).Methods("POST")
	router.HandleFunc("/email/update/{id}", emailHandler.Update).Methods("PUT")
	router.HandleFunc("/email/delete/{id}", emailHandler.Delete).Methods("DELETE")

	router.HandleFunc("/verify-auth", userHandler.VerifyAuth).Methods("GET")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/signup", userHandler.Signup).Methods("POST")
	router.HandleFunc("/logout", userHandler.Logout).Methods("POST")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	port := 8080
	fmt.Printf("The server is running on http://localhost:%d\n", port)
	fmt.Printf("Swagger is running on http://localhost:%d/swagger/index.html\n", port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), router)
	if err != nil {
		fmt.Println("Error when starting the server:", err)
	}
}
