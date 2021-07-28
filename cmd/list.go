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
	"github.com/RasaHQ/rasaxctl/pkg/rasaxctl"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {

	// cmd represents the open command
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list projects",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			rasaXCTL = &rasaxctl.RasaXCTL{}
			if err := rasaXCTL.InitClients(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := rasaXCTL.List(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}

			return nil
		},
	}

	return cmd
}

func init() {

	listCmd := listCmd()
	rootCmd.AddCommand(listCmd)
}
