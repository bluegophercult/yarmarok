// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kaznasho/yarmarok/service (interfaces: DonationStorage)

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDonationStorage is a mock of DonationStorage interface.
type MockDonationStorage struct {
	ctrl     *gomock.Controller
	recorder *MockDonationStorageMockRecorder
}

// MockDonationStorageMockRecorder is the mock recorder for MockDonationStorage.
type MockDonationStorageMockRecorder struct {
	mock *MockDonationStorage
}

// NewMockDonationStorage creates a new mock instance.
func NewMockDonationStorage(ctrl *gomock.Controller) *MockDonationStorage {
	mock := &MockDonationStorage{ctrl: ctrl}
	mock.recorder = &MockDonationStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDonationStorage) EXPECT() *MockDonationStorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockDonationStorage) Create(arg0 ParticipantStorage, arg1 PrizeStorage, arg2 *Donation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockDonationStorageMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDonationStorage)(nil).Create), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *MockDonationStorage) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDonationStorageMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDonationStorage)(nil).Delete), arg0)
}

// Get mocks base method.
func (m *MockDonationStorage) Get(arg0 string) (*Donation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*Donation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDonationStorageMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDonationStorage)(nil).Get), arg0)
}

// GetAll mocks base method.
func (m *MockDonationStorage) GetAll() ([]Donation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]Donation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockDonationStorageMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockDonationStorage)(nil).GetAll))
}

// Update mocks base method.
func (m *MockDonationStorage) Update(arg0 *Donation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockDonationStorageMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDonationStorage)(nil).Update), arg0)
}
