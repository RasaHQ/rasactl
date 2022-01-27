/*
Copyright © 2021 Rasa Technologies GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package k8s

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"golang.org/x/xerrors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"

	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/RasaHQ/rasactl/pkg/utils"
	"github.com/RasaHQ/rasactl/pkg/utils/cloud"
)

type KubernetesInterface interface {
	GetRasaXURL() (string, error)
	GetRasaXToken() (string, error)
	CreateNamespace() error
	IsRasaXRunning() (bool, error)
	GetPods() (*v1.PodList, error)
	DeleteRasaXPods() error
	GetPostgreSQLSvcNodePort() (int32, error)
	GetRasaXSvcNodePort() (int32, error)
	GetRabbitMqSvcNodePort() (int32, error)
	SaveSecretWithState(projectPath string) error
	UpdateRasaXConfig(token string) error
	ScaleDown() error
	ScaleUp() error
	UpdateSecretWithState(data ...interface{}) error
	ReadSecretWithState() (map[string][]byte, error)
	DeleteSecretWithState() error
	GetPostgreSQLCreds() (string, string, error)
	GetRabbitMqCreds() (string, string, error)
	IsNamespaceExist(namespace string) (bool, error)
	IsSecretWithStateExist() bool
	GetKindControlPlaneNode() (v1.Node, error)
	IsNamespaceManageable() bool
	AddNamespaceLabel() error
	DeleteNamespaceLabel() error
	DeleteNode(node string) error
	DeleteNamespace() error
	GetNamespaces() ([]string, error)
	PodStatus(conditions []v1.PodCondition) string
	CreateVolume(hostPath string) (string, error)
	DeleteVolume() error
	GetBackendType() types.KubernetesBackendType
	SetNamespace(namespace string)
	SetHelmValues(values map[string]interface{})
	SetHelmReleaseName(release string)
	GetCloudProvider() *cloud.Provider
	LoadConfig() (*rest.Config, error)
	GetLogs(pod string) *rest.Request
	GetPod(pod string) (*v1.Pod, error)
	GetServiceWithLabels(opts metav1.ListOptions) (*v1.ServiceList, error)
}

// Kubernetes represents Kubernetes client.
type Kubernetes struct {
	kubeconfig string

	clientset *kubernetes.Clientset

	// Namespace is a namepace name used by the client.
	Namespace string

	// Helm defines helm release configuration.
	Helm HelmSpec

	// Log defines logger.
	Log logr.Logger

	// BackendType stores a Kubernetes cluster type.
	BackendType types.KubernetesBackendType

	// CloudProvider defines a detected cloud provider.
	CloudProvider *cloud.Provider

	// Flags stores command flags used during the command execution.
	Flags *types.RasaCtlFlags
}

// HelmSpec stores data related to helm release.
type HelmSpec struct {
	// Values stores helm values for a given helm release.
	Values map[string]interface{}

	// ReleaseName is a helm release name used by the client.
	ReleaseName string
}

// New initializes a new Kubernetes client.
func New(client *Kubernetes) (KubernetesInterface, error) {
	client.Log.Info("Initializing Kubernetes client")

	config, err := client.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client.clientset = clientset

	client.BackendType = client.detectBackend()

	return client, nil
}

// GetBackendType returns the backend type.
func (k *Kubernetes) GetBackendType() types.KubernetesBackendType {
	return k.BackendType
}

// SetNamespace sets Namespace field.
func (k *Kubernetes) SetNamespace(namespace string) {
	k.Namespace = namespace
}

// SetHelmValues sets helm values.
func (k *Kubernetes) SetHelmValues(values map[string]interface{}) {
	k.Helm.Values = values
}

// SetHelmReleaseName sets helm release name.
func (k *Kubernetes) SetHelmReleaseName(release string) {
	k.Helm.ReleaseName = release
}

// GetCloudProvider returns CloudProvider field.
func (k *Kubernetes) GetCloudProvider() *cloud.Provider {
	return k.CloudProvider
}

// GetRasaXURL returns URL for a given deployment.
func (k *Kubernetes) GetRasaXURL() (string, error) {
	// Read URL from the environment variables
	if url := utils.GetRasaXURLEnv(k.Namespace); url != "" {
		k.Log.Info("Using Rasa X URL passed via the environment variables", "value", url)
		return url, nil
	}

	if k.Helm.Values == nil {
		return "", xerrors.Errorf("helm client requires values, %#v", k.Helm)
	}

	nginxValues := k.Helm.Values["nginx"]
	ingressValues := k.Helm.Values["ingress"]
	nginxIsEnabled := nginxValues.(map[string]interface{})["enabled"].(bool)
	ingressIsEnabled := ingressValues.(map[string]interface{})["enabled"].(bool)
	nginxServiceType := nginxValues.(map[string]interface{})["service"].(map[string]interface{})["type"].(string)
	rasaXScheme := k.Helm.Values["rasax"].(map[string]interface{})["scheme"].(string)

	url := "UNKNOWN"

	if nginxServiceType == "LoadBalancer" &&
		nginxIsEnabled &&
		(k.BackendType != types.KubernetesBackendLocal || k.CloudProvider.Name != types.CloudProviderUnknown) {

		nginxService, err := k.getNginxService()
		if err != nil {
			return "", err
		}

		service, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(), nginxService.Name, metav1.GetOptions{})
		if err != nil {
			return url, err
		}
		ipAddress := service.Status.LoadBalancer.Ingress[0].IP
		port := service.Spec.Ports[0].Port

		url = fmt.Sprintf("%s://%s:%d", rasaXScheme, ipAddress, port)

		return url, nil
	} else if nginxServiceType == "NodePort" &&
		nginxIsEnabled &&
		k.BackendType == types.KubernetesBackendLocal && k.CloudProvider.Name != types.CloudProviderUnknown {

		nginxService, err := k.getNginxService()
		if err != nil {
			return "", err
		}

		service, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(),
			nginxService.Name, metav1.GetOptions{})
		if err != nil {
			return url, err
		}
		port := service.Spec.Ports[0].NodePort
		ip := k.CloudProvider.ExternalIP

		if ip == "" {
			ip = "127.0.0.1"
			k.Log.Info("Can't get an external IP address, using localhost instead",
				"externalIP", k.CloudProvider.ExternalIP)
		}

		url = fmt.Sprintf("%s://%s:%d", rasaXScheme, ip, port)
	} else if ingressIsEnabled {
		// If a release name is different than "rasa-x" then a ingress name has
		// a different pattern (<release-name>-rasa-x).
		// Use labels to be sure that the correct ingress is read.
		labels := fmt.Sprintf("app.kubernetes.io/name=rasa-x,app.kubernetes.io/instance=%s",
			k.Helm.ReleaseName)

		ingresses, err := k.clientset.NetworkingV1().Ingresses(k.Namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels,
			Limit:         1,
		})
		if err != nil {
			return url, err
		}

		ingress, err := k.clientset.NetworkingV1().Ingresses(k.Namespace).Get(context.TODO(),
			ingresses.Items[0].Name, metav1.GetOptions{})
		if err != nil {
			return url, err
		}
		host := ingress.Spec.Rules[0].Host

		if len(ingress.Spec.TLS) != 0 {
			rasaXScheme = "https"
		}

		url = fmt.Sprintf("%s://%s", rasaXScheme, host)
	}

	return url, nil
}

func (k *Kubernetes) getNginxService() (v1.Service, error) {
	labels := fmt.Sprintf("app.kubernetes.io/component=nginx,app.kubernetes.io/instance=%s",
		k.Helm.ReleaseName)

	svc, err := k.GetServiceWithLabels(metav1.ListOptions{
		LabelSelector: labels,
		Limit:         1,
	})
	if err != nil {
		return v1.Service{}, err
	}

	return svc.Items[0], nil
}

// GetRasaXToken returns a Rasa X token that is stored in a Kubernetes secret.
func (k *Kubernetes) GetRasaXToken() (string, error) {

	secretName := fmt.Sprintf("%s-rasa", k.Helm.ReleaseName)
	keyName := "rasaXToken"

	secret, err := k.clientset.CoreV1().Secrets(k.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return string(secret.Data[keyName]), nil
}

// IsRasaXRunning checks if Rasa X deployment is running.
func (k *Kubernetes) IsRasaXRunning() (bool, error) {
	deployments, err := k.clientset.AppsV1().Deployments(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	if len(deployments.Items) == 0 {
		return false, nil
	}

	for _, deployment := range deployments.Items {
		if deployment.Status.Replicas == 0 || deployment.Status.ReadyReplicas == 0 {
			k.Log.V(1).Info("Deployment has replica number set to 0",
				"statefulset", deployment.Name, "replicas", deployment.Status.Replicas, "readyReplicas", deployment.Status.ReadyReplicas)
			return false, nil
		}
	}

	statefulsets, err := k.clientset.AppsV1().StatefulSets(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	if len(statefulsets.Items) == 0 {
		return false, nil
	}

	for _, statefulset := range statefulsets.Items {
		if statefulset.Status.Replicas == 0 || statefulset.Status.ReadyReplicas == 0 {
			k.Log.V(1).Info("Statefulset has replica number set to 0",
				"statefulset", statefulset.Name, "replicas", statefulset.Status.Replicas, "readyReplicas", statefulset.Status.ReadyReplicas)
			return false, nil
		}
	}

	return true, nil
}

// CreateNamespace creates a namespace.
func (k *Kubernetes) CreateNamespace() error {
	k.Log.V(1).Info("Creating namespace", "namespace", k.Namespace)

	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: k.Namespace,
			Labels: map[string]string{
				"rasactl": "true",
			},
		},
	}
	_, err := k.clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	switch t := err.(type) { //nolint:errorlint
	case *errors.StatusError:
		if t.ErrStatus.Code == 409 {
			k.Log.V(1).Info("Namespace already exists", "namespace", k.Namespace)
			return nil
		}
		return err
	default:
		return err
	}
}

// GetPods returns a list of pods for the active namespace.
func (k *Kubernetes) GetPods() (*v1.PodList, error) {
	pods, err := k.clientset.CoreV1().Pods(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func (k *Kubernetes) DeleteRasaXPods() error {
	pods, err := k.clientset.CoreV1().Pods(k.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "app.kubernetes.io/component=rasa-x"})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		if err := k.clientset.CoreV1().Pods(k.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	return nil
}

// GetPostgreSQLSvcPort returns a node port for the postgresql service.
func (k *Kubernetes) GetPostgreSQLSvcNodePort() (int32, error) {

	k.Log.V(1).Info("Getting a node port for the PostgreSQL service")

	svcName := fmt.Sprintf("%s-postgresql", k.Helm.ReleaseName)
	svc, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(), svcName, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}

	return svc.Spec.Ports[0].NodePort, nil
}

// GetRasaXSvcNodePort returns a node port for the postgresql service.
func (k *Kubernetes) GetRasaXSvcNodePort() (int32, error) {

	k.Log.V(1).Info("Getting a node port for the Rasa X service")

	// If a release name is different than "rasa-x" then a service name has
	// a different pattern (<release-name>-rasa-x-rasa-x).
	// Use labels to be sure that the correct service is read.
	lables := fmt.Sprintf("app.kubernetes.io/component=rasa-x,app.kubernetes.io/instance=%s",
		k.Helm.ReleaseName)

	svcs, err := k.clientset.CoreV1().Services(k.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: lables,
		Limit:         1,
	})
	if err != nil {
		return 0, err
	}

	svc, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(), svcs.Items[0].Name, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}

	return svc.Spec.Ports[0].NodePort, nil
}

// GetRabbitMqNodePort returns a node port for the rabbitmq service.
func (k *Kubernetes) GetRabbitMqSvcNodePort() (int32, error) {

	k.Log.V(1).Info("Getting a node port for the RabbitMQ service")

	helmValues := k.Helm.Values
	rabbitPort := helmValues["rabbitmq"].(map[string]interface{})["service"].(map[string]interface{})["port"].(float64)

	svcName := fmt.Sprintf("%s-rabbit", k.Helm.ReleaseName)
	svc, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(), svcName, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}

	for _, port := range svc.Spec.Ports {
		if port.Port == int32(rabbitPort) {
			return port.NodePort, nil
		}
	}

	return 0, xerrors.Errorf("can't determine a node port for the rabbitmq service")
}

// GetLogs returns the logs stream for a pod.
func (k *Kubernetes) GetLogs(pod string) *rest.Request {

	opts := v1.PodLogOptions{
		Previous: k.Flags.Logs.Previous,
		Follow:   k.Flags.Logs.Follow,
	}

	if k.Flags.Logs.TailLines > 0 {
		opts.TailLines = &k.Flags.Logs.TailLines
	}

	if k.Flags.Logs.Container != "" {
		opts.Container = k.Flags.Logs.Container
	}

	return k.clientset.CoreV1().
		Pods(k.Namespace).
		GetLogs(pod, &opts)
}

// GetPod returns a Pod object for a given pod.
func (k *Kubernetes) GetPod(pod string) (*v1.Pod, error) {
	return k.clientset.CoreV1().Pods(k.Namespace).Get(context.TODO(), pod, metav1.GetOptions{})
}
