package delivery_models

type Answer struct {
	ID         uint32 `json:"id,omitempty"`
	QuestionId uint32 `json:"question_id,omitempty"`
	Login      string `json:"login,omitempty"`
	Mark       uint32 `json:"mark,omitempty"`
	Text       string `json:"text,omitempty"`
}
