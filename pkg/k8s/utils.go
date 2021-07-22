package k8s

import (
	"net"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func (k *Kubernetes) detectBackend() (types.KubernetesBackendType, error) {

	var backend types.KubernetesBackendType

	host, _, err := net.SplitHostPort(k.clientset.RESTClient().Get().URL().Host)
	if err != nil {
		return "", err
	}

	if host == "127.0.0.1" || host == "localhost" {
		backend = types.KubernetesBackendLocal
	} else {
		backend = types.KubernetesBackendRemote
	}

	k.Log.V(1).Info("Detected Kubernetes backend", "type", backend)

	return backend, nil
}
