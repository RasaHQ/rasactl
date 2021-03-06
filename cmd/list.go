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
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

const (
	listDesc = `List all deployments.

The '*' in the 'CURRENT' field indicates a deployment that is used as default.
It means that every time when you execute 'rasactl' command without defining
the deployment name, the deployment marked with '*' is used.

A deployment is marked as 'CURRENT' if:

	- there is a '.rasactl' file that includes a deployment name in your current working directory.
	  The file is automatically created if you run the 'rasactl start' command with
		the '--project' or '--project-path' flag.
	- a default deployment is defined, e.g. via the 'rasactl config use-deployment' command.
  - there is only one deployment.

`
)

func listCmd() *cobra.Command {

	// cmd represents the list command
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list deployments",
		Long:    listDesc,
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := parseArgs(namespace, args, 0, 0, rasactlFlags); err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
			}

			if err := rasaCtl.List(); err != nil {
				return xerrors.Errorf(errorPrint.Sprintf("%s", err))
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
