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
	configUseDeploymentExample = `
	# Set the 'example' deployment as the current deployment.
	$ rasactl config use-deployment example
`
)

func configUseDeploymentCmd() *cobra.Command {

	// cmd represents the status command
	cmd := &cobra.Command{
		Use:     "use-deployment DEPLOYMENT-NAME",
		Short:   "set the current-deployment in the configuration file",
		Long:    "Sets the current-deployment in the configuration file.",
		Example: templates.Examples(configUseDeploymentExample),
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return checkIfNamespaceExists()
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := rasaCtl.ConfigUseDeployment(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			fmt.Println("Done.")
			return nil
		},
	}

	return cmd
}
