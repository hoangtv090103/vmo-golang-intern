// Code generated by MockGen. DO NOT EDIT.
// Source: internal/auth/repository/account_repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/auth/repository/account_repository.go -destination=internal/auth/mocks/mock_account_repository.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	entity "ecommerce/internal/auth/entity"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockIAccountRepository is a mock of IAccountRepository interface.
type MockIAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAccountRepositoryMockRecorder
}

// MockIAccountRepositoryMockRecorder is the mock recorder for MockIAccountRepository.
type MockIAccountRepositoryMockRecorder struct {
	mock *MockIAccountRepository
}

// NewMockIAccountRepository creates a new mock instance.
func NewMockIAccountRepository(ctrl *gomock.Controller) *MockIAccountRepository {
	mock := &MockIAccountRepository{ctrl: ctrl}
	mock.recorder = &MockIAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAccountRepository) EXPECT() *MockIAccountRepositoryMockRecorder {
	return m.recorder
}

// GetByUsername mocks base method.
func (m *MockIAccountRepository) GetByUsername(ctx context.Context, username string) (*entity.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsername", ctx, username)
	ret0, _ := ret[0].(*entity.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUsername indicates an expected call of GetByUsername.
func (mr *MockIAccountRepositoryMockRecorder) GetByUsername(ctx, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsername", reflect.TypeOf((*MockIAccountRepository)(nil).GetByUsername), ctx, username)
}

// Login mocks base method.
func (m *MockIAccountRepository) Login(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, account)
	ret0, _ := ret[0].(*entity.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockIAccountRepositoryMockRecorder) Login(ctx, account any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockIAccountRepository)(nil).Login), ctx, account)
}

// Register mocks base method.
func (m *MockIAccountRepository) Register(ctx context.Context, account *entity.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockIAccountRepositoryMockRecorder) Register(ctx, account any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockIAccountRepository)(nil).Register), ctx, account)
}
