package email

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"mail/pkg/domain/logger"
	domain "mail/pkg/domain/models"
	"mail/pkg/repository/converters"
	database "mail/pkg/repository/models"
	"time"
)

type EmailRepository struct {
	DB *sqlx.DB
}

var Logger = logger.InitializationBdLog()

func NewEmailRepository(db *sqlx.DB) *EmailRepository {
	return &EmailRepository{DB: db}
}

func (r *EmailRepository) Add(emailModelCore *domain.Email, requestID string) (int64, *domain.Email, error) {
	query := `
		INSERT INTO email (topic, text, date_of_dispatch, photoid, sender_email, recipient_email, read_status, deleted_status, draft_status, flag) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id
	`
	emailModelDb := converters.EmailConvertCoreInDb(*emailModelCore)
	format := "2006/01/02 15:04:05"
	args := []interface{}{emailModelDb.Topic, emailModelDb.Text, time.Now().Format(format), emailModelDb.PhotoID, emailModelDb.SenderEmail, emailModelDb.RecipientEmail, emailModelDb.ReadStatus, emailModelDb.Deleted, emailModelDb.DraftStatus, emailModelDb.Flag}
	var id int64
	start := time.Now()
	err := r.DB.QueryRow(query, args[0].(string), args[1].(string), args[2].(string), args[3].(string), args[4].(string), args[5].(string), args[6].(bool), args[7].(bool), args[8].(bool), args[9].(bool)).Scan(&id)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return 0, emailModelCore, fmt.Errorf("Email with id %d fail", id)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return id, emailModelCore, nil
}

func (r *EmailRepository) AddProfileEmail(email_id int64, sender, recipient, requestID string) error {
	query := `
		INSERT INTO profile_email (profile_id, email_id)
		VALUES ((SELECT id FROM profile WHERE login=$1), $3), ((SELECT id FROM profile WHERE login=$2), $3)
	`
	args := []interface{}{sender, recipient, email_id}
	start := time.Now()
	_, err := r.DB.Exec(query, sender, recipient, email_id)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return fmt.Errorf("Profile_email with profile_id=%d and  fail", email_id)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return nil

}

func (r *EmailRepository) FindEmail(login, requestID string) error {
	query := "SELECT * FROM profile WHERE login = $1"
	args := []interface{}{login}
	var userModelDb database.User
	start := time.Now()
	err := r.DB.Get(&userModelDb, query, login)

	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return fmt.Errorf("user with login = %d not found", login)
	}
	Logger.DbLog(query, requestID, 200, start, nil, args)
	return nil
}

func (r *EmailRepository) GetAllIncoming(login, requestID string, offset, limit int) ([]*domain.Email, error) {
	query := `
		SELECT * FROM email
		WHERE recipient_email = $1
		ORDER BY date_of_dispatch ASC
	`
	emailsModelDb := []database.Email{}

	var err error
	start := time.Now()
	args := []interface{}{}
	if offset >= 0 && limit > 0 {
		query += " OFFSET $2 LIMIT $3"
		args = []interface{}{login, offset, limit}
		err = r.DB.Select(&emailsModelDb, query, login, offset, limit)
	} else {
		args = []interface{}{login}
		err = r.DB.Select(&emailsModelDb, query, login)
	}

	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have emails")
		}
		return nil, err
	}

	var emailsModelCore []*domain.Email
	for _, e := range emailsModelDb {
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(e))
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return emailsModelCore, nil
}

func (r *EmailRepository) GetAllSent(login, requestID string, offset, limit int) ([]*domain.Email, error) {
	query := `
		SELECT * FROM email
		WHERE sender_email = $1
		ORDER BY date_of_dispatch ASC
	`
	emailsModelDb := []database.Email{}

	var err error
	start := time.Now()
	args := []interface{}{}
	if offset >= 0 && limit > 0 {
		query += " OFFSET $2 LIMIT $3"
		args = []interface{}{login, offset, limit}
		err = r.DB.Select(&emailsModelDb, query, login, offset, limit)
	} else {
		args = []interface{}{login}
		err = r.DB.Select(&emailsModelDb, query, login)
	}

	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have emails")
		}
		return nil, err
	}

	var emailsModelCore []*domain.Email
	for _, e := range emailsModelDb {
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(e))
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return emailsModelCore, nil
}

func (r *EmailRepository) GetByID(id uint64, login, requestID string) (*domain.Email, error) {
	query := `
		SELECT * FROM email WHERE id = $1 AND (recipient_email = $2 OR sender_email = &2)
	`
	args := []interface{}{id, login}

	var emailModelDb database.Email
	start := time.Now()
	err := r.DB.Get(&emailModelDb, query, id, login)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("email with id %d not found", id)
		}
		return nil, err
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return converters.EmailConvertDbInCore(emailModelDb), nil
}

func (r *EmailRepository) Update(newEmail *domain.Email, requestID string) (bool, error) {
	newEmailDb := converters.EmailConvertCoreInDb(*newEmail)

	query := `
        UPDATE email
        SET
            topic = $1, 
            text = $2, 
            photoid = $3,
            read_status = $4, 
            deleted_status = $5, 
            draft_status = $6, 
            reply_to_email_id = $7, 
            flag = $8
        WHERE
            id = $9 AND sender_email = $10
    `
	args := []interface{}{newEmailDb.Topic, newEmailDb.Text, newEmailDb.PhotoID, newEmailDb.ReadStatus, newEmailDb.Deleted, newEmailDb.DraftStatus, newEmailDb.ReplyToEmailID, newEmailDb.Flag, newEmailDb.ID, newEmailDb.RecipientEmail}

	start := time.Now()
	result, err := r.DB.Exec(query, newEmailDb.Topic, newEmailDb.Text, newEmailDb.PhotoID, newEmailDb.ReadStatus, newEmailDb.Deleted, newEmailDb.DraftStatus, newEmailDb.ReplyToEmailID, newEmailDb.Flag, newEmailDb.ID, newEmailDb.RecipientEmail)
	if err != nil {
		return false, fmt.Errorf("failed to update email: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return false, fmt.Errorf("email with id %d not found", newEmailDb.ID)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return true, nil
}

func (r *EmailRepository) Delete(id uint64, login, requestID string) (bool, error) {
	query := "DELETE FROM email WHERE id = $1 AND (recipient_email = $2 OR sender_email = &2)"

	args := []interface{}{id, login}
	start := time.Now()
	result, err := r.DB.Exec(query, id, login)
	if err != nil {
		return false, fmt.Errorf("failed to delete email: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return false, fmt.Errorf("email with id %d not found", id)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return true, nil
}
