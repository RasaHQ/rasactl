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
	modelDeleteDesc = `
Delete a model from Rasa X / Enterprise.
`

	modelDeleteExample = `
	# Delete the 'example-model' model (use the currently active deployment).
	$ rasactl model delete example-model

	# Delete the 'example-model' model for the 'my-deployment' deployment.
	$ rasactl model delete my-deployment example-model
`
)

func modelDeleteCmd() *cobra.Command {
	// cmd represents the model delete command
	cmd := &cobra.Command{
		Use:     "delete [DEPLOYMENT-NAME] MODEL-NAME",
		Short:   "delete a model from Rasa X / Enterprise",
		Long:    templates.LongDesc(modelDeleteDesc),
		Example: templates.Examples(modelDeleteExample),
		Args:    cobra.RangeArgs(1, 2),
		Aliases: []string{"del"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkIfDeploymentsExist(); err != nil {
				return err
			}

			args, err := parseArgs(namespace, args, 1, 2, rasactlFlags)
			if err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if err := checkIfNamespaceExists(); err != nil {
				return err
			}

			rasactlFlags.Model.Delete.Name = args[1]

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

			if err := rasaCtl.ModelDelete(); err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	return cmd
}
