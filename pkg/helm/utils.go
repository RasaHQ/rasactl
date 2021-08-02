package helm

import (
	"fmt"
	"io/ioutil"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"sigs.k8s.io/yaml"
)

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

func (h *Helm) GetStatus() (*release.Release, error) {
	client := action.NewStatus(h.ActionConfig)
	release, err := client.Run(h.Configuration.ReleaseName)
	if err != nil {
		return nil, err
	}

	return release, nil
}
