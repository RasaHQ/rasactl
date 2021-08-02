package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

func addStartUpgradeFlags(cmd *cobra.Command) {
	cmd.Flags().DurationVar(&helmConfiguration.Timeout, "wait-timeout", time.Minute*10, "time to wait for Rasa X to be ready")
	cmd.Flags().StringVar(&helmConfiguration.Version, "rasa-x-chart-version", "", "a helm chart version to use")
	cmd.Flags().StringVar(&helmConfiguration.ReleaseName, "rasa-x-release-name", "rasa-x", "a helm release name to manage")

	cmd.PersistentFlags().StringVar(&rasaxctlFlags.StartUpgrade.ValuesFile, "values-file", "", "absolute path to the values file")
}

func addStartFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&rasaxctlFlags.Start.ProjectPath, "project-path", "", "absolute path to the project directory directory mounted in kind")
	cmd.PersistentFlags().BoolVarP(&rasaxctlFlags.Start.Project, "project", "p", false, "use the current working directory as a project directory, the flag is ignored if the --project-path flag is used")
	cmd.PersistentFlags().StringVar(&rasaxctlFlags.Start.RasaXPassword, "rasa-x-password", "rasaxlocal", "Rasa X password")
}

func addUpgradeFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&helmConfiguration.Atomic, "atomic", false, "if set, upgrade process rolls back changes made in case of failed upgrade")
	cmd.Flags().BoolVar(&helmConfiguration.ReuseValues, "reuse-values", true, "when upgrading, reuse the last release's values and merge in any overrides")
}

func addDeleteFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&rasaxctlFlags.Delete.Force, "force", false, "if true, delete resources and ignore errors")
	cmd.PersistentFlags().BoolVar(&rasaxctlFlags.Delete.Prune, "prune", false, "if true, delete a namespace with a project")
}

func addStatusFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&rasaxctlFlags.Status.Details, "details", "d", false, "show detailed information, such as running pods, helm chart status")
}

func addAddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&helmConfiguration.ReleaseName, "rasa-x-release-name", "rasa-x", "a helm release name to manage")
}

func addConnectRasaFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&rasaxctlFlags.ConnectRasa.Port, "port", "p", 5005, "port to run the Rasa server at")
	cmd.Flags().BoolVar(&rasaxctlFlags.ConnectRasa.RunSeparateWorker, "run-saparate-worker", false, "runs a separate Rasa server for the worker environment")
	cmd.Flags().StringSliceVar(&rasaxctlFlags.ConnectRasa.ExtraArgs, "extra-args", nil, "extra arguments for Rasa server")
}
