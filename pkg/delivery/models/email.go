package models

import "time"

// Email represents the information about an email.
type Email struct {
	ID             uint64    `json:"id,omitempty"`             // ID is the unique identifier of the email in the database.
	Topic          string    `json:"topic"`                    // Topic is the subject of the email.
	Text           string    `json:"text"`                     // Text is the body of the email.
	AvatarID       string    `json:"avatarID"`                 // AvatarID avatar interlocutor.
	ReadStatus     bool      `json:"readStatus"`               // ReadStatus indicates whether the email has been read.
	Flag           bool      `json:"mark,omitempty"`           // Flag is a flag, such as marking the email as a favorite.
	Deleted        bool      `json:"deleted"`                  // Deleted indicates whether the email has been deleted.
	DateOfDispatch time.Time `json:"dateOfDispatch,omitempty"` // DateOfDispatch is the date when the email was sent.
	ReplyToEmailID uint64    `json:"replyToEmailId,omitempty"` // ReplyToEmailID is the ID of the email to which a reply can be sent.
	DraftStatus    bool      `json:"draftStatus"`              // DraftStatus indicates whether the email is a draft.
	SenderEmail    string    `json:"senderEmail"`              // SenderEmail is the email of the sender user
	RecipientEmail string    `json:"recipientEmail"`           // RecipientEmail is the email of the recipient user
}

//SenderPhotoID    string    `json:"senderPhotoID"`            // SenderPhotoID is the link to the photo attached to the email, if any.
//RecipientPhotoID string    `json:"recipientPhotoID"`         // RecipientPhotoID is the link to the photo attached to the email, if any.	ReadStatus     bool      `json:"readStatus"`               // ReadStatus indicates whether the email has been read.
