package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	domain "mail/internal/microservice/models/domain_models"
	converters "mail/internal/microservice/models/repository_converters"
	"mail/internal/microservice/models/repository_models"
	"mail/internal/pkg/logger"
	"time"
)

var requestIDContextKey interface{} = "requestID"

type QuestionAnswerRepository struct {
	DB *sqlx.DB
}

func NewQuestionRepository(db *sqlx.DB) *QuestionAnswerRepository {
	return &QuestionAnswerRepository{DB: db}
}

func (r *QuestionAnswerRepository) GetAllQuestions(ctx context.Context) ([]*domain.Question, error) {
	query := `
		SELECT question.id, question.text, question.min_text, question.max_text FROM question
	`

	questionsModelDb := []repository_models.Question{}

	var err error
	args := []interface{}{}
	start := time.Now()
	err = r.DB.Select(&questionsModelDb, query)

	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have questions")
		}
		return nil, err
	}

	var questionsModelCore []*domain.Question
	for _, e := range questionsModelDb {
		questionsModelCore = append(questionsModelCore, converters.QuestionConvertDbInCore(e))
	}

	return questionsModelCore, nil
}

func (r *QuestionAnswerRepository) GetAllAnswers(ctx context.Context) ([]*domain.Answer, error) {
	query := `
		SELECT answer.id, answer.question_id, answer.login, answer.mark, question.text 
		FROM answer
		LEFT JOIN question ON question.id = answer.question_id
	`

	answersModelDb := []repository_models.Answer{}

	var err error
	args := []interface{}{}
	start := time.Now()
	err = r.DB.Select(&answersModelDb, query)

	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have answers")
		}
		return nil, err
	}

	var answersModelCore []*domain.Answer
	for _, a := range answersModelDb {
		answersModelCore = append(answersModelCore, converters.AnswerConvertDbInCore(a))
	}

	return answersModelCore, nil
}

func (r *QuestionAnswerRepository) AddQuestion(newQuestion *domain.Question, ctx context.Context) (bool, error) {
	insertQuestionQuery := `
		INSERT INTO question (text, min_text, max_text)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	questionModelDb := converters.QuestionConvertCoreInDb(*newQuestion)

	var id uint64
	start := time.Now()
	err := r.DB.QueryRow(insertQuestionQuery, questionModelDb.Text, questionModelDb.MinResult, questionModelDb.MaxResult).Scan(&id)
	if err != nil {
		return false, fmt.Errorf("failed to add question: %v", err)
	}

	args := []interface{}{questionModelDb.Text, questionModelDb.MinResult, questionModelDb.MaxResult}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(insertQuestionQuery, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	return true, nil
}

func (r *QuestionAnswerRepository) AddAnswer(newAnswer *domain.Answer, ctx context.Context) (bool, error) {
	insertAnswerQuery := `
		INSERT INTO answer (question_id, login, mark)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	answerModelDb := converters.AnswerConvertCoreInDb(*newAnswer)

	var id uint64
	start := time.Now()
	err := r.DB.QueryRow(insertAnswerQuery, answerModelDb.QuestionID, answerModelDb.Login, answerModelDb.Mark).Scan(&id)
	if err != nil {
		return false, fmt.Errorf("failed to add answer: %v", err)
	}

	args := []interface{}{answerModelDb.QuestionID, answerModelDb.Login, answerModelDb.Mark}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(insertAnswerQuery, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	return true, nil
}
