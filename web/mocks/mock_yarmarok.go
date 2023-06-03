// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kaznasho/yarmarok/service (interfaces: YarmarokService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/kaznasho/yarmarok/service"
)

// MockYarmarokService is a mock of YarmarokService interface.
type MockYarmarokService struct {
	ctrl     *gomock.Controller
	recorder *MockYarmarokServiceMockRecorder
}

// MockYarmarokServiceMockRecorder is the mock recorder for MockYarmarokService.
type MockYarmarokServiceMockRecorder struct {
	mock *MockYarmarokService
}

// NewMockYarmarokService creates a new mock instance.
func NewMockYarmarokService(ctrl *gomock.Controller) *MockYarmarokService {
	mock := &MockYarmarokService{ctrl: ctrl}
	mock.recorder = &MockYarmarokServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockYarmarokService) EXPECT() *MockYarmarokServiceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockYarmarokService) Get(arg0 string) (*service.Yarmarok, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*service.Yarmarok)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockYarmarokServiceMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockYarmarokService)(nil).Get), arg0)
}

// Init mocks base method.
func (m *MockYarmarokService) Init(arg0 *service.YarmarokInitRequest) (*service.InitResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", arg0)
	ret0, _ := ret[0].(*service.InitResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Init indicates an expected call of Init.
func (mr *MockYarmarokServiceMockRecorder) Init(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockYarmarokService)(nil).Init), arg0)
}

// List mocks base method.
func (m *MockYarmarokService) List() ([]service.Yarmarok, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]service.Yarmarok)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockYarmarokServiceMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockYarmarokService)(nil).List))
}
