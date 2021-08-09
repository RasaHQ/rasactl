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

	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	deleteDesc = `
This command deletes a Rasa X deployment.

In a case where a Rasa X deployment is using kind and is deployed locally, you have to run the command with administrator permissions (sudo).
It's required because the command removes a hostname from /etc/hosts which was added during deployment creation.

`

	deleteExample = `
	# Delete the 'my-example' deployment.
	$ sudo rasactl delete my-example

	# Prune the 'my-example' deployment, execute the command with the --prune flag deletes the whole namespace.
	$ sudo rasactl delete my-example --prune
`
)

func deleteCmd() *cobra.Command {

	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:     "delete [DEPLOYMENT NAME]",
		Short:   "delete Rasa X deployment",
		Long:    deleteDesc,
		Example: examples(deleteExample),
		Aliases: []string{"del"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if namespace == "" {
				return errors.Errorf(errorPrint.Sprint("You have to pass a deployment name"))
			}

			if rasaCtl.KubernetesClient.BackendType == types.KubernetesBackendLocal {
				if os.Getuid() != 0 {
					return errors.Errorf(
						warnPrint.Sprintf(
							"Administrator permissions required, please run the command with sudo.\n%s needs administrator permissions to remove a hostname from /etc/hosts which was created during the creation of the deployment.",
							cmd.CommandPath(),
						),
					)
				}
			}

			stateData, err := rasaCtl.KubernetesClient.ReadSecretWithState()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			rasaCtl.HelmClient.Configuration = &types.HelmConfigurationSpec{
				ReleaseName: string(stateData[types.StateSecretHelmReleaseName]),
			}
			rasaCtl.KubernetesClient.Helm.ReleaseName = string(stateData[types.StateSecretHelmReleaseName])

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			isProjectExist, err := rasaCtl.KubernetesClient.IsNamespaceExist(rasaCtl.Namespace)
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isProjectExist {
				fmt.Printf("The %s project doesn't exist.\n", rasaCtl.Namespace)
				return nil
			}

			if !rasaCtl.KubernetesClient.IsNamespaceManageable() && !viper.GetBool("force") {
				return errors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasactl, can't continue :(", rasaCtl.Namespace))
			}

			if err := rasaCtl.Delete(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			defer rasaCtl.Spinner.Stop()
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
