// Code generated by MockGen. DO NOT EDIT.
// Source: ./server/internal/application/ports.go

// Package application is a generated GoMock package.
package application

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	gomoney "github.com/lordvidex/gomoney/pkg/gomoney"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(ctx context.Context, arg CreateUserArg) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, arg)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), ctx, arg)
}

// GetUserByPhone mocks base method.
func (m *MockUserRepository) GetUserByPhone(ctx context.Context, phone string) (gomoney.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByPhone", ctx, phone)
	ret0, _ := ret[0].(gomoney.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByPhone indicates an expected call of GetUserByPhone.
func (mr *MockUserRepositoryMockRecorder) GetUserByPhone(ctx, phone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByPhone", reflect.TypeOf((*MockUserRepository)(nil).GetUserByPhone), ctx, phone)
}

// MockAccountRepository is a mock of AccountRepository interface.
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder is the mock recorder for MockAccountRepository.
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock instance.
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockAccountRepository) CreateAccount(ctx context.Context, arg CreateAccountArg) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, arg)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockAccountRepositoryMockRecorder) CreateAccount(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountRepository)(nil).CreateAccount), ctx, arg)
}

// DeleteAccount mocks base method.
func (m *MockAccountRepository) DeleteAccount(ctx context.Context, accountID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", ctx, accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockAccountRepositoryMockRecorder) DeleteAccount(ctx, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockAccountRepository)(nil).DeleteAccount), ctx, accountID)
}

// GetAccountByID mocks base method.
func (m *MockAccountRepository) GetAccountByID(ctx context.Context, accountID int64) (*gomoney.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountByID", ctx, accountID)
	ret0, _ := ret[0].(*gomoney.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountByID indicates an expected call of GetAccountByID.
func (mr *MockAccountRepositoryMockRecorder) GetAccountByID(ctx, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountByID", reflect.TypeOf((*MockAccountRepository)(nil).GetAccountByID), ctx, accountID)
}

// GetAccountsForUser mocks base method.
func (m *MockAccountRepository) GetAccountsForUser(ctx context.Context, userID uuid.UUID) ([]gomoney.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountsForUser", ctx, userID)
	ret0, _ := ret[0].([]gomoney.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountsForUser indicates an expected call of GetAccountsForUser.
func (mr *MockAccountRepositoryMockRecorder) GetAccountsForUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountsForUser", reflect.TypeOf((*MockAccountRepository)(nil).GetAccountsForUser), ctx, userID)
}

// GetAllTransactions mocks base method.
func (m *MockAccountRepository) GetAllTransactions(ctx context.Context, accountID int64) ([]gomoney.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTransactions", ctx, accountID)
	ret0, _ := ret[0].([]gomoney.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTransactions indicates an expected call of GetAllTransactions.
func (mr *MockAccountRepositoryMockRecorder) GetAllTransactions(ctx, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTransactions", reflect.TypeOf((*MockAccountRepository)(nil).GetAllTransactions), ctx, accountID)
}

// GetLastNTransactions mocks base method.
func (m *MockAccountRepository) GetLastNTransactions(ctx context.Context, accountID int64, n int) ([]gomoney.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastNTransactions", ctx, accountID, n)
	ret0, _ := ret[0].([]gomoney.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastNTransactions indicates an expected call of GetLastNTransactions.
func (mr *MockAccountRepositoryMockRecorder) GetLastNTransactions(ctx, accountID, n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastNTransactions", reflect.TypeOf((*MockAccountRepository)(nil).GetLastNTransactions), ctx, accountID, n)
}

// Transfer mocks base method.
func (m *MockAccountRepository) Transfer(ctx context.Context, tx *gomoney.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transfer", ctx, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transfer indicates an expected call of Transfer.
func (mr *MockAccountRepositoryMockRecorder) Transfer(ctx, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transfer", reflect.TypeOf((*MockAccountRepository)(nil).Transfer), ctx, tx)
}

// MockTxLocker is a mock of TxLocker interface.
type MockTxLocker struct {
	ctrl     *gomock.Controller
	recorder *MockTxLockerMockRecorder
}

// MockTxLockerMockRecorder is the mock recorder for MockTxLocker.
type MockTxLockerMockRecorder struct {
	mock *MockTxLocker
}

// NewMockTxLocker creates a new mock instance.
func NewMockTxLocker(ctrl *gomock.Controller) *MockTxLocker {
	mock := &MockTxLocker{ctrl: ctrl}
	mock.recorder = &MockTxLockerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxLocker) EXPECT() *MockTxLockerMockRecorder {
	return m.recorder
}

// Lock mocks base method.
func (m *MockTxLocker) Lock(x any, y ...any) func() {
	m.ctrl.T.Helper()
	varargs := []interface{}{x}
	for _, a := range y {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Lock", varargs...)
	ret0, _ := ret[0].(func())
	return ret0
}

// Lock indicates an expected call of Lock.
func (mr *MockTxLockerMockRecorder) Lock(x interface{}, y ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{x}, y...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lock", reflect.TypeOf((*MockTxLocker)(nil).Lock), varargs...)
}
