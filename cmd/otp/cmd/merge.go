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
	"strconv"

	"github.com/mfojtik/dev-tools/pkg/cmds"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge [pull-id]",
	Short: "Tag a pull request in origin for merge",
	Long: `Tag a single pull request in origin for merge using CI.

For example:

$ otp merge 123 
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		hasError := false
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "Error: No pull-id specified, check usage.\n\n")
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
		if err := cmds.AddMergeComment(int(number)); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Adding merge comment failed, got %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(mergeCmd)
}
