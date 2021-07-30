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
)

func statusCmd() *cobra.Command {

	// cmd represents the status command
	cmd := &cobra.Command{
		Use:   "status [DEPLOYMENT NAME]",
		Short: "show deployment status",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			if namespace == "" {
				return errors.Errorf(errorPrint.Sprint("You have pass a deployment name"))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			isProjectExist, err := rasaXCTL.KubernetesClient.IsNamespaceExist(rasaXCTL.Namespace)
			if err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if !isProjectExist {
				fmt.Printf("The %s project doesn't exist.\n", rasaXCTL.Namespace)
				return nil
			}

			if !rasaXCTL.KubernetesClient.IsNamespaceManageable() {
				return errors.Errorf(errorPrint.Sprintf("The %s namespace exists but is not managed by rasaxctl, can't continue :(", rasaXCTL.Namespace))
			}

			if err := rasaXCTL.Status(); err != nil {
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
