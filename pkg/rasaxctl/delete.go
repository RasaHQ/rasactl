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

	r.Spinner.Message("Deleting Rasa X")

	state, err := r.KubernetesClient.ReadSecretWithState()
	if err != nil && !force {
		return err
	} else if err != nil && force {
		r.Log.Info("Can't read state secret", "error", err)
	}

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

	if r.DockerClient.Kind.ControlPlaneHost != "" && string(state[types.StateSecretProjectPath]) != "" || force {
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
		rasaxctlFile := fmt.Sprintf("%s/.rasaxctl", state[types.StateSecretProjectPath])
		os.Remove(rasaxctlFile)
	}

	if r.KubernetesClient.BackendType == types.KubernetesBackendLocal && r.CloudProvider.Name == types.CloudProviderUnknown {
		host := fmt.Sprintf("%s.rasaxctl.local.io", r.Namespace)
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

	r.Spinner.Message("Done!")
	r.Spinner.Stop()
	return nil
}
