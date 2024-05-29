package domain_models

import "time"

// User represents information about a user.
type User struct {
	ID          uint32     // ID uniquely identifies the user.
	FirstName   string     // FirstName stores the first name of the user.
	Surname     string     // Surname stores the last name of the user.
	Patronymic  string     // Patronymic stores the middle name of the user, if available.
	Gender      UserGender // Gender stores the gender of the user.
	Birthday    time.Time  // Birthday stores the birthdate of the user.
	Login       string     // Login is the username used for authentication.
	Password    string     // Password is the hashed password of the user.
	AvatarID    string     // AvatarID stores the identifier of the user's avatar image.
	PhoneNumber string     // PhoneNumber stores the phone number of the user.
	Description string     // Description stores additional information about the user.
	VKId        uint32     // VKId uniquely identifies the user in VK.
}
