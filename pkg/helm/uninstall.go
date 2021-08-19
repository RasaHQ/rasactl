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
	"helm.sh/helm/v3/pkg/action"
)

// Uninstall prepares and executes the uninstall.
func (h *Helm) Uninstall() error {

	client := action.NewUninstall(h.ActionConfig)
	client.Description = "rasactl"
	client.KeepHistory = false
	client.Timeout = h.Configuration.Timeout

	h.Log.V(1).Info("Helm client settings", "settings", client)

	msg := "Uninstalling Rasa X"
	h.Spinner.Message(msg)

	// Uninstall the chart
	rel, err := client.Run(h.Configuration.ReleaseName)
	if err != nil {
		return err
	}
	h.setCacheDirectory(cachePath)

	h.Log.Info(msg, "releaseName", rel.Release.Name, "namespace", h.Namespace)

	return nil
}
