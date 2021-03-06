package helm

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
import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// Upgrade prepares and executes the upgrade.
func (h *Helm) Upgrade() error {

	err := h.updateRepository()
	if err != nil {
		return err
	}

	co := action.ChartPathOptions{
		InsecureSkipTLSverify: false,
		RepoURL:               h.Repositories[0].URL,
		Version:               h.Configuration.Version,
	}

	chartPath, err := co.LocateChart(h.RasaXChartName, h.Settings)
	if err != nil {
		return err
	}

	// Creates a new upgrade object with a given configuration.
	client := action.NewUpgrade(h.ActionConfig)
	client.Namespace = h.Namespace
	client.Description = "rasactl"
	client.Wait = true
	client.Timeout = h.Configuration.Timeout
	client.Atomic = h.Configuration.Atomic
	client.ReuseValues = h.Configuration.ReuseValues
	client.Version = h.Configuration.Version
	client.MaxHistory = 10

	// In a case where deployment was stopped, and it's started again,
	// the upgrade action is run with the current values defined for a helm release.
	if h.Configuration.StartProject {
		client.ReuseValues = true
	}

	h.Log.V(1).Info("Helm client settings", "settings", client)

	if err := h.ReadValuesFile(); err != nil {
		return err
	}

	helmChart, err := loader.Load(chartPath)
	if err != nil {
		return err
	}

	if req := helmChart.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(helmChart, req); err != nil {
			return err
		}
	}

	// Upgrade the chart
	rel, err := client.Run(h.Configuration.ReleaseName, helmChart, h.Values)
	if err != nil {
		return err
	}
	h.setCacheDirectory(cachePath)

	var msg string
	if !h.Configuration.StartProject {
		msg = fmt.Sprintf("Upgrade has beed finished, status: %s", rel.Info.Status)
	} else {
		msg = fmt.Sprintf("Rasa X for the %s deployment is ready", h.Namespace)
	}
	h.Log.Info(msg, "releaseName", rel.Name, "namespace", client.Namespace)
	h.Log.V(1).Info(msg, "values", h.Values)
	h.Spinner.Message(msg)

	return nil
}
