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

	"github.com/mfojtik/dev-tools/pkg/api"
	"github.com/mfojtik/dev-tools/pkg/cmds"
	"github.com/spf13/cobra"
)

var rebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuilds a Docker image with openshift binary",
	Long: `Rebuilds a Docker image with openshift binary.
The image can be S2I, Docker or base image.

For example:

$ otp rebuild --image=openshift/origin-sti-builder
`,
	Run: func(cmd *cobra.Command, args []string) {
		rebuildImages := []string{}
		image, _ := cmd.Flags().GetString("image")
		isBuilders, _ := cmd.Flags().GetBool("builders")
		if isBuilders && len(image) == 0 {
			rebuildImages = api.OriginBuilders
		}
		if len(image) > 0 {
			rebuildImages = append(rebuildImages, image)
		}
		if len(rebuildImages) == 0 {
			fmt.Fprintf(os.Stderr, "Error: You must provide the image to rebuild, see --image")
			cmd.Usage()
			os.Exit(1)
		}
		for _, imageName := range rebuildImages {
			fmt.Fprintf(os.Stdout, "Rebuilding %q ...\n", imageName)
			if err := cmds.RebuildOpenShiftImage(imageName); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Rebuilding %q failed, got %v\n", imageName, err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(rebuildCmd)
	rebuildCmd.Flags().StringP("image", "i", "", "Image to rebuild")
	rebuildCmd.Flags().BoolP("builders", "b", false, "Rebuild all known builders")
}
