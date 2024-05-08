package server

import (
	"context"
	"fmt"
	"strconv"

	usecase "mail/internal/microservice/email/interface"
	"mail/internal/microservice/email/proto"
	converters "mail/internal/microservice/models/proto_converters"
)

type EmailServer struct {
	proto.UnimplementedEmailServiceServer
	EmailUseCase usecase.EmailUseCase
}

func NewEmailServer(emailUseCase usecase.EmailUseCase) *EmailServer {
	return &EmailServer{EmailUseCase: emailUseCase}
}

func (es *EmailServer) GetEmailByID(ctx context.Context, input *proto.EmailIdAndLogin) (*proto.Email, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid email id: %s", strconv.Itoa(int(input.Id)))
	}

	email, err := es.EmailUseCase.GetEmailByID(input.Id, input.Login, ctx)
	if err != nil {
		return nil, fmt.Errorf("email not found")
	}

	return converters.EmailConvertCoreInProto(*email), nil
}

func (es *EmailServer) GetAllIncoming(ctx context.Context, input *proto.LoginOffsetLimit) (*proto.Emails, error) {
	if input.Login == "" {
		return nil, fmt.Errorf("invalid email login: %s", input.Login)
	}

	emailsCore, err := es.EmailUseCase.GetAllEmailsIncoming(input.Login, input.Offset, input.Limit, ctx)
	if err != nil {
		return nil, fmt.Errorf("email not found")
	}

	emailsProto := make([]*proto.Email, len(emailsCore))
	for i, e := range emailsCore {
		emailsProto[i] = converters.EmailConvertCoreInProto(*e)
	}

	emailProto := new(proto.Emails)
	emailProto.Emails = emailsProto
	return emailProto, nil
}

func (es *EmailServer) GetAllSent(ctx context.Context, input *proto.LoginOffsetLimit) (*proto.Emails, error) {
	if input.Login == "" {
		return nil, fmt.Errorf("invalid email login: %s", input.Login)
	}

	emailsCore, err := es.EmailUseCase.GetAllEmailsSent(input.Login, input.Offset, input.Limit, ctx)
	if err != nil {
		return nil, fmt.Errorf("email not found")
	}

	emailsProto := make([]*proto.Email, len(emailsCore))
	for i, e := range emailsCore {
		emailsProto[i] = converters.EmailConvertCoreInProto(*e)
	}

	emailProto := new(proto.Emails)
	emailProto.Emails = emailsProto
	return emailProto, nil
}

func (es *EmailServer) GetDraftEmails(ctx context.Context, input *proto.LoginOffsetLimit) (*proto.Emails, error) {
	if input.Login == "" {
		return nil, fmt.Errorf("invalid email login: %s", input.Login)
	}

	emailsCore, err := es.EmailUseCase.GetAllDraftEmails(input.Login, input.Offset, input.Limit, ctx)
	if err != nil {
		return nil, fmt.Errorf("email not found")
	}

	emailsProto := make([]*proto.Email, len(emailsCore))
	for i, e := range emailsCore {
		emailsProto[i] = converters.EmailConvertCoreInProto(*e)
	}

	emailProto := new(proto.Emails)
	emailProto.Emails = emailsProto
	return emailProto, nil
}

func (es *EmailServer) GetSpamEmails(ctx context.Context, input *proto.LoginOffsetLimit) (*proto.Emails, error) {
	if input.Login == "" {
		return nil, fmt.Errorf("invalid email login: %s", input.Login)
	}

	emailsCore, err := es.EmailUseCase.GetAllSpamEmails(input.Login, input.Offset, input.Limit, ctx)
	if err != nil {
		return nil, fmt.Errorf("email not found")
	}

	emailsProto := make([]*proto.Email, len(emailsCore))
	for i, e := range emailsCore {
		emailsProto[i] = converters.EmailConvertCoreInProto(*e)
	}

	emailProto := new(proto.Emails)
	emailProto.Emails = emailsProto
	return emailProto, nil
}

func (es *EmailServer) CreateEmail(ctx context.Context, input *proto.Email) (*proto.EmailWithID, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid email format: %s", input)
	}

	id, email, err := es.EmailUseCase.CreateEmail(converters.EmailConvertProtoInCore(*input), ctx)
	if err != nil {
		return nil, fmt.Errorf("failed create email")
	}

	emailWithId := new(proto.EmailWithID)
	emailWithId.Id = id
	emailWithId.Email = converters.EmailConvertCoreInProto(*email)
	return emailWithId, nil
}

func (es *EmailServer) UpdateEmail(ctx context.Context, input *proto.Email) (*proto.StatusEmail, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid email format: %s", input)
	}

	okStatus, err := es.EmailUseCase.UpdateEmail(converters.EmailConvertProtoInCore(*input), ctx)

	if err != nil {
		return nil, fmt.Errorf("email not found")
	}

	protoStatusEmail := new(proto.StatusEmail)
	protoStatusEmail.Status = okStatus
	return protoStatusEmail, nil
}

func (es *EmailServer) DeleteEmail(ctx context.Context, input *proto.LoginWithID) (*proto.StatusEmail, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid email format: %s", input)
	}

	if input.Login == "" || input.Id <= 0 {
		return nil, fmt.Errorf("invalid input data")
	}

	okStatus, err := es.EmailUseCase.DeleteEmail(input.Id, input.Login, ctx)
	if err != nil {
		return nil, fmt.Errorf("email not found")
	}

	protoStatusEmail := new(proto.StatusEmail)
	protoStatusEmail.Status = okStatus
	return protoStatusEmail, nil
}

func (es *EmailServer) CreateProfileEmail(ctx context.Context, input *proto.IdSenderRecipient) (*proto.EmptyEmail, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid email id: %s", strconv.Itoa(int(input.Id)))
	}

	if input.Sender == "" || input.Recipient == "" {
		return nil, fmt.Errorf("invalid email sender or recipient login: %s, %s", input.Sender, input.Recipient)
	}

	err := es.EmailUseCase.CreateProfileEmail(input.Id, input.Sender, input.Recipient, ctx)
	if err != nil {
		return nil, fmt.Errorf("sender, recipient or id not found")
	}

	emailEmpty := new(proto.EmptyEmail)
	return emailEmpty, nil
}

func (es *EmailServer) CheckRecipientEmail(ctx context.Context, input *proto.Recipient) (*proto.EmptyEmail, error) {
	if input.Recipient == "" {
		return nil, fmt.Errorf("invalid recipient login: %s", input.Recipient)
	}

	err := es.EmailUseCase.CheckRecipientEmail(input.Recipient, ctx)
	if err != nil {
		return nil, fmt.Errorf("Recipient login not found")
	}

	emailEmpty := new(proto.EmptyEmail)
	return emailEmpty, nil
}
