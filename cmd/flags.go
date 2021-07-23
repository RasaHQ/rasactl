package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addStartUpgradeFlags(cmd *cobra.Command) {
	cmd.Flags().DurationVar(&helmConfiguration.Timeout, "wait-timeout", time.Minute*10, "time to wait for Rasa X to be ready")
	cmd.Flags().StringVar(&helmConfiguration.Version, "rasa-x-chart-version", "", "a helm chart version to use")
	cmd.Flags().StringVar(&helmConfiguration.ReleaseName, "rasa-x-release-name", "rasa-x", "a helm release name to manage")

	cmd.PersistentFlags().String("values-file", "", "absolute path to the values file")

	viper.BindPFlag("values-file", cmd.PersistentFlags().Lookup("values-file"))
}

func addStartFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("project-path", "", "absolute path to the project directory directory mounted in kind")
	cmd.PersistentFlags().String("rasa-x-password", "rasaxlocal", "Rasa X password")

	viper.BindPFlag("project-path", cmd.PersistentFlags().Lookup("project-path"))
	viper.BindPFlag("rasa-x-password", cmd.PersistentFlags().Lookup("rasa-x-password"))
}
