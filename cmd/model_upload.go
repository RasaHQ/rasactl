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
	"github.com/RasaHQ/rasactl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

const (
	modelUploadDesc = `
Upload a model to Rasa X / Enterprise.
`

	modelUploadExample = `
	# Upload the model.tar.gz model file to Rasa X / Enterprise.
	$ rasactl model upload model.tar.gz
`
)

func modelUploadCmd() *cobra.Command {
	// cmd represents the model upload command
	cmd := &cobra.Command{
		Use:     "upload [DEPLOYMENT NAME] MODEL-FILE",
		Short:   "upload model to Rasa X / Enterprise",
		Long:    templates.LongDesc(modelUploadDesc),
		Example: templates.Examples(modelUploadExample),
		Args:    cobra.RangeArgs(1, 2),
		Aliases: []string{"up"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			detectedNamespace := utils.GetActiveNamespace(log)
			modelFile, _, namespace, err := parseModelUpDownArgs(namespace, detectedNamespace, args)
			if err != nil {
				return err
			}
			if detectedNamespace != "" {
				rasaCtl.Namespace = namespace
				rasaCtl.KubernetesClient.Namespace = namespace
				rasaCtl.HelmClient.Namespace = namespace
				if err := rasaCtl.HelmClient.New(); err != nil {
					return err
				}
			}

			if err := checkIfNamespaceExists(); err != nil {
				return err
			}

			rasactlFlags.Model.Upload.File = modelFile

			stateData, err := rasaCtl.KubernetesClient.ReadSecretWithState()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
			rasaCtl.HelmClient.Configuration = &types.HelmConfigurationSpec{
				ReleaseName: string(stateData[types.StateHelmReleaseName]),
			}
			rasaCtl.KubernetesClient.Helm.ReleaseName = string(stateData[types.StateHelmReleaseName])

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !rasaCtl.KubernetesClient.IsNamespaceManageable() {
				return errors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasactl, can't continue :(", rasaCtl.Namespace))
			}

			if err := rasaCtl.ModelUpload(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	return cmd
}
