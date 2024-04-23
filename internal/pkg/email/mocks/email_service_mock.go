// Code generated by MockGen. DO NOT EDIT.
// Source: iemail_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	"mail/internal/microservice/models/domain_models"
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
func (m *MockEmailUseCase) CreateEmail(newEmail *domain_models.Email, ctx context.Context) (int64, *domain_models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEmail", newEmail, ctx)
	ret0, _ := ret[0].(int64)
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
func (m *MockEmailUseCase) CreateProfileEmail(email_id int64, sender, recipient string, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProfileEmail", email_id, sender, recipient, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProfileEmail indicates an expected call of CreateProfileEmail.
func (mr *MockEmailUseCaseMockRecorder) CreateProfileEmail(email_id, sender, recipient, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProfileEmail", reflect.TypeOf((*MockEmailUseCase)(nil).CreateProfileEmail), email_id, sender, recipient, ctx)
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

// GetAllEmailsIncoming mocks base method.
func (m *MockEmailUseCase) GetAllEmailsIncoming(login string, offset, limit int, ctx context.Context) ([]*domain_models.Email, error) {
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
func (m *MockEmailUseCase) GetAllEmailsSent(login string, offset, limit int, ctx context.Context) ([]*domain_models.Email, error) {
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