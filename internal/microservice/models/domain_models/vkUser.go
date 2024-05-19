package domain_models

import "time"

// VKUser represents information about a user.
type VKUser struct {
	ID        uint32     // ID uniquely identifies the user.
	FirstName string     // FirstName stores the first name of the user.
	Surname   string     // Surname stores the last name of the user.
	Gender    UserGender // Gender stores the gender of the user.
	Birthday  time.Time  // Birthday stores the birthdate of the user.
	VKId      uint32     // ID user in VK
	Login     string     // Login user
}
