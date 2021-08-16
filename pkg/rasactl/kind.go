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
)

// CreateAndJoinKindNode creates and joins a kind node.
func (r *RasaCtl) CreateAndJoinKindNode() error {
	nodeName := fmt.Sprintf("kind-%s", r.Namespace)
	if _, err := r.DockerClient.CreateKindNode(nodeName); err != nil {
		return err
	}
	return nil
}

// GetKindControlPlaneNodeInfo gets information about a kind control plane node
// and stores data in the DockerClient.Kind object.
func (r *RasaCtl) GetKindControlPlaneNodeInfo() error {
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
