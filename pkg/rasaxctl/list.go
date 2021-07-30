package rasaxctl

import (
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/status"
	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func (r *RasaXCTL) List() error {
	data := [][]string{}
	namespaces, err := r.KubernetesClient.GetNamespaces()
	if err != nil {
		return err
	}

	if len(namespaces) == 0 {
		fmt.Println("Nothing to show, use the start command to create a new project")
		return nil
	}

	for _, namespace := range namespaces {
		r.KubernetesClient.Namespace = namespace
		isRunning, err := r.KubernetesClient.IsRasaXRunning()
		if err != nil {
			return err
		}
		status := "Stopped"
		if isRunning {
			status = "Running"
		}

		stateData, err := r.KubernetesClient.ReadSecretWithState()
		if err != nil {
			r.Log.Info("Can't read a secret with state", "namespace", namespace, "error", err)
		}

		current := ""
		if namespace == r.Namespace {
			current = "*"
		}

		data = append(data, []string{current, namespace, status,
			string(stateData[types.StateSecretRasaWorkerVersion]),
			string(stateData[types.StateSecretEnterprise]),
			string(stateData[types.StateSecretRasaXVersion]),
		},
		)
	}

	status.PrintTable(
		[]string{"Current", "Name", "Status", "Rasa worker", "Enterprise", "Version"},
		data,
	)
	return nil
}