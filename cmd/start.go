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
package cmd

import (
	"fmt"
	"os"

	"github.com/RasaHQ/rasaxctl/pkg/helm"
	"github.com/RasaHQ/rasaxctl/pkg/k8s"
	"github.com/RasaHQ/rasaxctl/pkg/logger"
	"github.com/RasaHQ/rasaxctl/pkg/rasax"
	"github.com/RasaHQ/rasaxctl/pkg/status"
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils"
	"github.com/fatih/color"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	helmClient        *helm.Helm
	kubernetesClient  *k8s.Kubernetes
	spinnerMessage    *status.SpinnerMessage
	err               error
	helmConfiguration types.ConfigurationSpec = types.ConfigurationSpec{}
	errorPrint        *color.Color            = color.New(color.FgRed)
	log               logr.Logger
)

func startCmd() *cobra.Command {

	// cmd represents the start command
	cmd := &cobra.Command{
		Use:          "start [PROJECT NAME]",
		Short:        "Run Rasa X deployment",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			namespace := args[0]
			log = logger.New()
			spinnerMessage = &status.SpinnerMessage{}
			spinnerMessage.New()
			kubernetesClient = &k8s.Kubernetes{
				Namespace: namespace,
				Log:       log,
				Helm: k8s.HelmSpec{
					ReleaseName: helmConfiguration.ReleaseName,
				},
			}
			if err = kubernetesClient.New(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if kubernetesClient.BackendType == types.KubernetesBackendLocal {
				if os.Getuid() != 0 {
					return errors.Errorf(errorPrint.Sprint("Administrator permissions required, please run the command with sudo"))
				}
			}

			helmClient, err = helm.New(
				log,
				spinnerMessage,
				helmConfiguration,
				namespace,
			)
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			helmClient.KubernetesBackendType = kubernetesClient.BackendType

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			// Check if a Rasa X deployment is already installed and running
			isDeployed, err := helmClient.IsDeployed()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			isRasaXRunning, err := kubernetesClient.IsRasaXRunning()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if isDeployed && isRasaXRunning {
				fmt.Printf("Rasa X for the %s project is running.\n", helmClient.Namespace)
				return nil
			}

			// Install Rasa X
			if !isDeployed && !isRasaXRunning {
				spinnerMessage.Message("Deploying Rasa X")
				if viper.GetString("project-path") != "" {
					volume, err := kubernetesClient.CreateVolume(viper.GetString("project-path"))
					if err != nil {
						return errors.Errorf(errorPrint.Sprintf("%s", err))
					}
					helmClient.PVCName = volume
				}

				err = helmClient.Install()
				if err != nil {
					return errors.Errorf(errorPrint.Sprintf("%s", err))
				}
			} else if !isRasaXRunning {
				// Start Rasa X if deployments are scaled down to 0
				spinnerMessage.Message("Starting Rasa X")
			}

			allValues, err := helmClient.GetValues()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			kubernetesClient.Helm.Values = allValues
			url, err := kubernetesClient.GetRasaXURL()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			log.V(1).Info("Get Rasa X URL", "url", url)

			token, err := kubernetesClient.GetRasaXToken()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			// Wait for Rasa X to be ready
			rasaX := &rasax.RasaX{
				URL:            url,
				Token:          token,
				Log:            log,
				SpinnerMessage: spinnerMessage,
				WaitTimeout:    helmConfiguration.Timeout,
			}
			rasaX.New()
			err = rasaX.WaitForRasaX()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			spinnerMessage.Message("Ready!")

			if utils.IsDebugOrVerboseEnabled() {
				log.Info("Rasa X is ready", "url", url)
			}
			spinnerMessage.Stop()

			rasaXVersion, err := rasaX.GetVersionEndpoint()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isDeployed && !isRasaXRunning {
				// Print the status box only if it's a new Rasa X deployment
				status.PrintRasaXStatus(rasaXVersion, url)
			}

			return nil
		},
	}

	addStartUpgradeFlags(cmd)
	addStartFlags(cmd)

	return cmd
}

func init() {

	startCmd := startCmd()
	rootCmd.AddCommand(startCmd)
}
