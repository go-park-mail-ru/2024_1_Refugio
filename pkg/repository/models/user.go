package models

// User represents information about a user.
type User struct {
	ID       uint32 `json:"id,omitempty"`      // ID uniquely identifies the user.
	Name     string `json:"name,omitempty"`    // Name stores the first name of the user.
	Surname  string `json:"surname,omitempty"` // Surname stores the last name of the user.
	Login    string `json:"login"`             // Login is the username used for authentication.
	Password string `json:"password"`          // Password is the hashed password of the user.
	AvatarID string `json:"avatar,omitempty"`  // AvatarID stores the identifier of the user's avatar image.
}
