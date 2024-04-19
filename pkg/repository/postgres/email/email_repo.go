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
	"os"
	"time"
)

type EmailRepository struct {
	DB *sqlx.DB
}

var Logger = logger.InitializationEmptyLog()

func NewEmailRepository(db *sqlx.DB) *EmailRepository {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile in email_repo" + "log.txt")
	}
	Logger = logger.InitializationBdLog(f)
	return &EmailRepository{DB: db}
}

func (r *EmailRepository) Add(emailModelCore *domain.Email, requestID string) (int64, *domain.Email, error) {
	query := `
		INSERT INTO email (topic, text, date_of_dispatch, /*photoid,*/ sender_email, recipient_email, read_status, deleted_status, draft_status, reply_to_email_id, flag) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING id
	`

	/*query := `
		INSERT INTO profile_email (profile_id, email_id)
		VALUES ((SELECT id FROM profile WHERE login=$1), $3), ((SELECT id FROM profile WHERE login=$2), $3)
	`*/

	emailModelDb := converters.EmailConvertCoreInDb(*emailModelCore)
	format := "2006/01/02 15:04:05"
	args := []interface{}{emailModelDb.Topic, emailModelDb.Text, time.Now().Format(format), emailModelDb.SenderEmail, emailModelDb.RecipientEmail, emailModelDb.ReadStatus, emailModelDb.Deleted, emailModelDb.DraftStatus, emailModelDb.ReplyToEmailID, emailModelDb.Flag}
	var id int64
	start := time.Now()
	err := r.DB.QueryRow(query, emailModelDb.Topic, emailModelDb.Text, time.Now().Format(format), emailModelDb.SenderEmail, emailModelDb.RecipientEmail, emailModelDb.ReadStatus, emailModelDb.Deleted, emailModelDb.DraftStatus, emailModelDb.ReplyToEmailID, emailModelDb.Flag).Scan(&id)
	defer Logger.DbLog(query, requestID, start, &err, args)

	if err != nil {
		return 0, &domain.Email{}, fmt.Errorf("Email with id %d fail", id)
	}

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
	defer Logger.DbLog(query, requestID, start, &err, args)

	if err != nil {
		return fmt.Errorf("Profile_email with profile_id=%d and fail", email_id)
	}

	return nil

}

func (r *EmailRepository) FindEmail(login, requestID string) error {
	query := "SELECT * FROM profile WHERE login = $1"
	args := []interface{}{login}
	var userModelDb database.User
	start := time.Now()

	err := r.DB.Get(&userModelDb, query, login)
	defer Logger.DbLog(query, requestID, start, &err, args)

	if err != nil {
		return fmt.Errorf("user with login = %v not found", login)
	}

	return nil
}

func (r *EmailRepository) GetAllIncoming(login, requestID string, offset, limit int) ([]*domain.Email, error) {
	query := `
		SELECT email.id, email.topic, email.text, email.date_of_dispatch, email.sender_email, email.recipient_email, email.read_status, email.deleted_status, email.draft_status, email.reply_to_email_id, email.flag, profile.avatar_id
		FROM email
		JOIN profile ON email.sender_email = profile.login
		WHERE sender_email = $1
		ORDER BY date_of_dispatch DESC
	`
	/*query := `
		SELECT * FROM email
		WHERE recipient_email = $1
		ORDER BY date_of_dispatch ASC
	`*/
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
	defer Logger.DbLog(query, requestID, start, &err, args)

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

func (r *EmailRepository) GetAllSent(login, requestID string, offset, limit int) ([]*domain.Email, error) {
	query := `
		SELECT email.id, email.topic, email.text, email.date_of_dispatch, email.sender_email, email.recipient_email, email.read_status, email.deleted_status, email.draft_status, email.reply_to_email_id, email.flag, profile.avatar_id
		FROM email
		JOIN profile ON email.sender_email = profile.login
		WHERE sender_email = $1
		ORDER BY date_of_dispatch DESC
	`
	/*query := `
		SELECT * FROM email
		WHERE sender_email = $1
		ORDER BY date_of_dispatch ASC
	`*/
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
	defer Logger.DbLog(query, requestID, start, &err, args)

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

func (r *EmailRepository) GetByID(id uint64, login, requestID string) (*domain.Email, error) {
	query := `
		SELECT * FROM email WHERE id = $1 AND (recipient_email = $2 OR sender_email = $2)
	`
	args := []interface{}{id, login}

	var emailModelDb database.Email
	start := time.Now()
	err := r.DB.Get(&emailModelDb, query, int(id), login)
	defer Logger.DbLog(query, requestID, start, &err, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("email with id %d not found", id)
		}
		return nil, err
	}

	return converters.EmailConvertDbInCore(emailModelDb), nil
}

func (r *EmailRepository) Update(newEmail *domain.Email, requestID string) (bool, error) {
	newEmailDb := converters.EmailConvertCoreInDb(*newEmail)

	query := `
        UPDATE email
        SET
            topic = $1, 
            text = $2, 
            read_status = $3, 
            deleted_status = $4, 
            draft_status = $5, 
            reply_to_email_id = $6, 
            flag = $7
        WHERE
            id = $8 AND sender_email = $9
    `
	args := []interface{}{newEmailDb.Topic, newEmailDb.Text, newEmailDb.ReadStatus, newEmailDb.Deleted, newEmailDb.DraftStatus, newEmailDb.ReplyToEmailID, newEmailDb.Flag, newEmailDb.ID, newEmailDb.RecipientEmail}

	start := time.Now()
	result, err := r.DB.Exec(query, newEmailDb.Topic, newEmailDb.Text, newEmailDb.ReadStatus, newEmailDb.Deleted, newEmailDb.DraftStatus, newEmailDb.ReplyToEmailID, newEmailDb.Flag, newEmailDb.ID, newEmailDb.SenderEmail)
	defer Logger.DbLog(query, requestID, start, &err, args)
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

func (r *EmailRepository) Delete(id uint64, login, requestID string) (bool, error) {
	query := "DELETE FROM email WHERE id = $1 AND (recipient_email = $2 OR sender_email = $2)"

	args := []interface{}{id, login}
	start := time.Now()
	result, err := r.DB.Exec(query, id, login)
	defer Logger.DbLog(query, requestID, start, &err, args)
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
