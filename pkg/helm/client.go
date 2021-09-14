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
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/RasaHQ/rasactl/pkg/status"
	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/RasaHQ/rasactl/pkg/utils/cloud"
)

const (
	cachePath            = "/tmp/.helmcache"
	repositoryConfigPath = "/tmp/.helmrepo"
)

type Interface interface {
	SetNamespace(namespace string) error
	GetNamespace() string
	Install() error
	Uninstall() error
	Upgrade() error
	ReadValuesFile() error
	GetAllValues() (map[string]interface{}, error)
	IsDeployed() (bool, error)
	GetStatus() (*release.Release, error)
	SetConfiguration(config *types.HelmConfigurationSpec)
	GetConfiguration() *types.HelmConfigurationSpec
	GetValues() map[string]interface{}
	SetValues(values map[string]interface{})
	SetKubernetesBackendType(backend types.KubernetesBackendType)
	SetPersistanceVolumeClaimName(name string)
}

// Helm represents a helm client.
type Helm struct {
	// Settings describes all of the environment settings.
	Settings *cli.EnvSettings

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

	driver   string
	debugLog func(format string, v ...interface{})

	// RasaXChartName defines a helm chart name to be used.
	RasaXChartName string

	kubeConfig  string
	kubeContext string

	// Values store helm values that are used by the client.
	Values map[string]interface{}

	// CloudProvider stores information about a cloud provider.
	CloudProvider *cloud.Provider

	// Flags stores command flags and their values.
	Flags *types.RasaCtlFlags
}

// New initializes a new helm client.
func New(client *Helm) (Interface, error) {
	var driverIsSet bool

	client.Settings = cli.New()
	client.Settings.RepositoryCache = cachePath
	client.Settings.RepositoryConfig = repositoryConfigPath
	client.RasaXChartName = "rasa-x"
	client.kubeConfig = viper.GetString("kubeconfig")

	client.kubeContext = viper.GetString("kube-context")
	if client.kubeContext != "" {
		client.Settings.KubeContext = client.kubeContext
	}

	client.Log.Info("Initializing Helm client")

	client.Repositories = append(client.Repositories, types.RepositorySpec{
		Name: "rasa-x",
		URL:  "https://rasahq.github.io/rasa-x-helm",
	})

	client.ActionConfig = new(action.Configuration)

	client.driver, driverIsSet = os.LookupEnv("HELM_DRIVER")
	if !driverIsSet {
		client.driver = "secrets"
	}

	client.debugLog = func(format string, v ...interface{}) {
		client.Log.Info(fmt.Sprintf(format, v...))
	}

	genericcliopts := &genericclioptions.ConfigFlags{
		Namespace:  &client.Namespace,
		KubeConfig: &client.kubeConfig,
		Context:    &client.kubeContext,
	}

	err := client.ActionConfig.Init(genericcliopts, client.Namespace, client.driver, client.debugLog)
	return client, err
}

// SetNamespace sets namespace for initialized client.
func (h *Helm) SetNamespace(namespace string) error {
	h.Namespace = namespace
	genericcliopts := &genericclioptions.ConfigFlags{
		Namespace:  &h.Namespace,
		KubeConfig: &h.kubeConfig,
		Context:    &h.kubeContext,
	}

	return h.ActionConfig.Init(genericcliopts, h.Namespace, h.driver, h.debugLog)
}

// GetNamespace returns namespace name.
func (h *Helm) GetNamespace() string {
	return h.Namespace
}

func (h *Helm) addRepository() ([]*repo.ChartRepository, error) {
	var repos []*repo.ChartRepository

	for _, repEntry := range h.Repositories {
		rep := repo.Entry{
			Name: repEntry.Name,
			URL:  repEntry.URL,
		}
		r, err := repo.NewChartRepository(&rep, getter.All(h.Settings))
		r.CachePath = h.Settings.RepositoryCache
		if err != nil {
			return nil, err
		}

		if _, err := r.DownloadIndexFile(); err != nil {
			err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", rep.URL)
			return nil, err
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

	g, _ := errgroup.WithContext(context.Background())

	for _, re := range repos {
		re := re // create a new 're', https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			if _, err := re.DownloadIndexFile(); err != nil {
				return errors.Wrapf(err, "...Unable to get an update from the %q chart repository (%s):\n\t%s\n",
					re.Config.Name, re.Config.URL, err)
			}
			return nil
		})
	}

	return g.Wait()
}
