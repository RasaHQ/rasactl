package helm

import (
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

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

	chartPath, err := co.LocateChart(h.rasaXChartName, h.settings)
	if err != nil {
		return err
	}
	client := action.NewInstall(h.ActionConfig)
	client.Namespace = h.Namespace
	client.ReleaseName = h.Configuration.ReleaseName
	client.Description = "rasaxctl"
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

	h.Values = utils.MergeMaps(valuesDisableRasaProduction(), h.Values)

	// Add additional values for local PVC
	if (h.Flags.Start.ProjectPath != "" || h.Flags.Start.Project) && h.KubernetesBackendType == types.KubernetesBackendLocal {
		h.Values = utils.MergeMaps(valuesMountHostPath(h.PVCName), h.Values)
		h.Values = utils.MergeMaps(valuesUseDedicatedKindNode(h.Namespace), h.Values)
		h.Log.V(1).Info("Merging values", "result", h.Values)
	}

	// Configure ingress to use local hostname if Kubernetes backend is on a local machine
	if h.KubernetesBackendType == types.KubernetesBackendLocal && h.CloudProvider.Name == types.CloudProviderUnknown {
		host := fmt.Sprintf("%s.%s", h.Namespace, types.RasaXCtlLocalDomain)
		ip := "127.0.0.1"
		h.Values = utils.MergeMaps(valuesDisableNginx(), valuesSetupLocalIngress(host), h.Values)
		h.Log.V(1).Info("Merging values", "result", h.Values)

		// Add host to /etc/hosts - required sudo
		if err := utils.AddHostToEtcHosts(host, ip); err != nil {
			return err
		}
		h.Log.V(1).Info("Adding host", "host", host, "ip", ip)
	} else if h.KubernetesBackendType == types.KubernetesBackendLocal && h.CloudProvider.Name != types.CloudProviderUnknown {
		h.Values = utils.MergeMaps(h.Values, valuesEnableRasaProduction())
		h.Values = utils.MergeMaps(valuesNginxNodePort(), h.Values)
	}

	// Set Rasa X password
	h.Values = utils.MergeMaps(valuesSetRasaXPassword(h.Flags.Start.RasaXPassword), h.Values)
	h.Log.V(1).Info("Merging values", "result", h.Values)

	// install the chart
	rel, err := client.Run(helmChart, h.Values)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Installation has beed finished, status: %s", rel.Info.Status)
	h.Log.Info(msg, "releaseName", client.ReleaseName, "namespace", client.Namespace)
	h.Log.V(1).Info(msg, "values", h.Values)
	h.Spinner.Message(msg)

	return nil
}
