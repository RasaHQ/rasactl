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

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"

	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/RasaHQ/rasactl/pkg/utils"
)

// Install prepares and executes the installation.
func (h *Helm) Install() error {

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

	// Creates a new install object with a given configuration.
	client := action.NewInstall(h.ActionConfig)
	client.Namespace = h.Namespace
	client.ReleaseName = h.Configuration.ReleaseName
	client.Description = "rasactl"
	client.Wait = true
	client.DryRun = false
	client.Timeout = h.Configuration.Timeout

	h.Log.V(1).Info("Helm client settings", "settings", client)

	helmChart, err := loader.Load(chartPath)
	if err != nil {
		return err
	}

	if req := helmChart.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(helmChart, req); err != nil {
			return err
		}
	}

	// Merge helm values - set erlang cookie for rabbitmq.
	h.Values = utils.MergeMaps(valuesRabbitMQErlangCookie(), h.Values)

	// Use the latest edge release for Rasa X
	if h.Flags.Start.UseEdgeRelease {
		h.Values = utils.MergeMaps(valuesUseEdgeReleaseRasaX(), h.Values)
	}

	// Add additional values for local PVC
	if (h.Flags.Start.ProjectPath != "" || h.Flags.Start.Project) &&
		h.KubernetesBackendType == types.KubernetesBackendLocal {
		h.Values = utils.MergeMaps(valuesMountHostPath(h.PVCName), h.Values)
		h.Values = utils.MergeMaps(valuesUseDedicatedKindNode(h.Namespace), h.Values)
		h.Log.V(1).Info("Merging values", "result", h.Values)
	}

	// Configure ingress to use local hostname if Kubernetes backend is on a local machine
	if h.KubernetesBackendType == types.KubernetesBackendLocal && h.CloudProvider.Name == types.CloudProviderUnknown {
		host := fmt.Sprintf("%s.%s", h.Namespace, types.RasaCtlLocalDomain)
		ip := "127.0.0.1"
		h.Values = utils.MergeMaps(
			valuesDisableRasaProduction(),
			valuesDisableNginx(),
			valuesSetupLocalIngress(host),
			h.Values,
		)
		h.Log.V(1).Info("Merging values", "result", h.Values)

		// Add host to /etc/hosts - required sudo
		if err := utils.AddHostToEtcHosts(host, ip); err != nil {
			return err
		}
		h.Log.V(1).Info("Adding host", "host", host, "ip", ip)
	} else if h.KubernetesBackendType == types.KubernetesBackendLocal &&
		h.CloudProvider.Name != types.CloudProviderUnknown {
		h.Values = utils.MergeMaps(valuesEnableRasaProduction(), valuesNginxNodePort(), h.Values)
	}

	// Set Rasa X password
	h.Values = utils.MergeMaps(valuesSetRasaXPassword(h.Flags.Start.RasaXPassword), h.Values)
	h.Log.V(1).Info("Merging values", "result", h.Values)

	// Install the chart
	rel, err := client.Run(helmChart, h.Values)
	if err != nil {
		return err
	}
	h.setCacheDirectory(cachePath)

	msg := fmt.Sprintf("Installation has beed finished, status: %s", rel.Info.Status)
	h.Log.Info(msg, "releaseName", client.ReleaseName, "namespace", client.Namespace)
	h.Log.V(1).Info(msg, "values", h.Values)
	h.Spinner.Message(msg)

	return nil
}
