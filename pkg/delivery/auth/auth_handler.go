package auth

import (
	"encoding/json"
	"github.com/microcosm-cc/bluemonday"
	"mail/pkg/delivery"
	"mail/pkg/delivery/converters"
	api "mail/pkg/delivery/models"
	domainSession "mail/pkg/domain/delivery"
	domain "mail/pkg/domain/models"
	"mail/pkg/domain/usecase"
	"net/http"
	"regexp"
	"strings"
)

var (
	AHandler                        = &AuthHandler{}
	requestIDContextKey interface{} = "requestid"
)

// AuthHandler handles user-related HTTP requests.
type AuthHandler struct {
	UserUseCase usecase.UserUseCase
	Sessions    domainSession.SessionsManager
}

// InitializationAuthHandler initializes the auth handler with the provided user handler.
func InitializationAuthHandler(authHandler *AuthHandler) {
	AHandler = authHandler
}

// sanitizeString sanitizes the provided string using the UGCPolicy from the bluemonday package.
func sanitizeString(str string) string {
	p := bluemonday.UGCPolicy()

	return p.Sanitize(str)
}

// VerifyAuth verifies user authentication.
// @Summary Verify user authentication
// @Description Verify user authentication using sessions
// @Tags auth
// @Produce json
// @Success 200 {object} delivery.Response "OK"
// @Failure 401 {object} delivery.Response "Not Authorized"
// @Router /api/v1/verify-auth [get]
func (ah *AuthHandler) VerifyAuth(w http.ResponseWriter, r *http.Request) {
	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	sessionUser := ah.Sessions.GetSession(r, requestID)
	w.Header().Set("X-Csrf-Token", sessionUser.CsrfToken)

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "OK"})
}

// Login handles user login.
// @Summary User login
// @Description Handles user login.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body delivery.UserSwag true "User credentials for login"
// @Success 200 {object} delivery.Response "Login successful"
// @Failure 400 {object} delivery.ErrorResponse "Invalid request body"
// @Failure 401 {object} delivery.ErrorResponse "Invalid credentials"
// @Failure 500 {object} delivery.ErrorResponse "Failed to create session"
// @Router /api/v1/auth/login [post]
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials api.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	credentials.Login = sanitizeString(credentials.Login)
	credentials.Password = sanitizeString(credentials.Password)

	if isEmpty(credentials.Login) || isEmpty(credentials.Password) {
		delivery.HandleError(w, http.StatusInternalServerError, "All fields must be filled in")
		return
	}

	if !isValidEmailFormat(credentials.Login) {
		delivery.HandleError(w, http.StatusBadRequest, "Domain in the login is not suitable")
		return
	}

	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	ourUser, err := ah.UserUseCase.GetUserByLogin(credentials.Login, credentials.Password, requestID)
	if err != nil {
		delivery.HandleError(w, http.StatusUnauthorized, "Login failed")
		return
	}

	_, er := ah.Sessions.Create(w, ourUser.ID, requestID)
	if er != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Login successful"})
}

// Signup handles user signup.
// @Summary User signup
// @Description Handles user signup.
// @Tags auth
// @Accept json
// @Produce json
// @Param newUser body delivery.UserSwag true "New user details for signup"
// @Success 200 {object} delivery.Response "Signup successful"
// @Failure 400 {object} delivery.ErrorResponse "Invalid request body"
// @Failure 500 {object} delivery.ErrorResponse "Failed to add user"
// @Router /api/v1/auth/signup [post]
func (ah *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var newUser api.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newUser.Login = sanitizeString(newUser.Login)
	newUser.Password = sanitizeString(newUser.Password)
	newUser.FirstName = sanitizeString(newUser.FirstName)
	newUser.Surname = sanitizeString(newUser.Surname)

	if isEmpty(newUser.Login) || isEmpty(newUser.Password) || isEmpty(newUser.FirstName) || isEmpty(newUser.Surname) || !domain.IsValidGender(newUser.Gender) {
		delivery.HandleError(w, http.StatusBadRequest, "All fields must be filled in")
		return
	}

	if !isValidEmailFormat(newUser.Login) {
		delivery.HandleError(w, http.StatusBadRequest, "Domain in the login is not suitable")
		return
	}

	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	loginUnique, _ := ah.UserUseCase.IsLoginUnique(newUser.Login, requestID)
	if !loginUnique {
		delivery.HandleError(w, http.StatusBadRequest, "Such a login already exists")
		return
	}

	_, er := ah.UserUseCase.CreateUser(converters.UserConvertApiInCore(newUser), requestID)
	if er != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Failed to add user")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Signup successful"})
}

// Logout handles user logout.
// @Summary User logout
// @Description Handles user logout.
// @Tags auth
// @Produce json
// @Success 200 {object} delivery.Response "Logout successful"
// @Router /api/v1/auth/logout [post]
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	err := ah.Sessions.DestroyCurrent(w, r, requestID)
	if err != nil {
		delivery.HandleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Logout successful"})
}

// isEmpty checks if the given string is empty after trimming leading and trailing whitespace.
// Returns true if the string is empty, and false otherwise.
func isEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// isValidEmailFormat checks if the provided email string matches the specific format for emails ending with "@mailhub.ru".
func isValidEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@mailhub\.su$`)

	return emailRegex.MatchString(email)
}
