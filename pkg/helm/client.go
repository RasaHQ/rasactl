package helm

import (
	"fmt"
	"os"
	"sync"

	"github.com/RasaHQ/rasaxctl/pkg/status"
	"github.com/RasaHQ/rasaxctl/pkg/types"
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
	settings       *cli.EnvSettings
	ActionConfig   *action.Configuration
	Namespace      string
	Repositories   []types.RepositorySpec
	Configuration  types.ConfigurationSpec
	spinnerMessage *status.SpinnerMessage
	log            logr.Logger
	driver         string
	debugLog       func(format string, v ...interface{})
	rasaXChartName string
	kubeConfig     string
	values         map[string]interface{}
}

func New(log logr.Logger, spinnerMessage *status.SpinnerMessage, configuration types.ConfigurationSpec, namespace string) (*Helm, error) {
	var driverIsSet bool

	helmClient := &Helm{}
	helmClient.settings = cli.New()
	helmClient.log = log
	helmClient.spinnerMessage = spinnerMessage
	helmClient.rasaXChartName = "rasa-x"
	helmClient.kubeConfig = viper.GetString("kubeconfig")
	helmClient.Namespace = namespace
	helmClient.Configuration = configuration

	if err := helmClient.ReadValuesFile(); err != nil {
		return helmClient, err
	}

	helmClient.Repositories = append(helmClient.Repositories, types.RepositorySpec{
		Name: "rasa-x",
		URL:  "https://rasahq.github.io/rasa-x-helm",
	})

	helmClient.ActionConfig = new(action.Configuration)

	helmClient.driver, driverIsSet = os.LookupEnv("HELM_DRIVER")
	if !driverIsSet {
		helmClient.driver = "secrets"
	}

	helmClient.debugLog = func(format string, v ...interface{}) {
		log.Info(fmt.Sprintf(format, v...))
	}

	genericcliopts := &genericclioptions.ConfigFlags{
		Namespace:  &helmClient.Namespace,
		KubeConfig: &helmClient.kubeConfig,
	}

	if err := helmClient.ActionConfig.Init(genericcliopts, helmClient.Namespace, helmClient.driver, helmClient.debugLog); err != nil {
		return helmClient, err
	}

	log.Info("Initializing helm client")

	return helmClient, nil
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
