// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kaznasho/yarmarok/service (interfaces: RaffleStorage)

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRaffleStorage is a mock of RaffleStorage interface.
type MockRaffleStorage struct {
	ctrl     *gomock.Controller
	recorder *MockRaffleStorageMockRecorder
}

// MockRaffleStorageMockRecorder is the mock recorder for MockRaffleStorage.
type MockRaffleStorageMockRecorder struct {
	mock *MockRaffleStorage
}

// NewMockRaffleStorage creates a new mock instance.
func NewMockRaffleStorage(ctrl *gomock.Controller) *MockRaffleStorage {
	mock := &MockRaffleStorage{ctrl: ctrl}
	mock.recorder = &MockRaffleStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRaffleStorage) EXPECT() *MockRaffleStorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRaffleStorage) Create(arg0 *Raffle) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRaffleStorageMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRaffleStorage)(nil).Create), arg0)
}

// Get mocks base method.
func (m *MockRaffleStorage) Get(arg0 string) (*Raffle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*Raffle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRaffleStorageMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRaffleStorage)(nil).Get), arg0)
}

// GetAll mocks base method.
func (m *MockRaffleStorage) GetAll() ([]Raffle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]Raffle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockRaffleStorageMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockRaffleStorage)(nil).GetAll))
}

// ParticipantStorage mocks base method.
func (m *MockRaffleStorage) ParticipantStorage(arg0 string) ParticipantStorage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParticipantStorage", arg0)
	ret0, _ := ret[0].(ParticipantStorage)
	return ret0
}

// ParticipantStorage indicates an expected call of ParticipantStorage.
func (mr *MockRaffleStorageMockRecorder) ParticipantStorage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParticipantStorage", reflect.TypeOf((*MockRaffleStorage)(nil).ParticipantStorage), arg0)
}

// PrizeStorage mocks base method.
func (m *MockRaffleStorage) PrizeStorage(arg0 string) PrizeStorage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrizeStorage", arg0)
	ret0, _ := ret[0].(PrizeStorage)
	return ret0
}

// PrizeStorage indicates an expected call of PrizeStorage.
func (mr *MockRaffleStorageMockRecorder) PrizeStorage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrizeStorage", reflect.TypeOf((*MockRaffleStorage)(nil).PrizeStorage), arg0)
}

// Query mocks base method.
func (m *MockRaffleStorage) Query(arg0 Query) ([]Raffle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0)
	ret0, _ := ret[0].([]Raffle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockRaffleStorageMockRecorder) Query(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockRaffleStorage)(nil).Query), arg0)
}

// Update mocks base method.
func (m *MockRaffleStorage) Update(arg0 *Raffle) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRaffleStorageMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRaffleStorage)(nil).Update), arg0)
}
