package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	majorVersion = 0
	minorVersion = 1
	patchVersion = 0
)

var (
	fullVersion = fmt.Sprintf("%v.%v.%v", majorVersion, minorVersion, patchVersion)

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of dpl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("dev-pipeline %v\n", fullVersion)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
