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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

const (
	statusDesc = `
Show status of the deployment.
`

	statusExample = `
	# Show status for the 'example' deployment.
	$ rasactl status example

	# Show status for the 'example' deployment along with details.
	$ rasactl status example --details

`
)

func statusCmd() *cobra.Command {

	// cmd represents the status command
	cmd := &cobra.Command{
		Use:     "status [DEPLOYMENT NAME]",
		Short:   "show deployment status",
		Long:    templates.LongDesc(statusDesc),
		Example: templates.Examples(statusExample),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := checkIfNamespaceExists()
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !rasaCtl.KubernetesClient.IsNamespaceManageable() {
				return errors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasactl, can't continue :(", rasaCtl.Namespace))
			}

			if err := rasaCtl.Status(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	addStatusFlags(cmd)

	return cmd
}

func init() {

	statusCmd := statusCmd()
	rootCmd.AddCommand(statusCmd)
}
