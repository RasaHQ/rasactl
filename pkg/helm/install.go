package helm

import (
	"fmt"

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

	h.log.V(1).Info("Helm client settings", "settings", client)

	helmChart, err := loader.Load(chartPath)
	if err != nil {
		return err
	}

	if req := helmChart.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(helmChart, req); err != nil {
			return err
		}
	}

	// install the chart
	rel, err := client.Run(helmChart, h.values)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Installation has beed finished, status: %s", rel.Info.Status)
	h.log.Info(msg, "releaseName", client.ReleaseName, "namespace", client.Namespace)
	h.log.V(1).Info(msg, "values", h.values)
	h.spinnerMessage.Message(msg)

	return nil
}
