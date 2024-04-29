package repository_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

// FolderConvertDbInCore converts an folder model from database representation to core domain representation.
func FolderConvertDbInCore(folderModelDb database.Folder) *domain.Folder {
	return &domain.Folder{
		ID:        folderModelDb.ID,
		ProfileId: folderModelDb.ProfileId,
		Name:      folderModelDb.Name,
	}
}

// FolderConvertCoreInDb converts an folder model from core domain representation to database representation.
func FolderConvertCoreInDb(emailModelCore domain.Folder) *database.Folder {
	return &database.Folder{
		ID:        emailModelCore.ID,
		ProfileId: emailModelCore.ProfileId,
		Name:      emailModelCore.Name,
	}
}
