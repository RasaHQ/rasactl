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
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/RasaHQ/rasactl/pkg/types"
)

const (
	statusDesc = `
Show the status of a deployment.
`

	statusExample = `
	# Show status for the 'example' deployment.
	$ rasactl status example

	# Show status for the 'example' deployment along with details.
	$ rasactl status example --details

`
)

func statusCmd() *cobra.Command {

	// cmd represents the status command
	cmd := &cobra.Command{
		Use:     "status [DEPLOYMENT-NAME]",
		Short:   "show deployment status",
		Long:    templates.LongDesc(statusDesc),
		Example: templates.Examples(statusExample),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkIfDeploymentsExist(); err != nil {
				return err
			}

			if _, err := parseArgs(namespace, args, 1, 1, rasactlFlags); err != nil {
				return err
			}

			err := checkIfNamespaceExists()
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !rasaCtl.KubernetesClient.IsNamespaceManageable() {
				return xerrors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasactl, can't continue :(", rasaCtl.Namespace))
			}

			if rasaCtl.KubernetesClient.IsSecretWithStateExist() {
				stateData, err := rasaCtl.KubernetesClient.ReadSecretWithState()
				if err != nil {
					return xerrors.Errorf(errorPrint.Sprintf("%s", err))
				}

				helmConfiguration.ReleaseName = string(stateData[types.StateHelmReleaseName])
			}
			rasaCtl.KubernetesClient.SetHelmReleaseName(helmConfiguration.ReleaseName)
			rasaCtl.HelmClient.SetConfiguration(helmConfiguration)

			if err := rasaCtl.Status(); err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	addStatusFlags(cmd)

	return cmd
}

func init() {

	statusCmd := statusCmd()
	rootCmd.AddCommand(statusCmd)
}
