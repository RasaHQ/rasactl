package rasaxctl

import (
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/status"
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/spf13/viper"
)

func (r *RasaXCTL) Status() error {
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

	fmt.Printf("Name: %s\n", r.Namespace)
	fmt.Printf("Status: %s\n", statusProject)
	fmt.Printf("Version: %s\n", stateData[types.StateSecretRasaXVersion])
	fmt.Printf("Rasa worker version: %s\n", stateData[types.StateSecretRasaWorkerVersion])

	projectPath := "not defined"
	if string(stateData[types.StateSecretProjectPath]) != "" {
		projectPath = string(stateData[types.StateSecretProjectPath])
	}
	fmt.Printf("Project path: %s\n", projectPath)

	if viper.GetBool("details") {
		r.HelmClient.Configuration = &types.HelmConfigurationSpec{
			ReleaseName: string(stateData[types.StateSecretHelmReleaseName]),
		}
		release, err := r.HelmClient.GetStatus()
		if err != nil {
			return err
		}
		fmt.Printf("Helm chart: %s-%s\n", release.Chart.Name(), release.Chart.Metadata.Version)
		fmt.Printf("Helm release: %s\n", release.Name)
		fmt.Printf("Helm release status: %s\n", release.Info.Status)

		fmt.Println()

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
			fmt.Print("Pod details:\n\n")
		}

		status.PrintTable(
			[]string{"Name", "Condition", "Status"},
			data,
		)
	}

	fmt.Println()

	return nil
}
