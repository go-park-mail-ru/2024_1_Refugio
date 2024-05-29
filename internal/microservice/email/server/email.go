package server

import (
	"context"
	"fmt"
	"strconv"

	"mail/internal/microservice/email/proto"
	"mail/internal/pkg/utils/validators"

	usecase "mail/internal/microservice/email/interface"
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

	return converters.EmailConvertCoreInProto(email), nil
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
		emailsProto[i] = converters.EmailConvertCoreInProto(e)
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
		emailsProto[i] = converters.EmailConvertCoreInProto(e)
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
		emailsProto[i] = converters.EmailConvertCoreInProto(e)
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
		emailsProto[i] = converters.EmailConvertCoreInProto(e)
	}

	emailProto := new(proto.Emails)
	emailProto.Emails = emailsProto
	return emailProto, nil
}

func (es *EmailServer) CreateEmail(ctx context.Context, input *proto.Email) (*proto.EmailWithID, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid email format: %s", input)
	}

	id, email, err := es.EmailUseCase.CreateEmail(converters.EmailConvertProtoInCore(input), ctx)
	if err != nil {
		return nil, fmt.Errorf("failed create email")
	}

	emailWithId := new(proto.EmailWithID)
	emailWithId.Id = id
	emailWithId.Email = converters.EmailConvertCoreInProto(email)
	return emailWithId, nil
}

func (es *EmailServer) UpdateEmail(ctx context.Context, input *proto.Email) (*proto.StatusEmail, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid email format: %s", input)
	}

	okStatus, err := es.EmailUseCase.UpdateEmail(converters.EmailConvertProtoInCore(input), ctx)

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
		return nil, fmt.Errorf("recipient login not found")
	}

	emailEmpty := new(proto.EmptyEmail)
	return emailEmpty, nil
}

func (es *EmailServer) AddEmailDraft(ctx context.Context, input *proto.Email) (*proto.EmailWithID, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid email format: %s", input)
	}

	id, email, err := es.EmailUseCase.CreateEmail(converters.EmailConvertProtoInCore(input), ctx)
	if err != nil {
		return nil, fmt.Errorf("failed create email")
	}

	err = es.EmailUseCase.CreateProfileEmail(id, email.SenderEmail, "", ctx)
	if err != nil {
		return nil, fmt.Errorf("failed create profile email")
	}

	emailWithId := new(proto.EmailWithID)
	emailWithId.Id = id
	emailWithId.Email = converters.EmailConvertCoreInProto(email)
	return emailWithId, nil
}

func (es *EmailServer) AddAttachment(ctx context.Context, input *proto.AddAttachmentRequest) (*proto.AddAttachmentReply, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid file format: %s", input)
	}

	if validators.IsEmpty(input.FileId) || validators.IsEmpty(input.FileType) || validators.IsEmpty(input.FileName) || validators.IsEmpty(input.FileSize) {
		return nil, fmt.Errorf("file id or file type or file name or file size is empty")
	}

	if input.EmailId <= 0 {
		return nil, fmt.Errorf("invalid email id")
	}

	fileID, err := es.EmailUseCase.AddAttachment(input.FileId, input.FileType, input.FileName, input.FileSize, input.EmailId, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed add attachment")
	}

	return &proto.AddAttachmentReply{FileId: fileID}, nil
}

func (es *EmailServer) GetFileByID(ctx context.Context, input *proto.GetFileByIDRequest) (*proto.GetFileByIDReply, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid file format: %s", input)
	}

	if input.FileId <= 0 {
		return nil, fmt.Errorf("invalid file id")
	}

	file, err := es.EmailUseCase.GetFileByID(input.FileId, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get file by id")
	}

	return &proto.GetFileByIDReply{File: converters.FileConvertCoreInProto(file)}, nil
}

func (es *EmailServer) GetFilesByEmailID(ctx context.Context, input *proto.GetFilesByEmailIDRequest) (*proto.GetFilesByEmailIDReply, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid file format: %s", input)
	}

	if input.EmailId <= 0 {
		return nil, fmt.Errorf("invalid email id")
	}

	files, err := es.EmailUseCase.GetFilesByEmailID(input.EmailId, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get files by email id")
	}

	filesProto := make([]*proto.File, 0, len(files))
	for _, file := range files {
		filesProto = append(filesProto, converters.FileConvertCoreInProto(file))
	}

	return &proto.GetFilesByEmailIDReply{Files: filesProto}, nil
}

func (es *EmailServer) DeleteFileByID(ctx context.Context, input *proto.DeleteFileByIDRequest) (*proto.DeleteFileByIDReply, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid file format: %s", input)
	}

	if input.FileId <= 0 {
		return nil, fmt.Errorf("invalid file id")
	}

	deleted, err := es.EmailUseCase.DeleteFileByID(input.FileId, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed delete file by id")
	}
	if !deleted {
		return nil, fmt.Errorf("file not deleted")
	}

	return &proto.DeleteFileByIDReply{Status: deleted}, nil
}

func (es *EmailServer) UpdateFileByID(ctx context.Context, input *proto.UpdateFileByIDRequest) (*proto.UpdateFileByIDReply, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid file format: %s", input)
	}

	if validators.IsEmpty(input.NewFileId) || validators.IsEmpty(input.NewFileType) || validators.IsEmpty(input.NewFileName) || validators.IsEmpty(input.NewFileSize) {
		return nil, fmt.Errorf("file id or file type or file name or file size is empty")
	}

	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid file id")
	}

	updated, err := es.EmailUseCase.UpdateFileByID(input.Id, input.NewFileId, input.NewFileType, input.NewFileName, input.NewFileSize, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed update file by id")
	}
	if !updated {
		return nil, fmt.Errorf("file not updated")
	}

	return &proto.UpdateFileByIDReply{Status: updated}, nil
}

func (es *EmailServer) AddFile(ctx context.Context, input *proto.AddFileRequest) (*proto.AddFileReply, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid file format: %s", input)
	}

	if validators.IsEmpty(input.FileId) || validators.IsEmpty(input.FileType) || validators.IsEmpty(input.FileName) || validators.IsEmpty(input.FileSize) {
		return nil, fmt.Errorf("file id or file type or file name or file size is empty")
	}

	fileId, err := es.EmailUseCase.AddFile(input.FileId, input.FileType, input.FileName, input.FileSize, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed add attachment")
	}

	return &proto.AddFileReply{FileId: fileId}, nil
}

func (es *EmailServer) AddFileToEmail(ctx context.Context, input *proto.AddFileToEmailRequest) (*proto.AddFileToEmailReply, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid file format: %s", input)
	}

	if input.EmailId <= 0 || input.FileId <= 0 {
		return nil, fmt.Errorf("invalid file id")
	}

	err := es.EmailUseCase.AddFileToEmail(input.EmailId, input.FileId, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed add attachment")
	}

	return &proto.AddFileToEmailReply{Status: true}, nil
}
