package domain_models

// Answer represents the information.
type Answer struct {
	ID         uint64
	QuestionId string
	LoginUser  string
	Mark       uint64
}
