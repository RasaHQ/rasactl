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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/RasaHQ/rasactl/pkg/types"
)

const (
	stopDesc = `
Stop a given running Rasa X deployment.
`

	stopExample = `
	# Stop a Rasa X deployment with the 'my-deployment' name.
	$ rasactl stop my-deployment

	# Stop a currently active Rasa X deployment.
	# The command stops the currently active deployment.
	# You can use the 'rasactl list' command to check which deployment is currently used.
	$ rasactl stop
`
)

func stopCmd() *cobra.Command {

	// cmd represents the stop command
	cmd := &cobra.Command{
		Use:     "stop [DEPLOYMENT-NAME]",
		Short:   "stop Rasa X deployment",
		Long:    stopDesc,
		Example: templates.Examples(stopExample),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkIfDeploymentsExist(); err != nil {
				return err
			}

			if _, err := parseArgs(namespace, args, 1, 1, rasactlFlags); err != nil {
				return err
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
				},
			)
			rasaCtl.KubernetesClient.SetHelmReleaseName(string(stateData[types.StateHelmReleaseName]))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if a Rasa X deployment is already installed and running
			_, isRunning, err := rasaCtl.CheckDeploymentStatus()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isRunning {
				fmt.Printf("The %s deployment is not running.\n", rasaCtl.Namespace)
				return nil
			}
			defer rasaCtl.Spinner.Stop()
			if err := rasaCtl.Stop(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			return nil
		},
	}

	return cmd
}

func init() {

	stopCmd := stopCmd()
	rootCmd.AddCommand(stopCmd)
}
