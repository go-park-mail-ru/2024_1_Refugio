package delivery_models

// Question represents the information.
type Question struct {
	ID        uint64 `json:"id,omitempty"`
	Question  string `json:"question"`
	MinResult string `json:"minResult"`
	MaxResult string `json:"maxResult"`
}
