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
package rasactl

import (
	"fmt"

	"github.com/RasaHQ/rasactl/pkg/types"
)

// Stop stops a deployment.
func (r *RasaCtl) Stop() error {
	r.Spinner.Message("Stopping Rasa X")

	r.initRasaXClient()

	rasaXVersion, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		return err
	}

	if err := r.KubernetesClient.UpdateSecretWithState(rasaXVersion); err != nil {
		return err
	}

	if err := r.KubernetesClient.ScaleDown(); err != nil {
		return err
	}

	state, err := r.KubernetesClient.ReadSecretWithState()
	if err != nil {
		return err
	}

	if r.DockerClient.GetKind().ControlPlaneHost != "" && string(state[types.StateProjectPath]) != "" {
		nodeName := fmt.Sprintf("kind-%s", r.Namespace)
		if err := r.DockerClient.StopKindNode(nodeName); err != nil {
			return err
		}
	}
	r.Spinner.Stop()
	fmt.Printf("Rasa X for the %s deployment has been stopped\n", r.Namespace)
	return nil
}
