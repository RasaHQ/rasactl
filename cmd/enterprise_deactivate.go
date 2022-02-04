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
	enterpriseDeactivateDesc = `
	Deactivate an Enterprise license.
`

	enterpriseDeactivateExample = `
	# Deactivate an Enterprise license (use the currently active deployment).
	$ rasactl enterprise deactivate

	# Deactivate an Enterprise license for the 'my-deployment' deployment.
	$ rasactl enterprise deactivate my-deployment
`
)

func enterpriseDeactivateCmd() *cobra.Command {
	// cmd represents the enterprise deactivate command
	cmd := &cobra.Command{
		Use:     "deactivate [DEPLOYMENT-NAME]",
		Short:   "deactivate an Enterprise license",
		Long:    templates.LongDesc(enterpriseDeactivateDesc),
		Example: templates.Examples(enterpriseDeactivateExample),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkIfDeploymentsExist(); err != nil {
				return err
			}

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
				return xerrors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasactl, can't continue :(", rasaCtl.Namespace))
			}

			// Check if a Rasa X deployment is already installed and running
			_, isRunning, err := rasaCtl.CheckDeploymentStatus()
			if err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isRunning {
				fmt.Printf("Rasa X for the %s deployment is not running.\n", rasaCtl.Namespace)
				return nil
			}

			if err := rasaCtl.EnterpriseDeactivate(); err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	return cmd
}
