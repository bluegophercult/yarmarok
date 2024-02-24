// Code generated by MockGen. DO NOT EDIT.
// Source:  github.com/bluegophercult/yarmarok/service (interfaces: PrizeStorage)
//
// Generated by this command:
//
//	mockgen -destination=mock_prize_storage_test.go -package=service  github.com/bluegophercult/yarmarok/service PrizeStorage
//
// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockPrizeStorage is a mock of PrizeStorage interface.
type MockPrizeStorage struct {
	ctrl     *gomock.Controller
	recorder *MockPrizeStorageMockRecorder
}

// MockPrizeStorageMockRecorder is the mock recorder for MockPrizeStorage.
type MockPrizeStorageMockRecorder struct {
	mock *MockPrizeStorage
}

// NewMockPrizeStorage creates a new mock instance.
func NewMockPrizeStorage(ctrl *gomock.Controller) *MockPrizeStorage {
	mock := &MockPrizeStorage{ctrl: ctrl}
	mock.recorder = &MockPrizeStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrizeStorage) EXPECT() *MockPrizeStorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPrizeStorage) Create(arg0 *Prize) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockPrizeStorageMockRecorder) Create(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPrizeStorage)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockPrizeStorage) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPrizeStorageMockRecorder) Delete(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPrizeStorage)(nil).Delete), arg0)
}

// DonationStorage mocks base method.
func (m *MockPrizeStorage) DonationStorage(arg0 string) DonationStorage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DonationStorage", arg0)
	ret0, _ := ret[0].(DonationStorage)
	return ret0
}

// DonationStorage indicates an expected call of DonationStorage.
func (mr *MockPrizeStorageMockRecorder) DonationStorage(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DonationStorage", reflect.TypeOf((*MockPrizeStorage)(nil).DonationStorage), arg0)
}

// Get mocks base method.
func (m *MockPrizeStorage) Get(arg0 string) (*Prize, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*Prize)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPrizeStorageMockRecorder) Get(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPrizeStorage)(nil).Get), arg0)
}

// GetAll mocks base method.
func (m *MockPrizeStorage) GetAll() ([]Prize, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]Prize)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPrizeStorageMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPrizeStorage)(nil).GetAll))
}

// Update mocks base method.
func (m *MockPrizeStorage) Update(arg0 *Prize) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPrizeStorageMockRecorder) Update(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPrizeStorage)(nil).Update), arg0)
}
