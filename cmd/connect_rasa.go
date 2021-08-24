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
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/RasaHQ/rasactl/pkg/utils"
)

const (
	connectRasaDesc = `
Connect Rasa OSS (Open Source Server) to Rasa X deployment.

The command prepares a configuration that's required to connect Rasa X deployment and run a local Rasa server.

It's required to have the 'rasa' command accessible by rasactl.

The command works only if Rasa X deployment uses a local rasa project.
`

	connectRasaExample = `
	# Connect Rasa Server to Rasa X deployment.
	$ rasactl connect rasa

	# Run a separate rasa server for the Rasa X worker environment.
	$ rasactl connect rasa --run-separate-worker

	# Pass extra arguments to rasa server.
	$ rasactl connect rasa --extra-args="--debug"
`
)

func connectRasaCmd() *cobra.Command {

	// cmd represents the connect rasa command
	cmd := &cobra.Command{
		Use:     "rasa [DEPLOYMENT NAME]",
		Short:   "run Rasa OSS server and connect it to the Rasa X deployment",
		Long:    connectRasaDesc,
		Args:    cobra.MaximumNArgs(1),
		Example: templates.Examples(connectRasaExample),
		PreRunE: func(cmd *cobra.Command, args []string) error {

			if !utils.CommandExists("rasa") {
				return errors.Errorf(
					errorPrint.Sprint(
						"The 'rasa' command doesn't exist. Check out the docs to learn how to install rasa, https://rasa.com/docs/rasa/installation/",
					),
				)
			}

			if _, err := parseArgs(namespace, args, 1, 1, rasactlFlags); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if err := checkIfNamespaceExists(); err != nil {
				return err
			}

			stateData, err := rasaCtl.KubernetesClient.ReadSecretWithState()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			rasaCtl.HelmClient.SetConfiguration(
				&types.HelmConfigurationSpec{
					ReleaseName: string(stateData[types.StateHelmReleaseName]),
					ReuseValues: true,
					Timeout:     time.Minute * 10,
				},
			)

			rasaCtl.KubernetesClient.SetHelmReleaseName(string(stateData[types.StateHelmReleaseName]))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if !rasaCtl.KubernetesClient.IsNamespaceManageable() {
				return errors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasactl, can't continue :(", rasaCtl.Namespace))
			}

			// Check if a Rasa X deployment is already installed and running
			_, isRunning, err := rasaCtl.CheckDeploymentStatus()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isRunning {
				fmt.Printf("Rasa X for the %s deployment is not running.\n", rasaCtl.Namespace)
				return nil
			}

			if err := rasaCtl.ConnectRasa(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	addConnectRasaFlags(cmd)

	return cmd
}
