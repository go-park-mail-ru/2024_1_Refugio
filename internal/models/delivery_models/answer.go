package delivery_models

// Answer response to questionnaire
type Answer struct {
	ID         uint32 `json:"id,omitempty"`          // ID Unique identifier for the answer
	QuestionId uint32 `json:"question_id,omitempty"` // QuestionId Identifier for the question to which the answer belongs
	Login      string `json:"login,omitempty"`       // Login associated with the answer
	Mark       uint32 `json:"mark,omitempty"`        // Mark or score assigned to the answer
	Text       string `json:"text,omitempty"`        // Text content of the answer
}
