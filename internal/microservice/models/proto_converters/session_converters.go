package proto_converters

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	domain "mail/internal/microservice/models/domain_models"
	grpc "mail/internal/microservice/session/proto"
)

func SessionConvertCoreInProto(sessionModelCore domain.Session) *grpc.Session {
	return &grpc.Session{
		SessionId:    sessionModelCore.ID,
		UserId:       sessionModelCore.UserID,
		CreationDate: timestamppb.New(sessionModelCore.CreationDate),
		Device:       sessionModelCore.Device,
		LifeTime:     int32(sessionModelCore.LifeTime),
		CsrfToken:    sessionModelCore.CsrfToken,
	}
}

func SessionConvertProtoInCore(sessionModelProto grpc.Session) *domain.Session {
	return &domain.Session{
		ID:           sessionModelProto.SessionId,
		UserID:       sessionModelProto.UserId,
		CreationDate: sessionModelProto.CreationDate.AsTime(),
		Device:       sessionModelProto.Device,
		LifeTime:     int(sessionModelProto.LifeTime),
		CsrfToken:    sessionModelProto.CsrfToken,
	}
}
