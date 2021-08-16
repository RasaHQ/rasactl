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
package helm

import (
	"fmt"
	"io/ioutil"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"sigs.k8s.io/yaml"
)

// ReadValuesFiles reads the value file and store values in the Helm.Values object.
func (h *Helm) ReadValuesFile() error {
	file := h.Flags.StartUpgrade.ValuesFile

	if file != "" {
		h.Log.V(1).Info("Reading the values file", "file", file)
		valuesFile, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal([]byte(valuesFile), &h.Values)
		if err != nil {
			return err
		}
		h.Log.V(1).Info("Read values from the file",
			"file", file, "values", h.Values,
		)
	}
	return nil
}

// GetValues returns values for the active helm release.
func (h *Helm) GetValues() (map[string]interface{}, error) {
	client := action.NewGetValues(h.ActionConfig)
	client.AllValues = true

	if h.Configuration == nil {
		return nil, fmt.Errorf("helm client requires to define a release name: %#v", h)
	}

	values, err := client.Run(h.Configuration.ReleaseName)
	if err != nil {
		h.Log.V(1).Error(err, "Getting all values",
			"releaseName", h.Configuration.ReleaseName,
			"namespace", h.Namespace,
			"helmClient", fmt.Sprintf("%#v", h),
		)
		return values, err
	}
	h.Log.V(1).Info("Getting all values",
		"releaseName", h.Configuration.ReleaseName,
		"namespace", h.Namespace,
		"values", values,
	)
	return values, err
}

// IsDeployed checks if a given helm release is deployed.
// Return 'true' if release is found.
func (h *Helm) IsDeployed() (bool, error) {
	client := action.NewHistory(h.ActionConfig)
	client.Max = 1

	_, err := client.Run(h.Configuration.ReleaseName)
	if err == driver.ErrReleaseNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// GetStatus returns a helm relese status.
func (h *Helm) GetStatus() (*release.Release, error) {
	client := action.NewStatus(h.ActionConfig)
	release, err := client.Run(h.Configuration.ReleaseName)
	if err != nil {
		return nil, err
	}

	return release, nil
}
