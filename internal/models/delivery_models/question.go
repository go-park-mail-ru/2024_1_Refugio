package delivery_models

type Question struct {
	ID      uint32 `json:"id,omitempty"`
	Text    string `json:"text,omitempty"`
	MinText string `json:"min_text,omitempty"`
	MaxText string `json:"max_text,omitempty"`
}
