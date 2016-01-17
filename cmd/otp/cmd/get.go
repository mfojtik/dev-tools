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

	"github.com/mfojtik/dev-tools/pkg/cmds"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the pull request that has the current branch",
	Long: `Get the pull request which is based on the current branch.

For example:

$ otp get # => 12345
`,
	Run: func(cmd *cobra.Command, args []string) {
		user, _ := cmd.Flags().GetString("user")
		isQuiet, _ := cmd.Flags().GetBool("number")
		if err := cmds.GetPullRequest(user, isQuiet); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to get the pull request for current branch, got %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("user", "u", os.Getenv("USER"), "Github username to use for search")
	getCmd.Flags().BoolP("number", "n", false, "Print only pull requests numbers")
}
