package main

// EmailSMTP represents the information about an email.
type EmailSMTP struct {
	Topic          string `json:"topic"`          // Topic is the subject of the email.
	Text           string `json:"text"`           // Text is the body of the email.
	SenderEmail    string `json:"senderEmail"`    // SenderEmail is the email of the sender user
	RecipientEmail string `json:"recipientEmail"` // RecipientEmail is the email of the recipient user
}
