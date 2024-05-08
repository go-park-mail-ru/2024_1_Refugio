package domain_models

// Question represents information about a question.
type Question struct {
	ID          uint32 // ID is the unique identifier of the question.
	Text        string // Text contains the text of the question, describing the inquiry or topic.
	MinResult   string // MinResult holds the minimum expected result or outcome related to the question.
	MaxResult   string // MaxResult holds the maximum expected result or outcome related to the question.
	DopQuestion string // DopQuestion contains additional or supplementary questions that may be related to the main question.
}
