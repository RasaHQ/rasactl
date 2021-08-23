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

	"github.com/RasaHQ/rasactl/pkg/status"
	"github.com/RasaHQ/rasactl/pkg/types"
)

// Status prints status for a given deployment.
func (r *RasaCtl) Status() error {
	var d = [][]string{}
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
		ReleaseName: string(stateData[types.StateHelmReleaseName]),
	}
	r.KubernetesClient.SetHelmReleaseName(string(stateData[types.StateHelmReleaseName]))

	d = append(d, []string{"Name:", r.Namespace})
	d = append(d, []string{"Status:", statusProject})

	url, err := r.GetRasaXURL()
	if err != nil {
		return err
	}
	d = append(d, []string{"URL:", url})

	r.initRasaXClient()
	r.RasaXClient.URL = url

	versionEndpoint, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		d = append(d, []string{"Version:", string(stateData[types.StateRasaXVersion])})
		d = append(d, []string{"Enterprise:", string(stateData[types.StateEnterprise])})
	} else {
		enterprise := "inactive"
		if versionEndpoint.Enterprise {
			enterprise = "active"
		}
		d = append(d, []string{"Version:", versionEndpoint.RasaX})
		d = append(d, []string{"Enterprise:", enterprise})
		d = append(d, []string{"Rasa production version:", versionEndpoint.Rasa.Production})
		d = append(d, []string{"Rasa worker version:", versionEndpoint.Rasa.Worker})
	}

	projectPath := "not defined"
	if string(stateData[types.StateProjectPath]) != "" {
		projectPath = string(stateData[types.StateProjectPath])
	}
	d = append(d, []string{"Project path:", projectPath})

	if r.Flags.Status.Details {

		release, err := r.HelmClient.GetStatus()
		if err != nil {
			return err
		}
		d = append(d, []string{"Helm chart:", fmt.Sprintf("%s-%s", release.Chart.Name(), release.Chart.Metadata.Version)})
		d = append(d, []string{"Helm release:", release.Name})
		d = append(d, []string{"Helm release status", release.Info.Status.String()})

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
			status.PrintTableNoHeader(d)

			fmt.Println()

			status.PrintTable(
				[]string{"Name", "Condition", "Status"},
				data,
			)
			fmt.Println()
		} else {
			status.PrintTableNoHeader(d)
		}
		return nil
	}
	status.PrintTableNoHeader(d)

	return nil
}
