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

var requestIDContextKey interface{} = "requestID"

type FolderRepository struct {
	DB *sqlx.DB
}

func NewFolderRepository(db *sqlx.DB) *FolderRepository {
	return &FolderRepository{DB: db}
}

// Create adds a new folder to the storage and returns its assigned unique identifier.
func (r *FolderRepository) Create(folderModelCore *domain.Folder, ctx context.Context) (uint32, *domain.Folder, error) {
	insertFolderQuery := `
		INSERT INTO folder (profile_id, name)
		VALUES ($1, $2)
		RETURNING id
	`

	folderModelDb := converters.FolderConvertCoreInDb(folderModelCore)

	var id uint32

	start := time.Now()
	err := r.DB.QueryRow(insertFolderQuery, folderModelDb.ProfileId, folderModelDb.Name).Scan(&id)
	if err != nil {
		return 0, &domain.Folder{}, fmt.Errorf("failed to add folder: %v", err)
	}

	args := []interface{}{folderModelDb.ProfileId, folderModelDb.Name}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(insertFolderQuery, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	folderModelCore.ID = id
	return id, folderModelCore, nil
}

// GetAll get list folder user.
func (r *FolderRepository) GetAll(profileID uint32, offset, limit int64, ctx context.Context) ([]*domain.Folder, error) {
	query := `
		SELECT folder.id, folder.name, folder.profile_id FROM folder
		WHERE profile_id = $1
	`

	var foldersModelDb []repository_models.Folder

	var err error
	var args []interface{}
	start := time.Now()

	if offset >= 0 && limit > 0 {
		query += " OFFSET $2 LIMIT $3"
		args = []interface{}{profileID, offset, limit}
		err = r.DB.Select(&foldersModelDb, query, profileID, offset, limit)
	} else {
		args = []interface{}{profileID}
		err = r.DB.Select(&foldersModelDb, query, profileID)
	}

	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have folders")
		}
		return nil, err
	}

	var foldersModelCore []*domain.Folder
	for _, e := range foldersModelDb {
		foldersModelCore = append(foldersModelCore, converters.FolderConvertDbInCore(&e))
	}

	return foldersModelCore, nil
}

// Delete delete folder as user.
func (r *FolderRepository) Delete(folderID uint32, profileID uint32, ctx context.Context) (bool, error) {
	query := `
		DELETE FROM folder
		WHERE folder.id = $1 AND profile_id = $2
	`

	start := time.Now()
	result, err := r.DB.Exec(query, folderID, profileID)

	args := []interface{}{folderID, profileID}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return false, fmt.Errorf("failed to delete folder: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		err = fmt.Errorf("folderID=%d and profileID=%d dnot found", folderID, profileID)
		return false, err
	}

	return true, nil
}

// Update folder as user.
func (r *FolderRepository) Update(folderModelCore *domain.Folder, ctx context.Context) (bool, error) {
	query := `
        UPDATE folder
        SET
            name = $1
        WHERE
            folder.id = $2 AND folder.profile_id = $3
    `

	newUdFolderDb := converters.FolderConvertCoreInDb(folderModelCore)

	start := time.Now()
	result, err := r.DB.Exec(query, newUdFolderDb.Name, newUdFolderDb.ID, newUdFolderDb.ProfileId)

	args := []interface{}{newUdFolderDb.Name, newUdFolderDb.ID, newUdFolderDb.ProfileId}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return false, fmt.Errorf("failed to update folder: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		err = fmt.Errorf("folderID=%d and profileID=%d not found", newUdFolderDb.ID, newUdFolderDb.ProfileId)
		return false, err
	}

	return true, nil
}

// AddEmailFolder adds a new email in folder to the storage and returns its assigned unique identifier.
func (r *FolderRepository) AddEmailFolder(folderID, emailID uint32, ctx context.Context) (bool, error) {
	insertFolderEmailQuery := `
		INSERT INTO folder_email (folder_id, email_id)
		VALUES ($1, $2)
	`

	start := time.Now()
	_, err := r.DB.Exec(insertFolderEmailQuery, folderID, emailID)
	if err != nil {
		return false, fmt.Errorf("failed to add email in folder: %v", err)
	}

	args := []interface{}{folderID, emailID}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(insertFolderEmailQuery, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	return true, nil
}

// DeleteEmailFolder adds a new email in folder to the storage and returns its assigned unique identifier.
func (r *FolderRepository) DeleteEmailFolder(folderID uint32, emailID uint32, ctx context.Context) (bool, error) {
	query := `
		DELETE FROM folder_email
		WHERE folder_id = $1 AND email_id = $2
	`

	start := time.Now()
	result, err := r.DB.Exec(query, folderID, emailID)

	args := []interface{}{folderID, emailID}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		return false, fmt.Errorf("failed to delete email in folder: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		err = fmt.Errorf("folderID=%d and emailID=%d dnot found", folderID, emailID)
		return false, err
	}

	return true, nil
}

// CheckFolder checking that the folder belongs to the user.
func (r *FolderRepository) CheckFolder(folderID uint32, profileID uint32, ctx context.Context) (bool, error) {
	query := `
		SELECT id, profile_id, name 
		FROM folder 
		WHERE folder.id = $1 AND folder.profile_id = $2
	`

	folderModelDb := []repository_models.Folder{}

	start := time.Now()
	err := r.DB.Select(&folderModelDb, query, folderID, profileID)
	if err != nil || len(folderModelDb) == 0 {
		return false, fmt.Errorf("failed to check folder: %v", err)
	}

	args := []interface{}{folderID, profileID}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	return true, nil
}

// CheckEmail checking that the email belongs to the user.
func (r *FolderRepository) CheckEmail(emailID uint32, profileID uint32, ctx context.Context) (bool, error) {
	query := `
		SELECT profile_id, email_id 
		FROM profile_email 
		WHERE email_id = $1 AND profile_id = $2
	`

	folderProfileEmailModelDb := []repository_models.ProfileEmail{}

	start := time.Now()
	err := r.DB.Select(&folderProfileEmailModelDb, query, emailID, profileID)
	if err != nil || len(folderProfileEmailModelDb) == 0 {
		return false, fmt.Errorf("failed to check email: %v", err)
	}

	args := []interface{}{emailID, profileID}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	return true, nil
}

// GetAllEmails get list emails folder user.
func (r *FolderRepository) GetAllEmails(folderID, profileID, limit, offset uint32, ctx context.Context) ([]*domain.Email, error) {
	query := `
		SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important
		FROM email e
		JOIN profile_email pe ON e.id = pe.email_id
		JOIN profile p ON pe.profile_id = $1 
		LEFT JOIN folder_email ON e.id = folder_email.email_id
		WHERE folder_email.folder_id = $2
		ORDER BY e.date_of_dispatch DESC
	`

	var emailsModelDb []repository_models.Email

	var err error
	var args []interface{}
	start := time.Now()

	if offset > 0 && limit > 0 {
		query += " OFFSET $3 LIMIT $4"
		args = []interface{}{profileID, folderID, offset, limit}
		err = r.DB.Select(&emailsModelDb, query, profileID, folderID, offset, limit)
	} else {
		args = []interface{}{profileID, folderID}
		err = r.DB.Select(&emailsModelDb, query, profileID, folderID)
	}

	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have emails in folder")
		}
		return nil, err
	}

	var emailsModelCore []*domain.Email
	for _, e := range emailsModelDb {
		emailsModelCore = append(emailsModelCore, converters.EmailConvertDbInCore(&e))
	}

	return emailsModelCore, nil
}

// GetAvatarFileIDByLogin getting an avatar by login.
func (r *FolderRepository) GetAvatarFileIDByLogin(login string, ctx context.Context) (string, error) {
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

// GetAllFolderName retrieves the names of all folders associated with a given email ID.
func (r *FolderRepository) GetAllFolderName(emailID uint32, ctx context.Context) ([]*domain.Folder, error) {
	query := `
		SELECT f.id, f.name
		FROM folder f
		INNER JOIN folder_email fe ON f.id = fe.folder_id
		WHERE fe.email_id = $1
	`

	foldersModelDb := []repository_models.Folder{}

	var err error
	args := []interface{}{emailID}
	start := time.Now()
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	err = r.DB.Select(&foldersModelDb, query, emailID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("DB no have folders")
		}
		return nil, err
	}

	var foldersModelCore []*domain.Folder
	for _, e := range foldersModelDb {
		foldersModelCore = append(foldersModelCore, converters.FolderConvertDbInCore(&e))
	}

	return foldersModelCore, nil
}
