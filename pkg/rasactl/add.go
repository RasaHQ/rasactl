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

	"github.com/RasaHQ/rasactl/pkg/utils"
)

// Add adds an existing deployment.
func (r *RasaCtl) Add() error {
	r.Log.Info("Adding existing deployment",
		"namespace", r.Namespace, "releaseName", r.HelmClient.GetConfiguration().ReleaseName)

	release, err := r.HelmClient.GetStatus()
	if err != nil {
		return err
	}

	if err := utils.HelmChartVersionConstrains(release.Chart.Metadata.Version); err != nil {
		return err
	}

	url, err := r.GetRasaXURL()
	if err != nil {
		return err
	}
	r.initRasaXClient()
	r.RasaXClient.URL = url

	rasaXVersion, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		return err
	}
	if err := r.KubernetesClient.SaveSecretWithState(""); err != nil {
		return err
	}
	if err := r.KubernetesClient.UpdateSecretWithState(rasaXVersion, release); err != nil {
		return err
	}

	if err := r.KubernetesClient.AddNamespaceLabel(); err != nil {
		return err
	}

	fmt.Printf("The %s has been added as a deployment.\n", r.Namespace)

	return nil
}
