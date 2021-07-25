package rasaxctl

import "fmt"

func (r *RasaXCTL) CreateAndJoinKindNode() error {
	nodeName := fmt.Sprintf("kind-%s", r.Namespace)
	if _, err := r.DockerClient.CreateKindNode(nodeName); err != nil {
		return err
	}
	return nil
}
