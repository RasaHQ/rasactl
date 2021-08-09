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
	"bytes"
	"fmt"

	"github.com/RasaHQ/rasactl/pkg/status"
	"github.com/RasaHQ/rasactl/pkg/types"
)

func (r *RasaCtl) Status() error {
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
			fmt.Println()
		} else {
			fmt.Println(b.String())
		}
		return nil
	}
	fmt.Println(b.String())

	return nil
}
