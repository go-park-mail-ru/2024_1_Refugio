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

func (r *QuestionAnswerRepository) GetAllQuestions(ctx context.Context) ([]*domain.Email, error) {
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

	return emailsModelCore, nil
}

/*
func (r *QuestionAnswerRepository) Add(newQuestion *domain.Question, ctx context.Context) (uint64, *domain.Email, error) {
	insertEmailQuery := `
		INSERT INTO email (topic, text, date_of_dispatch, sender_email, recipient_email, isRead, isDeleted, isDraft, isSpam, reply_to_email_id, is_important)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	insertEmailFileQuery := `
		INSERT INTO email_file (email_id, file_id)
		SELECT $1, p.avatar_id
		FROM profile p
		WHERE p.login = $2
	`

	emailModelDb := converters.EmailConvertCoreInDb(*emailModelCore)
	format := "2006/01/02 15:04:05"

	var id uint64
	start := time.Now()

	err := r.DB.QueryRow(insertEmailQuery, emailModelDb.Topic, emailModelDb.Text, time.Now().Format(format), emailModelDb.SenderEmail, emailModelDb.RecipientEmail, emailModelDb.ReadStatus, emailModelDb.Deleted, emailModelDb.DraftStatus, emailModelCore.SpamStatus, emailModelDb.ReplyToEmailID, emailModelDb.Flag).Scan(&id)
	if err != nil {
		return 0, &domain.Email{}, fmt.Errorf("failed to add email: %v", err)
	}

	_, err = r.DB.Exec(insertEmailFileQuery, id, emailModelDb.SenderEmail)
	if err != nil {
		return 0, &domain.Email{}, fmt.Errorf("failed to add email file: %v", err)
	}

	args := []interface{}{emailModelDb.Topic, emailModelDb.Text, time.Now().Format(format), emailModelDb.SenderEmail, emailModelDb.RecipientEmail, emailModelDb.ReadStatus, emailModelDb.Deleted, emailModelDb.DraftStatus, emailModelDb.ReplyToEmailID, emailModelDb.Flag, emailModelDb.SenderEmail}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(insertEmailQuery, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	return id, emailModelCore, nil
}
*/
