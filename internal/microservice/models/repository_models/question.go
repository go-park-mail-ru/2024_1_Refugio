package repository_models

type Question struct {
	ID          uint32 `db:"id"`
	Text        string `db:"text"`
	MinResult   string `db:"min_text"`
	MaxResult   string `db:"max_text"`
	DopQuestion string `db:"dop_question"`
}
