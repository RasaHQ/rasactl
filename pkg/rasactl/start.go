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

import "github.com/RasaHQ/rasactl/pkg/utils"

// Start a Rasa X / Enterprise deployment.
func (r *RasaCtl) Start() error {

	if err := utils.HelmChartVersionConstrains(
		r.HelmClient.GetConfiguration().Version,
	); err != nil {
		return err
	}

	dockerVersion, err := r.DockerClient.GetServerVersion()
	if err != nil {
		return err
	}
	if err := utils.DockerVersionConstrains(
		dockerVersion,
	); err != nil {
		return err
	}

	r.Log.V(1).Info("Validating namespace name", "namespace", r.Namespace)
	if err := utils.ValidateName(r.HelmClient.GetNamespace()); err != nil {
		return err
	}

	if err := r.KubernetesClient.CreateNamespace(); err != nil {
		return err
	}

	if err := r.KubernetesClient.AddNamespaceLabel(); err != nil {
		return err
	}

	if err := r.startOrInstall(); err != nil {
		return err
	}

	// Init Rasa X client
	r.initRasaXClient()

	token, err := r.GetRasaXToken()
	if err != nil {
		return err
	}
	r.RasaXClient.Token = token

	return r.checkDeploymentStatus()
}
