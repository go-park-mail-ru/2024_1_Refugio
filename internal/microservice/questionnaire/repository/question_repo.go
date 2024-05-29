package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"mail/internal/microservice/models/repository_models"
	"mail/internal/pkg/logger"

	domain "mail/internal/microservice/models/domain_models"
	converters "mail/internal/microservice/models/repository_converters"
)

// requestIDContextKey is the context key for the request ID.
var requestIDContextKey interface{} = "requestID"

// QuestionAnswerRepository handles operations related to questions and answers in the database.
type QuestionAnswerRepository struct {
	DB *sqlx.DB
}

// NewQuestionRepository creates a new instance of QuestionAnswerRepository.
func NewQuestionRepository(db *sqlx.DB) *QuestionAnswerRepository {
	return &QuestionAnswerRepository{DB: db}
}

// GetAllQuestions retrieves all questions from the database.
func (r *QuestionAnswerRepository) GetAllQuestions(ctx context.Context) ([]*domain.Question, error) {
	query := `
		SELECT question.id, question.text, question.min_text, question.max_text, question.dop_question 
		FROM question
	`

	var err error
	questionsModelDb := []repository_models.Question{}

	start := time.Now()

	err = r.DB.Select(&questionsModelDb, query)

	args := []interface{}{}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have questions")
		}
		return nil, err
	}

	var questionsModelCore []*domain.Question
	for _, e := range questionsModelDb {
		questionsModelCore = append(questionsModelCore, converters.QuestionConvertDbInCore(&e))
	}

	return questionsModelCore, nil
}

// GetAllAnswers retrieves all answers from the database.
func (r *QuestionAnswerRepository) GetAllAnswers(ctx context.Context) ([]*domain.Answer, error) {
	query := `
		SELECT answer.id, answer.question_id, answer.login, answer.mark, question.text 
		FROM answer
		LEFT JOIN question ON question.id = answer.question_id
	`

	var err error
	answersModelDb := []repository_models.Answer{}

	start := time.Now()

	err = r.DB.Select(&answersModelDb, query)

	args := []interface{}{}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have answers")
		}
		return nil, err
	}

	var answersModelCore []*domain.Answer
	for _, a := range answersModelDb {
		answersModelCore = append(answersModelCore, converters.AnswerConvertDbInCore(&a))
	}

	return answersModelCore, nil
}

// AddQuestion adds a new question to the database.
func (r *QuestionAnswerRepository) AddQuestion(newQuestion *domain.Question, ctx context.Context) (bool, error) {
	insertQuestionQuery := `
		INSERT INTO question (text, min_text, max_text, dop_question)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id uint64
	questionModelDb := converters.QuestionConvertCoreInDb(newQuestion)

	start := time.Now()

	err := r.DB.QueryRow(insertQuestionQuery, questionModelDb.Text, questionModelDb.MinResult, questionModelDb.MaxResult, questionModelDb.DopQuestion).Scan(&id)

	args := []interface{}{questionModelDb.Text, questionModelDb.MinResult, questionModelDb.MaxResult, questionModelDb.DopQuestion}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(insertQuestionQuery, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return false, fmt.Errorf("failed to add question: %v", err)
	}

	return true, nil
}

// AddAnswer adds a new answer to the database.
func (r *QuestionAnswerRepository) AddAnswer(newAnswer *domain.Answer, ctx context.Context) (bool, error) {
	insertAnswerQuery := `
		INSERT INTO answer (question_id, login, mark, text)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id uint64
	answerModelDb := converters.AnswerConvertCoreInDb(newAnswer)

	start := time.Now()

	err := r.DB.QueryRow(insertAnswerQuery, answerModelDb.QuestionID, answerModelDb.Login, answerModelDb.Mark, answerModelDb.Text).Scan(&id)

	args := []interface{}{answerModelDb.QuestionID, answerModelDb.Login, answerModelDb.Mark, answerModelDb.Text}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(insertAnswerQuery, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return false, fmt.Errorf("failed to add answer: %v", err)
	}

	return true, nil
}
