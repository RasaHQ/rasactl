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

	"github.com/RasaHQ/rasaxctl/pkg/rasaxctl"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func upgradeCmd() *cobra.Command {

	// cmd represents the upgrade command
	cmd := &cobra.Command{
		Use:          "upgrade [PROJECT NAME]",
		Short:        "Upgrade/update Rasa X deployment",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			namespace := args[0]

			rasaXCTL = &rasaxctl.RasaXCTL{
				Namespace: namespace,
			}
			if err := rasaXCTL.InitClients(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			rasaXCTL.KubernetesClient.Helm.ReleaseName = helmConfiguration.ReleaseName
			rasaXCTL.HelmClient.Configuration = helmConfiguration

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			isProjectExist, err := rasaXCTL.KubernetesClient.IsNamespaceExist(rasaXCTL.Namespace)
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isProjectExist {
				fmt.Printf("The %s project doesn't exist.\n", rasaXCTL.Namespace)
				return nil
			}

			// Check if a Rasa X deployment is already installed and running
			_, isRunning, err := rasaXCTL.CheckDeploymentStatus()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isRunning {
				fmt.Printf("Rasa X for the %s project is not running.\n", rasaXCTL.Namespace)
				return nil
			}

			if err := rasaXCTL.Upgrade(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

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
