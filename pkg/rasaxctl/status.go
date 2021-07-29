package rasaxctl

import (
	"bytes"
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/status"
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/spf13/viper"
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
	fmt.Fprintf(&b, "Version: %s\n", stateData[types.StateSecretRasaXVersion])
	fmt.Fprintf(&b, "Rasa worker version: %s\n", stateData[types.StateSecretRasaWorkerVersion])

	projectPath := "not defined"
	if string(stateData[types.StateSecretProjectPath]) != "" {
		projectPath = string(stateData[types.StateSecretProjectPath])
	}
	fmt.Fprintf(&b, "Project path: %s\n", projectPath)

	url, err := r.GetRasaXURL()
	if err != nil {
		return err
	}
	fmt.Fprintf(&b, "URL: %s", url)

	if viper.GetBool("details") {

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
			fmt.Fprintf(&b, "Pod details:\n\n")
		}

		fmt.Println(b.String())

		status.PrintTable(
			[]string{"Name", "Condition", "Status"},
			data,
		)
	}
	fmt.Println(b.String())
	fmt.Println()

	return nil
}
