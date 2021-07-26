package k8s

import (
	"context"
	"net"

	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
