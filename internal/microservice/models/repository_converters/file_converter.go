package repository_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

// FileConvertDbInCore converts a file model from database representation to core domain representation.
func FileConvertDbInCore(fileModelDb *database.File) *domain.File {
	return &domain.File{
		ID:       fileModelDb.ID,
		FileId:   fileModelDb.FileId,
		FileType: fileModelDb.FileType,
		FileName: fileModelDb.FileName,
		FileSize: fileModelDb.FileSize,
	}
}

// FileConvertCoreInDb converts a file model from core domain representation to database representation.
func FileConvertCoreInDb(fileModelCore *domain.File) *database.File {
	return &database.File{
		ID:       fileModelCore.ID,
		FileId:   fileModelCore.FileId,
		FileType: fileModelCore.FileType,
		FileName: fileModelCore.FileName,
		FileSize: fileModelCore.FileSize,
	}
}
