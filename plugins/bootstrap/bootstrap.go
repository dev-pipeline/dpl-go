package bootstrap

import (
	"github.com/spf13/cobra"

	"github.com/dev-pipeline/dpl-go/cmd"
	icmd "github.com/dev-pipeline/dpl-go/internal/cmd"
	"github.com/dev-pipeline/dpl-go/internal/common"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/build"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/scm"
)

var (
	bootstrapCommon common.Args

	bootstrapCmd = &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap a dpl project",
		Run:   doBootstrap,
	}
)

func doBootstrap(cmd *cobra.Command, args []string) {
	common.DoCommand(args, bootstrapCommon, []common.Task{
		scm.CheckoutTask,
		build.BuildTask,
	})
}

func init() {
	icmd.AddCommonArgs(bootstrapCmd, &bootstrapCommon)
	cmd.AddCommand(bootstrapCmd)
}
