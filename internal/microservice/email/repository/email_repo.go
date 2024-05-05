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

type EmailRepository struct {
	DB *sqlx.DB
}

func NewEmailRepository(db *sqlx.DB) *EmailRepository {
	/*f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile in email_repo" + "log.txt")
	}
	Logger = logger.InitializationBdLog(f)*/
	return &EmailRepository{DB: db}
}

func (r *EmailRepository) Add(emailModelCore *domain.Email, ctx context.Context) (uint64, *domain.Email, error) {
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

func (r *EmailRepository) AddProfileEmail(email_id uint64, sender, recipient string, ctx context.Context) error {
	var query string
	var err error
	var args []interface{}
	start := time.Now()
	if sender == recipient {
		query = `
			INSERT INTO profile_email (profile_id, email_id)
			VALUES ((SELECT id FROM profile WHERE login=$1), $2)
		`
		_, err = r.DB.Exec(query, sender, email_id)
		args = []interface{}{sender, email_id}
	} else {
		query = `
			INSERT INTO profile_email (profile_id, email_id)
			VALUES ((SELECT id FROM profile WHERE login=$1), $3), ((SELECT id FROM profile WHERE login=$2), $3)
		`
		_, err = r.DB.Exec(query, sender, recipient, email_id)
		args = []interface{}{sender, recipient, email_id}
	}

	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return fmt.Errorf("Profile_email with email_id=%d and fail", email_id)
	}

	return nil
}

func (r *EmailRepository) AddProfileEmailMyself(email_id uint64, sender, recipient string, ctx context.Context) error {
	query := `
		INSERT INTO profile_email (profile_id, email_id)
		VALUES ((SELECT id FROM profile WHERE login=$1), $2)
	`

	start := time.Now()
	_, err := r.DB.Exec(query, sender, email_id)

	args := []interface{}{sender, email_id}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return fmt.Errorf("Profile_email with profile_id=%d and fail", email_id)
	}

	return nil

}

func (r *EmailRepository) FindEmail(login string, ctx context.Context) error {
	query := "SELECT * FROM profile WHERE login = $1"

	var userModelDb repository_models.User

	start := time.Now()
	err := r.DB.Get(&userModelDb, query, login)

	args := []interface{}{login}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return fmt.Errorf("user with login = %v not found", login)
	}

	return nil
}

func (r *EmailRepository) GetAllIncoming(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important, f.file_id AS photoid
		FROM email e
		LEFT JOIN email_file ef ON e.id = ef.email_id
		LEFT JOIN file f ON ef.file_id = f.id
		JOIN profile_email pe ON e.id = pe.email_id
		JOIN profile p ON pe.profile_id = (
			SELECT id FROM profile WHERE login = $1
		)
		WHERE e.recipient_email = $1 AND e.isSpam = false AND e.isDraft = false
		ORDER BY e.date_of_dispatch DESC
	`

	emailsModelDb := []repository_models.Email{}

	var err error
	args := []interface{}{}
	start := time.Now()
	if offset >= 0 && limit > 0 {
		query += " OFFSET $2 LIMIT $3"
		args = []interface{}{login, offset, limit}
		err = r.DB.Select(&emailsModelDb, query, login, offset, limit)
	} else {
		args = []interface{}{login}
		err = r.DB.Select(&emailsModelDb, query, login)
	}

	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have emails")
		}
		return nil, err
	}

	var emailsModelCore []*domain.Email
	for _, e := range emailsModelDb {
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(e))
	}

	return emailsModelCore, nil
}

func (r *EmailRepository) GetAllSent(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important, f.file_id AS photoid
		FROM email e
		LEFT JOIN email_file ef ON e.id = ef.email_id
		LEFT JOIN file f ON ef.file_id = f.id
		JOIN profile_email pe ON e.id = pe.email_id
		JOIN profile p ON pe.profile_id = (
			SELECT id FROM profile WHERE login = $1
		)
		WHERE e.sender_email = $1 AND e.isSpam = false AND e.isDraft = false
		ORDER BY e.date_of_dispatch DESC
	`

	emailsModelDb := []repository_models.Email{}

	var err error
	args := []interface{}{}
	start := time.Now()
	if offset >= 0 && limit > 0 {
		query += " OFFSET $2 LIMIT $3"
		args = []interface{}{login, offset, limit}
		err = r.DB.Select(&emailsModelDb, query, login, offset, limit)
	} else {
		args = []interface{}{login}
		err = r.DB.Select(&emailsModelDb, query, login)
	}

	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have emails")
		}
		return nil, err
	}

	var emailsModelCore []*domain.Email
	for _, e := range emailsModelDb {
		e.ReadStatus = true
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(e))
	}

	return emailsModelCore, nil
}

func (r *EmailRepository) GetAllDraft(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important, f.file_id AS photoid
		FROM email e
		LEFT JOIN email_file ef ON e.id = ef.email_id
		LEFT JOIN file f ON ef.file_id = f.id
		JOIN profile_email pe ON e.id = pe.email_id
		JOIN profile p ON pe.profile_id = (
			SELECT id FROM profile WHERE login = $1
		)
		WHERE e.sender_email = $1 AND e.isDraft = true
		ORDER BY e.date_of_dispatch DESC
	`

	emailsModelDb := []repository_models.Email{}

	var err error
	args := []interface{}{}
	start := time.Now()
	if offset >= 0 && limit > 0 {
		query += " OFFSET $2 LIMIT $3"
		args = []interface{}{login, offset, limit}
		err = r.DB.Select(&emailsModelDb, query, login, offset, limit)
	} else {
		args = []interface{}{login}
		err = r.DB.Select(&emailsModelDb, query, login)
	}

	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have emails")
		}
		return nil, err
	}

	var emailsModelCore []*domain.Email
	for _, e := range emailsModelDb {
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(e))
	}

	return emailsModelCore, nil
}

func (r *EmailRepository) GetAllSpam(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important, f.file_id AS photoid
		FROM email e
		LEFT JOIN email_file ef ON e.id = ef.email_id
		LEFT JOIN file f ON ef.file_id = f.id
		JOIN profile_email pe ON e.id = pe.email_id
		JOIN profile p ON pe.profile_id = (
			SELECT id FROM profile WHERE login = $1
		)
		WHERE e.sender_email = $1 AND e.isSpam = true
		ORDER BY e.date_of_dispatch DESC
	`

	emailsModelDb := []repository_models.Email{}

	var err error
	args := []interface{}{}
	start := time.Now()
	if offset >= 0 && limit > 0 {
		query += " OFFSET $2 LIMIT $3"
		args = []interface{}{login, offset, limit}
		err = r.DB.Select(&emailsModelDb, query, login, offset, limit)
	} else {
		args = []interface{}{login}
		err = r.DB.Select(&emailsModelDb, query, login)
	}

	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have emails")
		}
		return nil, err
	}

	var emailsModelCore []*domain.Email
	for _, e := range emailsModelDb {
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(e))
	}

	return emailsModelCore, nil
}

func (r *EmailRepository) GetByID(id uint64, login string, ctx context.Context) (*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important, f.file_id AS photoid
		FROM email e
		LEFT JOIN email_file ef ON e.id = ef.email_id
		LEFT JOIN file f ON ef.file_id = f.id
		JOIN profile_email pe ON e.id = pe.email_id
		JOIN profile p ON pe.profile_id = (
			SELECT id FROM profile WHERE login = $2
		)
		WHERE e.id = $1
	`

	var emailModelDb repository_models.Email
	start := time.Now()
	err := r.DB.Get(&emailModelDb, query, id, login)

	args := []interface{}{id, login}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("email with id %d not found", id)
		}
		return nil, err
	}

	return converters.EmailConvertDbInCore(emailModelDb), nil
}

func (r *EmailRepository) Update(newEmail *domain.Email, ctx context.Context) (bool, error) {
	newEmailDb := converters.EmailConvertCoreInDb(*newEmail)

	query := `
        UPDATE email
        SET
            topic = $1, 
            text = $2, 
            isread = $3, 
            isdeleted = $4, 
            isdraft = $5, 
            isspam = $6,
            reply_to_email_id = $7, 
            is_important = $8
        WHERE
            id = $9 AND sender_email = $10
    `

	start := time.Now()
	result, err := r.DB.Exec(query, newEmailDb.Topic, newEmailDb.Text, newEmailDb.ReadStatus, newEmailDb.Deleted, newEmailDb.DraftStatus, newEmail.SpamStatus, newEmailDb.ReplyToEmailID, newEmailDb.Flag, newEmailDb.ID, newEmailDb.SenderEmail)

	args := []interface{}{newEmailDb.Topic, newEmailDb.Text, newEmailDb.ReadStatus, newEmailDb.Deleted, newEmailDb.DraftStatus, newEmail.SpamStatus, newEmailDb.ReplyToEmailID, newEmailDb.Flag, newEmailDb.ID, newEmailDb.RecipientEmail}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return false, fmt.Errorf("failed to update email: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		err = fmt.Errorf("email with id %d not found", newEmailDb.ID)
		return false, err
	}

	return true, nil
}

func (r *EmailRepository) Delete(id uint64, login string, ctx context.Context) (bool, error) {
	query := `
		DELETE FROM profile_email
		WHERE profile_id = (
			SELECT profile_id 
			FROM profile_email pe
			JOIN profile p ON pe.profile_id = p.id
			WHERE email_id = $1 AND p.login = $2
		)
		AND email_id = $1;
	`

	start := time.Now()
	result, err := r.DB.Exec(query, id, login)

	args := []interface{}{id, login}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return false, fmt.Errorf("failed to delete email: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		err = fmt.Errorf("email with id %d not found", id)
		return false, err
	}

	return true, nil
}
