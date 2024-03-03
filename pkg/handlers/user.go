package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mail/pkg/session"
	"mail/pkg/user"

	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	UserRepository user.UserRepository
	Sessions       *session.SessionsManager
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	b, _ := bcrypt.GenerateFromPassword([]byte("1234"), 14)
	fmt.Println("1234 ---> ", string(b))
	return string(bytes), err
}

// VerifyAuth verifies user authentication.
// @Summary Verify user authentication
// @Description Verify user authentication using sessions
// @Tags users
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 401 {string} string "Not Authorized"
// @Router /verify-auth [get]
func (uh *UserHandler) VerifyAuth(w http.ResponseWriter, r *http.Request) {
	_, err := uh.Sessions.Check(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Login handles user login.
// @Summary User login
// @Description Handles user login.
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body user.User true "User credentials for login"
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Failed to create session"
// @Router /login [post]
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials user.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Data validation
	users, err := uh.UserRepository.GetAll()
	ourUser, ourUserDefault := user.User{}, user.User{}
	for _, u := range users {
		if u.Login == credentials.Login {
			if user.CheckPasswordHash(credentials.Password, u.Password) {
				ourUser = *u
				break
			} else {
				break
			}
		}
	}
	if ourUser == ourUserDefault {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Login failed"))
		return
	}
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	_, er := uh.Sessions.Create(w, ourUser.ID)
	if er != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// Signup handles user signup.
// @Summary User signup
// @Description Handles user signup.
// @Tags users
// @Accept json
// @Produce json
// @Param newUser body user.User true "New user details for signup"
// @Success 200 {string} string "Signup successful"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to add user"
// @Router /signup [post]
func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var newUser user.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Data validation
	if newUser.Name == "" || newUser.Surname == "" || newUser.Login == "" || newUser.Password == "" {
		http.Error(w, `All fields must be filled in`, http.StatusBadRequest)
		return
	}
	users, er := uh.UserRepository.GetAll()
	if er != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, u := range users {
		if u.Login == newUser.Login {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Such a login already exists"))
			return
		}
	}

	newUser.Password, err = HashPassword(newUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, erro := uh.UserRepository.Add(&newUser)
	if erro != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Signup successful"))
}

// Logout handles user logout.
// @Summary User logout
// @Description Handles user logout.
// @Tags users
// @Produce json
// @Success 200 {string} string "Logout successful"
// @Router /logout [post]
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := uh.Sessions.DestroyCurrent(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Logged out"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}
