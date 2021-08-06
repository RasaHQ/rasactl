package rasaxctl

import (
	"bytes"
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/status"
	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func (r *RasaXCTL) Status() error {
	var b bytes.Buffer
	isRunning, err := r.KubernetesClient.IsRasaXRunning()
	if err != nil {
		return err
	}
	statusProject := "Stopped"
	if isRunning {
		statusProject = "Running"
	}

	stateData, err := r.KubernetesClient.ReadSecretWithState()
	if err != nil {
		return err
	}
	r.HelmClient.Configuration = &types.HelmConfigurationSpec{
		ReleaseName: string(stateData[types.StateSecretHelmReleaseName]),
	}
	r.KubernetesClient.Helm.ReleaseName = string(stateData[types.StateSecretHelmReleaseName])

	fmt.Fprintf(&b, "Name: %s\n", r.Namespace)
	fmt.Fprintf(&b, "Status: %s\n", statusProject)

	url, err := r.GetRasaXURL()
	if err != nil {
		return err
	}
	fmt.Fprintf(&b, "URL: %s\n", url)

	r.initRasaXClient()
	r.RasaXClient.URL = url

	versionEndpoint, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		fmt.Fprintf(&b, "Version: %s\n", stateData[types.StateSecretRasaXVersion])
		fmt.Fprintf(&b, "Enterprise: %s\n", stateData[types.StateSecretEnterprise])
		fmt.Fprintf(&b, "Rasa worker version: %s\n", stateData[types.StateSecretRasaWorkerVersion])
	} else {
		fmt.Fprintf(&b, "Version: %s\n", versionEndpoint.RasaX)
		fmt.Fprintf(&b, "Enterprise: %t\n", versionEndpoint.Enterprise)
		fmt.Fprintf(&b, "Rasa production version: %s\n", versionEndpoint.Rasa.Production)
		fmt.Fprintf(&b, "Rasa worker version: %s\n", versionEndpoint.Rasa.Worker)
	}

	projectPath := "not defined"
	if string(stateData[types.StateSecretProjectPath]) != "" {
		projectPath = string(stateData[types.StateSecretProjectPath])
	}
	fmt.Fprintf(&b, "Project path: %s\n", projectPath)

	if r.Flags.Status.Details {

		release, err := r.HelmClient.GetStatus()
		if err != nil {
			return err
		}
		fmt.Fprintf(&b, "Helm chart: %s-%s\n", release.Chart.Name(), release.Chart.Metadata.Version)
		fmt.Fprintf(&b, "Helm release: %s\n", release.Name)
		fmt.Fprintf(&b, "Helm release status: %s\n\n", release.Info.Status)

		pods, err := r.KubernetesClient.GetPods()
		if err != nil {
			return err
		}

		data := [][]string{}
		for _, pod := range pods.Items {
			data = append(data,
				[]string{
					pod.Name,
					r.KubernetesClient.PodStatus(pod.Status.Conditions),
					string(pod.Status.Phase),
				},
			)
		}

		if len(pods.Items) != 0 {
			fmt.Fprintf(&b, "Pod details:\n")

			fmt.Println(b.String())

			status.PrintTable(
				[]string{"Name", "Condition", "Status"},
				data,
			)
		} else {
			fmt.Println(b.String())
		}
		fmt.Println()
		return nil
	}
	fmt.Println(b.String())

	return nil
}
