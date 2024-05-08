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

type UserGenderSwag string

const (
	Male   UserGenderSwag = "Male"
	Female UserGenderSwag = "Female"
	Other  UserGenderSwag = "Other"
)

type UserSwag struct {
	ID          uint32         `json:"id,omitempty"`
	FirstName   string         `json:"firstname,omitempty"`
	Surname     string         `json:"surname,omitempty"`
	Patronymic  string         `json:"middlename,omitempty"`
	Gender      UserGenderSwag `json:"gender,omitempty"`
	Birthday    time.Time      `json:"birthday,omitempty"`
	Login       string         `json:"login"`
	Password    string         `json:"password"`
	AvatarID    string         `json:"avatar,omitempty"`
	PhoneNumber string         `json:"phonenumber,omitempty"`
	Description string         `json:"description,omitempty"`
}

type EmailSwag struct {
	ID             uint64    `json:"id,omitempty"`
	Topic          string    `json:"topic"`
	Text           string    `json:"text"`
	ReadStatus     bool      `json:"readStatus"`
	Flag           bool      `json:"mark,omitempty"`
	Deleted        bool      `json:"deleted"`
	DateOfDispatch time.Time `json:"dateOfDispatch,omitempty"`
	ReplyToEmailID uint64    `json:"replyToEmailId,omitempty"`
	DraftStatus    bool      `json:"draftStatus"`
	SpamStatus     bool      `json:"spamStatus"`
	SenderEmail    string    `json:"senderEmail"`
	RecipientEmail string    `json:"recipientEmail"`
}

type FolderSwag struct {
	Name string `json:"name"`
}

type FolderEmailSwag struct {
	FolderID uint32 `json:"folderId"`
	EmailID  uint32 `json:"emailId"`
}

type QuestionSwag struct {
	ID          uint32 `json:"id,omitempty"`
	Text        string `json:"text,omitempty"`
	MinText     string `json:"min_text,omitempty"`
	MaxText     string `json:"max_text,omitempty"`
	DopQuestion string `json:"dop_question,omitempty"`
}

type AnswerSwag struct {
	ID         uint32 `json:"id,omitempty"`
	QuestionId uint32 `json:"question_id,omitempty"`
	Login      string `json:"login,omitempty"`
	Mark       uint32 `json:"mark,omitempty"`
	Text       string `json:"text,omitempty"`
}
