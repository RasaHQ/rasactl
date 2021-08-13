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
	"github.com/RasaHQ/rasactl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	modelDownloadDesc = `
`

	modelDownloadExample = `
`
)

func modelDownloadCmd() *cobra.Command {
	// cmd represents the status command
	cmd := &cobra.Command{
		Use:     "download [DEPLOYMENT NAME] MODEL-NAME [STORE-PATH]",
		Short:   "download a model from Rasa X / Enterprise",
		Long:    modelDownloadDesc,
		Example: examples(modelDownloadExample),
		Args:    maximumNArgs(3),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			detectedNamespace := utils.GetActiveNamespace(log)
			modelName, modelPath, namespace, err := parseModelDownloadArgs(namespace, detectedNamespace, args)
			if err != nil {
				return err
			}
			rasaCtl.Namespace = namespace

			isProjectExist, err := rasaCtl.KubernetesClient.IsNamespaceExist(rasaCtl.Namespace)
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isProjectExist {
				fmt.Printf("The %s project doesn't exist.\n", rasaCtl.Namespace)
				return nil
			}
			rasactlFlags.Model.Download.Name = modelName
			rasactlFlags.Model.Download.FilePath = modelPath

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

			if err := rasaCtl.ModelDownload(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	return cmd
}
