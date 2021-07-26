package helm

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

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

	chartPath, err := co.LocateChart(h.rasaXChartName, h.settings)
	if err != nil {
		return err
	}

	client := action.NewUpgrade(h.ActionConfig)
	client.Namespace = h.Namespace
	client.Description = "rasaxctl"
	client.Wait = true
	client.Timeout = h.Configuration.Timeout
	client.Atomic = h.Configuration.Atomic
	client.ReuseValues = h.Configuration.ReuseValues

	if h.Configuration.StartProject {
		client.ReuseValues = true
	}

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

	// upgrade the chart
	rel, err := client.Run(h.Configuration.ReleaseName, helmChart, h.values)
	if err != nil {
		return err
	}

	var msg string
	if h.Configuration.StartProject {
		msg = fmt.Sprintf("Upgrade has beed finished, status: %s", rel.Info.Status)
	} else {
		msg = fmt.Sprintf("Rasa X for the %s project is ready", h.Namespace)
	}
	h.Log.Info(msg, "releaseName", rel.Name, "namespace", client.Namespace)
	h.Log.V(1).Info(msg, "values", h.values)
	h.Spinner.Message(msg)

	return nil
}
