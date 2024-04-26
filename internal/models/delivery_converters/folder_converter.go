package delivery_converters

import (
	folderCore "mail/internal/microservice/models/domain_models"
	folderApi "mail/internal/models/delivery_models"
)

// FolderConvertCoreInApi converts an folder model from the core package to the API representation.
func FolderConvertCoreInApi(folderModelDb folderCore.Folder) *folderApi.Folder {
	return &folderApi.Folder{
		ID:        folderModelDb.ID,
		ProfileId: folderModelDb.ProfileId,
		Name:      folderModelDb.Name,
	}
}

// FolderConvertApiInCore converts an folder model from the API representation to the core package.
func FolderConvertApiInCore(folderModelApi folderApi.Folder) *folderCore.Folder {
	return &folderCore.Folder{
		ID:        folderModelApi.ID,
		ProfileId: folderModelApi.ProfileId,
		Name:      folderModelApi.Name,
	}
}
