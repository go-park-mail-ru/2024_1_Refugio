package repository_models

// Question represents information about a question.
type Question struct {
	ID          uint32 `db:"id"`           // ID is the unique identifier of the question.
	Text        string `db:"text"`         // Text contains the text of the question, describing the inquiry or topic.
	MinResult   string `db:"min_text"`     // MinResult holds the minimum expected result or outcome related to the question.
	MaxResult   string `db:"max_text"`     // MaxResult holds the maximum expected result or outcome related to the question.
	DopQuestion string `db:"dop_question"` // DopQuestion contains additional or supplementary questions that may be related to the main question.
}
