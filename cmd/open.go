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

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func openCmd() *cobra.Command {

	// cmd represents the open command
	cmd := &cobra.Command{
		Use:   "open [DEPLOYMENT NAME]",
		Short: "open Rasa X in a web browser",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := parseArgs(args, 1, 1); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if err := checkIfNamespaceExists(); err != nil {
				return err
			}
			rasaCtl.KubernetesClient.SetHelmReleaseName(helmConfiguration.ReleaseName)
			rasaCtl.HelmClient.Configuration = helmConfiguration

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if a Rasa X deployment is already installed and running
			_, isRunning, err := rasaCtl.CheckDeploymentStatus()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isRunning {
				fmt.Printf("The %s deployment is not running.\n", rasaCtl.Namespace)
				return nil
			}

			url, err := rasaCtl.GetRasaXURL()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if err := browser.OpenURL(url); err != nil {
				return errors.Errorf(errorPrint.Sprintf("Can't open the URL %s in your web browser: %s", url, err))
			}

			fmt.Printf("The %s URL has been opened in your web browser\n", url)
			rasaCtl.Spinner.Stop()
			return nil
		},
	}

	return cmd
}

func init() {

	openCmd := openCmd()
	rootCmd.AddCommand(openCmd)
}
