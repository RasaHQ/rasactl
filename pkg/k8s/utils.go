package k8s

import (
	"context"
	"encoding/json"
	"net"

	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ktypes "k8s.io/apimachinery/pkg/types"
)

func (k *Kubernetes) detectBackend() (types.KubernetesBackendType, error) {

	var backend types.KubernetesBackendType

	host, _, err := net.SplitHostPort(k.clientset.RESTClient().Get().URL().Host)
	if err != nil {
		return "", err
	}

	ip := net.ParseIP(host)
	if utils.IsPrivate(ip) || ip.IsLoopback() {
		backend = types.KubernetesBackendLocal
	} else {
		backend = types.KubernetesBackendRemote
	}

	k.Log.V(1).Info("Detected Kubernetes backend", "type", backend, "host", host)

	return backend, nil
}

func (k *Kubernetes) IsNamespaceExist(namespace string) (bool, error) {

	_, err := k.clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	switch t := err.(type) {
	case *errors.StatusError:
		if t.ErrStatus.Code == 404 {
			k.Log.V(1).Info("Namespace not found", "namespace", k.Namespace)
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (k *Kubernetes) GetKindControlPlaneNode() (v1.Node, error) {

	nodes, err := k.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: "node-role.kubernetes.io/control-plane="})
	if err != nil {
		return v1.Node{}, err
	}

	for _, node := range nodes.Items {
		return node, nil
	}

	return v1.Node{}, nil
}

func (k *Kubernetes) IsNamespaceManageable() bool {
	namespace, err := k.clientset.CoreV1().Namespaces().Get(context.TODO(), k.Namespace, metav1.GetOptions{})
	if err != nil {
		return false
	}
	if namespace.Labels["rasaxctl"] == "true" {
		return true
	}
	return false
}

func (k *Kubernetes) AddNamespaceLabel() error {
	type patch struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}

	payload := []patch{{
		Op:    "add",
		Path:  "/metadata/labels/rasaxctl",
		Value: "true",
	}}

	payloadBytes, _ := json.Marshal(payload)
	k.Log.V(1).Info("Adding label", "namespace", k.Namespace, "payload", string(payloadBytes))
	if _, err := k.clientset.CoreV1().Namespaces().Patch(context.TODO(), k.Namespace, ktypes.JSONPatchType, payloadBytes, metav1.PatchOptions{}); err != nil {
		return err
	}
	return nil
}

func (k *Kubernetes) DeleteNamespaceLabel() error {
	type patch struct {
		Op   string `json:"op"`
		Path string `json:"path"`
	}

	payload := []patch{{
		Op:   "remove",
		Path: "/metadata/labels/rasaxctl",
	}}

	payloadBytes, _ := json.Marshal(payload)
	k.Log.V(1).Info("Deleting label", "namespace", k.Namespace, "payload", string(payloadBytes))
	if _, err := k.clientset.CoreV1().Namespaces().Patch(context.TODO(), k.Namespace, ktypes.JSONPatchType, payloadBytes, metav1.PatchOptions{}); err != nil {
		return err
	}
	return nil
}

func (k *Kubernetes) DeleteNode(node string) error {
	if err := k.clientset.CoreV1().Nodes().Delete(context.TODO(), node, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func (k *Kubernetes) DeleteNamespace() error {
	if err := k.clientset.CoreV1().Namespaces().Delete(context.TODO(), k.Namespace, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func (k *Kubernetes) GetNamespaces() ([]string, error) {
	result := []string{}
	namespaces, err := k.clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{LabelSelector: "rasaxctl=true"})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaces.Items {
		result = append(result, namespace.Name)
	}

	return result, nil
}
