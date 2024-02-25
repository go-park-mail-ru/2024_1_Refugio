package main

import (
	"fmt"
	"net/http"

	"mail/pkg/email"
	"mail/pkg/handlers"
	"mail/pkg/user"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "mail/docs"
)

// @title API Mail
// @version 1.0
// @description API server for mail

// @host localhost:8080
// @BasePath /
func main() {
	emailRepository := email.NewEmailMemoryRepository()
	userRepository := user.NewInMemoryUserRepository()

	emailHandler := &handlers.EmailHandler{
		EmailRepository: emailRepository,
	}

	userHandler := &handlers.UserHandler{
		UserRepository: userRepository,
	}

	router := mux.NewRouter()
	router.HandleFunc("/emails", emailHandler.List).Methods("GET")
	router.HandleFunc("/email/{id}", emailHandler.GetByID).Methods("GET")
	router.HandleFunc("/email/add", emailHandler.Add).Methods("POST")
	router.HandleFunc("/email/edit/{id}", emailHandler.Update).Methods("PUT")
	router.HandleFunc("/email/delete/{id}", emailHandler.Delete).Methods("DELETE")

	router.HandleFunc("/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/singup", userHandler.Signup).Methods("POST")
	router.HandleFunc("/logout", userHandler.Logout).Methods("POST")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	port := 8080
	fmt.Printf("The server is running on http://localhost:%d\n", port)
	fmt.Printf("Swagger is running on http://localhost:%d/swagger/index.html\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
