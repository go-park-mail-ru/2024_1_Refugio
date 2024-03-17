package main

import (
	"fmt"
	"github.com/rs/cors"
	"mail/pkg/delivery/session"
	"net/http"

	emailHand "mail/pkg/delivery/email"
	userHand "mail/pkg/delivery/user"
	emailRepo "mail/pkg/repository/email"
	userRepo "mail/pkg/repository/user"
	emailUc "mail/pkg/usecase/email"
	userUc "mail/pkg/usecase/user"

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
	sessionsManager := session.NewSessionsManager()

	emailRepository := emailRepo.NewEmailMemoryRepository()
	emailUseCase := emailUc.NewEmailUseCase(emailRepository)

	userRepository := userRepo.NewInMemoryUserRepository()
	userUseCase := userUc.NewUserUseCase(userRepository)

	emailHandler := &emailHand.EmailHandler{
		EmailUseCase: emailUseCase,
		Sessions:     sessionsManager,
	}

	userHandler := &userHand.UserHandler{
		UserUseCase: userUseCase,
		Sessions:    sessionsManager,
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
		AllowedOrigins:   []string{"http://127.0.0.1:8081", "http://89.208.223.140:8081", "http://localhost:8080", "http://localhost:8081", "http://89.208.223.140:8080"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
	})

	corsHandler := c.Handler(router)

	port := 8080
	fmt.Printf("The server is running on http://localhost:%d\n", port)
	fmt.Printf("Swagger is running on http://localhost:%d/swagger/index.html\n", port)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), corsHandler)
	if err != nil {
		fmt.Println("Error when starting the server:", err)
	}
	// 89.208.223.140
}
