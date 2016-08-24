// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

	"github.com/mfojtik/dev-tools/pkg/server"
	"github.com/spf13/cobra"
)

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstraps the OpenShift development server",
	Long: `Bootstraps the OpenShift development server

For example:

$ otp bootstrap
`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("master-config-dir")
		netInterface, _ := cmd.Flags().GetString("net-interface")
		server, err := server.NewOpenShift(netInterface, dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize server: %q\n", err)
			cmd.Usage()
			os.Exit(1)
		}
		server, err = server.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run server: %q\n", err)
			os.Exit(1)
		}
		fmt.Println(server)
		server.Stop()
	},
}

func init() {
	bootstrapCmd.Flags().StringP("master-config-dir", "c", "/var/lib/openshift", "Master config directory")
	bootstrapCmd.Flags().StringP("net-interface", "n", "enp0s8", "Network interface to use")
	RootCmd.AddCommand(bootstrapCmd)
}
