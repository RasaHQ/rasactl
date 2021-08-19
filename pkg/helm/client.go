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
	"os"
	"sync"

	"github.com/RasaHQ/rasactl/pkg/status"
	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/RasaHQ/rasactl/pkg/utils/cloud"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	cachePath            = "/tmp/.helmcache"
	repositoryConfigPath = "/tmp/.helmrepo"
)

// Helm represents a helm client.
type Helm struct {
	settings *cli.EnvSettings

	// ActionConfig injects the dependencies that all actions share.
	ActionConfig *action.Configuration

	// Namespace is a namespace name that is used for the current client.
	Namespace string

	// PVCName defines a persistent volume claim name that is used to create a PVC.
	PVCName string

	// KubernetesBackendType defines a Kubernetes cluster type.
	KubernetesBackendType types.KubernetesBackendType

	// Repositories store slices of helm repository that used used by the client.
	Repositories []types.RepositorySpec

	// Configuration defines configuration for the client.
	Configuration *types.HelmConfigurationSpec

	// Spinner stores a spinner client.
	Spinner *status.SpinnerMessage

	// Log defines logger.
	Log logr.Logger

	driver         string
	debugLog       func(format string, v ...interface{})
	rasaXChartName string
	kubeConfig     string

	// Values store helm values that are used by the client.
	Values map[string]interface{}

	// CloudProvider stores information about a cloud provider.
	CloudProvider *cloud.Provider

	// Flags stores command flags and their values.
	Flags *types.RasaCtlFlags
}

// New initializes a new helm client.
func (h *Helm) New() error {
	var driverIsSet bool

	h.settings = cli.New()
	h.settings.RepositoryCache = cachePath
	h.settings.RepositoryConfig = repositoryConfigPath
	h.rasaXChartName = "rasa-x"
	h.kubeConfig = viper.GetString("kubeconfig")

	h.Log.Info("Initializing Helm client")
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

	err := h.ActionConfig.Init(genericcliopts, h.Namespace, h.driver, h.debugLog)
	return err
}

func (h *Helm) addRepository() ([]*repo.ChartRepository, error) {
	var repos []*repo.ChartRepository

	for _, repEntry := range h.Repositories {
		rep := repo.Entry{
			Name: repEntry.Name,
			URL:  repEntry.URL,
		}
		r, err := repo.NewChartRepository(&rep, getter.All(h.settings))
		r.CachePath = h.settings.RepositoryCache
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
				errResult = errors.Wrapf(err, "...Unable to get an update from the %q chart repository (%s):\n\t%s\n",
					re.Config.Name, re.Config.URL, err)
			}
		}(re)
	}
	wg.Wait()

	return errResult
}
