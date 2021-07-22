package helm

import (
	"io/ioutil"

	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/storage/driver"
	"sigs.k8s.io/yaml"
)

func (h *Helm) ReadValuesFile() error {
	file := viper.GetString("values-file")

	if file != "" {
		h.log.V(1).Info("Reading the values file", "file", file)
		valuesFile, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal([]byte(valuesFile), &h.values)
		if err != nil {
			return err
		}
		h.log.V(1).Info("Read values from the file",
			"file", file, "values", h.values,
		)
	}
	return nil
}

func (h *Helm) GetValues() (map[string]interface{}, error) {
	client := action.NewGetValues(h.ActionConfig)
	client.AllValues = true

	values, err := client.Run(h.Configuration.ReleaseName)
	if err != nil {
		h.log.V(1).Error(err, "Getting all values",
			"releaseName", h.Configuration.ReleaseName,
			"namespace", h.Namespace,
		)
		return values, err
	}
	h.log.V(1).Info("Getting all values",
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
