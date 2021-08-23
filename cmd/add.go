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
)

const (
	addDesc = `
This command adds existing Rasa X deployment to rasactl.

If you already have a Rasa X deployment that uses the rasa-x-helm chart you can add the deployment and control it via rasactl.
`

	addExample = `
	# Add a Rasa X deployment that is deployed in the 'my-test' namespace.
	$ rasactl add my-test

	# Add a Rasa X deployment that is deployed in the 'my-test' namespace and
	# a helm release name for the deployment is 'rasa-x-example'.
	$ rasactl add my-test --rasa-x-release-name rasa-x-example
`
)

func addCmd() *cobra.Command {
	// cmd represents the add command
	cmd := &cobra.Command{
		Use:     "add NAMESPACE",
		Short:   "add existing Rasa X deployment",
		Long:    addDesc,
		Example: templates.Examples(addExample),
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			rasaCtl.KubernetesClient.SetHelmReleaseName(helmConfiguration.ReleaseName)
			rasaCtl.HelmClient.Configuration = helmConfiguration

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := parseArgs(args, 1, 1); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if err := checkIfNamespaceExists(); err != nil {
				return err
			}

			if rasaCtl.KubernetesClient.IsNamespaceManageable() {
				fmt.Println("Already added")
				return nil
			}

			if err := rasaCtl.Add(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	addAddFlags(cmd)

	return cmd
}

func init() {

	addCmd := addCmd()
	rootCmd.AddCommand(addCmd)
}
