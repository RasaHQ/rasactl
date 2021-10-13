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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/RasaHQ/rasactl/pkg/types"
)

const (
	deleteDesc = `
This command deletes a Rasa X deployment.
`

	deleteExample = `
	# Delete the 'my-example' deployment.
	$ rasactl delete my-example

	# Prune the 'my-example' deployment, execute the command with the --prune flag deletes the whole namespace.
	$ rasactl delete my-example --prune
`
)

func deleteCmd() *cobra.Command {

	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:     "delete DEPLOYMENT-NAME",
		Short:   "delete Rasa X deployment",
		Long:    deleteDesc,
		Example: templates.Examples(deleteExample),
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"del"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
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
			if !rasaCtl.KubernetesClient.IsNamespaceManageable() && !viper.GetBool("force") {
				return errors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasactl, can't continue :(", rasaCtl.Namespace))
			}

			defer rasaCtl.Spinner.Stop()
			if err := rasaCtl.Delete(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			return nil
		},
	}

	addDeleteFlags(cmd)

	return cmd
}

func init() {

	deleteCmd := deleteCmd()
	rootCmd.AddCommand(deleteCmd)
}
