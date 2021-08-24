// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/RasaHQ/rasactl/pkg/docker (interfaces: Interface)

// Package fake is a generated GoMock package.
package fake

import (
	reflect "reflect"

	container "github.com/docker/docker/api/types/container"
	gomock "github.com/golang/mock/gomock"

	docker "github.com/RasaHQ/rasactl/pkg/docker"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// CreateKindNode mocks base method.
func (m *MockInterface) CreateKindNode(arg0 string) (container.ContainerCreateCreatedBody, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateKindNode", arg0)
	ret0, _ := ret[0].(container.ContainerCreateCreatedBody)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateKindNode indicates an expected call of CreateKindNode.
func (mr *MockInterfaceMockRecorder) CreateKindNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateKindNode", reflect.TypeOf((*MockInterface)(nil).CreateKindNode), arg0)
}

// DeleteKindNode mocks base method.
func (m *MockInterface) DeleteKindNode(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteKindNode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteKindNode indicates an expected call of DeleteKindNode.
func (mr *MockInterfaceMockRecorder) DeleteKindNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteKindNode", reflect.TypeOf((*MockInterface)(nil).DeleteKindNode), arg0)
}

// GetKind mocks base method.
func (m *MockInterface) GetKind() docker.KindSpec {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKind")
	ret0, _ := ret[0].(docker.KindSpec)
	return ret0
}

// GetKind indicates an expected call of GetKind.
func (mr *MockInterfaceMockRecorder) GetKind() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKind", reflect.TypeOf((*MockInterface)(nil).GetKind))
}

// SetKind mocks base method.
func (m *MockInterface) SetKind(arg0 docker.KindSpec) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetKind", arg0)
}

// SetKind indicates an expected call of SetKind.
func (mr *MockInterfaceMockRecorder) SetKind(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetKind", reflect.TypeOf((*MockInterface)(nil).SetKind), arg0)
}

// SetNamespace mocks base method.
func (m *MockInterface) SetNamespace(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNamespace", arg0)
}

// SetNamespace indicates an expected call of SetNamespace.
func (mr *MockInterfaceMockRecorder) SetNamespace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNamespace", reflect.TypeOf((*MockInterface)(nil).SetNamespace), arg0)
}

// SetProjectPath mocks base method.
func (m *MockInterface) SetProjectPath(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetProjectPath", arg0)
}

// SetProjectPath indicates an expected call of SetProjectPath.
func (mr *MockInterfaceMockRecorder) SetProjectPath(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProjectPath", reflect.TypeOf((*MockInterface)(nil).SetProjectPath), arg0)
}

// StartKindNode mocks base method.
func (m *MockInterface) StartKindNode(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartKindNode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartKindNode indicates an expected call of StartKindNode.
func (mr *MockInterfaceMockRecorder) StartKindNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartKindNode", reflect.TypeOf((*MockInterface)(nil).StartKindNode), arg0)
}

// StopKindNode mocks base method.
func (m *MockInterface) StopKindNode(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopKindNode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopKindNode indicates an expected call of StopKindNode.
func (mr *MockInterfaceMockRecorder) StopKindNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopKindNode", reflect.TypeOf((*MockInterface)(nil).StopKindNode), arg0)
}