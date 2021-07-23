package helm

import (
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils"
	"github.com/spf13/viper"
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

	// Add additional values for local PVC
	if viper.GetString("project-path") != "" {
		h.values = utils.MergeMaps(valuesMountHostPath(h.PVCName), h.values)
		h.Log.V(1).Info("Merging values", "result", h.values)
	}

	// Configure ingress to use local hostname if Kubernetes backend is on a local machine
	if h.KubernetesBackendType == types.KubernetesBackendLocal {
		host := fmt.Sprintf("%s.rasaxctl.local.io", h.Namespace)
		ip := "127.0.0.1"
		h.values = utils.MergeMaps(valuesDisableNginx(), valuesSetupLocalIngress(host), h.values)
		h.Log.V(1).Info("Merging values", "result", h.values)

		// Add host to /etc/hosts - required sudo
		if err := utils.AddHostToEtcHosts(host, ip); err != nil {
			return err
		}
		h.Log.V(1).Info("Adding host", "host", host, "ip", ip)
	}

	// Set Rasa X password
	h.values = utils.MergeMaps(valuesSetRasaXPassword(viper.GetString("rasa-x-password")), h.values)
	h.Log.V(1).Info("Merging values", "result", h.values)

	// install the chart
	rel, err := client.Run(helmChart, h.values)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Installation has beed finished, status: %s", rel.Info.Status)
	h.Log.Info(msg, "releaseName", client.ReleaseName, "namespace", client.Namespace)
	h.Log.V(1).Info(msg, "values", h.values)
	h.Spinner.Message(msg)

	return nil
}
