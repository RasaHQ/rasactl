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
	"encoding/json"
	"net"

	"github.com/spf13/viper"
	"golang.org/x/xerrors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ktypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/RasaHQ/rasactl/pkg/types"
)

func (k *Kubernetes) detectBackend() types.KubernetesBackendType {

	var backend types.KubernetesBackendType

	host, _, err := net.SplitHostPort(k.clientset.RESTClient().Get().URL().Host)
	if err != nil {
		host = k.clientset.RESTClient().Get().URL().Host
		k.Log.Info("Can't parse Kubernetes server host", "error", err)
	}

	ip := net.ParseIP(host)
	if ip.IsLoopback() {
		backend = types.KubernetesBackendLocal
	} else {
		backend = types.KubernetesBackendRemote
	}

	k.Log.V(1).Info("Detected Kubernetes backend", "type", backend, "host", host)

	return backend
}

// IsNamespaceExist checks if a namespace exists.
func (k *Kubernetes) IsNamespaceExist(namespace string) (bool, error) {

	_, err := k.clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})

	if err != nil {
		k.Log.V(1).Info("Namespace is invalid", "namespace", namespace, "error", err)
		return false, nil
	} else if errors.IsNotFound(err) {
		k.Log.V(1).Info("Namespace not found", "namespace", namespace)
		return false, nil
	}

	return true, nil

}

// GetKindControlPlaneNode returns v1.Node object that defines a kind control plane node.
func (k *Kubernetes) GetKindControlPlaneNode() (v1.Node, error) {

	nodes, err := k.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: "node-role.kubernetes.io/control-plane="})
	if err != nil {
		return v1.Node{}, err
	}

	if len(nodes.Items) != 0 {
		return nodes.Items[0], nil
	}

	return v1.Node{}, nil
}

// IsNamespaceManageable checks if a given namespace is managed by rasactl and returns `true` if it is.
func (k *Kubernetes) IsNamespaceManageable() bool {
	namespace, err := k.clientset.CoreV1().Namespaces().Get(context.TODO(), k.Namespace, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return namespace.Labels["rasactl"] == "true"
}

// AddNamespaceLabels adds an extra label to a given namespace that indicates that the namespace
// is managed by rasactl.
func (k *Kubernetes) AddNamespaceLabel() error {

	ns, err := k.clientset.CoreV1().Namespaces().Get(context.TODO(), k.Namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	var mergedLabels = make(map[string]string)
	for key, value := range ns.Labels {
		mergedLabels[key] = value
	}
	mergedLabels["rasactl"] = "true"
	ns.Labels = mergedLabels

	if _, err := k.clientset.CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

// DeleteNamespaceLabel deletes a label that indicates if a given namespaces is managed by rasactl.
func (k *Kubernetes) DeleteNamespaceLabel() error {
	type patch struct {
		Op   string `json:"op"`
		Path string `json:"path"`
	}

	payload := []patch{{
		Op:   "remove",
		Path: "/metadata/labels/rasactl",
	}}

	payloadBytes, _ := json.Marshal(payload)
	k.Log.V(1).Info("Deleting label", "namespace", k.Namespace, "payload", string(payloadBytes))
	if _, err := k.clientset.CoreV1().Namespaces().Patch(context.TODO(), k.Namespace,
		ktypes.JSONPatchType, payloadBytes, metav1.PatchOptions{}); err != nil {
		return err
	}
	return nil
}

// DeleteNode deletes a given Kubernetes node.
func (k *Kubernetes) DeleteNode(node string) error {
	return k.clientset.CoreV1().Nodes().Delete(context.TODO(), node, metav1.DeleteOptions{})
}

// DeleteNamespace deletes the active namespace.
func (k *Kubernetes) DeleteNamespace() error {
	return k.clientset.CoreV1().Namespaces().Delete(context.TODO(), k.Namespace, metav1.DeleteOptions{})
}

// GetNamespaces returns namespaces that are managed by rasactl.
func (k *Kubernetes) GetNamespaces() ([]string, error) {
	result := []string{}
	namespaces, err := k.clientset.CoreV1().Namespaces().List(context.TODO(),
		metav1.ListOptions{LabelSelector: "rasactl=true"})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaces.Items {
		if namespace.Status.Phase != v1.NamespaceActive {
			continue
		}
		result = append(result, namespace.Name)
	}

	return result, nil
}

// PodStatus returns a pod condition.
func (k *Kubernetes) PodStatus(conditions []v1.PodCondition) string {
	for _, c := range conditions {
		if c.Status != v1.ConditionTrue {
			return "NotReady"
		}
	}

	return "Ready"
}

// LoadConfig loads the kubeconfig file and returns a complete client config.
func (k *Kubernetes) LoadConfig() (*rest.Config, error) {
	context := viper.GetString("kube-context")
	k.kubeconfig = viper.GetString("kubeconfig")

	rawConfig, err := clientcmd.LoadFromFile(k.kubeconfig)
	if err != nil {
		return nil, err
	}

	if rawConfig.CurrentContext == "" {
		return nil, xerrors.Errorf("kubeconfig: the current context is empty, use the --kube-context flag or kubectl to set a context")
	}

	if _, ok := rawConfig.Contexts[context]; !ok && context != "" {
		return nil, xerrors.Errorf("kubeconfig: the %s context doesn't exist", context)
	}

	if context != "" {
		rawConfig.CurrentContext = context
	}

	return clientcmd.NewDefaultClientConfig(*rawConfig, nil).ClientConfig()
}

// GetServiceWithLabels returns a list of services with a given labels.
func (k *Kubernetes) GetServiceWithLabels(opts metav1.ListOptions) (*v1.ServiceList, error) {
	return k.clientset.CoreV1().Services(k.Namespace).List(context.TODO(), opts)
}
