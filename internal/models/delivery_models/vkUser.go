package delivery_models

import (
	domain "mail/internal/microservice/models/domain_models"
	"time"
)

// User represents information about a user.
type VKUser struct {
	ID        uint32            `json:"id,omitempty"`        // ID uniquely identifies the user.
	FirstName string            `json:"firstname,omitempty"` // FirstName stores the first name of the user.
	Surname   string            `json:"surname,omitempty"`   // Surname stores the last name of the user.
	Gender    domain.UserGender `json:"gender,omitempty"`    // Gender stores the gender of the user.
	Birthday  time.Time         `json:"birthday,omitempty"`  // Birthday stores the birthdate of the user.
	VKId      uint32            `json:"vkId"`                // ID user in VK
	Login     string            `json:"login,omitempty"`     // Login user
}
