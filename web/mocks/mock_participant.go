// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kaznasho/yarmarok/service (interfaces: ParticipantService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/kaznasho/yarmarok/service"
)

// MockParticipantService is a mock of ParticipantService interface.
type MockParticipantService struct {
	ctrl     *gomock.Controller
	recorder *MockParticipantServiceMockRecorder
}

// MockParticipantServiceMockRecorder is the mock recorder for MockParticipantService.
type MockParticipantServiceMockRecorder struct {
	mock *MockParticipantService
}

// NewMockParticipantService creates a new mock instance.
func NewMockParticipantService(ctrl *gomock.Controller) *MockParticipantService {
	mock := &MockParticipantService{ctrl: ctrl}
	mock.recorder = &MockParticipantServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParticipantService) EXPECT() *MockParticipantServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockParticipantService) Create(arg0 *service.ParticipantRequest) (*service.CreateResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*service.CreateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockParticipantServiceMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockParticipantService)(nil).Create), arg0)
}

// Edit mocks base method.
func (m *MockParticipantService) Edit(arg0 *service.ParticipantEditRequest) (*service.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Edit", arg0)
	ret0, _ := ret[0].(*service.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Edit indicates an expected call of Edit.
func (mr *MockParticipantServiceMockRecorder) Edit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockParticipantService)(nil).Edit), arg0)
}

// List mocks base method.
func (m *MockParticipantService) List() (*service.ParticipantListResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].(*service.ParticipantListResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockParticipantServiceMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockParticipantService)(nil).List))
}
