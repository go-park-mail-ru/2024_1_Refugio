package domain_models

// Answer represents information about an answer.
type Answer struct {
	ID         uint32 // ID is the unique identifier of the answer.
	QuestionID uint32 // QuestionID holds the identifier of the question to which this answer belongs.
	Login      string // Login represents the username of the user who provided the answer.
	Mark       uint32 // Mark denotes the rating given to the answer.
	Text       string // Text contains the text of the answer.
}
