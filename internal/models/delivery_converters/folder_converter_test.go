package delivery_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	api "mail/internal/models/delivery_models"
	"reflect"
	"testing"
)

func TestFolderConvertCoreInApi(t *testing.T) {
	folderModelCore := domain.Folder{
		ID:        1,
		ProfileId: 1,
		Name:      "folder",
	}

	folderModelApi := FolderConvertCoreInApi(folderModelCore)

	expectedFolderModelApi := &api.Folder{
		ID:        folderModelCore.ID,
		ProfileId: folderModelCore.ProfileId,
		Name:      folderModelCore.Name,
	}

	if !reflect.DeepEqual(folderModelApi, expectedFolderModelApi) {
		t.Errorf("FolderConvertCoreInApi() = %v, want %v", folderModelApi, expectedFolderModelApi)
	}
}

func TestFolderConvertApiInCore(t *testing.T) {
	folderModelApi := api.Folder{
		ID:        1,
		ProfileId: 1,
		Name:      "folder",
	}

	folderModelCore := FolderConvertApiInCore(folderModelApi)

	expectedFolderModelCore := &domain.Folder{
		ID:        folderModelApi.ID,
		ProfileId: folderModelApi.ProfileId,
		Name:      folderModelApi.Name,
	}

	if !reflect.DeepEqual(folderModelCore, expectedFolderModelCore) {
		t.Errorf("FolderConvertApiInCore() = %v, want %v", folderModelCore, expectedFolderModelCore)
	}
}
