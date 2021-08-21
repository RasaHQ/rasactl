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
	"github.com/RasaHQ/rasactl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

const (
	startDesc = `This command creates a Rasa X deployment or starts stopped deployment if a given deployment already exists.

If the --project or --project-path is used, a Rasa X deployment will be using a local directory with Rasa project.

If a deployment name is not defined, a random name is generated and used as a deployment name.

If there is no existing deployment or you use the --project or --project-path flag a new deployment will be created,
otherwise, you have to use the --create flags to create a deployment.
`

	startExample = `
	# Create a Rasa X deployment.
	$ rasactl start

	# Create a Rasa X deployment with custom configuration, e.g the following configuration changes a Rasa X version.
	# All available values: https://github.com/RasaHQ/rasa-x-helm/blob/main/charts/rasa-x/values.yaml
	$ rasactl start --values-file custom-configuration.yaml

	# Create a Rasa X deployment with a defined password.
	$ rasactl start --rasa-x-password mypassword

	# Create a Rasa X deployment that uses a local Rasa project.
	# The command is executed in a Rasa project directory.
	$ rasactl start --project

	# Create a Rasa X deployment with a defined name.
	$ rasactl start my-deployment

	# Create a new deployment if there is already one or more deployments.
	# rasactl start --create
`
)

func startCmd() *cobra.Command {

	// cmd represents the start command
	cmd := &cobra.Command{
		Use:     "start [DEPLOYMENT NAME]",
		Short:   "start a Rasa X deployment",
		Long:    startDesc,
		Example: templates.Examples(startExample),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			rasaCtl.KubernetesClient.Helm.ReleaseName = helmConfiguration.ReleaseName
			rasaCtl.HelmClient.Configuration = helmConfiguration

			if rasactlFlags.Start.RasaXPasswordStdin {
				password, err := utils.GetPasswordStdin()
				if err != nil {
					return err
				}
				rasactlFlags.Start.RasaXPassword = password
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := parseArgs(args, 1, 1); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			// Get list of namespaces (deployments)
			namespaces, err := rasaCtl.KubernetesClient.GetNamespaces()
			if err != nil {
				return errors.Errorf(errorPrint.Sprint(err))
			}

			// Check if namespace exists only if the number of namespaces >= 2
			// and a new deployment wasn't not requested
			if len(namespaces) != 0 && !rasactlFlags.Start.Create &&
				!rasactlFlags.Start.Project && rasactlFlags.Start.ProjectPath == "" {
				if err := checkIfNamespaceExists(); err != nil {
					return err
				}
			}

			isDeployed, isRunning, err := rasaCtl.CheckDeploymentStatus()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isDeployed {
				if rasaCtl.KubernetesClient.BackendType == types.KubernetesBackendLocal &&
					rasaCtl.KubernetesClient.CloudProvider.Name == types.CloudProviderUnknown {
					if os.Getuid() != 0 {
						return errors.Errorf(
							warnPrint.Sprintf(
								"Administrator permissions required, please run the command with sudo.\n%s needs "+
									"administrator permissions to add a hostname to /etc/hosts so that a connection to your deployment is possible.",
								cmd.CommandPath(),
							),
						)
					}
				}
			}

			if isRunning {
				fmt.Printf("Rasa X for the %s namespace is running.\n", rasaCtl.HelmClient.Namespace)
				return nil
			}

			if err := rasaCtl.Start(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			defer rasaCtl.Spinner.Stop()
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
