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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/RasaHQ/rasactl/pkg/types"
)

const (
	authLoginDesc = `
	Log in to Rasa X / Enterprise.
	It stores credentials in an external credentials store, such as the native keychain of the operating system.

	- Apple macOS Keychain Access for macOS
	- pass for Linux
	- Microsoft Windows Credential Manager for Windows
`

	authLoginExample = `
	# Login to Rasa X / Enterprise.
	$ rasactl auth login

	# Provide a password using STDIN.
	# You can login non-interactively by using the --password-stdin flag to provide a password through STDIN.
	# Using STDIN prevents the password from ending up in the shell’s history.
	$ rasactl auth login --username me --password-stdin
`
)

func authLoginCmd() *cobra.Command {

	// cmd represents the status command
	cmd := &cobra.Command{
		Use:     "login [DEPLOYMENT NAME]",
		Short:   "authorize rasactl to access the Rasa X / Enterprise with user credentials",
		Long:    templates.LongDesc(authLoginDesc),
		Example: templates.Examples(authLoginExample),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := parseArgs(namespace, args, 1, 1, rasactlFlags); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if err := checkIfNamespaceExists(); err != nil {
				return err
			}

			stateData, err := rasaCtl.KubernetesClient.ReadSecretWithState()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
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
			if !rasaCtl.KubernetesClient.IsNamespaceManageable() {
				return errors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasactl, can't continue :(", rasaCtl.Namespace))
			}

			// Check if a Rasa X deployment is already installed and running
			_, isRunning, err := rasaCtl.CheckDeploymentStatus()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isRunning {
				fmt.Printf("Rasa X for the %s deployment is not running.\n", rasaCtl.Namespace)
				return nil
			}

			if err := rasaCtl.AuthLogin(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	addAuthLoginFlags(cmd)

	return cmd
}
