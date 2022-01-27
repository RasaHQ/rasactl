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
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/go-logr/logr"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/RasaHQ/rasactl/pkg/logger"
	"github.com/RasaHQ/rasactl/pkg/rasactl"
	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/RasaHQ/rasactl/pkg/utils"
	"github.com/RasaHQ/rasactl/pkg/version"
)

var (
	cfgFile           string
	helmConfiguration *types.HelmConfigurationSpec = &types.HelmConfigurationSpec{}
	errorPrint        *color.Color                 = color.New(color.FgRed)
	rasaCtl           *rasactl.RasaCtl
	log               logr.Logger
	namespace         string
	rasactlFlags      *types.RasaCtlFlags = &types.RasaCtlFlags{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "rasactl",
	Short:   "rasactl provisions and manages Rasa X deployments.",
	Long:    `rasactl provisions and manages Rasa X deployments.`,
	Version: version.VERSION,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		rasaCtl = &rasactl.RasaCtl{
			Log:   log,
			Flags: rasactlFlags,
		}

		if !strings.Contains(cmd.CommandPath(), "help") && !strings.Contains(cmd.CommandPath(), "completion") {
			if err := rasaCtl.InitClients(); err != nil {
				return errors.Errorf(errorPrint.Sprintf("%s", err))
			}
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLog, initConfig, getNamespace)

	home, _ := homedir.Dir()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rasactl.yaml)")
	rootCmd.PersistentFlags().BoolVar(&rasactlFlags.Global.Verbose, "verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&rasactlFlags.Global.Debug, "debug", false, "enable debug output")
	rootCmd.PersistentFlags().String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	rootCmd.PersistentFlags().String("kube-context", "", "name of the kubeconfig context to use")

	//nolint:golint,errcheck
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	//nolint:golint,errcheck
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	//nolint:golint,errcheck
	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	//nolint:golint,errcheck
	viper.BindPFlag("kube-context", rootCmd.PersistentFlags().Lookup("kube-context"))
}

func initLog() {
	log = logger.New(rasactlFlags)
}

func getNamespace() {
	namespace = utils.GetActiveNamespace(log)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		if rasactlFlags.Config.CreateFile {
			f, _ := os.Create(cfgFile)
			f.Close()
		}

	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".rasactl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rasactl")

		if rasactlFlags.Config.CreateFile {

			file := strings.Join([]string{home, ".rasactl.yaml"}, "/")
			f, _ := os.Create(file)
			f.Close()
		}
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("rasactl")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config", "file", viper.ConfigFileUsed())
	}
}
