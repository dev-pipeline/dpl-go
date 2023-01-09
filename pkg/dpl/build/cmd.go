package build

import (
	"github.com/spf13/cobra"

	"github.com/dev-pipeline/dpl-go/cmd"
	icmd "github.com/dev-pipeline/dpl-go/internal/cmd"
	"github.com/dev-pipeline/dpl-go/internal/common"
)

var (
	buildCommon common.Args

	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build a dpl project",
		Run:   doBuildCmd,
	}
)

func doBuildCmd(cmd *cobra.Command, args []string) {
	common.DoCommand(args, buildCommon, []common.Task{
		BuildTask,
	})
}

func init() {
	icmd.AddCommonArgs(buildCmd, &buildCommon)
	cmd.AddCommand(buildCmd)
}
