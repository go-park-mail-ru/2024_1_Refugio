package domain_models

// Answer represents the information.
type Answer struct {
	ID         uint32
	QuestionID uint32
	Login      string
	Mark       uint64
}
