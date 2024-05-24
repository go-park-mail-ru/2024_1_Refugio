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

// EmailRepository represents a repository for managing email data in the database.
type EmailRepository struct {
	DB *sqlx.DB
}

// NewEmailRepository creates a new instance of EmailRepository with the given database connection.
func NewEmailRepository(db *sqlx.DB) *EmailRepository {
	return &EmailRepository{DB: db}
}

// Add adds a new email to the storage and returns its assigned unique identifier.
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

	var id uint64
	var err error
	start := time.Now()
	emailModelDb := converters.EmailConvertCoreInDb(emailModelCore)
	format := "2006/01/02 15:04:05"

	err = r.DB.QueryRow(insertEmailQuery, emailModelDb.Topic, emailModelDb.Text, time.Now().Format(format), emailModelDb.SenderEmail, emailModelDb.RecipientEmail, emailModelDb.ReadStatus, emailModelDb.Deleted, emailModelDb.DraftStatus, emailModelCore.SpamStatus, emailModelDb.ReplyToEmailID, emailModelDb.Flag).Scan(&id)

	args := []interface{}{emailModelDb.Topic, emailModelDb.Text, time.Now().Format(format), emailModelDb.SenderEmail, emailModelDb.RecipientEmail, emailModelDb.ReadStatus, emailModelDb.Deleted, emailModelDb.DraftStatus, emailModelDb.ReplyToEmailID, emailModelDb.Flag, emailModelDb.SenderEmail}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(insertEmailQuery, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return 0, &domain.Email{}, fmt.Errorf("failed to add email: %v", err)
	}

	_, err = r.DB.Exec(insertEmailFileQuery, id, emailModelDb.SenderEmail)
	if err != nil {
		return 0, &domain.Email{}, fmt.Errorf("failed to add email file: %v", err)
	}

	return id, emailModelCore, nil
}

// AddProfileEmail links an email to one or more profiles based on sender and recipient information.
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
		return fmt.Errorf("profile_email with email_id=%d and fail", email_id)
	}

	return nil
}

// AddProfileEmailMyself links an email to the profile corresponding to the sender (when sender and recipient are the same).
func (r *EmailRepository) AddProfileEmailMyself(email_id uint64, login string, ctx context.Context) error {
	query := `
		INSERT INTO profile_email (profile_id, email_id)
		VALUES ((SELECT id FROM profile WHERE login=$1), $2)
	`

	start := time.Now()
	_, err := r.DB.Exec(query, login, email_id)

	args := []interface{}{login, email_id}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return fmt.Errorf("Profile_email with profile_id=%d and fail", email_id)
	}

	return nil

}

// FindEmail searches for a user in the database based on their login.
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

// GetAllIncoming returns all emails incoming from the storage.
func (r *EmailRepository) GetAllIncoming(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
		FROM email e
		JOIN profile_email pe ON e.id = pe.email_id
		JOIN profile p ON pe.profile_id = (
			SELECT id FROM profile WHERE login = $1
		)
		WHERE e.recipient_email = $1 AND e.isSpam = false AND e.isDraft = false
		ORDER BY e.date_of_dispatch DESC
	`

	var emailsModelDb []repository_models.Email

	var err error
	var args []interface{}
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
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(&e))
	}

	return emailsModelCore, nil
}

// GetAllSent returns all emails sent from the storage.
func (r *EmailRepository) GetAllSent(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
		FROM email e
		JOIN profile_email pe ON e.id = pe.email_id
		JOIN profile p ON pe.profile_id = (
			SELECT id FROM profile WHERE login = $1
		)
		WHERE e.sender_email = $1 AND e.isSpam = false AND e.isDraft = false
		ORDER BY e.date_of_dispatch DESC
	`

	var emailsModelDb []repository_models.Email

	var err error
	var args []interface{}
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
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(&e))
	}

	return emailsModelCore, nil
}

// GetAllDraft returns all draft emails from the storage.
func (r *EmailRepository) GetAllDraft(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
		FROM email e
		JOIN profile_email pe ON e.id = pe.email_id
		JOIN profile p ON pe.profile_id = (
			SELECT id FROM profile WHERE login = $1
		)
		WHERE e.sender_email = $1 AND e.isDraft = true
		ORDER BY e.date_of_dispatch DESC
	`

	var emailsModelDb []repository_models.Email

	var err error
	var args []interface{}
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
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(&e))
	}

	return emailsModelCore, nil
}

// GetAllSpam returns all draft emails from the storage.
func (r *EmailRepository) GetAllSpam(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
		FROM email e
		JOIN profile_email pe ON e.id = pe.email_id
		JOIN profile p ON pe.profile_id = (
			SELECT id FROM profile WHERE login = $1
		)
		WHERE e.recipient_email = $1 AND e.isSpam = true
		ORDER BY e.date_of_dispatch DESC
	`

	var emailsModelDb []repository_models.Email

	var err error
	var args []interface{}
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
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(&e))
	}

	return emailsModelCore, nil
}

// GetByID returns the email by its unique identifier.
func (r *EmailRepository) GetByID(id uint64, login string, ctx context.Context) (*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
		FROM email e
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

	return converters.EmailConvertDbInCore(&emailModelDb), nil
}

// GetAvatarFileIDByLogin getting an avatar by login.
func (r *EmailRepository) GetAvatarFileIDByLogin(login string, ctx context.Context) (string, error) {
	query := `
		SELECT f.file_id
		FROM file f
		LEFT JOIN profile p ON p.avatar_id = f.id
		WHERE p.login = $1
	`

	var fileID sql.NullString
	start := time.Now()
	err := r.DB.Get(&fileID, query, login)

	args := []interface{}{login}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no avatar found for login %s", login)
		}
		return "", err
	}

	if !fileID.Valid {
		return "", nil
	}

	return fileID.String, nil
}

// Update updates the information of an email in the storage based on the provided new email.
func (r *EmailRepository) Update(newEmail *domain.Email, ctx context.Context) (bool, error) {
	newEmailDb := converters.EmailConvertCoreInDb(newEmail)

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

// Delete removes the email from the storage by its unique identifier.
func (r *EmailRepository) Delete(id uint64, login string, ctx context.Context) (bool, error) {
	query := `
		DELETE FROM profile_email
		WHERE profile_id = (
			SELECT profile_id 
			FROM profile_email pe
			JOIN profile p ON pe.profile_id = p.id
			WHERE email_id = $1 AND p.login = $2
		)
		AND email_id = $1
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

// AddFile adds a file entry to the database with the provided file ID, file type, file name and file size.
func (r *EmailRepository) AddFile(fileID string, fileType string, fileName string, fileSize string, ctx context.Context) (uint64, error) {
	query := `
        INSERT INTO file (file_id, file_type, file_name, file_size)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var id uint64
	start := time.Now()
	err := r.DB.QueryRowContext(ctx, query, fileID, fileType, fileName, fileSize).Scan(&id)

	args := []interface{}{fileID, fileType, fileName, fileSize}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return 0, fmt.Errorf("failed to add file: %v", err)
	}

	return id, nil
}

// AddAttachment links a file to an email by inserting a record into the email_file table.
func (r *EmailRepository) AddAttachment(emailID uint64, fileID uint64, ctx context.Context) error {
	query := `
        INSERT INTO email_file (email_id, file_id)
        VALUES ($1, $2)
    `

	start := time.Now()
	_, err := r.DB.ExecContext(ctx, query, emailID, fileID)

	args := []interface{}{emailID, fileID}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return fmt.Errorf("failed to add attachment: %v", err)
	}

	return nil
}

// GetFileByID retrieves file information based on the provided file ID.
func (r *EmailRepository) GetFileByID(id uint64, ctx context.Context) (*domain.File, error) {
	query := `
        SELECT file_id, file_type, file_name, file_size
        FROM file
        WHERE id = $1
    `

	var fileID string
	var fileType string
	var fileName string
	var fileSize string
	start := time.Now()
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&fileID, &fileType, &fileName, &fileSize)

	args := []interface{}{id}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	return &domain.File{ID: id, FileId: fileID, FileType: fileType, FileName: fileName, FileSize: fileSize}, nil
}

// GetFilesByEmailID retrieves all files associated with a given email ID.
func (r *EmailRepository) GetFilesByEmailID(emailID uint64, ctx context.Context) ([]*domain.File, error) {
	query := `
        SELECT f.id, f.file_id, f.file_type, f.file_name, f.file_size
        FROM file f
        JOIN email_file ef ON f.id = ef.file_id
        WHERE ef.email_id = $1
    `

	var filesModelDb []*repository_models.File
	start := time.Now()
	rows, err := r.DB.QueryContext(ctx, query, emailID)

	args := []interface{}{emailID}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return nil, fmt.Errorf("failed to get files")
	}

	for rows.Next() {
		var file repository_models.File
		err := rows.Scan(&file.ID, &file.FileId, &file.FileType, &file.FileName, &file.FileSize)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %v", err)
		}
		if file.FileType == "PHOTO" {
			continue
		}
		filesModelDb = append(filesModelDb, &file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over files: %v", err)
	}

	var filesModelCore []*domain.File
	for _, f := range filesModelDb {
		filesModelCore = append(filesModelCore, converters.FileConvertDbInCore(f))
	}

	return filesModelCore, nil
}

// DeleteFileByID deletes a file entry from the database based on the provided file ID.
func (r *EmailRepository) DeleteFileByID(fileID uint64, ctx context.Context) error {
	query := `
        DELETE FROM file
        WHERE id = $1
    `

	start := time.Now()
	_, err := r.DB.ExecContext(ctx, query, fileID)

	args := []interface{}{fileID}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// UpdateFileByID updates the file ID, file type, file name and file size of a file entry in the database based on the provided file ID.
func (r *EmailRepository) UpdateFileByID(fileID uint64, newFileID string, newFileType string, newFileName string, newFileSize string, ctx context.Context) error {
	query := `
        UPDATE file
        SET file_id = $1, file_type = $2, file_name = $3, file_size = $3
        WHERE id = $4
    `

	start := time.Now()
	_, err := r.DB.ExecContext(ctx, query, newFileID, newFileType, newFileName, newFileSize, fileID)

	args := []interface{}{newFileID, newFileType, newFileName, newFileSize, fileID}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return fmt.Errorf("failed to update file: %v", err)
	}

	return nil
}
