// Copyright © 2016 Michal Fojtik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/mfojtik/dev-tools/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "otp",
	Short: "OpenShift tag pull request",
	Long: `OpenShift tag pull request is a developer helper tool to use
the different tags that OpenShift team has available to control the Jenkins CI.
`,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.otp.yaml)")
	apiKeyDesc := "github api key"
	if len(os.Getenv("GITHUB_API_KEY")) > 0 {
		apiKeyDesc += " (using GITHUB_API_KEY)"
	}
	RootCmd.PersistentFlags().StringP("api-key", "", "", apiKeyDesc)
	RootCmd.PersistentFlags().StringVarP(&api.OriginRepoName, "repository", "", api.OriginRepoName, "github repository name")
	RootCmd.PersistentFlags().StringVarP(&api.OriginRepoOwner, "owner", "", api.OriginRepoOwner, "github repository owner")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if len(cfgFile) != 0 {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".otp-config")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	apiKey, _ := RootCmd.Flags().GetString("api-key")
	if len(apiKey) == 0 && len(os.Getenv("GITHUB_API_KEY")) > 0 {
		RootCmd.Flags().Set("api-key", os.Getenv("GITHUB_API_KEY"))
	}
	if len(apiKey) > 0 {
		if err := os.Setenv("GITHUB_API_KEY", apiKey); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to set GITHUB_API_KEY\n")
			os.Exit(1)
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
