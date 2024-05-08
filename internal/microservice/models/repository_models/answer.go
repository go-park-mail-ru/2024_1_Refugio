package repository_models

// Answer represents information about an answer.
type Answer struct {
	ID         uint32 `db:"id"`          // ID is the unique identifier of the answer.
	QuestionID uint32 `db:"question_id"` // QuestionID holds the identifier of the question to which this answer belongs.
	Login      string `db:"login"`       // Login represents the username of the user who provided the answer.
	Mark       uint32 `db:"mark"`        // Mark denotes the rating given to the answer.
	Text       string `db:"text"`        // Text contains the text of the answer.
}
