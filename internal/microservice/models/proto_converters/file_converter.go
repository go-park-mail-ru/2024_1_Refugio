package proto_converters

import (
	grpc "mail/internal/microservice/email/proto"
	domain "mail/internal/microservice/models/domain_models"
)

// FileConvertCoreInProto converts a file model from the application core to the gRPC format.
func FileConvertCoreInProto(fileModelCore *domain.File) *grpc.File {
	return &grpc.File{
		Id:       fileModelCore.ID,
		FileId:   fileModelCore.FileId,
		FileType: fileModelCore.FileType,
		FileName: fileModelCore.FileName,
		FileSize: fileModelCore.FileSize,
	}
}

// FileConvertProtoInCore converts a file model from the gRPC format to the application core.
func FileConvertProtoInCore(fileModelProto *grpc.File) *domain.File {
	return &domain.File{
		ID:       fileModelProto.Id,
		FileId:   fileModelProto.FileId,
		FileType: fileModelProto.FileType,
		FileName: fileModelProto.FileName,
		FileSize: fileModelProto.FileSize,
	}
}
