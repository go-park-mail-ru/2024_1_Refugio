// Code generated by MockGen. DO NOT EDIT.
// Source: ./iemail_service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	domain_models "mail/internal/microservice/models/domain_models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEmailUseCase is a mock of EmailUseCase interface.
type MockEmailUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockEmailUseCaseMockRecorder
}

// MockEmailUseCaseMockRecorder is the mock recorder for MockEmailUseCase.
type MockEmailUseCaseMockRecorder struct {
	mock *MockEmailUseCase
}

// NewMockEmailUseCase creates a new mock instance.
func NewMockEmailUseCase(ctrl *gomock.Controller) *MockEmailUseCase {
	mock := &MockEmailUseCase{ctrl: ctrl}
	mock.recorder = &MockEmailUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailUseCase) EXPECT() *MockEmailUseCaseMockRecorder {
	return m.recorder
}

// AddAttachment mocks base method.
func (m *MockEmailUseCase) AddAttachment(fileID, fileType, fileName, fileSize string, emailID uint64, ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAttachment", fileID, fileType, fileName, fileSize, emailID, ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAttachment indicates an expected call of AddAttachment.
func (mr *MockEmailUseCaseMockRecorder) AddAttachment(fileID, fileType, fileName, fileSize, emailID, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAttachment", reflect.TypeOf((*MockEmailUseCase)(nil).AddAttachment), fileID, fileType, fileName, fileSize, emailID, ctx)
}

// AddFile mocks base method.
func (m *MockEmailUseCase) AddFile(fileID, fileType, fileName, fileSize string, ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFile", fileID, fileType, fileName, fileSize, ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddFile indicates an expected call of AddFile.
func (mr *MockEmailUseCaseMockRecorder) AddFile(fileID, fileType, fileName, fileSize, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFile", reflect.TypeOf((*MockEmailUseCase)(nil).AddFile), fileID, fileType, fileName, fileSize, ctx)
}

// AddFileToEmail mocks base method.
func (m *MockEmailUseCase) AddFileToEmail(emailID, fileID uint64, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFileToEmail", emailID, fileID, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFileToEmail indicates an expected call of AddFileToEmail.
func (mr *MockEmailUseCaseMockRecorder) AddFileToEmail(emailID, fileID, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFileToEmail", reflect.TypeOf((*MockEmailUseCase)(nil).AddFileToEmail), emailID, fileID, ctx)
}

// CheckRecipientEmail mocks base method.
func (m *MockEmailUseCase) CheckRecipientEmail(recipient string, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRecipientEmail", recipient, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckRecipientEmail indicates an expected call of CheckRecipientEmail.
func (mr *MockEmailUseCaseMockRecorder) CheckRecipientEmail(recipient, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRecipientEmail", reflect.TypeOf((*MockEmailUseCase)(nil).CheckRecipientEmail), recipient, ctx)
}

// CreateEmail mocks base method.
func (m *MockEmailUseCase) CreateEmail(newEmail *domain_models.Email, ctx context.Context) (uint64, *domain_models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEmail", newEmail, ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(*domain_models.Email)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateEmail indicates an expected call of CreateEmail.
func (mr *MockEmailUseCaseMockRecorder) CreateEmail(newEmail, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEmail", reflect.TypeOf((*MockEmailUseCase)(nil).CreateEmail), newEmail, ctx)
}

// CreateProfileEmail mocks base method.
func (m *MockEmailUseCase) CreateProfileEmail(emailId uint64, sender, recipient string, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProfileEmail", emailId, sender, recipient, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProfileEmail indicates an expected call of CreateProfileEmail.
func (mr *MockEmailUseCaseMockRecorder) CreateProfileEmail(emailId, sender, recipient, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProfileEmail", reflect.TypeOf((*MockEmailUseCase)(nil).CreateProfileEmail), emailId, sender, recipient, ctx)
}

// DeleteEmail mocks base method.
func (m *MockEmailUseCase) DeleteEmail(id uint64, login string, ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEmail", id, login, ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteEmail indicates an expected call of DeleteEmail.
func (mr *MockEmailUseCaseMockRecorder) DeleteEmail(id, login, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEmail", reflect.TypeOf((*MockEmailUseCase)(nil).DeleteEmail), id, login, ctx)
}

// DeleteFileByID mocks base method.
func (m *MockEmailUseCase) DeleteFileByID(fileID uint64, ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFileByID", fileID, ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFileByID indicates an expected call of DeleteFileByID.
func (mr *MockEmailUseCaseMockRecorder) DeleteFileByID(fileID, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFileByID", reflect.TypeOf((*MockEmailUseCase)(nil).DeleteFileByID), fileID, ctx)
}

// GetAllDraftEmails mocks base method.
func (m *MockEmailUseCase) GetAllDraftEmails(login string, offset, limit int64, ctx context.Context) ([]*domain_models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDraftEmails", login, offset, limit, ctx)
	ret0, _ := ret[0].([]*domain_models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDraftEmails indicates an expected call of GetAllDraftEmails.
func (mr *MockEmailUseCaseMockRecorder) GetAllDraftEmails(login, offset, limit, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDraftEmails", reflect.TypeOf((*MockEmailUseCase)(nil).GetAllDraftEmails), login, offset, limit, ctx)
}

// GetAllEmailsIncoming mocks base method.
func (m *MockEmailUseCase) GetAllEmailsIncoming(login string, offset, limit int64, ctx context.Context) ([]*domain_models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllEmailsIncoming", login, offset, limit, ctx)
	ret0, _ := ret[0].([]*domain_models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllEmailsIncoming indicates an expected call of GetAllEmailsIncoming.
func (mr *MockEmailUseCaseMockRecorder) GetAllEmailsIncoming(login, offset, limit, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllEmailsIncoming", reflect.TypeOf((*MockEmailUseCase)(nil).GetAllEmailsIncoming), login, offset, limit, ctx)
}

// GetAllEmailsSent mocks base method.
func (m *MockEmailUseCase) GetAllEmailsSent(login string, offset, limit int64, ctx context.Context) ([]*domain_models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllEmailsSent", login, offset, limit, ctx)
	ret0, _ := ret[0].([]*domain_models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllEmailsSent indicates an expected call of GetAllEmailsSent.
func (mr *MockEmailUseCaseMockRecorder) GetAllEmailsSent(login, offset, limit, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllEmailsSent", reflect.TypeOf((*MockEmailUseCase)(nil).GetAllEmailsSent), login, offset, limit, ctx)
}

// GetAllSpamEmails mocks base method.
func (m *MockEmailUseCase) GetAllSpamEmails(login string, offset, limit int64, ctx context.Context) ([]*domain_models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSpamEmails", login, offset, limit, ctx)
	ret0, _ := ret[0].([]*domain_models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSpamEmails indicates an expected call of GetAllSpamEmails.
func (mr *MockEmailUseCaseMockRecorder) GetAllSpamEmails(login, offset, limit, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSpamEmails", reflect.TypeOf((*MockEmailUseCase)(nil).GetAllSpamEmails), login, offset, limit, ctx)
}

// GetEmailByID mocks base method.
func (m *MockEmailUseCase) GetEmailByID(id uint64, login string, ctx context.Context) (*domain_models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmailByID", id, login, ctx)
	ret0, _ := ret[0].(*domain_models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmailByID indicates an expected call of GetEmailByID.
func (mr *MockEmailUseCaseMockRecorder) GetEmailByID(id, login, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmailByID", reflect.TypeOf((*MockEmailUseCase)(nil).GetEmailByID), id, login, ctx)
}

// GetFileByID mocks base method.
func (m *MockEmailUseCase) GetFileByID(fileID uint64, ctx context.Context) (*domain_models.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileByID", fileID, ctx)
	ret0, _ := ret[0].(*domain_models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileByID indicates an expected call of GetFileByID.
func (mr *MockEmailUseCaseMockRecorder) GetFileByID(fileID, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileByID", reflect.TypeOf((*MockEmailUseCase)(nil).GetFileByID), fileID, ctx)
}

// GetFilesByEmailID mocks base method.
func (m *MockEmailUseCase) GetFilesByEmailID(emailID uint64, ctx context.Context) ([]*domain_models.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilesByEmailID", emailID, ctx)
	ret0, _ := ret[0].([]*domain_models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilesByEmailID indicates an expected call of GetFilesByEmailID.
func (mr *MockEmailUseCaseMockRecorder) GetFilesByEmailID(emailID, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilesByEmailID", reflect.TypeOf((*MockEmailUseCase)(nil).GetFilesByEmailID), emailID, ctx)
}

// UpdateEmail mocks base method.
func (m *MockEmailUseCase) UpdateEmail(updatedEmail *domain_models.Email, ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmail", updatedEmail, ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEmail indicates an expected call of UpdateEmail.
func (mr *MockEmailUseCaseMockRecorder) UpdateEmail(updatedEmail, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmail", reflect.TypeOf((*MockEmailUseCase)(nil).UpdateEmail), updatedEmail, ctx)
}

// UpdateFileByID mocks base method.
func (m *MockEmailUseCase) UpdateFileByID(fileID uint64, newFileID, newFileType, newFileName, newFileSize string, ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFileByID", fileID, newFileID, newFileType, newFileName, newFileSize, ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFileByID indicates an expected call of UpdateFileByID.
func (mr *MockEmailUseCaseMockRecorder) UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFileByID", reflect.TypeOf((*MockEmailUseCase)(nil).UpdateFileByID), fileID, newFileID, newFileType, newFileName, newFileSize, ctx)
}
