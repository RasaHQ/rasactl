package rasaxctl

import (
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func (r *RasaXCTL) Stop() error {
	r.Spinner.Message("Stopping Rasa X")
	if err := r.KubernetesClient.ScaleDown(); err != nil {
		return err
	}

	state, err := r.KubernetesClient.ReadSecretWithState()
	if err != nil {
		return err
	}

	if r.DockerClient.Kind.ControlPlaneHost != "" && string(state[types.StateSecretProjectPath]) != "" {
		nodeName := fmt.Sprintf("kind-%s", r.Namespace)
		if err := r.DockerClient.StopKindNode(nodeName); err != nil {
			return err
		}
	}
	r.Spinner.Message(fmt.Sprintf("Rasa X for the %s project has been stopped", r.Namespace))
	r.Spinner.Stop()
	return nil
}
