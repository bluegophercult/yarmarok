// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kaznasho/yarmarok/service (interfaces: PrizeService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/kaznasho/yarmarok/service"
)

// MockPrizeService is a mock of PrizeService interface.
type MockPrizeService struct {
	ctrl     *gomock.Controller
	recorder *MockPrizeServiceMockRecorder
}

// MockPrizeServiceMockRecorder is the mock recorder for MockPrizeService.
type MockPrizeServiceMockRecorder struct {
	mock *MockPrizeService
}

// NewMockPrizeService creates a new mock instance.
func NewMockPrizeService(ctrl *gomock.Controller) *MockPrizeService {
	mock := &MockPrizeService{ctrl: ctrl}
	mock.recorder = &MockPrizeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrizeService) EXPECT() *MockPrizeServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPrizeService) Create(arg0 *service.PrizeRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPrizeServiceMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPrizeService)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockPrizeService) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPrizeServiceMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPrizeService)(nil).Delete), arg0)
}

// Edit mocks base method.
func (m *MockPrizeService) Edit(arg0 string, arg1 *service.PrizeRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Edit", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Edit indicates an expected call of Edit.
func (mr *MockPrizeServiceMockRecorder) Edit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockPrizeService)(nil).Edit), arg0, arg1)
}

// Get mocks base method.
func (m *MockPrizeService) Get(arg0 string) (*service.Prize, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*service.Prize)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPrizeServiceMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPrizeService)(nil).Get), arg0)
}

// List mocks base method.
func (m *MockPrizeService) List() ([]service.Prize, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]service.Prize)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockPrizeServiceMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPrizeService)(nil).List))
}
