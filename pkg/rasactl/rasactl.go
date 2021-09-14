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
package rasactl

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	"github.com/RasaHQ/rasactl/pkg/docker"
	"github.com/RasaHQ/rasactl/pkg/helm"
	"github.com/RasaHQ/rasactl/pkg/k8s"
	"github.com/RasaHQ/rasactl/pkg/rasax"
	"github.com/RasaHQ/rasactl/pkg/status"
	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/RasaHQ/rasactl/pkg/utils/cloud"
)

// RasaCtl defines the rasactl client.
type RasaCtl struct {
	// KubernetesClient defines the Kubernetes client.
	KubernetesClient k8s.KubernetesInterface

	// HelmClient defines the helm client.
	HelmClient helm.Interface

	// RasaXClient defines the Rasa X client.
	RasaXClient *rasax.RasaX

	// DockerClient defines the Docker client.
	DockerClient docker.Interface

	// Log defines logger.
	Log logr.Logger

	// Spinner stores a spinner object.
	Spinner *status.SpinnerMessage

	// Namespace is a namespace name.
	Namespace string

	isRasaXRunning  bool
	isRasaXDeployed bool

	// CloudProvider stores a type of a detected cloud provider.
	CloudProvider *cloud.Provider

	// Flags stores the command flags.
	Flags *types.RasaCtlFlags
}

// InitClients initializes clients.
func (r *RasaCtl) InitClients() error {
	r.Spinner = status.NewSpinner()

	cloudProvider := &cloud.Provider{Log: r.Log}
	cloudProvider.New()
	r.CloudProvider = cloudProvider

	kubernetesClient, err := k8s.New(
		&k8s.Kubernetes{
			Namespace:     r.Namespace,
			Log:           r.Log,
			CloudProvider: r.CloudProvider,
			Flags:         r.Flags,
		},
	)
	if err != nil {
		return err
	}
	r.KubernetesClient = kubernetesClient

	helmClient, err := helm.New(
		&helm.Helm{
			Log:           r.Log,
			Namespace:     r.Namespace,
			Spinner:       r.Spinner,
			CloudProvider: r.CloudProvider,
			Flags:         r.Flags,
		},
	)
	if err != nil {
		return err
	}
	helmClient.SetKubernetesBackendType(r.KubernetesClient.GetBackendType())
	r.HelmClient = helmClient

	dockerClient, err := docker.New(
		&docker.Docker{
			Namespace: r.Namespace,
			Log:       r.Log,
			Spinner:   r.Spinner,
			Flags:     r.Flags,
		},
	)
	if err != nil {
		return err
	}
	r.DockerClient = dockerClient

	err = r.GetKindControlPlaneNodeInfo()
	return err
}

// SetNamespaceClients sets namespace for initialized clients.
func (r *RasaCtl) SetNamespaceClients(namespace string) error {
	r.Log.V(1).Info("Setting namespace for clients", "namespace", namespace)
	r.KubernetesClient.SetNamespace(namespace)
	r.DockerClient.SetNamespace(namespace)

	err := r.HelmClient.SetNamespace(namespace)
	return err
}

// CheckDeploymentStatus checks if a given deployment is deployed and running.
//
// It returns 'true' as the first value if the deployment is deployed, and 'true'
// as the second value if the deployment is running.
func (r *RasaCtl) CheckDeploymentStatus() (bool, bool, error) {
	// Check if a Rasa X deployment is already installed and running
	isRasaXDeployed, err := r.HelmClient.IsDeployed()
	if err != nil {
		return false, false, err
	}
	r.isRasaXDeployed = isRasaXDeployed

	isRasaXRunning, err := r.KubernetesClient.IsRasaXRunning()
	if err != nil {
		return false, false, err
	}

	r.isRasaXRunning = isRasaXRunning

	return isRasaXDeployed, isRasaXRunning, nil
}

func (r *RasaCtl) startOrInstall() error {
	projectPath := r.Flags.Start.ProjectPath
	// Install Rasa X
	if !r.isRasaXDeployed && !r.isRasaXRunning {
		if projectPath != "" || r.Flags.Start.Project {
			if r.DockerClient.GetKind().ControlPlaneHost != "" {
				if !r.Flags.Start.Project {
					// check if the project path exists
					if path, err := os.Stat(projectPath); err != nil {
						if os.IsNotExist(err) {
							return err
						} else if !path.IsDir() {
							return errors.Errorf("The %s path can't point to a file, it has to be a directory", projectPath)
						}
						return err
					}
					r.DockerClient.SetProjectPath(projectPath)
				} else {
					// use a current working directory
					wd, err := os.Getwd()
					if err != nil {
						return err
					}
					r.DockerClient.SetProjectPath(wd)
					projectPath = wd
				}

				r.Spinner.Message("Creating and joining a kind node")
				if err := r.CreateAndJoinKindNode(); err != nil {
					return err
				}
				volume, err := r.KubernetesClient.CreateVolume(projectPath)
				if err != nil {
					return err
				}
				r.HelmClient.SetPersistanceVolumeClaimName(volume)

			} else {
				return errors.Errorf("It looks like you don't use kind as a current Kubernetes context, the project-path flag is supported only with kind.")
			}

			if err := r.writeStatusFile(projectPath); err != nil {
				return err
			}
		}

		if err := r.KubernetesClient.SaveSecretWithState(projectPath); err != nil {
			return err
		}

		r.Spinner.Message("Deploying Rasa X")
		if err := r.HelmClient.Install(); err != nil {
			return err
		}
	} else if !r.isRasaXRunning {
		state, err := r.KubernetesClient.ReadSecretWithState()
		if err != nil {
			return err
		}
		// Start Rasa X if deployments are scaled down to 0
		msg := "Starting Rasa X"
		r.Spinner.Message(msg)
		r.Log.Info(msg)

		if string(state[types.StateProjectPath]) != "" {
			if r.DockerClient.GetKind().ControlPlaneHost != "" {
				nodeName := fmt.Sprintf("kind-%s", r.Namespace)
				if err := r.DockerClient.StartKindNode(nodeName); err != nil {
					return err
				}
			}
		}
		// Set configuration used for starting a stopped project.
		helmConfig := r.HelmClient.GetConfiguration()
		helmConfig.StartProject = true
		r.HelmClient.SetConfiguration(helmConfig)

		if err := r.HelmClient.Upgrade(); err != nil {
			return err
		}

		if err := r.KubernetesClient.ScaleUp(); err != nil {
			return err
		}
	}
	return nil
}

// GetAllHelmValues gets all evaluated values for a given helm release,
// and stores it in KubernetesClient.Helm.Values, HelmClient.Values.
func (r *RasaCtl) GetAllHelmValues() error {
	allValues, err := r.HelmClient.GetAllValues()
	if err != nil {
		return err
	}
	r.KubernetesClient.SetHelmValues(allValues)
	r.HelmClient.SetValues(allValues)

	return nil
}

// GetRasaXURL returns a Rasa X URL.
func (r *RasaCtl) GetRasaXURL() (string, error) {
	if err := r.GetAllHelmValues(); err != nil {
		return "", err
	}

	url, err := r.KubernetesClient.GetRasaXURL()
	if err != nil {
		return url, err
	}

	r.Log.V(1).Info("Get Rasa X URL", "url", url)
	return url, nil
}

func (r *RasaCtl) GetRasaXToken() (string, error) {
	token, err := r.KubernetesClient.GetRasaXToken()
	if err != nil {
		return token, err
	}

	return token, nil
}

func (r *RasaCtl) initRasaXClient() {
	url, _ := r.GetRasaXURL()
	r.RasaXClient = &rasax.RasaX{
		Log:            r.Log,
		SpinnerMessage: r.Spinner,
		WaitTimeout:    r.HelmClient.GetConfiguration().Timeout,
		Flags:          r.Flags,
	}
	r.RasaXClient.New()
	r.RasaXClient.URL = url
}

func (r *RasaCtl) checkDeploymentStatus() error {
	err := r.RasaXClient.WaitForRasaX()
	if err != nil {
		return err
	}

	r.Log.Info("Rasa X is ready", "url", r.RasaXClient.URL, "password", r.Flags.Start.RasaXPassword)
	r.Spinner.Message("Ready!")
	r.Spinner.Stop()

	rasaXVersion, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		return err
	}

	helmRelease, err := r.HelmClient.GetStatus()
	if err != nil {
		return err
	}

	if err := r.KubernetesClient.UpdateSecretWithState(rasaXVersion, helmRelease); err != nil {
		return err
	}

	if !r.isRasaXDeployed && !r.isRasaXRunning {
		// Print the status box only if it's a new Rasa X deployment
		status.PrintRasaXStatus(
			rasaXVersion,
			r.RasaXClient.URL,
			r.Flags,
		)
	}
	return nil
}
