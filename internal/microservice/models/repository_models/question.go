package repository_models

type Question struct {
	ID        uint64 `db:"id"`
	Text      string `db:"text"`
	MinResult string `db:"min_text"`
	MaxResult string `db:"max_text"`
}
