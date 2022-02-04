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
	logsDesc = `
Print the logs for a container in a pod. If the pod has only one container, the container name is
optional.
`

	logsExample = `
	# Choose a pod and show logs for it (use the currently active deployment).
	$ rasactl logs

	# Show logs from pod rasa-x (use the currently active deployment).
	$ rasactl logs rasa-x

	# Show logs from pod rasa-x for the 'my-deployment' deployment.
	$ rasactl logs my-deployment rasa-x

	# Display only the most recent 10 lines of output in pod rasa-x
	$ rasactl logs rasa-x --tail=10

	# Return snapshot of previous terminated nginx container logs from pod rasa
	$ rasactl logs -p -c nginx rasa

	# Begin streaming the logs from pod rasa-x
  $ rasactl logs -f rasa-x
`
)

func logsCmd() *cobra.Command {
	parsedArgs := []string{}
	// cmd represents the logs command
	cmd := &cobra.Command{
		Use:     "logs [DEPLOYMENT-NAME] [POD]",
		Short:   "print the logs for a container in a pod",
		Long:    templates.LongDesc(logsDesc),
		Example: templates.Examples(logsExample),
		Args:    cobra.RangeArgs(0, 2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkIfDeploymentsExist(); err != nil {
				return err
			}

			args, err := parseArgs(namespace, args, 1, 2, rasactlFlags)
			if err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}
			parsedArgs = args

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

			if err := rasaCtl.Logs(parsedArgs); err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	logsFlags(cmd)

	return cmd
}

func init() {

	logsCmd := logsCmd()
	rootCmd.AddCommand(logsCmd)
}
