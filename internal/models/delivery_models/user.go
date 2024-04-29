package delivery_models

import (
	domain "mail/internal/microservice/models/domain_models"
	"time"
)

// User represents information about a user.
type User struct {
	ID          uint32            `json:"id,omitempty"`          // ID uniquely identifies the user.
	FirstName   string            `json:"firstname,omitempty"`   // FirstName stores the first name of the user.
	Surname     string            `json:"surname,omitempty"`     // Surname stores the last name of the user.
	Patronymic  string            `json:"middlename,omitempty"`  // Patronymic stores the middle name of the user, if available.
	Gender      domain.UserGender `json:"gender,omitempty"`      // Gender stores the gender of the user.
	Birthday    time.Time         `json:"birthday,omitempty"`    // Birthday stores the birthdate of the user.
	Login       string            `json:"login"`                 // Login is the username used for authentication.
	Password    string            `json:"password"`              // Password is the hashed password of the user.
	AvatarID    string            `json:"avatar,omitempty"`      // AvatarID stores the identifier of the user's avatar image.
	PhoneNumber string            `json:"phonenumber,omitempty"` // PhoneNumber stores the phone number of the user.
	Description string            `json:"description,omitempty"` // Description stores additional information about the user.
}
