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

	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/RasaHQ/rasactl/pkg/utils/cloud"
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetes struct {
	kubeconfig    string
	clientset     *kubernetes.Clientset
	Namespace     string
	Helm          HelmSpec
	Log           logr.Logger
	BackendType   types.KubernetesBackendType
	CloudProvider *cloud.Provider
	Flags         *types.RasaCtlFlags
}

type HelmSpec struct {
	Values      map[string]interface{}
	ReleaseName string
}

func (k *Kubernetes) New() error {
	k.kubeconfig = viper.GetString("kubeconfig")
	config, err := clientcmd.BuildConfigFromFlags("", k.kubeconfig)
	if err != nil {
		return err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	k.clientset = clientset

	backendType, err := k.detectBackend()
	if err != nil {
		return err
	}
	k.BackendType = backendType

	k.Log.Info("Initializing Kubernetes client")

	return nil
}

func (k *Kubernetes) GetRasaXURL() (string, error) {

	if k.Helm.Values == nil {
		return "", fmt.Errorf("helm client requires values, %#v", k.Helm)
	}

	nginxValues := k.Helm.Values["nginx"]
	ingressValues := k.Helm.Values["ingress"]
	nginxIsEnabled := nginxValues.(map[string]interface{})["enabled"].(bool)
	ingressIsEnabled := ingressValues.(map[string]interface{})["enabled"].(bool)
	nginxServiceType := nginxValues.(map[string]interface{})["service"].(map[string]interface{})["type"].(string)
	rasaXScheme := k.Helm.Values["rasax"].(map[string]interface{})["scheme"].(string)
	serviceName := fmt.Sprintf("%s-nginx", k.Helm.ReleaseName)

	url := "UNKNOWN"

	if nginxServiceType == "LoadBalancer" && nginxIsEnabled && (k.BackendType != types.KubernetesBackendLocal || k.CloudProvider.Name != types.CloudProviderUnknown) {
		service, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err != nil {
			return url, err
		}
		ipAddress := service.Status.LoadBalancer.Ingress[0].IP
		port := service.Spec.Ports[0].Port

		url = fmt.Sprintf("%s://%s:%d", rasaXScheme, ipAddress, port)

		return url, nil
	} else if nginxServiceType == "NodePort" && nginxIsEnabled && k.BackendType == types.KubernetesBackendLocal && k.CloudProvider.Name != types.CloudProviderUnknown {

		service, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err != nil {
			return url, err
		}
		port := service.Spec.Ports[0].NodePort

		url = fmt.Sprintf("%s://%s:%d", rasaXScheme, k.CloudProvider.ExternalIP, port)
	} else if ingressIsEnabled {
		ingress, err := k.clientset.NetworkingV1().Ingresses(k.Namespace).Get(context.TODO(), k.Helm.ReleaseName, metav1.GetOptions{})
		if err != nil {
			return url, err
		}
		host := ingress.Spec.Rules[0].Host

		if len(ingress.Spec.TLS) != 0 {
			rasaXScheme = "https://"
		}

		url = fmt.Sprintf("%s://%s", rasaXScheme, host)
	}

	return url, nil
}

func (k *Kubernetes) GetRasaXToken() (string, error) {
	var token string

	secretName := fmt.Sprintf("%s-rasa", k.Helm.ReleaseName)
	keyName := "rasaXToken"

	secret, err := k.clientset.CoreV1().Secrets(k.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return token, err
	}

	return string(secret.Data[keyName]), nil
}

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

func (k *Kubernetes) CreateNamespace() error {
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: k.Namespace,
			Labels: map[string]string{
				"rasactl": "true",
			},
		},
	}
	_, err := k.clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	switch t := err.(type) {
	case *errors.StatusError:
		if t.ErrStatus.Code == 409 {
			k.Log.V(1).Info("Namespace already exists", "namespace", k.Namespace)
			return nil
		}
		return err
	}

	k.Log.V(1).Info("Create namespace", "namespace", k.Namespace)
	return nil
}

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

func (k *Kubernetes) GetPostgreSQLSvcNodePort() (int32, error) {

	svcName := fmt.Sprintf("%s-postgresql", k.Helm.ReleaseName)
	svc, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(), svcName, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}

	return svc.Spec.Ports[0].NodePort, nil
}

func (k *Kubernetes) GetRabbitMqSvcNodePort() (int32, error) {

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

	return 0, nil
}
