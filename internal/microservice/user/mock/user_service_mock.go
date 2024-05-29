// Code generated by MockGen. DO NOT EDIT.
// Source: ./iuser_service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	domain_models "mail/internal/microservice/models/domain_models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserUseCase is a mock of UserUseCase interface.
type MockUserUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUseCaseMockRecorder
}

// MockUserUseCaseMockRecorder is the mock recorder for MockUserUseCase.
type MockUserUseCaseMockRecorder struct {
	mock *MockUserUseCase
}

// NewMockUserUseCase creates a new mock instance.
func NewMockUserUseCase(ctrl *gomock.Controller) *MockUserUseCase {
	mock := &MockUserUseCase{ctrl: ctrl}
	mock.recorder = &MockUserUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUseCase) EXPECT() *MockUserUseCaseMockRecorder {
	return m.recorder
}

// AddAvatar mocks base method.
func (m *MockUserUseCase) AddAvatar(id uint32, fileID string, ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAvatar", id, fileID, ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAvatar indicates an expected call of AddAvatar.
func (mr *MockUserUseCaseMockRecorder) AddAvatar(id, fileID, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAvatar", reflect.TypeOf((*MockUserUseCase)(nil).AddAvatar), id, fileID, ctx)
}

// CreateUser mocks base method.
func (m *MockUserUseCase) CreateUser(user *domain_models.User, ctx context.Context) (*domain_models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user, ctx)
	ret0, _ := ret[0].(*domain_models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserUseCaseMockRecorder) CreateUser(user, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserUseCase)(nil).CreateUser), user, ctx)
}

// DeleteAvatarByUserID mocks base method.
func (m *MockUserUseCase) DeleteAvatarByUserID(userID uint32, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAvatarByUserID", userID, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAvatarByUserID indicates an expected call of DeleteAvatarByUserID.
func (mr *MockUserUseCaseMockRecorder) DeleteAvatarByUserID(userID, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAvatarByUserID", reflect.TypeOf((*MockUserUseCase)(nil).DeleteAvatarByUserID), userID, ctx)
}

// DeleteUserByID mocks base method.
func (m *MockUserUseCase) DeleteUserByID(id uint32, ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserByID", id, ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUserByID indicates an expected call of DeleteUserByID.
func (mr *MockUserUseCaseMockRecorder) DeleteUserByID(id, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserByID", reflect.TypeOf((*MockUserUseCase)(nil).DeleteUserByID), id, ctx)
}

// GetAllUsers mocks base method.
func (m *MockUserUseCase) GetAllUsers(ctx context.Context) ([]*domain_models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", ctx)
	ret0, _ := ret[0].([]*domain_models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockUserUseCaseMockRecorder) GetAllUsers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockUserUseCase)(nil).GetAllUsers), ctx)
}

// GetUserByID mocks base method.
func (m *MockUserUseCase) GetUserByID(id uint32, ctx context.Context) (*domain_models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", id, ctx)
	ret0, _ := ret[0].(*domain_models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserUseCaseMockRecorder) GetUserByID(id, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserUseCase)(nil).GetUserByID), id, ctx)
}

// GetUserByLogin mocks base method.
func (m *MockUserUseCase) GetUserByLogin(login, password string, ctx context.Context) (*domain_models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", login, password, ctx)
	ret0, _ := ret[0].(*domain_models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockUserUseCaseMockRecorder) GetUserByLogin(login, password, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockUserUseCase)(nil).GetUserByLogin), login, password, ctx)
}

// GetUserByOnlyLogin mocks base method.
func (m *MockUserUseCase) GetUserByOnlyLogin(login string, ctx context.Context) (*domain_models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByOnlyLogin", login, ctx)
	ret0, _ := ret[0].(*domain_models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByOnlyLogin indicates an expected call of GetUserByOnlyLogin.
func (mr *MockUserUseCaseMockRecorder) GetUserByOnlyLogin(login, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByOnlyLogin", reflect.TypeOf((*MockUserUseCase)(nil).GetUserByOnlyLogin), login, ctx)
}

// GetUserVkID mocks base method.
func (m *MockUserUseCase) GetUserVkID(vkId uint32, ctx context.Context) (*domain_models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserVkID", vkId, ctx)
	ret0, _ := ret[0].(*domain_models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserVkID indicates an expected call of GetUserVkID.
func (mr *MockUserUseCaseMockRecorder) GetUserVkID(vkId, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserVkID", reflect.TypeOf((*MockUserUseCase)(nil).GetUserVkID), vkId, ctx)
}

// IsLoginUnique mocks base method.
func (m *MockUserUseCase) IsLoginUnique(login string, ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLoginUnique", login, ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLoginUnique indicates an expected call of IsLoginUnique.
func (mr *MockUserUseCaseMockRecorder) IsLoginUnique(login, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLoginUnique", reflect.TypeOf((*MockUserUseCase)(nil).IsLoginUnique), login, ctx)
}

// UpdateUser mocks base method.
func (m *MockUserUseCase) UpdateUser(userNew *domain_models.User, ctx context.Context) (*domain_models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", userNew, ctx)
	ret0, _ := ret[0].(*domain_models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserUseCaseMockRecorder) UpdateUser(userNew, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserUseCase)(nil).UpdateUser), userNew, ctx)
}
