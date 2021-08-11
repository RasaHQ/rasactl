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

func (r *RasaCtl) Start() error {

	r.Log.V(1).Info("Validating namespace name", "namespace", r.Namespace)
	if err := utils.ValidateName(r.HelmClient.Namespace); err != nil {
		return err
	}

	if err := r.KubernetesClient.CreateNamespace(); err != nil {
		return err
	}

	if err := r.KubernetesClient.AddNamespaceLabel(); err != nil {
		return err
	}

	// Init Rasa X client
	r.initRasaXClient()

	if err := r.startOrInstall(); err != nil {
		return err
	}

	url, err := r.GetRasaXURL()
	if err != nil {
		return err
	}
	r.RasaXClient.URL = url

	token, err := r.GetRasaXToken()
	if err != nil {
		return err
	}
	r.RasaXClient.Token = token

	if err := r.checkDeploymentStatus(); err != nil {
		return err
	}

	return nil
}