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

	"github.com/RasaHQ/rasaxctl/pkg/rasaxctl"
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func startCmd() *cobra.Command {

	// cmd represents the start command
	cmd := &cobra.Command{
		Use:          "start [PROJECT NAME]",
		Short:        "Run Rasa X deployment",
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

			if rasaXCTL.KubernetesClient.BackendType == types.KubernetesBackendLocal {
				if os.Getuid() != 0 {
					return errors.Errorf(errorPrint.Sprint("Administrator permissions required, please run the command with sudo"))
				}
			}

			rasaXCTL.KubernetesClient.Helm.ReleaseName = helmConfiguration.ReleaseName
			rasaXCTL.HelmClient.Configuration = helmConfiguration

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if a Rasa X deployment is already installed and running
			_, isRunning, err := rasaXCTL.CheckDeploymentStatus()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if isRunning {
				fmt.Printf("Rasa X for the %s project is running.\n", rasaXCTL.HelmClient.Namespace)
				return nil
			}

			if err := rasaXCTL.Start(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			rasaXCTL.Spinner.Stop()
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
