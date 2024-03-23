package models

// User represents information about a user.
type User struct {
	ID       uint32 // ID uniquely identifies the user.
	Name     string // Name stores the first name of the user.
	Surname  string // Surname stores the last name of the user.
	Login    string // Login is the username used for authentication.
	Password string // Password is the hashed password of the user.
	AvatarID string // AvatarID stores the identifier of the user's avatar image.
}
