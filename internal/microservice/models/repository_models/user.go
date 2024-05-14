package repository_models

import (
	"time"

	domain "mail/internal/microservice/models/domain_models"
)

// User represents information about a user.
type User struct {
	ID               uint32            `db:"id"`                // ID uniquely identifies the user.
	FirstName        string            `db:"firstname"`         // FirstName stores the first name of the user.
	Surname          string            `db:"surname"`           // Surname stores the last name of the user.
	Patronymic       string            `db:"patronymic"`        // Patronymic stores the middle name of the user, if available.
	Gender           domain.UserGender `db:"gender"`            // Gender stores the gender of the user.
	Birthday         time.Time         `db:"birthday"`          // Birthday stores the birthdate of the user.
	RegistrationDate time.Time         `db:"registration_date"` // RegistrationDate stores the date when the user registered.
	Login            string            `db:"login"`             // Login is the username used for authentication.
	Password         string            `db:"password_hash"`     // Password is the hashed password of the user.
	AvatarID         *string           `db:"avatar_id"`         // AvatarID stores the identifier of the user's avatar image.
	PhoneNumber      string            `db:"phone_number"`      // PhoneNumber stores the phone number of the user.
	Description      string            `db:"description"`       // Description stores additional information about the user.
	VKId             uint32            `db:"vkid"`              // ID uniquely identifies the user in VK.
}
