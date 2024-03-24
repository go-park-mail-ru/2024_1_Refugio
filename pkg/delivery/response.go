package delivery

import (
	"encoding/json"
	"mail/pkg/domain/models"
	"net/http"
	"time"

	api "mail/pkg/delivery/models"
)

// Response represents the response format.
type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

// ErrorResponse represents the error response format.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HandleSuccess is a utility function to handle successful responses uniformly in the API.
func HandleSuccess(w http.ResponseWriter, status int, body interface{}) {
	response := Response{
		Status: status,
		Body:   body,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// HandleError is a utility function to handle errors uniformly in the API responses.
func HandleError(w http.ResponseWriter, status int, message string) {
	response := Response{
		Status: status,
		Body:   ErrorResponse{Error: message},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

type Email = api.Email

// UserSwag represents information about a user.
type UserSwag struct {
	ID          uint32            `json:"id,omitempty"`          // ID uniquely identifies the user.
	FirstName   string            `json:"firstname,omitempty"`   // FirstName stores the first name of the user.
	Surname     string            `json:"surname,omitempty"`     // Surname stores the last name of the user.
	Patronymic  string            `json:"middlename,omitempty"`  // Patronymic stores the middle name of the user, if available.
	Gender      models.UserGender `json:"gender,omitempty"`      // Gender stores the gender of the user.
	Birthday    time.Time         `json:"birthday,omitempty"`    // Birthday stores the birthdate of the user.
	Login       string            `json:"login"`                 // Login is the username used for authentication.
	Password    string            `json:"password"`              // Password is the hashed password of the user.
	AvatarID    string            `json:"avatar,omitempty"`      // AvatarID stores the identifier of the user's avatar image.
	PhoneNumber string            `json:"phonenumber,omitempty"` // PhoneNumber stores the phone number of the user.
	Description string            `json:"description,omitempty"` // Description stores additional information about the user.
}
