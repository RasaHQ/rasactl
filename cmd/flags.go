package cmd

import (
	"time"

	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/spf13/cobra"
)

func addStartUpgradeFlags(cmd *cobra.Command) {
	cmd.Flags().DurationVar(&helmConfiguration.Timeout, "wait-timeout", time.Minute*15, "time to wait for Rasa X to be ready")
	cmd.Flags().StringVar(&helmConfiguration.Version, "rasa-x-chart-version", types.HelmChartVersionRasaX, "a helm chart version to use")
	cmd.PersistentFlags().StringVar(&rasactlFlags.StartUpgrade.ValuesFile, "values-file", "", "absolute path to the values file")
}

func addStartFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&helmConfiguration.ReleaseName, "rasa-x-release-name", "rasa-x", "a helm release name to manage")

	cmd.PersistentFlags().StringVar(&rasactlFlags.Start.ProjectPath, "project-path", "",
		"absolute path to the project directory mounted in kind")

	cmd.PersistentFlags().BoolVarP(&rasactlFlags.Start.Project, "project", "p", false,
		"use the current working directory as a project directory, the flag is ignored if --project-path is used")

	cmd.PersistentFlags().StringVar(&rasactlFlags.Start.RasaXPassword, "rasa-x-password", "rasaxlocal", "Rasa X password")
	cmd.PersistentFlags().BoolVar(&rasactlFlags.Start.RasaXPasswordStdin, "rasa-x-password-stdin", false, "read the Rasa X password from stdin")
	cmd.Flags().BoolVar(&rasactlFlags.Start.UseEdgeRelease, "rasa-x-edge-release", false, "use the latest edge release of Rasa X")
	cmd.Flags().BoolVar(&rasactlFlags.Start.Create, "create", false,
		"create a new deployment. If --project or --project-path is set, or there is no existing deployment,"+
			" the flag is not required to create a new deployment")

}

func addUpgradeFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&helmConfiguration.Atomic, "atomic", false, "if set, upgrade process rolls back changes made in case of failed upgrade")
	cmd.Flags().BoolVar(&helmConfiguration.ReuseValues, "reuse-values", true,
		"when upgrading, reuse the last release's values and merge in any overrides")
}

func addDeleteFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&rasactlFlags.Delete.Force, "force", false, "if true, delete resources and ignore errors")
	cmd.PersistentFlags().BoolVar(&rasactlFlags.Delete.Prune, "prune", false, "if true, delete a namespace with a project")
}

func addStatusFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&rasactlFlags.Status.Details, "details", "d", false,
		"show detailed information, such as running pods, helm chart status")
	cmd.PersistentFlags().StringVarP(&rasactlFlags.Status.Output, "output", "o", "table",
		"output format. One of: json|table")
}

func addAddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&helmConfiguration.ReleaseName, "rasa-x-release-name", "rasa-x", "a helm release name to manage")
}

func addConnectRasaFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&rasactlFlags.ConnectRasa.Port, "port", "p", 5005, "port to run the Rasa server at")
	cmd.Flags().BoolVar(&rasactlFlags.ConnectRasa.RunSeparateWorker, "run-separate-worker", false,
		"runs a separate Rasa server for the worker environment")
	cmd.Flags().StringSliceVar(&rasactlFlags.ConnectRasa.ExtraArgs, "extra-args", nil, "extra arguments for Rasa server")
}

func addAuthLoginFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&rasactlFlags.Auth.Login.PasswordStdin, "password-stdin", false, "read the password from stdin")
	cmd.PersistentFlags().StringVarP(&rasactlFlags.Auth.Login.Username, "username", "u", "", "username")
	cmd.PersistentFlags().StringVarP(&rasactlFlags.Auth.Login.Password, "password", "p", "", "password")
}

func configFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&rasactlFlags.Config.CreateFile, "create-file", false, "create the configuration file if it doesn't exist")
}

func enterpriseActivateFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&rasactlFlags.Enterprise.Activate.LicenseStdin, "license-stdin", false, "read an Enterprise license from stdin")
	cmd.PersistentFlags().StringVarP(&rasactlFlags.Enterprise.Activate.License, "license", "l", "", "an Enterprise license")
}

func logsFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&rasactlFlags.Logs.Follow, "follow", "f", false, "specify if the logs should be streamed")
	cmd.PersistentFlags().BoolVarP(&rasactlFlags.Logs.Previous, "previous", "p", false, "print the logs for the previous instance of "+
		"the container in a pod if it exists")
	cmd.PersistentFlags().Int64Var(&rasactlFlags.Logs.TailLines, "tail", -1, "lines of recent log file to display. Defaults to -1 showing all log lines")
	cmd.PersistentFlags().StringVarP(&rasactlFlags.Logs.Container, "container", "c", "", "a container name")
}
