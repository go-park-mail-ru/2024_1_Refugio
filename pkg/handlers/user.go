package handlers

import (
	"encoding/json"
	"net/http"
	//"mail/pkg/session"
	"mail/pkg/user"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	UserRepository user.UserRepository
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userRepo user.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepository: userRepo,
	}
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

	// Assume you have a function that checks the credentials against the UserRepository
	// and returns the user ID if successful.
	//userID, err := uh.UserRepository.VerifyCredentials(credentials.Login, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create a new session
	//sess, err := sessionsManager.Create(w, userID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Return success response or handle errors.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))

	// You can also include additional information in the response, such as the user ID or session ID.
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
	// Parse request body
	var newUser user.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Assume you have a function that adds the new user to the UserRepository
	// and returns the assigned user ID.
	// userID, err := uh.UserRepository.Add(&newUser)
	if err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}

	// Create a new session
	//sess, err := sessionsManager.Create(w, userID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Return success response or handle errors.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Signup successful"))

	// You can also include additional information in the response, such as the user ID or session ID.
}

// Logout handles user logout.
// @Summary User logout
// @Description Handles user logout.
// @Tags users
// @Produce json
// @Success 200 {string} string "Logout successful"
// @Router /logout [post]
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	/*err := sessionsManager.DestroyCurrent(w, r)
	if err != nil {
		http.Error(w, "Failed to destroy session", http.StatusInternalServerError)
		return
	}*/

	// Return success response or handle errors.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}
