package helm

import (
	"fmt"
	"os"
	"sync"

	"github.com/RasaHQ/rasaxctl/pkg/status"
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils/cloud"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Helm struct {
	settings              *cli.EnvSettings
	ActionConfig          *action.Configuration
	Namespace             string
	PVCName               string
	KubernetesBackendType types.KubernetesBackendType
	Repositories          []types.RepositorySpec
	Configuration         *types.HelmConfigurationSpec
	Spinner               *status.SpinnerMessage
	Log                   logr.Logger
	driver                string
	debugLog              func(format string, v ...interface{})
	rasaXChartName        string
	kubeConfig            string
	values                map[string]interface{}
	CloudProvider         *cloud.Provider
}

func (h *Helm) New() error {
	var driverIsSet bool

	h.settings = cli.New()
	h.rasaXChartName = "rasa-x"
	h.kubeConfig = viper.GetString("kubeconfig")

	if err := h.ReadValuesFile(); err != nil {
		return err
	}

	h.Repositories = append(h.Repositories, types.RepositorySpec{
		Name: "rasa-x",
		URL:  "https://rasahq.github.io/rasa-x-helm",
	})

	h.ActionConfig = new(action.Configuration)

	h.driver, driverIsSet = os.LookupEnv("HELM_DRIVER")
	if !driverIsSet {
		h.driver = "secrets"
	}

	h.debugLog = func(format string, v ...interface{}) {
		h.Log.Info(fmt.Sprintf(format, v...))
	}

	genericcliopts := &genericclioptions.ConfigFlags{
		Namespace:  &h.Namespace,
		KubeConfig: &h.kubeConfig,
	}

	if err := h.ActionConfig.Init(genericcliopts, h.Namespace, h.driver, h.debugLog); err != nil {
		return err
	}

	h.Log.Info("Initializing Helm client")

	return nil
}

func (h *Helm) addRepository() ([]*repo.ChartRepository, error) {
	var repos []*repo.ChartRepository

	for _, repEntry := range h.Repositories {
		rep := repo.Entry{
			Name: repEntry.Name,
			URL:  repEntry.URL,
		}
		r, err := repo.NewChartRepository(&rep, getter.All(h.settings))
		if err != nil {
			return repos, err
		}

		if _, err := r.DownloadIndexFile(); err != nil {
			err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", rep.URL)
			return repos, err
		}
		repos = append(repos, r)
	}

	return repos, nil
}

func (h *Helm) updateRepository() error {
	repos, err := h.addRepository()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var errResult error
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				errResult = errors.Wrapf(err, "...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			}
		}(re)
	}
	wg.Wait()

	return errResult
}
