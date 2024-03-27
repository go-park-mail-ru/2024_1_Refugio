// Code generated by MockGen. DO NOT EDIT.
// Source: ./isession_repo.go

// Package mock is a generated GoMock package.
package mock

import (
	models "mail/pkg/domain/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSessionRepository is a mock of SessionRepository interface.
type MockSessionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepositoryMockRecorder
}

// MockSessionRepositoryMockRecorder is the mock recorder for MockSessionRepository.
type MockSessionRepositoryMockRecorder struct {
	mock *MockSessionRepository
}

// NewMockSessionRepository creates a new mock instance.
func NewMockSessionRepository(ctrl *gomock.Controller) *MockSessionRepository {
	mock := &MockSessionRepository{ctrl: ctrl}
	mock.recorder = &MockSessionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRepository) EXPECT() *MockSessionRepositoryMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockSessionRepository) CreateSession(userID uint32, device string, lifeTime int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", userID, device, lifeTime)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockSessionRepositoryMockRecorder) CreateSession(userID, device, lifeTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSessionRepository)(nil).CreateSession), userID, device, lifeTime)
}

// DeleteExpiredSessions mocks base method.
func (m *MockSessionRepository) DeleteExpiredSessions() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExpiredSessions")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExpiredSessions indicates an expected call of DeleteExpiredSessions.
func (mr *MockSessionRepositoryMockRecorder) DeleteExpiredSessions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExpiredSessions", reflect.TypeOf((*MockSessionRepository)(nil).DeleteExpiredSessions))
}

// DeleteSessionByID mocks base method.
func (m *MockSessionRepository) DeleteSessionByID(sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSessionByID", sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSessionByID indicates an expected call of DeleteSessionByID.
func (mr *MockSessionRepositoryMockRecorder) DeleteSessionByID(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSessionByID", reflect.TypeOf((*MockSessionRepository)(nil).DeleteSessionByID), sessionID)
}

// GetSessionByID mocks base method.
func (m *MockSessionRepository) GetSessionByID(sessionID string) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionByID", sessionID)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionByID indicates an expected call of GetSessionByID.
func (mr *MockSessionRepositoryMockRecorder) GetSessionByID(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionByID", reflect.TypeOf((*MockSessionRepository)(nil).GetSessionByID), sessionID)
}