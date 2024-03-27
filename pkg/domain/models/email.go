package models

import "time"

// Email represents the information about an email.
type Email struct {
	ID             uint64    // ID is the unique identifier of the email in the database.
	Topic          string    // Topic is the subject of the email.
	Text           string    // Text is the body of the email.
	PhotoID        string    // PhotoID is the link to the photo attached to the email, if any.
	ReadStatus     bool      // ReadStatus indicates whether the email has been read.
	Flag           bool      // Mark is a flag, such as marking the email as a favorite.
	Deleted        bool      // Deleted indicates whether the email has been deleted.
	DateOfDispatch time.Time // DateOfDispatch is the date when the email was sent.
	ReplyToEmailID uint64    // ReplyToEmailID is the ID of the email to which a reply can be sent.
	DraftStatus    bool      // DraftStatus indicates whether the email is a draft.
	SenderID       uint64    // SenderID is the ID of the sender user
	RecipientID    uint64    // RecipientID is the ID of the recipient user
}
