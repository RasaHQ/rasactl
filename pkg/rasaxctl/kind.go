package rasaxctl

import (
	"fmt"

	"github.com/pkg/errors"
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
		return errors.Errorf("Can't find kind control plane. Are you sure that the current Kubernetes context is kind?")
	}

	r.DockerClient.Kind.ControlPlaneHost = node.Name
	r.DockerClient.Kind.Version = node.Status.NodeInfo.KubeletVersion
	return nil
}
