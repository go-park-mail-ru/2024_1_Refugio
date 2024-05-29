package repository_converters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

func TestFileConvertDbInCore(t *testing.T) {
	fileModelDb := database.File{
		ID:       123,
		FileId:   "url",
		FileType: "Photo",
	}

	expectedCore := &domain.File{
		ID:       123,
		FileId:   "url",
		FileType: "Photo",
	}

	actualCore := FileConvertDbInCore(&fileModelDb)
	assert.Equal(t, expectedCore, actualCore)
}

func TestFileConvertCoreInDb(t *testing.T) {
	fileModelCore := domain.File{
		ID:       123,
		FileId:   "url",
		FileType: "Photo",
	}

	expectedDb := &database.File{
		ID:       123,
		FileId:   "url",
		FileType: "Photo",
	}

	actualDb := FileConvertCoreInDb(&fileModelCore)
	assert.Equal(t, expectedDb, actualDb)
}
