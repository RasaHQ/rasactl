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
	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

const (
	modelDeleteDesc = `
`

	modelDeleteExample = `
`
)

func modelDeleteCmd() *cobra.Command {
	// cmd represents the status command
	cmd := &cobra.Command{
		Use:     "delete [DEPLOYMENT NAME] MODEL-NAME",
		Short:   "delete a model from Rasa X / Enterprise",
		Long:    modelDeleteDesc,
		Example: templates.Examples(modelDeleteExample),
		Args:    cobra.MaximumNArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {

			if err := checkIfNamespaceExists(); err != nil {
				return err
			}

			var modelName string
			if namespace == "" {
				return errors.Errorf(errorPrint.Sprint("You have to pass a deployment name"))
			} else if len(args) == 1 {
				modelName = args[0]
			} else if len(args) == 2 {
				modelName = args[1]
			}

			rasactlFlags.Model.Delete.Name = modelName

			stateData, err := rasaCtl.KubernetesClient.ReadSecretWithState()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			rasaCtl.HelmClient.Configuration = &types.HelmConfigurationSpec{
				ReleaseName: string(stateData[types.StateSecretHelmReleaseName]),
			}
			rasaCtl.KubernetesClient.Helm.ReleaseName = string(stateData[types.StateSecretHelmReleaseName])

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !rasaCtl.KubernetesClient.IsNamespaceManageable() {
				return errors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasactl, can't continue :(", rasaCtl.Namespace))
			}

			if err := rasaCtl.ModelDelete(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	return cmd
}
