package repository_models

type Answer struct {
	ID         uint32 `db:"id"`
	QuestionID uint32 `db:"question_id"`
	Login      string `db:"login"`
	Mark       uint64 `db:"mark"`
}
