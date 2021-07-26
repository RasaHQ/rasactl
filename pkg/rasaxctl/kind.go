package rasaxctl

import (
	"fmt"
)

func (r *RasaXCTL) CreateAndJoinKindNode() error {
	nodeName := fmt.Sprintf("kind-%s", r.Namespace)
	if _, err := r.DockerClient.CreateKindNode(nodeName); err != nil {
		return err
	}
	return nil
}

func (r *RasaXCTL) GetKindControlPlaneNodeInfo() error {
	node, err := r.KubernetesClient.GetKindControlPlaneNode()
	if err != nil {
		return err
	}

	if node.Name == "" {
		r.Log.Info("Can't find kind control plane. Are you sure that the current Kubernetes context is kind?")
	}

	r.DockerClient.Kind.ControlPlaneHost = node.Name
	r.DockerClient.Kind.Version = node.Status.NodeInfo.KubeletVersion

	return nil
}
