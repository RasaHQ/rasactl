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
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func openCmd() *cobra.Command {

	// cmd represents the open command
	cmd := &cobra.Command{
		Use:   "open [PROJECT NAME]",
		Short: "Open Rasa X in a web browser",
		Args:  cobra.MinimumNArgs(1),
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

			err := rasaXCTL.CreateAndJoinKindNode()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

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

			url, err := rasaXCTL.GetRasaXURL()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if err := browser.OpenURL(url); err != nil {
				return errors.Errorf(errorPrint.Sprintf("Can't open the URL %s in your web browser: %s", url, err))
			}

			fmt.Printf("The URL %s has been opened in your web browser", url)

			return nil
		},
	}

	return cmd
}

func init() {

	openCmd := openCmd()
	rootCmd.AddCommand(openCmd)
}
