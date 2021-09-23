// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/RasaHQ/rasactl/pkg/k8s (interfaces: KubernetesInterface)

// Package fake is a generated GoMock package.
package fake

import (
	reflect "reflect"

	types "github.com/RasaHQ/rasactl/pkg/types"
	cloud "github.com/RasaHQ/rasactl/pkg/utils/cloud"
	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
	rest "k8s.io/client-go/rest"
)

// MockKubernetesInterface is a mock of KubernetesInterface interface.
type MockKubernetesInterface struct {
	ctrl     *gomock.Controller
	recorder *MockKubernetesInterfaceMockRecorder
}

// MockKubernetesInterfaceMockRecorder is the mock recorder for MockKubernetesInterface.
type MockKubernetesInterfaceMockRecorder struct {
	mock *MockKubernetesInterface
}

// NewMockKubernetesInterface creates a new mock instance.
func NewMockKubernetesInterface(ctrl *gomock.Controller) *MockKubernetesInterface {
	mock := &MockKubernetesInterface{ctrl: ctrl}
	mock.recorder = &MockKubernetesInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKubernetesInterface) EXPECT() *MockKubernetesInterfaceMockRecorder {
	return m.recorder
}

// AddNamespaceLabel mocks base method.
func (m *MockKubernetesInterface) AddNamespaceLabel() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNamespaceLabel")
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNamespaceLabel indicates an expected call of AddNamespaceLabel.
func (mr *MockKubernetesInterfaceMockRecorder) AddNamespaceLabel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNamespaceLabel", reflect.TypeOf((*MockKubernetesInterface)(nil).AddNamespaceLabel))
}

// CreateNamespace mocks base method.
func (m *MockKubernetesInterface) CreateNamespace() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNamespace")
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNamespace indicates an expected call of CreateNamespace.
func (mr *MockKubernetesInterfaceMockRecorder) CreateNamespace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNamespace", reflect.TypeOf((*MockKubernetesInterface)(nil).CreateNamespace))
}

// CreateVolume mocks base method.
func (m *MockKubernetesInterface) CreateVolume(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVolume", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVolume indicates an expected call of CreateVolume.
func (mr *MockKubernetesInterfaceMockRecorder) CreateVolume(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVolume", reflect.TypeOf((*MockKubernetesInterface)(nil).CreateVolume), arg0)
}

// DeleteNamespace mocks base method.
func (m *MockKubernetesInterface) DeleteNamespace() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNamespace")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNamespace indicates an expected call of DeleteNamespace.
func (mr *MockKubernetesInterfaceMockRecorder) DeleteNamespace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNamespace", reflect.TypeOf((*MockKubernetesInterface)(nil).DeleteNamespace))
}

// DeleteNamespaceLabel mocks base method.
func (m *MockKubernetesInterface) DeleteNamespaceLabel() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNamespaceLabel")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNamespaceLabel indicates an expected call of DeleteNamespaceLabel.
func (mr *MockKubernetesInterfaceMockRecorder) DeleteNamespaceLabel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNamespaceLabel", reflect.TypeOf((*MockKubernetesInterface)(nil).DeleteNamespaceLabel))
}

// DeleteNode mocks base method.
func (m *MockKubernetesInterface) DeleteNode(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNode indicates an expected call of DeleteNode.
func (mr *MockKubernetesInterfaceMockRecorder) DeleteNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNode", reflect.TypeOf((*MockKubernetesInterface)(nil).DeleteNode), arg0)
}

// DeleteRasaXPods mocks base method.
func (m *MockKubernetesInterface) DeleteRasaXPods() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRasaXPods")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRasaXPods indicates an expected call of DeleteRasaXPods.
func (mr *MockKubernetesInterfaceMockRecorder) DeleteRasaXPods() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRasaXPods", reflect.TypeOf((*MockKubernetesInterface)(nil).DeleteRasaXPods))
}

// DeleteSecretWithState mocks base method.
func (m *MockKubernetesInterface) DeleteSecretWithState() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecretWithState")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecretWithState indicates an expected call of DeleteSecretWithState.
func (mr *MockKubernetesInterfaceMockRecorder) DeleteSecretWithState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecretWithState", reflect.TypeOf((*MockKubernetesInterface)(nil).DeleteSecretWithState))
}

// DeleteVolume mocks base method.
func (m *MockKubernetesInterface) DeleteVolume() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVolume")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVolume indicates an expected call of DeleteVolume.
func (mr *MockKubernetesInterfaceMockRecorder) DeleteVolume() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVolume", reflect.TypeOf((*MockKubernetesInterface)(nil).DeleteVolume))
}

// GetBackendType mocks base method.
func (m *MockKubernetesInterface) GetBackendType() types.KubernetesBackendType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBackendType")
	ret0, _ := ret[0].(types.KubernetesBackendType)
	return ret0
}

// GetBackendType indicates an expected call of GetBackendType.
func (mr *MockKubernetesInterfaceMockRecorder) GetBackendType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBackendType", reflect.TypeOf((*MockKubernetesInterface)(nil).GetBackendType))
}

// GetCloudProvider mocks base method.
func (m *MockKubernetesInterface) GetCloudProvider() *cloud.Provider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCloudProvider")
	ret0, _ := ret[0].(*cloud.Provider)
	return ret0
}

// GetCloudProvider indicates an expected call of GetCloudProvider.
func (mr *MockKubernetesInterfaceMockRecorder) GetCloudProvider() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCloudProvider", reflect.TypeOf((*MockKubernetesInterface)(nil).GetCloudProvider))
}

// GetKindControlPlaneNode mocks base method.
func (m *MockKubernetesInterface) GetKindControlPlaneNode() (v1.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKindControlPlaneNode")
	ret0, _ := ret[0].(v1.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKindControlPlaneNode indicates an expected call of GetKindControlPlaneNode.
func (mr *MockKubernetesInterfaceMockRecorder) GetKindControlPlaneNode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKindControlPlaneNode", reflect.TypeOf((*MockKubernetesInterface)(nil).GetKindControlPlaneNode))
}

// GetLogs mocks base method.
func (m *MockKubernetesInterface) GetLogs(arg0 string) *rest.Request {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogs", arg0)
	ret0, _ := ret[0].(*rest.Request)
	return ret0
}

// GetLogs indicates an expected call of GetLogs.
func (mr *MockKubernetesInterfaceMockRecorder) GetLogs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockKubernetesInterface)(nil).GetLogs), arg0)
}

// GetNamespaces mocks base method.
func (m *MockKubernetesInterface) GetNamespaces() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespaces")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNamespaces indicates an expected call of GetNamespaces.
func (mr *MockKubernetesInterfaceMockRecorder) GetNamespaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespaces", reflect.TypeOf((*MockKubernetesInterface)(nil).GetNamespaces))
}

// GetPod mocks base method.
func (m *MockKubernetesInterface) GetPod(arg0 string) (*v1.Pod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPod", arg0)
	ret0, _ := ret[0].(*v1.Pod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPod indicates an expected call of GetPod.
func (mr *MockKubernetesInterfaceMockRecorder) GetPod(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPod", reflect.TypeOf((*MockKubernetesInterface)(nil).GetPod), arg0)
}

// GetPods mocks base method.
func (m *MockKubernetesInterface) GetPods() (*v1.PodList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPods")
	ret0, _ := ret[0].(*v1.PodList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPods indicates an expected call of GetPods.
func (mr *MockKubernetesInterfaceMockRecorder) GetPods() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPods", reflect.TypeOf((*MockKubernetesInterface)(nil).GetPods))
}

// GetPostgreSQLCreds mocks base method.
func (m *MockKubernetesInterface) GetPostgreSQLCreds() (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostgreSQLCreds")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPostgreSQLCreds indicates an expected call of GetPostgreSQLCreds.
func (mr *MockKubernetesInterfaceMockRecorder) GetPostgreSQLCreds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostgreSQLCreds", reflect.TypeOf((*MockKubernetesInterface)(nil).GetPostgreSQLCreds))
}

// GetPostgreSQLSvcNodePort mocks base method.
func (m *MockKubernetesInterface) GetPostgreSQLSvcNodePort() (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostgreSQLSvcNodePort")
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostgreSQLSvcNodePort indicates an expected call of GetPostgreSQLSvcNodePort.
func (mr *MockKubernetesInterfaceMockRecorder) GetPostgreSQLSvcNodePort() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostgreSQLSvcNodePort", reflect.TypeOf((*MockKubernetesInterface)(nil).GetPostgreSQLSvcNodePort))
}

// GetRabbitMqCreds mocks base method.
func (m *MockKubernetesInterface) GetRabbitMqCreds() (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRabbitMqCreds")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetRabbitMqCreds indicates an expected call of GetRabbitMqCreds.
func (mr *MockKubernetesInterfaceMockRecorder) GetRabbitMqCreds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRabbitMqCreds", reflect.TypeOf((*MockKubernetesInterface)(nil).GetRabbitMqCreds))
}

// GetRabbitMqSvcNodePort mocks base method.
func (m *MockKubernetesInterface) GetRabbitMqSvcNodePort() (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRabbitMqSvcNodePort")
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRabbitMqSvcNodePort indicates an expected call of GetRabbitMqSvcNodePort.
func (mr *MockKubernetesInterfaceMockRecorder) GetRabbitMqSvcNodePort() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRabbitMqSvcNodePort", reflect.TypeOf((*MockKubernetesInterface)(nil).GetRabbitMqSvcNodePort))
}

// GetRasaXToken mocks base method.
func (m *MockKubernetesInterface) GetRasaXToken() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRasaXToken")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRasaXToken indicates an expected call of GetRasaXToken.
func (mr *MockKubernetesInterfaceMockRecorder) GetRasaXToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRasaXToken", reflect.TypeOf((*MockKubernetesInterface)(nil).GetRasaXToken))
}

// GetRasaXURL mocks base method.
func (m *MockKubernetesInterface) GetRasaXURL() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRasaXURL")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRasaXURL indicates an expected call of GetRasaXURL.
func (mr *MockKubernetesInterfaceMockRecorder) GetRasaXURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRasaXURL", reflect.TypeOf((*MockKubernetesInterface)(nil).GetRasaXURL))
}

// IsNamespaceExist mocks base method.
func (m *MockKubernetesInterface) IsNamespaceExist(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNamespaceExist", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsNamespaceExist indicates an expected call of IsNamespaceExist.
func (mr *MockKubernetesInterfaceMockRecorder) IsNamespaceExist(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNamespaceExist", reflect.TypeOf((*MockKubernetesInterface)(nil).IsNamespaceExist), arg0)
}

// IsNamespaceManageable mocks base method.
func (m *MockKubernetesInterface) IsNamespaceManageable() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNamespaceManageable")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNamespaceManageable indicates an expected call of IsNamespaceManageable.
func (mr *MockKubernetesInterfaceMockRecorder) IsNamespaceManageable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNamespaceManageable", reflect.TypeOf((*MockKubernetesInterface)(nil).IsNamespaceManageable))
}

// IsRasaXRunning mocks base method.
func (m *MockKubernetesInterface) IsRasaXRunning() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRasaXRunning")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRasaXRunning indicates an expected call of IsRasaXRunning.
func (mr *MockKubernetesInterfaceMockRecorder) IsRasaXRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRasaXRunning", reflect.TypeOf((*MockKubernetesInterface)(nil).IsRasaXRunning))
}

// LoadConfig mocks base method.
func (m *MockKubernetesInterface) LoadConfig() (*rest.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadConfig")
	ret0, _ := ret[0].(*rest.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadConfig indicates an expected call of LoadConfig.
func (mr *MockKubernetesInterfaceMockRecorder) LoadConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadConfig", reflect.TypeOf((*MockKubernetesInterface)(nil).LoadConfig))
}

// PodStatus mocks base method.
func (m *MockKubernetesInterface) PodStatus(arg0 []v1.PodCondition) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PodStatus", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// PodStatus indicates an expected call of PodStatus.
func (mr *MockKubernetesInterfaceMockRecorder) PodStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PodStatus", reflect.TypeOf((*MockKubernetesInterface)(nil).PodStatus), arg0)
}

// ReadSecretWithState mocks base method.
func (m *MockKubernetesInterface) ReadSecretWithState() (map[string][]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSecretWithState")
	ret0, _ := ret[0].(map[string][]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSecretWithState indicates an expected call of ReadSecretWithState.
func (mr *MockKubernetesInterfaceMockRecorder) ReadSecretWithState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSecretWithState", reflect.TypeOf((*MockKubernetesInterface)(nil).ReadSecretWithState))
}

// SaveSecretWithState mocks base method.
func (m *MockKubernetesInterface) SaveSecretWithState(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSecretWithState", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSecretWithState indicates an expected call of SaveSecretWithState.
func (mr *MockKubernetesInterfaceMockRecorder) SaveSecretWithState(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSecretWithState", reflect.TypeOf((*MockKubernetesInterface)(nil).SaveSecretWithState), arg0)
}

// ScaleDown mocks base method.
func (m *MockKubernetesInterface) ScaleDown() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScaleDown")
	ret0, _ := ret[0].(error)
	return ret0
}

// ScaleDown indicates an expected call of ScaleDown.
func (mr *MockKubernetesInterfaceMockRecorder) ScaleDown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScaleDown", reflect.TypeOf((*MockKubernetesInterface)(nil).ScaleDown))
}

// ScaleUp mocks base method.
func (m *MockKubernetesInterface) ScaleUp() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScaleUp")
	ret0, _ := ret[0].(error)
	return ret0
}

// ScaleUp indicates an expected call of ScaleUp.
func (mr *MockKubernetesInterfaceMockRecorder) ScaleUp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScaleUp", reflect.TypeOf((*MockKubernetesInterface)(nil).ScaleUp))
}

// SetHelmReleaseName mocks base method.
func (m *MockKubernetesInterface) SetHelmReleaseName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHelmReleaseName", arg0)
}

// SetHelmReleaseName indicates an expected call of SetHelmReleaseName.
func (mr *MockKubernetesInterfaceMockRecorder) SetHelmReleaseName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHelmReleaseName", reflect.TypeOf((*MockKubernetesInterface)(nil).SetHelmReleaseName), arg0)
}

// SetHelmValues mocks base method.
func (m *MockKubernetesInterface) SetHelmValues(arg0 map[string]interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHelmValues", arg0)
}

// SetHelmValues indicates an expected call of SetHelmValues.
func (mr *MockKubernetesInterfaceMockRecorder) SetHelmValues(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHelmValues", reflect.TypeOf((*MockKubernetesInterface)(nil).SetHelmValues), arg0)
}

// SetNamespace mocks base method.
func (m *MockKubernetesInterface) SetNamespace(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNamespace", arg0)
}

// SetNamespace indicates an expected call of SetNamespace.
func (mr *MockKubernetesInterfaceMockRecorder) SetNamespace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNamespace", reflect.TypeOf((*MockKubernetesInterface)(nil).SetNamespace), arg0)
}

// UpdateRasaXConfig mocks base method.
func (m *MockKubernetesInterface) UpdateRasaXConfig(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRasaXConfig", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRasaXConfig indicates an expected call of UpdateRasaXConfig.
func (mr *MockKubernetesInterfaceMockRecorder) UpdateRasaXConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRasaXConfig", reflect.TypeOf((*MockKubernetesInterface)(nil).UpdateRasaXConfig), arg0)
}

// UpdateSecretWithState mocks base method.
func (m *MockKubernetesInterface) UpdateSecretWithState(arg0 ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateSecretWithState", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSecretWithState indicates an expected call of UpdateSecretWithState.
func (mr *MockKubernetesInterfaceMockRecorder) UpdateSecretWithState(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecretWithState", reflect.TypeOf((*MockKubernetesInterface)(nil).UpdateSecretWithState), arg0...)
}
