package proto_converters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	grpc "mail/internal/microservice/email/proto"
	domain "mail/internal/microservice/models/domain_models"
)

func TestFileConvertCoreInProto(t *testing.T) {
	fileModelCore := domain.File{
		ID:       123,
		FileId:   "file_id_123",
		FileType: "pdf",
	}

	fileModelProto := FileConvertCoreInProto(&fileModelCore)

	assert.Equal(t, fileModelCore.ID, fileModelProto.Id)
	assert.Equal(t, fileModelCore.FileId, fileModelProto.FileId)
	assert.Equal(t, fileModelCore.FileType, fileModelProto.FileType)
}

func TestFileConvertProtoInCore(t *testing.T) {
	fileModelProto := grpc.File{
		Id:       123,
		FileId:   "file_id_123",
		FileType: "pdf",
	}

	fileModelCore := FileConvertProtoInCore(&fileModelProto)

	assert.Equal(t, fileModelProto.Id, fileModelCore.ID)
	assert.Equal(t, fileModelProto.FileId, fileModelCore.FileId)
	assert.Equal(t, fileModelProto.FileType, fileModelCore.FileType)
}
