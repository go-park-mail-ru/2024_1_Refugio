// Code generated by MockGen. DO NOT EDIT.
// Source: ./imanager.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	delivery_models "mail/internal/models/delivery_models"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSessionsManager is a mock of SessionsManager interface.
type MockSessionsManager struct {
	ctrl     *gomock.Controller
	recorder *MockSessionsManagerMockRecorder
}

// MockSessionsManagerMockRecorder is the mock recorder for MockSessionsManager.
type MockSessionsManagerMockRecorder struct {
	mock *MockSessionsManager
}

// NewMockSessionsManager creates a new mock instance.
func NewMockSessionsManager(ctrl *gomock.Controller) *MockSessionsManager {
	mock := &MockSessionsManager{ctrl: ctrl}
	mock.recorder = &MockSessionsManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionsManager) EXPECT() *MockSessionsManagerMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockSessionsManager) Check(r *http.Request, ctx context.Context) (*delivery_models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", r, ctx)
	ret0, _ := ret[0].(*delivery_models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockSessionsManagerMockRecorder) Check(r, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockSessionsManager)(nil).Check), r, ctx)
}

// CheckLogin mocks base method.
func (m *MockSessionsManager) CheckLogin(login string, r *http.Request, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckLogin", login, r, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckLogin indicates an expected call of CheckLogin.
func (mr *MockSessionsManagerMockRecorder) CheckLogin(login, r, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLogin", reflect.TypeOf((*MockSessionsManager)(nil).CheckLogin), login, r, ctx)
}

// Create mocks base method.
func (m *MockSessionsManager) Create(w http.ResponseWriter, userID uint32, ctx context.Context) (*delivery_models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", w, userID, ctx)
	ret0, _ := ret[0].(*delivery_models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionsManagerMockRecorder) Create(w, userID, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionsManager)(nil).Create), w, userID, ctx)
}

// DestroyCurrent mocks base method.
func (m *MockSessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyCurrent", w, r, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DestroyCurrent indicates an expected call of DestroyCurrent.
func (mr *MockSessionsManagerMockRecorder) DestroyCurrent(w, r, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyCurrent", reflect.TypeOf((*MockSessionsManager)(nil).DestroyCurrent), w, r, ctx)
}

// GetLoginBySession mocks base method.
func (m *MockSessionsManager) GetLoginBySession(r *http.Request, ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoginBySession", r, ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoginBySession indicates an expected call of GetLoginBySession.
func (mr *MockSessionsManagerMockRecorder) GetLoginBySession(r, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoginBySession", reflect.TypeOf((*MockSessionsManager)(nil).GetLoginBySession), r, ctx)
}

// GetProfileIDBySessionID mocks base method.
func (m *MockSessionsManager) GetProfileIDBySessionID(r *http.Request, ctx context.Context) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileIDBySessionID", r, ctx)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileIDBySessionID indicates an expected call of GetProfileIDBySessionID.
func (mr *MockSessionsManagerMockRecorder) GetProfileIDBySessionID(r, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileIDBySessionID", reflect.TypeOf((*MockSessionsManager)(nil).GetProfileIDBySessionID), r, ctx)
}

// GetSession mocks base method.
func (m *MockSessionsManager) GetSession(r *http.Request, ctx context.Context) *delivery_models.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", r, ctx)
	ret0, _ := ret[0].(*delivery_models.Session)
	return ret0
}

// GetSession indicates an expected call of GetSession.
func (mr *MockSessionsManagerMockRecorder) GetSession(r, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockSessionsManager)(nil).GetSession), r, ctx)
}

// SetSession mocks base method.
func (m *MockSessionsManager) SetSession(sessionId string, w http.ResponseWriter, r *http.Request, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSession", sessionId, w, r, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSession indicates an expected call of SetSession.
func (mr *MockSessionsManagerMockRecorder) SetSession(sessionId, w, r, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSession", reflect.TypeOf((*MockSessionsManager)(nil).SetSession), sessionId, w, r, ctx)
}
