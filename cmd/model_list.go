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

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/RasaHQ/rasactl/pkg/types"
)

const (
	modelListDesc = `
	List all models stored in Rasa X / Enterprise.
`

	modelListExample = `
	# List all models (use the currently active deployment).
	$ rasactl model list

	# List all models for the 'my-deployment' deployment.
	$ rasactl model list my-deployment
`
)

func modelListCmd() *cobra.Command {
	// cmd represents the model list command
	cmd := &cobra.Command{
		Use:     "list [DEPLOYMENT-NAME]",
		Short:   "list models stored in Rasa X / Enterprise",
		Long:    templates.LongDesc(modelListDesc),
		Example: templates.Examples(modelListExample),
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"ls"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkIfDeploymentsExist(); err != nil {
				return err
			}

			if _, err := parseArgs(namespace, args, 1, 1, rasactlFlags); err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if err := checkIfNamespaceExists(); err != nil {
				return err
			}

			stateData, err := rasaCtl.KubernetesClient.ReadSecretWithState()
			if err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
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

			// Check if a Rasa X deployment is running
			_, isRunning, err := rasaCtl.CheckDeploymentStatus()
			if err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isRunning {
				fmt.Printf("The %s deployment is not running.\n", rasaCtl.Namespace)
				return nil
			}

			if !rasaCtl.KubernetesClient.IsNamespaceManageable() {
				return xerrors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasactl, can't continue :(", rasaCtl.Namespace))
			}

			if err := rasaCtl.ModelList(); err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	return cmd
}
