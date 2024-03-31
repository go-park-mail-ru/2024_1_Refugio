package email

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	domain "mail/pkg/domain/models"
	"mail/pkg/repository/converters"
	database "mail/pkg/repository/models"
	"time"
)

type EmailRepository struct {
	DB *sqlx.DB
}

func NewEmailRepository(db *sqlx.DB) *EmailRepository {
	return &EmailRepository{DB: db}
}

func (r *EmailRepository) Add(emailModelCore *domain.Email) (*domain.Email, error) {
	query := `
		INSERT INTO email (topic, text, date_of_dispatch, photoid, sender_email, recipient_email, read_status, deleted_status, draft_status, flag)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	emailModelDb := converters.EmailConvertCoreInDb(*emailModelCore)
	_, err := r.DB.Exec(query, emailModelDb.Topic, emailModelDb.Text, time.Now(), emailModelDb.PhotoID, emailModelDb.SenderEmail, emailModelDb.RecipientEmail, emailModelDb.ReadStatus, emailModelDb.Deleted, emailModelDb.DraftStatus, emailModelDb.Flag)
	if err != nil {
		return &domain.Email{}, fmt.Errorf("Email with id %d fail", emailModelDb.ID)
	}

	return emailModelCore, nil
}

func (r *EmailRepository) GetAll(offset, limit int) ([]*domain.Email, error) {
	query := `SELECT * FROM email`
	emailsModelDb := []database.Email{}

	var err error
	if offset >= 0 && limit > 0 {
		query += " OFFSET $1 LIMIT $2"
		err = r.DB.Select(&emailsModelDb, query, offset, limit)
	} else {
		err = r.DB.Select(&emailsModelDb, query)
	}
	//err := r.DB.Select(&emailsModelDb, query)
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

func (r *EmailRepository) GetByID(id uint64) (*domain.Email, error) {
	query := "SELECT * FROM email WHERE id = $1"

	var emailModelDb database.Email
	err := r.DB.Get(&emailModelDb, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("email with id %d not found", id)
		}
		return nil, err
	}

	return converters.EmailConvertDbInCore(emailModelDb), nil
}

func (r *EmailRepository) Update(newEmail *domain.Email) (bool, error) {
	newEmailDb := converters.EmailConvertCoreInDb(*newEmail)

	query := `
        UPDATE email
        SET
            topic = $1, 
            text = $2, 
            photoid = $3, 
            /*sender_email = $5, 
            recipient_email = $6, 
		    date_of_dispatch = $3, */
            read_status = $4, 
            deleted_status = $5, 
            draft_status = $6, 
            reply_to_email_id = $7, 
            flag = $8
        WHERE
            id = $9
    `

	result, err := r.DB.Exec(
		query,
		newEmailDb.Topic,
		newEmailDb.Text,
		newEmailDb.PhotoID,
		newEmailDb.ReadStatus,
		newEmailDb.Deleted,
		newEmailDb.DraftStatus,
		newEmailDb.ReplyToEmailID,
		newEmailDb.Flag,
		newEmailDb.ID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to update email: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("email with id %d not found", newEmailDb.ID)
	}

	return true, nil
}

func (r *EmailRepository) Delete(id uint64) (bool, error) {
	query := "DELETE FROM email WHERE id = $1"

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete email: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("email with id %d not found", id)
	}

	return true, nil
}
