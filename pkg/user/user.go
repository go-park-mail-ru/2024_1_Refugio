package user

// User represents information about a user.
type User struct {
	ID       uint32 `json:"id,omitempty"`      // Unique identifier of the user.
	Name     string `json:"name,omitempty"`    // User's first name.
	Surname  string `json:"surname,omitempty"` // User's last name.
	Login    string `json:"login"`             // User's login.
	Password string `json:"password"`          // User's password.
	AvatarId string `json:"avatar,omitempty"`  // User's avatar.
}

// UserRepository represents the interface for working with users.
type UserRepository interface {
	// GetAll returns all users from the storage.
	GetAll() ([]*User, error)

	// GetByID returns the user by its unique identifier.
	GetByID(id uint32) (*User, error)

	// Add adds a new user to the storage and returns its assigned unique identifier.
	Add(user *User) (uint32, error)

	// Update updates the information of a user in the storage based on the provided new user.
	Update(newUser *User) (bool, error)

	// Delete removes the user from the storage by its unique identifier.
	Delete(id uint32) (bool, error)
}
