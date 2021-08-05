package rasaxctl

import (
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/status"
	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func (r *RasaXCTL) List() error {
	data := [][]string{}
	header := []string{"Current", "Name", "Status", "Rasa production", "Rasa worker", "Enterprise", "Version"}
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

		r.HelmClient.Configuration = &types.HelmConfigurationSpec{
			ReleaseName: string(stateData[types.StateSecretHelmReleaseName]),
		}
		r.KubernetesClient.Helm.ReleaseName = string(stateData[types.StateSecretHelmReleaseName])
		url, err := r.GetRasaXURL()
		if err != nil {
			return err
		}
		r.initRasaXClient()
		r.RasaXClient.URL = url

		versionEndpoint, err := r.RasaXClient.GetVersionEndpoint()
		if err == nil {
			enterprise := "inactive"
			if versionEndpoint.Enterprise {
				enterprise = "active"
			}

			data = append(data, []string{current, namespace, status,
				versionEndpoint.Rasa.Production,
				versionEndpoint.Rasa.Worker,
				enterprise,
				versionEndpoint.RasaX,
			},
			)
		} else {
			data = append(data, []string{current, namespace, status,
				"0.0.0",
				string(stateData[types.StateSecretRasaWorkerVersion]),
				string(stateData[types.StateSecretEnterprise]),
				string(stateData[types.StateSecretRasaXVersion]),
			},
			)
		}
	}

	status.PrintTable(
		header,
		data,
	)
	return nil
}
