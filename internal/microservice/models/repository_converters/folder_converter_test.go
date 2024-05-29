package repository_converters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

func TestFolderConvertDbInCore(t *testing.T) {
	folderModelDb := database.Folder{
		ID:        123,
		ProfileId: 456,
		Name:      "Inbox",
	}

	expectedCore := &domain.Folder{
		ID:        123,
		ProfileId: 456,
		Name:      "Inbox",
	}

	actualCore := FolderConvertDbInCore(&folderModelDb)
	assert.Equal(t, expectedCore, actualCore)
}

func TestFolderConvertCoreInDb(t *testing.T) {
	folderModelCore := domain.Folder{
		ID:        123,
		ProfileId: 456,
		Name:      "Inbox",
	}

	expectedDb := &database.Folder{
		ID:        123,
		ProfileId: 456,
		Name:      "Inbox",
	}

	actualDb := FolderConvertCoreInDb(&folderModelCore)
	assert.Equal(t, expectedDb, actualDb)
}
