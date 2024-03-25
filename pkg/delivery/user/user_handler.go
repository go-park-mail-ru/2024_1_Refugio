package user

import (
	"encoding/json"
	"mail/pkg/delivery"
	"mail/pkg/delivery/converters"
	"mail/pkg/delivery/models"
	"mail/pkg/delivery/session"
	"mail/pkg/domain/usecase"
	"net/http"
	"strings"
)

var (
	UHandler = &UserHandler{}
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	UserUseCase usecase.UserUseCase
	Sessions    *session.SessionsManager
}

func InitializationEmailHandler(userHandler *UserHandler) {
	UHandler = userHandler
}

// VerifyAuth verifies user authentication.
// @Summary Verify user authentication
// @Description Verify user authentication using sessions
// @Tags users
// @Produce json
// @Param X-CSRF-Token header string true "CSRF Token"
// @Success 200 {object} delivery.Response "OK"
// @Failure 401 {object} delivery.Response "Not Authorized"
// @Router /api/v1/auth/verify-auth [get]
func (uh *UserHandler) VerifyAuth(w http.ResponseWriter, r *http.Request) {
	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "OK"})
}

// Login handles user login.
// @Summary User login
// @Description Handles user login.
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body delivery.UserSwag true "User credentials for login"
// @Success 200 {object} delivery.Response "Login successful"
// @Failure 400 {object} delivery.ErrorResponse "Invalid request body"
// @Failure 401 {object} delivery.ErrorResponse "Invalid credentials"
// @Failure 500 {object} delivery.ErrorResponse "Failed to create session"
// @Router /api/v1/login [post]
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if isEmpty(credentials.Login) || isEmpty(credentials.Password) {
		delivery.HandleError(w, http.StatusInternalServerError, "All fields must be filled in")
		return
	}

	ourUser, err := uh.UserUseCase.GetUserByLogin(credentials.Login, credentials.Password)
	if err != nil {
		delivery.HandleError(w, http.StatusUnauthorized, "Login failed")
		return
	}

	_, er := uh.Sessions.Create(w, ourUser.ID)
	if er != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Login successful"})
}

// Signup handles user signup.
// @Summary User signup
// @Description Handles user signup.
// @Tags users
// @Accept json
// @Produce json
// @Param newUser body delivery.UserSwag true "New user details for signup"
// @Success 200 {object} delivery.Response "Signup successful"
// @Failure 400 {object} delivery.ErrorResponse "Invalid request body"
// @Failure 500 {object} delivery.ErrorResponse "Failed to add user"
// @Router /api/v1/signup [post]
func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if isEmpty(newUser.FirstName) || isEmpty(newUser.Surname) || isEmpty(newUser.Login) || isEmpty(newUser.Password) {
		delivery.HandleError(w, http.StatusBadRequest, "All fields must be filled in")
		return
	}

	loginUnique, _ := uh.UserUseCase.IsLoginUnique(newUser.Login)
	if !loginUnique {
		delivery.HandleError(w, http.StatusBadRequest, "Such a login already exists")
		return
	}

	_, er := uh.UserUseCase.CreateUser(converters.UserConvertApiInCore(newUser))
	if er != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Failed to add user")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Signup successful"})
}

// Logout handles user logout.
// @Summary User logout
// @Description Handles user logout.
// @Tags users
// @Produce json
// @Success 200 {object} delivery.Response "Logout successful"
// @Router /api/v1/logout [post]
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := uh.Sessions.DestroyCurrent(w, r)
	if err != nil {
		delivery.HandleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Logout successful"})
}

// GetUserBySession retrieves the user associated with the current session.
// @Summary Get user by session
// @Description Retrieve the user associated with the current session
// @Tags users
// @Produce json
// @Param X-CSRF-Token header string true "CSRF Token"
// @Success 200 {object} delivery.Response "User details"
// @Failure 401 {object} delivery.Response "Not Authorized"
// @Failure 500 {object} delivery.ErrorResponse "Internal Server Error"
// @Router /api/v1/auth/get-user [get]
func (uh *UserHandler) GetUserBySession(w http.ResponseWriter, r *http.Request) {
	sessionUser := uh.Sessions.GetSession(r)
	userData, err := uh.UserUseCase.GetUserByID(sessionUser.UserID)
	if err != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"user": converters.UserConvertCoreInApi(*userData)})
}

// isEmpty checks if the given string is empty after trimming leading and trailing whitespace.
// Returns true if the string is empty, and false otherwise.
func isEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}
