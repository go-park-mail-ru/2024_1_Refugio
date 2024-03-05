package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"mail/pkg/session"
	"mail/pkg/user"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	UserRepository user.UserRepository
	Sessions       *session.SessionsManager
}

// VerifyAuth verifies user authentication.
// @Summary Verify user authentication
// @Description Verify user authentication using sessions
// @Tags users
// @Produce json
// @Success 200 {object} Response "OK"
// @Failure 401 {object} Response "Not Authorized"
// @Router /api/v1/verify-auth [get]
func (uh *UserHandler) VerifyAuth(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)

	_, err := uh.Sessions.Check(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	handleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "OK"})
}

// Login handles user login.
// @Summary User login
// @Description Handles user login.
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body user.User true "User credentials for login"
// @Success 200 {object} Response "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Failed to create session"
// @Router /api/v1/login [post]
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)

	var credentials user.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if isEmpty(credentials.Login) || isEmpty(credentials.Password) {
		handleError(w, http.StatusInternalServerError, "All fields must be filled in")
		return
	}

	users, _ := uh.UserRepository.GetAll()
	fmt.Printf("users:%d\n", users)
	fmt.Printf("credentials:%d\n", credentials)
	ourUser, ourUserDefault := user.User{}, user.User{}
	for _, u := range users {
		if u.Login == credentials.Login {
			if user.CheckPasswordHash(credentials.Password, u.Password) {
				ourUser = *u
				fmt.Printf("u:%d\n", *u)
				break
			} else {
				break
			}
		}
	}
	if ourUser.Login == "" || ourUser == ourUserDefault {
		fmt.Printf("ourUser:%d\n", ourUser)
		fmt.Printf("ourUserDefault:%d\n", ourUserDefault)
		handleError(w, http.StatusUnauthorized, "Login failed")
		return
	}

	_, er := uh.Sessions.Create(w, ourUser.ID)
	if er != nil {
		handleError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	handleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Login successful"})
}

// Signup handles user signup.
// @Summary User signup
// @Description Handles user signup.
// @Tags users
// @Accept json
// @Produce json
// @Param newUser body user.User true "New user details for signup"
// @Success 200 {object} Response "Signup successful"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to add user"
// @Router /api/v1/signup [post]
func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)

	var newUser user.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if isEmpty(newUser.Name) || isEmpty(newUser.Surname) || isEmpty(newUser.Login) || isEmpty(newUser.Password) {
		handleError(w, http.StatusBadRequest, "All fields must be filled in")
	}

	users, _ := uh.UserRepository.GetAll()
	for _, u := range users {
		if u.Login == newUser.Login {
			handleError(w, http.StatusBadRequest, "Such a login already exists")
			return
		}
	}

	_, er := uh.UserRepository.Add(&newUser)
	if er != nil {
		handleError(w, http.StatusInternalServerError, "Failed to add user")
	}

	handleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Signup successful"})
}

// Logout handles user logout.
// @Summary User logout
// @Description Handles user logout.
// @Tags users
// @Produce json
// @Success 200 {object} Response "Logout successful"
// @Router /api/v1/logout [post]
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)

	err := uh.Sessions.DestroyCurrent(w, r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	handleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Logout successful"})
}

// GetUserBySession retrieves the user associated with the current session.
// @Summary Get user by session
// @Description Retrieve the user associated with the current session
// @Tags users
// @Produce json
// @Success 200 {object} Response "User details"
// @Failure 401 {object} Response "Not Authorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/get-user [get]
func (uh *UserHandler) GetUserBySession(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)

	sessionUser, err := uh.Sessions.Check(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	userData, err := uh.UserRepository.GetByID(sessionUser.UserID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userData.Password = ""
	handleSuccess(w, http.StatusOK, map[string]interface{}{"user": userData})
}

// isEmpty checks if the given string is empty after trimming leading and trailing whitespace.
// Returns true if the string is empty, and false otherwise.
func isEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}
