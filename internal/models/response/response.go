package response

import (
	"encoding/json"
	"net/http"
	"time"
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

// UserSwag represents information about a user.
type UserSwag struct {
	ID          uint32         `json:"id,omitempty"`          // ID uniquely identifies the user.
	FirstName   string         `json:"firstname,omitempty"`   // FirstName stores the first name of the user.
	Surname     string         `json:"surname,omitempty"`     // Surname stores the last name of the user.
	Patronymic  string         `json:"middlename,omitempty"`  // Patronymic stores the middle name of the user, if available.
	Gender      UserGenderSwag `json:"gender,omitempty"`      // Gender stores the gender of the user.
	Birthday    time.Time      `json:"birthday,omitempty"`    // Birthday stores the birthdate of the user.
	Login       string         `json:"login"`                 // Login is the username used for authentication.
	Password    string         `json:"password"`              // Password is the hashed password of the user.
	AvatarID    string         `json:"avatar,omitempty"`      // AvatarID stores the identifier of the user's avatar image.
	PhoneNumber string         `json:"phonenumber,omitempty"` // PhoneNumber stores the phone number of the user.
	Description string         `json:"description,omitempty"` // Description stores additional information about the user.
}

// EmailSwag represents information about a email.
type EmailSwag struct {
	ID             uint64    `json:"id,omitempty"`             // ID is the unique identifier of the email in the database.
	Topic          string    `json:"topic"`                    // Topic is the subject of the email.
	Text           string    `json:"text"`                     // Text is the body of the email.
	PhotoID        string    `json:"photoId,omitempty"`        // PhotoID is the link to the photo attached to the email, if any.
	ReadStatus     bool      `json:"readStatus"`               // ReadStatus indicates whether the email has been read.
	Flag           bool      `json:"mark,omitempty"`           // Mark is a flag, such as marking the email as a favorite.
	Deleted        bool      `json:"deleted"`                  // Deleted indicates whether the email has been deleted.
	DateOfDispatch time.Time `json:"dateOfDispatch,omitempty"` // DateOfDispatch is the date when the email was sent.
	ReplyToEmailID uint64    `json:"replyToEmailId,omitempty"` // ReplyToEmailID is the ID of the email to which a reply can be sent.
	DraftStatus    bool      `json:"draftStatus"`              // DraftStatus indicates whether the email is a draft.
	SpamStatus     bool      `json:"spamStatus"`               // SpamStatus indicates whether the email is a spam
	SenderEmail    string    `json:"senderEmail"`              // SenderEmail is the Email of the sender user
	RecipientEmail string    `json:"recipientEmail"`           // RecipientEmail is the Email of the recipient user
}

type FolderSwag struct {
	ID   uint64 `json:"id,omitempty"` // ID he unique ID of the folder in the database.
	Name string `json:"name"`         // Name the name of the folder.
}

type UserGenderSwag string

const (
	Male   UserGenderSwag = "Male"
	Female UserGenderSwag = "Female"
	Other  UserGenderSwag = "Other"
)
