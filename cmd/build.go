package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dev-pipeline/dpl-go/internal/build"
	"github.com/dev-pipeline/dpl-go/internal/common"
)

var (
	buildCommon common.Args

	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build a dpl project",
		Run:   doBuild,
	}
)

func doBuild(cmd *cobra.Command, args []string) {
	common.DoCommand(args, buildCommon, []common.Task{
		build.BuildTask,
	})
}

func init() {
	addCommonArgs(buildCmd, &buildCommon)
	rootCmd.AddCommand(buildCmd)
}
