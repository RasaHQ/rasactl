package helm

import (
	"helm.sh/helm/v3/pkg/action"
)

func (h *Helm) Uninstall() error {

	client := action.NewUninstall(h.ActionConfig)
	client.Description = "rasaxctl"
	client.KeepHistory = false
	client.Timeout = h.Configuration.Timeout

	h.Log.V(1).Info("Helm client settings", "settings", client)

	// uninstall the chart
	rel, err := client.Run(h.Configuration.ReleaseName)
	if err != nil {
		return err
	}

	msg := "Uninstalling Rasa X"
	h.Log.Info(msg, "releaseName", rel.Release.Name, "namespace", h.Namespace)
	h.Spinner.Message(msg)

	return nil
}
