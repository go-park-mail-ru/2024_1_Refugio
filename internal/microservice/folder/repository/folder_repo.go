package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	domain "mail/internal/microservice/models/domain_models"
	converters "mail/internal/microservice/models/repository_converters"
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

func (r *FolderRepository) CreateFolder(folderModelCore *domain.Folder, ctx context.Context) (uint64, *domain.Folder, error) {
	insertEmailQuery := `
		INSERT INTO folder (profile_id, name)
		VALUES ($1, $2)
		RETURNING id
	`

	folderModelDb := converters.FolderConvertCoreInDb(*folderModelCore)

	var id uint64

	start := time.Now()
	err := r.DB.QueryRow(insertEmailQuery, folderModelDb.ProfileId, folderModelDb.Name).Scan(&id)
	if err != nil {
		return 0, &domain.Folder{}, fmt.Errorf("failed to add email: %v", err)
	}

	args := []interface{}{folderModelDb.ProfileId, folderModelDb.Name}
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(insertEmailQuery, ctx.Value(requestIDContextKey).([]string)[0], start, &err, args)

	return id, folderModelCore, nil
}
