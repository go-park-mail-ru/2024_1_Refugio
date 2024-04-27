package domain_models

// Answer represents the information.
type Answer struct {
	ID         uint32
	QuestionId uint32
	LoginUser  string
	Mark       uint32
}
