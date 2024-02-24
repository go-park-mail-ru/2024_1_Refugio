package email

import "time"

// Email represents the information about an email.
type Email struct {
	ID             uint64    `json:"id,omitempty"`             // Unique identifier of the email in the database.
	Topic          string    `json:"topic"`                    // Subject of the email.
	Text           string    `json:"text"`                     // Text body of the email.
	PhotoID        string    `json:"photoId,omitempty"`        // Link to the photo attached to the email, if any.
	ReadStatus     bool      `json:"readStatus"`               // Status indicating whether the email has been read.
	Mark           string    `json:"mark,omitempty"`           // A flag, for example, marking the email as a favorite.
	Deleted        bool      `json:"deleted"`                  // Status indicating whether the email has been deleted.
	DateOfDispatch time.Time `json:"dateOfDispatch,omitempty"` // Date when the email was sent.
	ReplyToEmailID uint64    `json:"replyToEmailId,omitempty"` // ID of the email to which a reply can be sent.
	DraftStatus    bool      `json:"draftStatus"`              // Status indicating that the email is a draft.
}

// EmailRepository represents the interface for working with emails.
type EmailRepository interface {
	// GetAll returns all emails from the storage.
	GetAll() ([]*Email, error)

	// GetByID returns the email by its unique identifier.
	GetByID(id uint64) (*Email, error)

	// Add adds a new email to the storage and returns its assigned unique identifier.
	Add(email *Email) (uint64, error)

	// Update updates the information of an email in the storage based on the provided new email.
	Update(newEmail *Email) (bool, error)

	// Delete removes the email from the storage by its unique identifier.
	Delete(id uint64) (bool, error)
}
