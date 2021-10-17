package cmd

import (
	"github.com/spf13/cobra"
)

var (
	buildCommon commonArgs

	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build a dpl project",
		Run:   doBuild,
	}
)

func doBuild(cmd *cobra.Command, args []string) {
}

func init() {
	addCommonArgs(buildCmd, &buildCommon)
	rootCmd.AddCommand(buildCmd)
}
