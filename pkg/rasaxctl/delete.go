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
package rasaxctl

import (
	"fmt"
	"os"

	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils"
)

func (r *RasaXCTL) Delete() error {
	force := r.Flags.Delete.Force
	prune := r.Flags.Delete.Prune

	if prune {
		aYes, err := utils.AskForConfirmation("You're about to delete the namespace with all resources in it, are you sure?", 5, os.Stdin)
		if err != nil {
			return err
		}

		if !aYes {
			return nil
		}
	}

	msg := "Deleting Rasa X"
	r.Spinner.Message(msg)
	r.Log.Info(msg, "namespace", r.Namespace)

	state, err := r.KubernetesClient.ReadSecretWithState()
	if err != nil && !force {
		return err
	} else if err != nil && force {
		r.Log.Info("Can't read state secret", "error", err)
	}
	rasaxctlFile := fmt.Sprintf("%s/.rasaxctl", state[types.StateSecretProjectPath])

	if !prune {
		if err := r.HelmClient.Uninstall(); err != nil && !force {
			return err
		} else if err != nil && force {
			r.Log.Info("Can't uninstall helm chart", "error", err)
		}

		msgDelSec := "Deleting secret with rasaxctl state"
		r.Spinner.Message(msgDelSec)
		r.Log.Info(msgDelSec)
		if err := r.KubernetesClient.DeleteSecretWithState(); err != nil && !force {
			return err
		} else if err != nil && force {
			r.Log.Info("Can't delete secret with state", "error", err)
		}

		if err := r.KubernetesClient.DeleteNamespaceLabel(); err != nil && !force {
			return err
		} else if err != nil && force {
			r.Log.Info("Can't delete label", "error", err)
		}
	}

	if (r.DockerClient.Kind.ControlPlaneHost != "" && string(state[types.StateSecretProjectPath]) != "") || force {
		r.Spinner.Message("Deleting persistent volume")
		if err := r.KubernetesClient.DeleteVolume(); err != nil && !force {
			return err
		} else if err != nil && force {
			r.Log.Info("Can't delete persistent volume", "error", err)
		}

		r.Spinner.Message("Deleting a kind node")
		nodeName := fmt.Sprintf("kind-%s", r.Namespace)
		r.Log.Info("Deleting a kind node", "node", nodeName)
		if err := r.DockerClient.DeleteKindNode(nodeName); err != nil && !force {
			return err
		} else if err != nil && force {
			r.Log.Info("Can't delete a kind node", "node", nodeName, "error", err)
		}
		if err := r.KubernetesClient.DeleteNode(nodeName); err != nil && !force {
			return err
		} else if err != nil && force {
			r.Log.Info("Can't delete a Kubernetes node", "node", nodeName, "error", err)
		}
	}

	if r.KubernetesClient.BackendType == types.KubernetesBackendLocal && r.CloudProvider.Name == types.CloudProviderUnknown {
		host := fmt.Sprintf("%s.%s", r.Namespace, types.RasaXCtlLocalDomain)
		err := utils.DeleteHostToEtcHosts(host)
		if err != nil && !force {
			return err
		} else if err != nil && force {
			r.Log.Info("Can't delete host entry", "error", err)
		}
	}

	if prune {
		r.Log.Info("Deleting namespace", "namespace", r.Namespace)
		if err := r.KubernetesClient.DeleteNamespace(); err != nil && !force {
			return err
		} else if err != nil && force {
			r.Log.Info("Can't delete namespace", "namespace", r.Namespace, "error", err)
		}
	}

	if string(state[types.StateSecretProjectPath]) != "" {
		if err := os.Remove(rasaxctlFile); err != nil && !force {
			return err
		}
	}

	r.Spinner.Message("Done!")
	r.Spinner.Stop()
	return nil
}
