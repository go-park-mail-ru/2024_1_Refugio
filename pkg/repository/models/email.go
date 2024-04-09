package models

import "time"

// Email represents the information about an email.
type Email struct {
	ID             uint64      `db:"id"`                   // ID is the unique identifier of the email in the database.
	Topic          string      `db:"topic"`                // Topic is the subject of the email.
	Text           string      `db:"text"`                 // Text is the body of the email.
	PhotoID        string      `db:"photoid"`              // PhotoID is the link to the photo attached to the email, if any.
	ReadStatus     bool        `db:"isRead"`               // ReadStatus indicates whether the email has been read.
	Flag           bool        `db:"flag_important_email"` // Flag is a flag, such as marking the email as a favorite.
	Deleted        bool        `db:"isDeleted"`            // Deleted indicates whether the email has been deleted.
	DateOfDispatch time.Time   `db:"date_of_dispatch"`     // DateOfDispatch is the date when the email was sent.
	ReplyToEmailID interface{} `db:"reply_to_email_id"`    // ReplyToEmailID is the ID of the email to which a reply can be sent.
	DraftStatus    bool        `db:"isDraft"`              // DraftStatus indicates whether the email is a draft.
	SenderEmail    string      `db:"sender_email"`         // SenderEmail is the email of the sender user
	RecipientEmail string      `db:"recipient_email"`      // RecipientEmail is the email of the recipient user
}
