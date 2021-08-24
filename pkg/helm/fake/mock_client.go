// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/RasaHQ/rasactl/pkg/helm (interfaces: Interface)

// Package fake is a generated GoMock package.
package fake

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	release "helm.sh/helm/v3/pkg/release"

	types "github.com/RasaHQ/rasactl/pkg/types"
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

// GetAllValues mocks base method.
func (m *MockInterface) GetAllValues() (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllValues")
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllValues indicates an expected call of GetAllValues.
func (mr *MockInterfaceMockRecorder) GetAllValues() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllValues", reflect.TypeOf((*MockInterface)(nil).GetAllValues))
}

// GetConfiguration mocks base method.
func (m *MockInterface) GetConfiguration() *types.HelmConfigurationSpec {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfiguration")
	ret0, _ := ret[0].(*types.HelmConfigurationSpec)
	return ret0
}

// GetConfiguration indicates an expected call of GetConfiguration.
func (mr *MockInterfaceMockRecorder) GetConfiguration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfiguration", reflect.TypeOf((*MockInterface)(nil).GetConfiguration))
}

// GetNamespace mocks base method.
func (m *MockInterface) GetNamespace() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespace")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetNamespace indicates an expected call of GetNamespace.
func (mr *MockInterfaceMockRecorder) GetNamespace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespace", reflect.TypeOf((*MockInterface)(nil).GetNamespace))
}

// GetStatus mocks base method.
func (m *MockInterface) GetStatus() (*release.Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus")
	ret0, _ := ret[0].(*release.Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockInterfaceMockRecorder) GetStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockInterface)(nil).GetStatus))
}

// GetValues mocks base method.
func (m *MockInterface) GetValues() map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValues")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// GetValues indicates an expected call of GetValues.
func (mr *MockInterfaceMockRecorder) GetValues() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValues", reflect.TypeOf((*MockInterface)(nil).GetValues))
}

// Install mocks base method.
func (m *MockInterface) Install() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Install")
	ret0, _ := ret[0].(error)
	return ret0
}

// Install indicates an expected call of Install.
func (mr *MockInterfaceMockRecorder) Install() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Install", reflect.TypeOf((*MockInterface)(nil).Install))
}

// IsDeployed mocks base method.
func (m *MockInterface) IsDeployed() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDeployed")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsDeployed indicates an expected call of IsDeployed.
func (mr *MockInterfaceMockRecorder) IsDeployed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDeployed", reflect.TypeOf((*MockInterface)(nil).IsDeployed))
}

// ReadValuesFile mocks base method.
func (m *MockInterface) ReadValuesFile() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadValuesFile")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadValuesFile indicates an expected call of ReadValuesFile.
func (mr *MockInterfaceMockRecorder) ReadValuesFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadValuesFile", reflect.TypeOf((*MockInterface)(nil).ReadValuesFile))
}

// SetConfiguration mocks base method.
func (m *MockInterface) SetConfiguration(arg0 *types.HelmConfigurationSpec) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetConfiguration", arg0)
}

// SetConfiguration indicates an expected call of SetConfiguration.
func (mr *MockInterfaceMockRecorder) SetConfiguration(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConfiguration", reflect.TypeOf((*MockInterface)(nil).SetConfiguration), arg0)
}

// SetKubernetesBackendType mocks base method.
func (m *MockInterface) SetKubernetesBackendType(arg0 types.KubernetesBackendType) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetKubernetesBackendType", arg0)
}

// SetKubernetesBackendType indicates an expected call of SetKubernetesBackendType.
func (mr *MockInterfaceMockRecorder) SetKubernetesBackendType(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetKubernetesBackendType", reflect.TypeOf((*MockInterface)(nil).SetKubernetesBackendType), arg0)
}

// SetNamespace mocks base method.
func (m *MockInterface) SetNamespace(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNamespace", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetNamespace indicates an expected call of SetNamespace.
func (mr *MockInterfaceMockRecorder) SetNamespace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNamespace", reflect.TypeOf((*MockInterface)(nil).SetNamespace), arg0)
}

// SetPersistanceVolumeClaimName mocks base method.
func (m *MockInterface) SetPersistanceVolumeClaimName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPersistanceVolumeClaimName", arg0)
}

// SetPersistanceVolumeClaimName indicates an expected call of SetPersistanceVolumeClaimName.
func (mr *MockInterfaceMockRecorder) SetPersistanceVolumeClaimName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPersistanceVolumeClaimName", reflect.TypeOf((*MockInterface)(nil).SetPersistanceVolumeClaimName), arg0)
}

// SetValues mocks base method.
func (m *MockInterface) SetValues(arg0 map[string]interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetValues", arg0)
}

// SetValues indicates an expected call of SetValues.
func (mr *MockInterfaceMockRecorder) SetValues(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetValues", reflect.TypeOf((*MockInterface)(nil).SetValues), arg0)
}

// Uninstall mocks base method.
func (m *MockInterface) Uninstall() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Uninstall")
	ret0, _ := ret[0].(error)
	return ret0
}

// Uninstall indicates an expected call of Uninstall.
func (mr *MockInterfaceMockRecorder) Uninstall() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Uninstall", reflect.TypeOf((*MockInterface)(nil).Uninstall))
}

// Upgrade mocks base method.
func (m *MockInterface) Upgrade() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upgrade")
	ret0, _ := ret[0].(error)
	return ret0
}

// Upgrade indicates an expected call of Upgrade.
func (mr *MockInterfaceMockRecorder) Upgrade() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upgrade", reflect.TypeOf((*MockInterface)(nil).Upgrade))
}