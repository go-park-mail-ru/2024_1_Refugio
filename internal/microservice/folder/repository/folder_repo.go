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

type FolderRepository struct {
	DB *sqlx.DB
}

func NewFolderRepository(db *sqlx.DB) *FolderRepository {
	return &FolderRepository{DB: db}
}

func (r *FolderRepository) CreateFolder(folderModelCore *domain.Folder, ctx context.Context) (uint32, *domain.Folder, error) {
	insertFolderQuery := `
		INSERT INTO folder (profile_id, name)
		VALUES ($1, $2)
		RETURNING id
	`

	folderModelDb := converters.FolderConvertCoreInDb(*folderModelCore)

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

func (r *FolderRepository) GetAll(profileID uint32, offset, limit int64, ctx context.Context) ([]*domain.Folder, error) {
	query := `
		SELECT folder.id, folder.name, folder.profile_id FROM folder
		WHERE profile_id = $1
	`

	foldersModelDb := []repository_models.Folder{}

	var err error
	args := []interface{}{}
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
		foldersModelCore = append(foldersModelCore, converters.FolderConvertDbInCore(e))
	}

	return foldersModelCore, nil
}
