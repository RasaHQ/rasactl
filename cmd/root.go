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

	"github.com/RasaHQ/rasactl/pkg/logger"
	"github.com/RasaHQ/rasactl/pkg/rasactl"
	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/RasaHQ/rasactl/pkg/utils"
	"github.com/RasaHQ/rasactl/pkg/version"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/fatih/color"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile           string
	helmConfiguration *types.HelmConfigurationSpec = &types.HelmConfigurationSpec{}
	errorPrint        *color.Color                 = color.New(color.FgRed)
	warnPrint         *color.Color                 = color.New(color.FgYellow)
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

		if len(args) != 0 {
			namespace = args[0]
		}

		if namespace == "" && cmd.CalledAs() == "start" {
			namespace = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
		}

		rasaCtl = &rasactl.RasaCtl{
			Log:       log,
			Namespace: namespace,
			Flags:     rasactlFlags,
		}
		if err := rasaCtl.InitClients(); err != nil {
			return errors.Errorf(errorPrint.Sprintf("%s", err))
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
	cobra.OnInitialize(initConfig, initLog)

	home, _ := homedir.Dir()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rasactl.yaml)")
	rootCmd.PersistentFlags().BoolVar(&rasactlFlags.Global.Verbose, "verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&rasactlFlags.Global.Debug, "debug", false, "enable debug output")
	rootCmd.PersistentFlags().String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
}

func initLog() {
	log = logger.New(rasactlFlags)
	namespace = utils.GetActiveNamespace(log)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
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
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
