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
	modelTagDesc = `
	Rasa Enterprise allows multiple versions of an assistant to be run simultaneously and served to different users.
	By default, two environments are defined:

	- production
	- worker

	If you want to activate a model you have to tag it as 'production'.

	Learn more: https://rasa.com/docs/rasa-x/enterprise/deployment-environments/
`

	modelTagExample = `
	# Tag a model as 'production'
	$ rasactl model tag my-model production
`
)

func modelTagCmd() *cobra.Command {
	// cmd represents the status command
	cmd := &cobra.Command{
		Use:     "tag [DEPLOYMENT NAME] MODEL-NAME TAG",
		Short:   "tag a model in Rasa X / Enterprise",
		Long:    templates.LongDesc(modelTagDesc),
		Example: templates.Examples(modelTagExample),
		Args:    cobra.RangeArgs(2, 3),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			detectedNamespace := utils.GetActiveNamespace(log)
			modelName, modelTag, namespace, err := parseModelTagArgs(namespace, detectedNamespace, args)
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

			rasactlFlags.Model.Tag.Model = modelName
			rasactlFlags.Model.Tag.Name = modelTag

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

			if err := rasaCtl.ModelTag(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	return cmd
}
