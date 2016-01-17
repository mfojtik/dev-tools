// Copyright © 2016 Michal Fojtik <mi@mifo.sk>
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
	"strconv"

	"github.com/mfojtik/dev-tools/pkg/cmds"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [pull-id]",
	Short: "Tag a pull request in origin for testing",
	Long: `Tag a single pull requrst in origin for testing in CI.

For example:

Run unit tests and integration tests for pull request #123:
$ otp test 123

Run the extended tests for pull request #123:
$ otp test 123 --extended

Run only extended tests that has "builds" in the description in "core" group:
$ otp test 123 --only-extended --focus "builds" --group "core"

`,
	PreRun: func(cmd *cobra.Command, args []string) {
		hasError := false
		flags := cmd.Flags()
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "Error: No pull-id specified, check usage.\n\n")
			hasError = true
		}
		isExtended, _ := flags.GetBool("extended")
		isOnlyExtended, _ := flags.GetBool("only-extended")
		focus, _ := flags.GetString("focus")
		group, _ := flags.GetString("group")
		if isExtended && isOnlyExtended {
			fmt.Fprintf(os.Stderr, "Error: Test must be extended or only-extended.\n\n")
			hasError = true
		}
		if len(focus) > 0 && !isOnlyExtended {
			fmt.Fprintf(os.Stderr, "Error: Focus is supported for only-extended.\n\n")
			hasError = true
		}
		if isOnlyExtended && len(group) == 0 {
			fmt.Fprintf(os.Stderr, "Error: Group must be set.\n\n")
			hasError = true
		}
		if hasError {
			cmd.Usage()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		number, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: The pull-id must be a number, got: %v\n", args[0])
		}

		flags := cmd.Flags()
		isExtended, _ := flags.GetBool("extended")
		isOnlyExtended, _ := flags.GetBool("only-extended")
		focus, _ := flags.GetString("focus")
		group, _ := flags.GetString("group")

		if err := cmds.AddTestComment(int(number), isExtended, isOnlyExtended, focus, group); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Adding merge comment failed, got %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)

	testCmd.Flags().BoolP("extended", "e", false, "Run extended tests")
	testCmd.Flags().BoolP("only-extended", "o", false, "Run only extended tests")
	testCmd.Flags().StringP("focus", "f", "", "Set ginkgo focus when running only extended tests")
	testCmd.Flags().StringP("group", "g", "core", "Set different group when running extended tests")
}
