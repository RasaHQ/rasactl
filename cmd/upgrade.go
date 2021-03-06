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
	"github.com/RasaHQ/rasactl/pkg/utils"
)

const (
	upgradeDesc = `
This command upgrades Rasa X / Enterprise deployment.

The upgrade command upgrades or changes configuration for Rasa X / Enterprise deployment.

You can specify a values file with you custom configuration. The values file has the same form as a values file for helm chart.
Here you can find all available values that can be configured: https://github.com/RasaHQ/rasa-x-helm/blob/main/charts/rasa-x/values.yaml
`

	upgradeExample = `
	# Change configuration for Rasa X / Enterprise deployment by passing a custom configuration.
	$ rasactl upgrade my-deployment --values-file my-custom-values.yaml
`
)

func upgradeCmd() *cobra.Command {

	// cmd represents the upgrade command
	cmd := &cobra.Command{
		Use:     "upgrade [DEPLOYMENT-NAME]",
		Short:   "upgrade Rasa X deployment",
		Long:    upgradeDesc,
		Example: templates.Examples(upgradeExample),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			utils.CheckHelmChartDir()
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

			if helmConfiguration.Version == "" {
				helmConfiguration.Version = string(stateData[types.StateHelmChartVersion])
			}

			helmConfiguration.ReleaseName = string(stateData[types.StateHelmReleaseName])
			rasaCtl.HelmClient.SetConfiguration(helmConfiguration)
			rasaCtl.KubernetesClient.SetHelmReleaseName(string(stateData[types.StateHelmReleaseName]))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if a Rasa X deployment is already installed and running
			_, isRunning, err := rasaCtl.CheckDeploymentStatus()
			if err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isRunning {
				fmt.Printf("The %s deployment is not running.\n", rasaCtl.Namespace)
				return nil
			}

			if err := rasaCtl.Upgrade(); err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}
			rasaCtl.Spinner.Message("Ready!")
			defer rasaCtl.Spinner.Stop()
			return nil
		},
	}

	addStartUpgradeFlags(cmd)
	addUpgradeFlags(cmd)

	return cmd
}

func init() {

	upgradeCmd := upgradeCmd()
	rootCmd.AddCommand(upgradeCmd)
}
