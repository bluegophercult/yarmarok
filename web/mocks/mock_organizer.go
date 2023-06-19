// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kaznasho/yarmarok/service (interfaces: OrganizerService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/kaznasho/yarmarok/service"
)

// MockOrganizerService is a mock of OrganizerService interface.
type MockOrganizerService struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizerServiceMockRecorder
}

// MockOrganizerServiceMockRecorder is the mock recorder for MockOrganizerService.
type MockOrganizerServiceMockRecorder struct {
	mock *MockOrganizerService
}

// NewMockOrganizerService creates a new mock instance.
func NewMockOrganizerService(ctrl *gomock.Controller) *MockOrganizerService {
	mock := &MockOrganizerService{ctrl: ctrl}
	mock.recorder = &MockOrganizerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizerService) EXPECT() *MockOrganizerServiceMockRecorder {
	return m.recorder
}

// InitOrganizerIfNotExists mocks base method.
func (m *MockOrganizerService) InitOrganizerIfNotExists(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitOrganizerIfNotExists", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitOrganizerIfNotExists indicates an expected call of InitOrganizerIfNotExists.
func (mr *MockOrganizerServiceMockRecorder) InitOrganizerIfNotExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitOrganizerIfNotExists", reflect.TypeOf((*MockOrganizerService)(nil).InitOrganizerIfNotExists), arg0)
}

// YarmarokService mocks base method.
func (m *MockOrganizerService) YarmarokService(arg0 string) service.YarmarokService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "YarmarokService", arg0)
	ret0, _ := ret[0].(service.YarmarokService)
	return ret0
}

// YarmarokService indicates an expected call of YarmarokService.
func (mr *MockOrganizerServiceMockRecorder) YarmarokService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "YarmarokService", reflect.TypeOf((*MockOrganizerService)(nil).YarmarokService), arg0)
}