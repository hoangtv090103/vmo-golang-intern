// Code generated by MockGen. DO NOT EDIT.
// Source: internal/order/repository/order_repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/order/repository/order_repository.go -destination=internal/order/mocks/mock_order_repository.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	entity "ecommerce/internal/order/entity"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockIOrderRepository is a mock of IOrderRepository interface.
type MockIOrderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIOrderRepositoryMockRecorder
}

// MockIOrderRepositoryMockRecorder is the mock recorder for MockIOrderRepository.
type MockIOrderRepositoryMockRecorder struct {
	mock *MockIOrderRepository
}

// NewMockIOrderRepository creates a new mock instance.
func NewMockIOrderRepository(ctrl *gomock.Controller) *MockIOrderRepository {
	mock := &MockIOrderRepository{ctrl: ctrl}
	mock.recorder = &MockIOrderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIOrderRepository) EXPECT() *MockIOrderRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIOrderRepositoryMockRecorder) Create(ctx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIOrderRepository)(nil).Create), ctx, order)
}

// Delete mocks base method.
func (m *MockIOrderRepository) Delete(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIOrderRepositoryMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIOrderRepository)(nil).Delete), ctx, id)
}

// GetAll mocks base method.
func (m *MockIOrderRepository) GetAll(ctx context.Context) ([]*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockIOrderRepositoryMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockIOrderRepository)(nil).GetAll), ctx)
}

// GetByID mocks base method.
func (m *MockIOrderRepository) GetByID(ctx context.Context, id int) (*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIOrderRepositoryMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIOrderRepository)(nil).GetByID), ctx, id)
}

// GetInvoice mocks base method.
func (m *MockIOrderRepository) GetInvoice(ctx context.Context, orderID int) ([]*entity.InvoiceData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoice", ctx, orderID)
	ret0, _ := ret[0].([]*entity.InvoiceData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvoice indicates an expected call of GetInvoice.
func (mr *MockIOrderRepositoryMockRecorder) GetInvoice(ctx, orderID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoice", reflect.TypeOf((*MockIOrderRepository)(nil).GetInvoice), ctx, orderID)
}

// GetUserOrders mocks base method.
func (m *MockIOrderRepository) GetUserOrders(ctx context.Context, username string) ([]*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserOrders", ctx, username)
	ret0, _ := ret[0].([]*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserOrders indicates an expected call of GetUserOrders.
func (mr *MockIOrderRepositoryMockRecorder) GetUserOrders(ctx, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserOrders", reflect.TypeOf((*MockIOrderRepository)(nil).GetUserOrders), ctx, username)
}

// Update mocks base method.
func (m *MockIOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIOrderRepositoryMockRecorder) Update(ctx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIOrderRepository)(nil).Update), ctx, order)
}
