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
package types

import "time"

const (
	// HelmChartNameRasaX stores a name of helm chart used to deploy Rasa X / Enterprise.
	HelmChartNameRasaX string = "rasa-x"

	//HelmChartVersionRasaX storage a version of helm chart used to deploy Rasa X / Enterprise.
	HelmChartVersionRasaX string = "4.3.1"
)

// RepositorySpec stores data related to a helm repository.
type RepositorySpec struct {
	Name string
	URL  string
}

// HelmConfigurationSpec stores a configuration for the helm client.
type HelmConfigurationSpec struct {
	Timeout      time.Duration
	ReleaseName  string
	Version      string
	ReuseValues  bool
	StartProject bool
	Atomic       bool
	Wait         bool
}
