package main

type EmailSMTP struct {
	Topic          string `json:"topic"`
	Text           string `json:"text"`
	SenderEmail    string `json:"senderEmail"`
	RecipientEmail string `json:"recipientEmail"`
}
