package delivery_models

// Question represents a structure to hold information about a question.
type Question struct {
	ID          uint32 `json:"id,omitempty"`           // ID is the unique identifier of the question.
	Text        string `json:"text,omitempty"`         // Text contains the main text of the question.
	MinText     string `json:"min_text,omitempty"`     // MinText is the minimum allowable text for the question.
	MaxText     string `json:"max_text,omitempty"`     // MaxText is the maximum allowable text for the question.
	DopQuestion string `json:"dop_question,omitempty"` // DopQuestion contains additional question information.
}
