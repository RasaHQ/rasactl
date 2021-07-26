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

	"github.com/RasaHQ/rasaxctl/pkg/rasaxctl"
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile           string
	helmConfiguration *types.HelmConfigurationSpec = &types.HelmConfigurationSpec{}
	errorPrint        *color.Color                 = color.New(color.FgRed)
	rasaXCTL          *rasaxctl.RasaXCTL
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rasaxctl",
	Short: "A tools to manage Rasa X deployments",
	Long:  `rasaxctl helps you to manage Rasa X deployments.`,
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
	cobra.OnInitialize(initConfig)

	home, _ := homedir.Dir()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rasaxctl.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug output")
	rootCmd.PersistentFlags().String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
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

		// Search config in home directory with name ".rasaxctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rasaxctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
