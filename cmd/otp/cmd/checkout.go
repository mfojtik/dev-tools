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

// checkoutCmd represents the get command
var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Checkout the local branch associated with the pull request",
	Long: `Checkout the local branch associalted with the pull request.

For example:

$ otp checkout 1234
feature-branch
$
`,
	Run: func(cmd *cobra.Command, args []string) {
		user, _ := cmd.Flags().GetString("user")
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Must specify the pull request number")
			os.Exit(1)
		}
		if err := cmds.CheckoutPullRequest(user, args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to fetch pull request %q got %v\n", args[0], err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(checkoutCmd)

	checkoutCmd.Flags().StringP("user", "u", os.Getenv("USER"), "Github username to use for search")
}
