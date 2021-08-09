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

	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func upgradeCmd() *cobra.Command {

	// cmd represents the upgrade command
	cmd := &cobra.Command{
		Use:          "upgrade [DEPLOYMENT NAME]",
		Short:        "upgrade/update Rasa X deployment",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if namespace == "" {
				return errors.Errorf(errorPrint.Sprint("You have to pass a deployment name"))
			}

			stateData, err := rasaCtl.KubernetesClient.ReadSecretWithState()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			rasaCtl.HelmClient.Configuration = helmConfiguration
			rasaCtl.HelmClient.Configuration.ReleaseName = string(stateData[types.StateSecretHelmReleaseName])
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

			// Check if a Rasa X deployment is already installed and running
			_, isRunning, err := rasaCtl.CheckDeploymentStatus()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isRunning {
				fmt.Printf("Rasa X for the %s project is not running.\n", rasaCtl.Namespace)
				return nil
			}

			if err := rasaCtl.Upgrade(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
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
