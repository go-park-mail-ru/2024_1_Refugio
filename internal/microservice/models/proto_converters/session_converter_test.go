package proto_converters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "mail/internal/microservice/models/domain_models"
	grpc "mail/internal/microservice/session/proto"
)

func TestSessionConvertCoreInProto(t *testing.T) {
	creationDate := time.Now()
	sessionModelCore := domain.Session{
		ID:           "session123",
		UserID:       42,
		CreationDate: creationDate,
		Device:       "mobile",
		LifeTime:     3600,
		CsrfToken:    "csrf_token",
	}

	expectedProto := &grpc.Session{
		SessionId:    "session123",
		UserId:       42,
		CreationDate: timestamppb.New(creationDate),
		Device:       "mobile",
		LifeTime:     3600,
		CsrfToken:    "csrf_token",
	}

	actualProto := SessionConvertCoreInProto(&sessionModelCore)
	assert.Equal(t, expectedProto, actualProto)
}

func TestSessionConvertProtoInCore(t *testing.T) {
	sessionModelProto := grpc.Session{
		SessionId: "session123",
		UserId:    42,
		Device:    "mobile",
		LifeTime:  3600,
		CsrfToken: "csrf_token",
	}

	actualCore := SessionConvertProtoInCore(&sessionModelProto)
	assert.NotNil(t, actualCore)
}
