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
	"k8s.io/kubectl/pkg/util/templates"
)

const (
	modelDownloadDesc = `
Download a model from Rasa X / Enterprise to your local machine.
`

	modelDownloadExample = `
	# Download the 'model' model.
	# If the destination is not defined, the model will be stored in a current working directory.
	$ rasactl model download deployment-name model

	# Download the 'model' model and store it in the /tmp directory.
	$ rasactl model download deployment-name model /tmp/model.tar.gz
`
)

func modelDownloadCmd() *cobra.Command {
	// cmd represents the status command
	cmd := &cobra.Command{
		Use:     "download [DEPLOYMENT NAME] MODEL-NAME [DESTINATION]",
		Short:   "download a model from Rasa X / Enterprise",
		Long:    templates.LongDesc(modelDownloadDesc),
		Example: templates.Examples(modelDownloadExample),
		Args:    cobra.RangeArgs(1, 3),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := parseNamespaceModelDownloadCommand(args); err != nil {
				return err
			}

			detectedNamespace := utils.GetActiveNamespace(log)
			modelName, modelPath, namespace, err := parseModelUpDownArgs(namespace, detectedNamespace, args)
			if err != nil {
				return err
			}
			if detectedNamespace != "" {
				rasaCtl.Namespace = namespace
				if err := rasaCtl.SetNamespaceClients(namespace); err != nil {
					return err
				}
			}

			if err := checkIfNamespaceExists(); err != nil {
				return err
			}

			rasactlFlags.Model.Download.Name = modelName
			rasactlFlags.Model.Download.FilePath = modelPath

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
			// Check if a Rasa X deployment is running
			_, isRunning, err := rasaCtl.CheckDeploymentStatus()
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isRunning {
				fmt.Printf("The %s deployment is not running.\n", rasaCtl.Namespace)
				return nil
			}

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
