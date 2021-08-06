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

	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	startDesc = `
This command creates a Rasa X deployment or starts stopped deployment if a given deployment already exists.

If the --project or --project-path is used, a Rasa X deployment will be using a local directory with Rasa project.

If a deployment name is not defined, a random name is generated and used as a deployment name.
`

	startExample = `
	# Create a Rasa X deployment.
	$ rasaxctl start

	# Create a Rasa X deployment with custom configuration, e.g the following configuration changes a Rasa X version.
	# All available values: https://github.com/RasaHQ/rasa-x-helm/blob/main/charts/rasa-x/values.yaml
	$ rasaxctl start --values-file custom-configuration.yaml

	# Create a Rasa X deployment with a defined password.
	$ rasaxctl start --rasa-x-password mypassword

	# Create a Rasa X deployment that uses a local Rasa project.
	# The command is executed in a Rasa project directory.
	$ rasaxctl start --project

	# Create a Rasa X deployment with a defined name.
	$ rasaxctl start my-deployment

`
)

func startCmd() *cobra.Command {

	// cmd represents the start command
	cmd := &cobra.Command{
		Use:          "start [DEPLOYMENT NAME]",
		Short:        "start a Rasa X deployment",
		Long:         startDesc,
		Example:      examples(startExample),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			rasaXCTL.KubernetesClient.Helm.ReleaseName = helmConfiguration.ReleaseName
			rasaXCTL.HelmClient.Configuration = helmConfiguration

			if rasaxctlFlags.Start.RasaXPasswordStdin {
				password, err := getRasaXPasswordStdin()
				if err != nil {
					return err
				}
				rasaxctlFlags.Start.RasaXPassword = password
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			isDeployed, isRunning, err := rasaXCTL.CheckDeploymentStatus()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isDeployed {
				if rasaXCTL.KubernetesClient.BackendType == types.KubernetesBackendLocal && rasaXCTL.KubernetesClient.CloudProvider.Name == types.CloudProviderUnknown {
					if os.Getuid() != 0 {
						return errors.Errorf(
							warnPrint.Sprintf(
								"Administrator permissions required, please run the command with sudo.\n%s needs administrator permissions to add a hostname to /etc/hosts so that a connection to your deployment is possible.",
								cmd.CommandPath(),
							),
						)
					}
				}
			}

			if isRunning {
				fmt.Printf("Rasa X for the %s namespace is running.\n", rasaXCTL.HelmClient.Namespace)
				return nil
			}

			if err := rasaXCTL.Start(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			defer rasaXCTL.Spinner.Stop()
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
