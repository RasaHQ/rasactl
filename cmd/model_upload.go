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
	modelUploadDesc = `
`

	modelUploadExample = `
`
)

func modelUploadCmd() *cobra.Command {
	// cmd represents the status command
	cmd := &cobra.Command{
		Use:     "upload [DEPLOYMENT NAME] MODEL-FILE",
		Short:   "upload model to Rasa X / Enterprise",
		Long:    templates.LongDesc(modelUploadDesc),
		Example: templates.Examples(modelUploadExample),
		Args:    maximumNArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkIfNamespaceExists(); err != nil {
				return err
			}

			var modelFile string
			if namespace == "" {
				return errors.Errorf(errorPrint.Sprint("You have to pass a deployment name"))
			} else if len(args) == 1 {
				modelFile = args[0]
			} else if len(args) == 2 {
				modelFile = args[1]
			}

			rasactlFlags.Model.Upload.File = modelFile

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

			if err := rasaCtl.ModelUpload(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	return cmd
}
