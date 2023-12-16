// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kaznasho/yarmarok/service (interfaces: OrganizerStorage)

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockOrganizerStorage is a mock of OrganizerStorage interface.
type MockOrganizerStorage struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizerStorageMockRecorder
}

// MockOrganizerStorageMockRecorder is the mock recorder for MockOrganizerStorage.
type MockOrganizerStorageMockRecorder struct {
	mock *MockOrganizerStorage
}

// NewMockOrganizerStorage creates a new mock instance.
func NewMockOrganizerStorage(ctrl *gomock.Controller) *MockOrganizerStorage {
	mock := &MockOrganizerStorage{ctrl: ctrl}
	mock.recorder = &MockOrganizerStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizerStorage) EXPECT() *MockOrganizerStorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockOrganizerStorage) Create(arg0 *Organizer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockOrganizerStorageMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOrganizerStorage)(nil).Create), arg0)
}

// Exists mocks base method.
func (m *MockOrganizerStorage) Exists(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockOrganizerStorageMockRecorder) Exists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockOrganizerStorage)(nil).Exists), arg0)
}

// RaffleStorage mocks base method.
func (m *MockOrganizerStorage) RaffleStorage(arg0 string) RaffleStorage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RaffleStorage", arg0)
	ret0, _ := ret[0].(RaffleStorage)
	return ret0
}

// RaffleStorage indicates an expected call of RaffleStorage.
func (mr *MockOrganizerStorageMockRecorder) RaffleStorage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RaffleStorage", reflect.TypeOf((*MockOrganizerStorage)(nil).RaffleStorage), arg0)
}
